package eventlog

import (
	"context"
	"errors"
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

const lgfFileName = "1Cv8.lgf"

type ManagerOptions struct {
	Folder             string
	PoolSize           int
	IdleCheckFrequency time.Duration
}

type journal struct {
	File   string
	Offset int64
	EventReader

	lgfReader *LgfReader
}

func NewManager(opt ManagerOptions) *Manager {

	p := &Manager{

		Folder:     opt.Folder,
		queue:      make(chan struct{}, opt.PoolSize),
		journals:   make(map[string]journal),
		muJournals: sync.Mutex{},
		stop:       make(chan struct{}),
	}

	//p.connsMu.Lock()
	//p.checkMinIdleConns()
	//p.connsMu.Unlock()

	if opt.IdleCheckFrequency > 0 {
		go p.reaper(opt.IdleCheckFrequency)
	}

	return p
}

type Manager struct {
	poolSize int    // Лимит обновременных экспортеров
	Folder   string // Каталог жерналов регистрации
	LiveMode bool

	journals   map[string]journal
	muJournals sync.Mutex

	exporters map[string]*Exporter

	queue chan struct{}

	storage ExporterStorage
	stop    chan struct{}
}

func (m *Manager) Watch(ctx context.Context, ticker time.Duration) error {

	fileWatcher := watcher.New()
	if err := fileWatcher.AddRecursive(m.Folder); err != nil {
		return err
	}
	fileWatcher.SetMaxEvents(1)
	fileWatcher.FilterOps(watcher.Create, watcher.Write, watcher.Remove)

	errChan := make(chan error)

	go func() {
		errChan <- fileWatcher.Start(ticker)
	}()

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			fileWatcher.Close()
			m.stopExporters()
			return ctx.Err()
		case e := <-fileWatcher.Event:
			switch e.Op {

			case watcher.Write, watcher.Create:
				m.addExporter(ctx, e)
			case watcher.Remove:
				m.removeExporter(ctx, e)
			}

		}
	}
}

func (m *Manager) stopExporters() {

}

func (m *Manager) processExporters() {

	for {
		select {
		case <-m.stop:

		case info := <-m.queue:

		}
	}
}

func (m *Manager) addExporter(ctx context.Context, event watcher.Event) {

	if info.IsDir() {
		return
	}

	fileSize := event.Size()
	fileName := event.Name()

	_ = m.waitTurn(ctx)

	m.muJournals.Lock()
	defer m.muJournals.Unlock()

	// TODO Подумать над циклом чтения и записи offset
	journ, ok := m.journals[fileName]
	_, exporter := m.exporters[fileName]

	if offset == fileSize {
		return
	}

}

var ErrLgfNotFound = errors.New("lgf not found")

func (m *Manager) doExporter(journal journal) error {

	lgfDir := filepath.Dir(journal.File)
	LgfFile := filepath.Join(lgfDir, lgfFileName)

	lgfStream, err := os.OpenFile(LgfFile, os.O_RDONLY, 644)
	if err != nil {
		return ErrLgfNotFound
	}

	lgpOpts := LgpReaderOptions{
		LgfDir:    lgfDir,
		LgfFile:   LgfFile,
		LgfStream: lgfStream,
		Offset:    journal.Offset,
	}

	reader, err := NewLgpReader(journal.File, lgpOpts)

	if err != nil {
		return err
	}

	exporter := NewExporter(reader, m.storage)
	exporter.Poller = &LongPoller{
		Limit:   10000,           // TODO
		Timeout: 1 * time.Second, // TODO
	}
	exporter.TZ = time.Local

	m.exporters[journal.File] = exporter

	go func() {

		_ = m.waitTurn()
		exporter.Start()
		m.freeTurn()

	}()

	return nil
}

func (m *Manager) removeExporter(ctx context.Context, info os.FileInfo) {
	m.muJournals.Lock()
	defer m.muJournals.Unlock()

}

func (m *Manager) getTurn() {
	m.queue <- struct{}{}
}

func (p *Manager) waitTurn(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.queue <- struct{}{}:
		return nil
	}
}

func (p *Manager) freeTurn() {
	<-p.queue
}

func (m *Manager) reaper(frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		select {
		case <-m.stop:
			log.Println("reaper stopped")
			break
		case <-ticker.C:

			log.Println("TODO reaper ticked")

			//_, err := p.ReapStaleConns()
			//if err != nil {
			//	continue
			//}
		}
	}
}
