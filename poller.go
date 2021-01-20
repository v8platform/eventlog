package eventlog

import (
	"github.com/xelaj/go-dry"
	"io"
	"time"
)

type Poller interface {
	Poll(r EventReader, dest chan Event, stop chan struct{})
}

type LongPoller struct {
	Limit   int
	Timeout time.Duration

	// AllowedSeverity contains the event types
	// Possible values:
	//	SeverityInfo  => "Информация"
	//	SeverityError => "Ошибка"
	//	SeverityWarn  => "Предупреждение"
	//	SeverityNote  => "Примечание"
	//
	AllowedSeverity []SeverityType

	LastErr error
}

// Poll does long polling.
func (p *LongPoller) Poll(r EventReader, dest chan Event, stop chan struct{}) {

	p.LastErr = nil

	for {

		select {
		case <-stop:
			return
		default:
		}

		events, err := r.Read(p.Limit, p.Timeout)
		p.LastErr = err

		if err != nil && err != io.EOF {
			return
		}

		p.pushEvents(events, dest)

		if p.LastErr == io.EOF {
			return
		}
	}
}

func (p *LongPoller) pushEvents(events []Event, dest chan Event) {

	for _, event := range events {

		if len(p.AllowedSeverity) > 0 &&
			!dry.SliceContains(p.AllowedSeverity, event.Severity) {
			continue
		}

		dest <- event
	}
}
