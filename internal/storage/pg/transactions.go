package pg

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Nexadis/Storage/internal/storage"
)

var TransactionsScheme = `CREATE TABLE TRANSACTIONS(
ID SERIAL PRIMARY KEY,
USERNAME VARCHAR(256),
TYPE INT NOT NULL,
KEY VARCHAR(1024) NOT NULL,
VALUE VARCHAR(1024) NOT NULL
);
`

type PostgreTransactionLogger struct {
	events chan<- storage.Event
	errors <-chan error
	db     *sql.DB
}

func (p *PostgreTransactionLogger) WritePut(user, key, value string) {
	p.events <- storage.Event{
		User:      user,
		EventType: storage.EventPut,
		Key:       key,
		Value:     value,
	}
}

func (p *PostgreTransactionLogger) WriteDelete(user, key string) {
	p.events <- storage.Event{
		User:      user,
		EventType: storage.EventDelete,
		Key:       key,
	}
}

func (p *PostgreTransactionLogger) Err() <-chan error {
	return p.errors
}

func (p *PostgreTransactionLogger) ReadEvents() (<-chan storage.Event, <-chan error) {
	errors := make(chan error, 1)
	events := make(chan storage.Event)

	go func() {
		var e storage.Event

		defer close(errors)
		defer close(events)

		q, err := p.db.Prepare(`SELECT * FROM TRANSACTIONS`)
		if err != nil {
			errors <- err
			return
		}

		rows, err := q.Query()
		if err != nil {
			errors <- err
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&e.ID, &e.User, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				errors <- err
				return
			}
			log.Printf("Read Transaction %v", e)
			events <- e
		}

		err = rows.Err()
		if err != nil {
			errors <- err
			return
		}
	}()

	return events, errors
}

func (p *PostgreTransactionLogger) Run() {
	events := make(chan storage.Event, 16)
	p.events = events

	errors := make(chan error, 1)
	p.errors = errors

	go func() {
		for err := range errors {
			log.Printf("Logger err: %v", err)
		}
	}()
	insQ, err := p.db.Prepare("INSERT INTO TRANSACTIONS (\"username\",\"type\",\"key\",\"value\") VALUES($1,$2,$3,$4)")
	if err != nil {
		errors <- err
		return
	}
	delQ, err := p.db.Prepare("DELETE FROM TRANSACTIONS WHERE \"key\"=$1 AND \"username\"=$2")
	if err != nil {
		errors <- err
		return
	}
	go func() {
		defer close(errors)
		for e := range events {
			e := e
			switch {
			case e.EventType == storage.EventPut:
				log.Printf("Write Transaction %v", e)
				_, err := insQ.Exec(&e.User, &e.EventType, &e.Key, &e.Value)
				if err != nil {
					errors <- err
					return
				}
			case e.EventType == storage.EventDelete:
				log.Printf("Delete Transaction %v", e)
				_, err := delQ.Exec(&e.Key, &e.User)
				if err != nil {
					errors <- err
					return
				}
			}
		}
	}()
}

func NewPostgreTransactionLogger(conn string) (storage.TransactionLogger, error) {
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	pl := &PostgreTransactionLogger{
		db: db,
	}
	_, err = pl.db.Exec(TransactionsScheme)
	if err != nil {
		log.Printf("Error while scheme creating:%v", err)
	}
	return pl, nil
}
