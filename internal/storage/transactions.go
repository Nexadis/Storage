package storage

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