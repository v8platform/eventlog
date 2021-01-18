package eventlog

import (
	"context"
	"time"
)

type EventsStream chan Event

type EventReader interface {
	SetOffset(offset int64) (int64, error)
	Offset() int64
	Read() *Event
	Stream(ctx context.Context, events EventsStream)
	StreamRead(ctx context.Context, bufSize ...int) EventsStream
}

type Event struct {
	Date                  time.Time
	TransactionStatus     TransactionStatusType
	TransactionDate       time.Time
	TransactionNumber     int64
	UserUuid              string
	User                  string
	Computer              string
	Application           ApplicationType
	Connection            int64
	Event                 EventType
	Severity              SeverityType
	Comment               string
	MetadataUuid          string
	Metadata              string
	Data                  interface{}
	DataPresentation      string
	Server                string
	MainPort              string
	AddPort               string
	Session               int64
	SessionDataSeparators []RefObject
}

type Objects interface {
	ReferencedObjectValue(objectType int, id ...int) (value, uuid string)
	ObjectValue(objectType int, id ...int) (value string)
}
