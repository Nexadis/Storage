package storage

import (
	"bufio"
	"fmt"
	"os"
)

var saveFormat = "%d\t%d\t%s\t%s"

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

type EventType byte

const (
	EventPut = iota
	EventDelete
)

type Event struct {
	ID        uint64
	EventType EventType
	Key       string
	Value     string
}

type FileTransactionLogger struct {
	events chan<- Event
	errors <-chan error
	lastID uint64
	file   *os.File
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- Event{
		EventType: EventPut,
		Key:       key,
		Value:     value,
	}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- Event{
		EventType: EventDelete,
		Key:       key,
	}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastID++
			_, err := fmt.Fprintf(
				l.file,
				saveFormat+"\n",
				l.lastID, e.EventType, e.Key, e.Value,
			)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)
	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if _, err := fmt.Sscanf(line, saveFormat,
				&e.ID, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}
			if l.lastID >= e.ID {
				outError <- fmt.Errorf("transaction numbers out of sequence")

				return
			}
			l.lastID = e.ID
			outEvent <- e
		}
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}

func RestoreTransactions(s Storage, l TransactionLogger) error {
	var err error

	events, errors := l.ReadEvents()
	e, ok := Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = s.Delete(e.Key)
			case EventPut:
				err = s.Put(e.Key, e.Value)
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("can't open tansaction log file: %w", err)
	}
	return &FileTransactionLogger{file: file}, nil
}
