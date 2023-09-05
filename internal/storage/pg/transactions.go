package pg

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Nexadis/Storage/internal/storage"
)

var TransactionsScheme = `CREATE TABLE TRANSACTIONS(
ID SERIAL PRIMARY KEY,
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

func (p *PostgreTransactionLogger) WritePut(key, value string) {
	p.events <- storage.Event{
		EventType: storage.EventPut,
		Key:       key,
		Value:     value,
	}
}

func (p *PostgreTransactionLogger) WriteDelete(key string) {
	p.events <- storage.Event{
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
			err := rows.Scan(&e.ID, &e.EventType, &e.Key, &e.Value)
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
	q, err := p.db.Prepare("INSERT INTO TRANSACTIONS (\"type\",\"key\",\"value\") VALUES($1,$2,$3)")
	if err != nil {
		errors <- err
		return
	}
	go func() {
		for e := range events {
			e := e
			log.Printf("Write Transaction %v", e)
			_, err := q.Exec(&e.EventType, &e.Key, &e.Value)
			if err != nil {
				errors <- err
				return
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
