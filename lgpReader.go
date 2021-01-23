package eventlog

import (
	"bufio"
	"bytes"
	"context"
	"github.com/v8platform/brackets"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var defaultOptions = LgpReaderOptions{}

var _ EventReader = (*LgpReader)(nil)

type LgpReaderOptions struct {
	LgfDir    string
	LgfFile   string
	LgfStream io.ReadSeekCloser
	LgfOffset int64
	Offset    int64
}

type LgpReader struct {
	stream  io.ReadSeekCloser
	parser  *brackets.Parser
	objects Objects
	offset  int64
	Uuid    string
	Version string
}

func (r *LgpReader) Close() error {
	err := r.stream.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *LgpReader) Seek(offset int64) (int64, error) {

	if r.offset == offset {
		return 0, nil
	}

	n, err := r.stream.Seek(offset, io.SeekStart)
	if err != nil {
		return n, err
	}

	r.offset = offset

	return n, nil
}

func (r *LgpReader) Offset() int64 {

	return r.offset
}

func (r *LgpReader) readMetadata() {

	br := bufio.NewReader(r.stream)

	versionBytes, _ := br.ReadBytes('\n')
	uuidString, _ := br.ReadString('\n')
	versionBytes = bytes.Trim(versionBytes, "\xef\xbb\xbf")

	r.Version = strings.TrimSpace(string(versionBytes))
	r.Uuid = strings.TrimSpace(uuidString)

	r.offset, _ = r.stream.Seek(0, io.SeekCurrent)

}

func (r *LgpReader) Read(limit int, timeout time.Duration) (items []Event, err error) {

	return r.read(context.Background(), limit, timeout)
}

func (r *LgpReader) ReadCtx(ctx context.Context, limit int, timeout time.Duration) (items []Event, err error) {
	return r.read(ctx, limit, timeout)
}

func (r *LgpReader) read(ctx context.Context, limit int, timeout time.Duration) (items []Event, err error) {

	if limit < 1 {
		// Указывать лимит считывания обязательно
		// уменьшает нагрузку на ЦП и память
		return nil, nil
	}

	var timeoutC <-chan time.Time

	if timeout > 0 {
		timeoutC = time.After(timeout)
	}

	var count int
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	//limiter := make(chan struct{}, 10)
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return items, ctx.Err()
		case <-timeoutC:
			wg.Wait()
			return
		default:

			if limit > 0 && count == limit {
				wg.Wait()
				return items, nil
			}

			//limiter <-empty
			node, n := r.parser.NextNode()
			start := r.offset
			r.offset += int64(n)
			if node == nil {
				wg.Wait()
				return items, io.EOF
			}

			count++
			wg.Add(1)

			go func(n brackets.Node, offset, size int64) {

				event := &Event{
					Offset: offset,
					Size:   size,
				}

				parseEventLogItemData(event, n, r.objects)
				mu.Lock()
				defer mu.Unlock()
				defer wg.Done()
				items = append(items, *event)
				//<-limiter
			}(node, start, int64(n))

		}
	}
}

//NewLgpReader создает новый читатель журнала регистрации 1С
func NewLgpReader(path string, opts ...LgpReaderOptions) (*LgpReader, error) {

	lgpStream, err := os.OpenFile(path, os.O_RDONLY, 644)
	if err != nil {
		return nil, err
	}

	options := defaultOptions

	if len(opts) > 0 {
		options = opts[0]
	}

	if len(options.LgfDir) == 0 {
		options.LgfDir = filepath.Dir(path)
	}

	lgfStream, err := getLgfFile(options)

	if err != nil {
		return nil, err
	}

	reader := &LgpReader{
		stream:  lgpStream,
		parser:  brackets.NewParser(lgpStream),
		objects: NewLgfReader(lgfStream),
	}

	reader.readMetadata()

	if options.Offset > 0 {
		if _, err := lgpStream.Seek(options.Offset, io.SeekStart); err != nil {
			return nil, err
		}
		reader.offset = options.Offset
	}

	return reader, nil

}

