package eventlog

import (
	"context"
	"errors"
	"github.com/radovskyb/watcher"
	"github.com/xelaj/go-dry"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

/*
Логика работы менеджера жерналов регистрации 1С. Предприятие

1. Добавление каталога с файлом
2. Созмадие вотчера на файлы каталога
3. Получение событий от вотчера
4. При получении события создание читателя файлов
5. Проверка перед запуском, что такой читатель еще не работает.
6. Если читатель работает то постановка в очередь запуска с проверкой оффсета для файла
7. По завершенни работы сохранение читателя в базу или кэш
8. Цикл начиная с пункта 3
9. Безконечный цикл проверки читателей на актуальность ??? надо ли


*/
const lgfFileName = "1Cv8.lgf"

type ManagerOptions struct {
	Timeout            time.Duration
	Folder             []string
	PoolSize           int
	IdleCheckFrequency time.Duration
	JournalStorage     JournalStorage
	Exporters          []ExporterStorage
	BulkSize           int
}

func createLgpExporter(file string, offset int64, storage []ExporterStorage, poller Poller, tz *time.Location) (*Exporter, error) {

	lgfDir := filepath.Dir(file)
	LgfFile := filepath.Join(lgfDir, lgfFileName)
	lgfStream, err := os.OpenFile(LgfFile, os.O_RDONLY, 644)

	if err != nil {
		return nil, ErrLgfNotFound
	}

	lgpOpts := LgpReaderOptions{
		LgfDir:    lgfDir,
		LgfFile:   LgfFile,
		LgfStream: lgfStream,
		Offset:    offset,
	}

	reader, err := NewLgpReader(file, lgpOpts)

	if err != nil {
		return nil, err
	}

	exporter := NewExporter(reader, storage)
	exporter.Poller = poller
	exporter.TZ = tz

	return exporter, nil

}

func NewManager(ctx context.Context, opt ManagerOptions) *Manager {

	p := &Manager{
		queue:       make(chan struct{}, opt.PoolSize),
		fileWatcher: watcher.New(),
		exporters:   map[string]*Exporter{},
		mu:          sync.Mutex{},
		stop:        make(chan struct{}),
		journals:    NewInMemoryJournal(),
	}

	if opt.JournalStorage != nil {
		p.journals = opt.JournalStorage
	}

	go p.process(ctx)

	return p
}

type JournalStorage interface {
	GetOffset(file string) int64
	SetOffset(file string, off int64)
}

var _ JournalStorage = (*InMemoryJournal)(nil)

func NewInMemoryJournal() *InMemoryJournal {
	return &InMemoryJournal{
		data: &sync.Map{},
	}
}

type InMemoryJournal struct {
	data *sync.Map
}

func (i InMemoryJournal) GetOffset(file string) int64 {
	if value, ok := i.data.Load(file); ok {
		return value.(int64)
	}
	return 0
}

func (i InMemoryJournal) SetOffset(file string, off int64) {
	i.data.Store(file, off)
}

func extFilterHook(ext ...string) watcher.FilterFileHookFunc {
	return func(info os.FileInfo, fullPath string) error {

		if !info.IsDir() && dry.StringListContains(ext, filepath.Ext(info.Name())) {
			return nil
		}

		// No match.
		return watcher.ErrSkip
	}
}

var filterFiles = extFilterHook(".lgp", ".lgd")

// Manager основной объект выполнения чтения и экспорта журналов регистрации
type Manager struct {
	poolSize int      // Лимит обновременных экспортеров
	Folder   []string // Каталог жерналов регистрации
	LiveMode bool     // Auto use live mode
	BulkSize int
	Timeout  time.Duration
	Ticker   time.Duration
	TZ       *time.Location

	fileWatcher *watcher.Watcher

	journals  JournalStorage
	mu        sync.Mutex
	exporters map[string]*Exporter

	queue chan struct{}

	storage []ExporterStorage
	stop    chan struct{}

	running bool
}

func (m *Manager) Stop() {
	close(m.stop)
	m.fileWatcher.Close()
}

func (m *Manager) Running() bool {
	return m.running
}

func (m *Manager) Watch(folder string) error {

	if err := m.fileWatcher.AddRecursive(folder); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Unwatch(folder string) error {

	if err := m.fileWatcher.RemoveRecursive(folder); err != nil {
		return err
	}

	return nil
}

var ErrLgfNotFound = errors.New("lgf not found")

func (m *Manager) getPoller() Poller {
	poller := &LongPoller{
		Limit:   m.BulkSize,
		Timeout: m.Timeout,
	}
	return poller
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

func (m *Manager) removeExporter(ctx context.Context, e watcher.Event) {

}

func (m *Manager) process(ctx context.Context) {

	fileWatcher := m.fileWatcher
	fileWatcher.AddFilterHook(filterFiles)
	//fileWatcher.SetMaxEvents(1)
	fileWatcher.FilterOps(watcher.Create, watcher.Write, watcher.Remove)

	var err error
	go func() {
		err = fileWatcher.Start(m.Ticker)
	}()

	fileWatcher.Wait()

	if err != nil {
		return
	}

	m.running = true

	for {
		select {
		case <-m.stop:
			return
		case e := <-m.fileWatcher.Event:
			switch e.Op {
			case watcher.Write:
				m.writeWatcherHook(ctx, e)
			case watcher.Create:
				m.createWatcherHook(ctx, e)
			case watcher.Remove:
				m.removeWatcherHook(ctx, e)
			}

		}
	}
}

func (m *Manager) writeWatcherHook(ctx context.Context, event watcher.Event) {

	fileName := event.Path

	_, exporterInWork := m.exporters[fileName]
	if exporterInWork {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO Подумать над циклом чтения и записи offset
	offset := m.journals.GetOffset(fileName)
	exporter, err := createLgpExporter(fileName, offset, m.storage, m.getPoller(), m.TZ)

	if err != nil {
		log.Print(err)
		return
	}

	go func(key string) {
		err := m.waitTurn(ctx)
		if err != nil {
			return
		}

		m.getTurn()
		exporter.Start()
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.exporters, key)
		m.journals.SetOffset(fileName, exporter.eventReader.Offset())
		m.freeTurn()

	}(fileName)
}

func (m *Manager) createWatcherHook(ctx context.Context, e watcher.Event) {

}

func (m *Manager) removeWatcherHook(ctx context.Context, e watcher.Event) {

}
