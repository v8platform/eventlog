package eventlog

import (
	"context"
	"errors"
	"github.com/radovskyb/watcher"
	"github.com/xelaj/go-dry"
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
	Cache              CacheStorage
	Exporters          []ExporterStorage
	BulkSize           int
}

const (
	JournalFormatLGP = iota
	JournalFormatLGD
)

type lgpJournal struct {
	File   string
	UUID   string
	Offset int64
	EventReader

	lgfReader *LgfReader
}

type lgdJournal struct {
	File   string
	Offset int64
	EventReader
}

type EventJournal interface {
	CreateExporter(storage ExporterStorage, poller Poller, tz *time.Location) (*Exporter, error)
	AboveOffset(off int64) bool
}

func (j *lgpJournal) AboveOffset(off int64) bool {
	return j.Offset < off
}

func (j *lgpJournal) CreateExporter(storage ExporterStorage, poller Poller, tz *time.Location) (*Exporter, error) {

	lgfDir := filepath.Dir(j.File)
	LgfFile := filepath.Join(lgfDir, lgfFileName)
	lgfStream, err := os.OpenFile(LgfFile, os.O_RDONLY, 644)

	if err != nil {
		return nil, ErrLgfNotFound
	}

	lgpOpts := LgpReaderOptions{
		LgfDir:    lgfDir,
		LgfFile:   LgfFile,
		LgfStream: lgfStream,
		Offset:    j.Offset,
	}

	reader, err := NewLgpReader(j.File, lgpOpts)

	if err != nil {
		return nil, err
	}

	exporter := NewExporter(reader, storage)
	exporter.Poller = poller
	exporter.TZ = tz

	return exporter, nil

}

func NewManager(opt ManagerOptions) *Manager {

	p := &Manager{

		Folder:    opt.Folder,
		queue:     make(chan struct{}, opt.PoolSize),
		journals:  make(map[string]EventJournal),
		watchers:  make(map[string]chan struct{}),
		exporters: map[string]*Exporter{},
		mu:        sync.Mutex{},
		stop:      make(chan struct{}),
	}

	return p
}

type CacheStorage interface {
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

	watchers map[string]chan struct{}

	journals  map[string]EventJournal
	mu        sync.Mutex
	exporters map[string]*Exporter

	queue chan struct{}

	cache   CacheStorage
	storage []ExporterStorage
	stop    chan struct{}
}

func (m *Manager) Stop() {

}

func (m *Manager) Watch(ctx context.Context, folder string, ticker time.Duration) error {

	fileWatcher := watcher.New()
	fileWatcher.AddFilterHook(filterFiles)

	if err := fileWatcher.AddRecursive(folder); err != nil {
		return err
	}

	fileWatcher.SetMaxEvents(1)
	fileWatcher.FilterOps(watcher.Create, watcher.Write, watcher.Remove)

	errChan := make(chan error)

	go func() {
		errChan <- fileWatcher.Start(ticker)
	}()

	stopWatcher := make(chan struct{})
	m.mu.Lock()
	m.watchers[folder] = stopWatcher
	m.mu.Unlock()

	deleteWatcher := func() {
		fileWatcher.Close()
		m.mu.Lock()
		delete(m.watchers, folder)
		m.mu.Unlock()
	}

	for {
		select {
		case err := <-errChan:
			deleteWatcher()
			return err
		case <-stopWatcher:
			deleteWatcher()
			return nil
		case <-m.stop:
			deleteWatcher()
			return nil
		case <-ctx.Done():
			deleteWatcher()
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

func (m *Manager) addExporter(ctx context.Context, event watcher.Event) {

	fileSize := event.Size()
	fileName := event.Path

	m.mu.Lock()
	defer m.mu.Unlock()

	// TODO Подумать над циклом чтения и записи offset
	eventJournal, ok := m.journals[fileName]
	if ok && eventJournal.AboveOffset(fileSize) {
		return
	}

	_, exporterInWork := m.exporters[fileName]
	if exporterInWork {
		return
	}
	_ = m.waitTurn(ctx)

	poller := &LongPoller{
		Limit:   m.BulkSize,
		Timeout: m.Timeout,
	}

	exporter, _ := eventJournal.CreateExporter(m.storage[0], poller, time.Local)
	m.exporters[fileName] = exporter

	go func(key string) {

		exporter.Start()
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.exporters, key)
		m.freeTurn()

	}(fileName)

}

var ErrLgfNotFound = errors.New("lgf not found")

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