func getLgfFile(opt LgpReaderOptions) (io.ReadSeekCloser, error) {

	switch {

	case opt.LgfStream != nil:
		return opt.LgfStream, nil

	case len(opt.LgfFile) > 0:

		lgfStream, err := os.OpenFile(opt.LgfFile, os.O_RDONLY, 644)
		if err != nil {
			return nil, err
		}

		return lgfStream, nil
	default:

		LgfFile := filepath.Join(opt.LgfDir, "1Cv8.lgf")

		lgfStream, err := os.OpenFile(LgfFile, os.O_RDONLY, 644)
		if err != nil {
			return nil, err
		}

		return lgfStream, nil
	}

}

func parseEventLogItemData(event *Event, parsedData brackets.Node, objects Objects) {

	event.Date, _ = time.Parse(`20060102150405`, parsedData.Get(0))

	event.TransactionStatus = TransactionStatusType(parsedData.Get(1))
	event.TransactionNumber, event.TransactionDate = getTransactionData(parsedData.GetNode(2))

	event.User, event.UserUuid = objects.ReferencedObjectValue(ObjectTypeUsers, parsedData.Int(3))

	event.Computer = objects.ObjectValue(ObjectTypeComputers, parsedData.Int(4))
	event.Application = ApplicationType(objects.ObjectValue(ObjectTypeApplications, parsedData.Int(5)))

	event.Connection = parsedData.Int64(6)
	event.Event = EventType(objects.ObjectValue(ObjectTypeEvents, parsedData.Int(7)))
	event.Severity = SeverityType(parsedData.Get(8))

	event.Comment = parsedData.Get(9)

	event.Metadata, event.MetadataUuid = objects.ReferencedObjectValue(ObjectTypeMetadata, parsedData.Int(10))

	event.Data = getData(parsedData.GetNode(11), event.Event)
	event.DataPresentation = parsedData.Get(12)

	event.Server = objects.ObjectValue(ObjectTypeServers, parsedData.Int(13))
	event.MainPort = objects.ObjectValue(ObjectTypeMainPorts, parsedData.Int(14))
	event.AddPort = objects.ObjectValue(ObjectTypeAddPorts, parsedData.Int(15))
	event.Session = parsedData.Int64(16)

	//sessionDataSeparators := getSessionDataSeparators(parsedData.GetNode(18), objects)
	//
	//if len(sessionDataSeparators) > 0 {
	//	event.SessionDataSeparators = sessionDataSeparators
	//}

}

func getSessionDataSeparators(node brackets.Node, objects Objects) []RefObject {

	count := node.Int(0)
	if count == 0 {
		return nil
	}

	var dataSeparators []RefObject

	for i := 1; i <= count*2; i = i + 2 {

		key := node.Int(i)
		name, uuid := objects.ReferencedObjectValue(ObjectTypeSessionDataSeparator, key)
		value := objects.ObjectValue(ObjectTypeSessionDataSeparatorValue, key, node.Int(i+1))
		dataSeparators = append(dataSeparators, RefObject{
			Name:  name,
			Uuid:  uuid,
			Value: value,
		})
	}

	return dataSeparators
}

func getData(node brackets.Node, eventType EventType) interface{} {
	dataType := node.Get(0)
	switch dataType {
	case "R": // Reference

		return RefObject{
			Name: node.Get(0),
			Uuid: node.Get(1),
		}

	case "U": // Undefined
		return ""
	case "O": // object
		return RefObject{
			Name: node.Get(1, 1),
			Uuid: node.Get(1, 2),
		}
	case "A": // array

		count := node.Int(0)
		if count == 0 {
			return nil
		}
		var arr []interface{}

		for i := count; i > 0; i-- {

			arr = append(arr, getData(node.GetNode(i+1), eventType))

		}

		return arr

	case "S": // String
		return node.Get(1)
	case "B": // Boolean
		return node.Bool(1)
	case "P": // Complex data

		subDataNode := node.GetNode(1)
		subDataType := ComplexDataType(subDataNode.Int(0))

		parser := subDataType.Parser()

		if parser == nil {
			return nil
		}

		parser.Parse(subDataNode, eventType)

		if d, ok := parser.(*complexDataMapParser); ok {
			return d.Data
		}

		return parser

	default:
		return ""
	}
}

// Конвертация во время далась очень сложно только через unix
func getTransactionData(data brackets.Node) (int64, time.Time) {

	seconds := int(From16To10(data.Get(0)) / 10000)

	transactionDate := SecondsToUnixTime(seconds)
	transactionNumber := From16To10(data.Get(1))

	return transactionNumber, transactionDate
}
