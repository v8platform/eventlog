package eventlog

import (
	"time"
)

type EventsStream chan Event

type EventReader interface {
	Read() *Event
	StreamRead(bufSize ...int) EventsStream
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
	MainPort              string // TODO to int
	AddPort               string // TODO to int
	Session               int64
	SessionDataSeparators []RefObject
}

type Objects interface {
	ReferencedObjectValue(objectType int, id ...int) (value, uuid string)
	ObjectValue(objectType int, id ...int) (value string)
}
