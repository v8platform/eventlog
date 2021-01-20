package eventlog

import (
	"context"
	"io"
	"time"
)

type EventReader interface {
	io.Closer
	seekReader
	Read(limit int, timeout time.Duration) (items []Event, err error)
}

type CtxEventReader interface {
	io.Closer
	seekReader
	ReadCtx(ctx context.Context, limit int, timeout time.Duration) (items []Event, err error)
}

type seekReader interface {
	Offset() int64
	Seek(offset int64) (int64, error)
}

type Event struct {
	Date              time.Time
	TransactionStatus TransactionStatusType
	TransactionDate   time.Time
	TransactionNumber int64
	UserUuid          string
	User              string
	Computer          string
	Application       ApplicationType
	Connection        int64
	Event             EventType
	Severity          SeverityType
	Comment           string
	MetadataUuid      string
	Metadata          string
	Data              interface{}
	DataPresentation  string
	Server            string
	MainPort          string
	AddPort           string
	Session           int64
	//SessionDataSeparators []

	Offset int64
}

type Objects interface {
	ReferencedObjectValue(objectType int, id ...int) (value, uuid string)
	ObjectValue(objectType int, id ...int) (value string)
}

var empty = struct{}{}

type EventManager struct {
	Events chan Event
	Poller Poller
	reader EventReader

	stop chan struct{}
}

func (m *EventManager) Start() {
	if m.Poller == nil {
		panic("manager: can't start without a poller")
	}
}

func (m *EventManager) Stop() {
	m.stop <- empty
}

func (m *EventManager) Poll(dest chan Event, stop chan struct{}) {
	if m.Poller == nil {
		panic("manager: can't start without a poller")
	}
}
