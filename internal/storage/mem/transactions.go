package mem

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Nexadis/Storage/internal/storage"
)

var saveFormat = "%d\t%d\t%s\t%s"

type FileTransactionLogger struct {
	events chan<- storage.Event
	errors <-chan error
	lastID uint64
	file   *os.File
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- storage.Event{
		EventType: storage.EventPut,
		Key:       key,
		Value:     value,
	}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- storage.Event{
		EventType: storage.EventDelete,
		Key:       key,
	}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) Run() {
	events := make(chan storage.Event, 16)
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

func (l *FileTransactionLogger) ReadEvents() (<-chan storage.Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan storage.Event)
	outError := make(chan error, 1)
	go func() {
		var e storage.Event

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

func NewFileTransactionLogger(filename string) (storage.TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("can't open tansaction log file: %w", err)
	}
	return &FileTransactionLogger{file: file}, nil
}
