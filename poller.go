package eventlog

import (
	"github.com/xelaj/go-dry"
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
}

// Poll does long polling.
func (p *LongPoller) Poll(r EventReader, dest chan Event, stop chan struct{}) {

	for {
		select {
		case <-stop:
			return
		default:
		}

		events, err := r.Read(p.Limit, p.Timeout)
		if err != nil {
			continue
		}

		for _, event := range events {

			if !dry.SliceContains(p.AllowedSeverity, event.Severity) {
				continue
			}

			dest <- event
		}
	}
}
