package eventlog

import "time"

type ExporterStorage interface {
	Push(event Event)
}

type ExporterConfig struct {
	TZ      *time.Location // Временная зона для времени логово
	Timeout time.Duration  // timeout чтения
	Poller  Poller         // Читатель данных из файлов журнала регистрации

}

func NewExporter(eventReader EventReader, storage []ExporterStorage, config ...ExporterConfig) *Exporter {

	cfg := ExporterConfig{}

	if len(config) > 0 {
		cfg = config[0]
	}

	tz := time.Local
	if cfg.TZ != nil {
		tz = cfg.TZ
	}
	timeout := 1 * time.Second
	if cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}
	poller := cfg.Poller

	exporter := &Exporter{
		TZ:          tz,
		Timeout:     timeout,
		Events:      make(chan Event),
		Poller:      poller,
		eventReader: eventReader,
		storage:     storage,
		stop:        make(chan struct{}),
	}

	return exporter

}

type Exporter struct {
	File    string         // Файл логов
	TZ      *time.Location // Временная зона для времени логово
	Timeout time.Duration  // timeout чтения
	Events  chan Event
	Poller  Poller
	//Offset 		int64
	eventReader EventReader

	storage []ExporterStorage
	stop    chan struct{}
}

//LiveMode bool          // Чтение онлайн данных
func (e *Exporter) Start() {

	if e.Poller == nil {
		panic("exporter: can't start without a poller")
	}
	finished := make(chan struct{})
	stop := make(chan struct{})
	go e.Poller.Poll(e.eventReader, e.Events, stop)

	for {
		select {
		// handle incoming updates
		case event, closed := <-e.Events:
			if !closed {
				close(finished)
				close(stop)
				return
			}
			e.process(event)
		// call to stop polling
		case <-e.stop:
			close(finished)
			close(stop)

			// TODO Надо дочитать последню партию событий
			// TODO зафиксировать текущую позицию или не надо т.к. это можно зафиксировать на уровень уже
			return
		}
	}

}

// Stop gracefully shuts the poller down.
func (e *Exporter) Stop() error {
	e.stop <- struct{}{}
	err := e.eventReader.Close()
	return err
}

func (e *Exporter) process(event Event) {

	for _, storage := range e.storage {
		storage.Push(event)
	}

}
