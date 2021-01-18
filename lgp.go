package eventlog

import (
	"bufio"
	"bytes"
	"context"
	"github.com/k0kubun/pp"
	"github.com/v8platform/brackets"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultBufSize = 10

var defaultOptions = LgpReaderOptions{}

type LgpReaderOptions struct {
	DirLgpFile string
	LgfFile    string
	LgfStream  io.ReadSeekCloser
	LgfOffset  int64
	Offset     int64
}

type LgpReader struct {
	stream   io.ReadSeekCloser
	parser   *brackets.Parser
	objects  Objects
	offset   int64
	stopChan chan struct{}

	uuid    string
	version string
}

func (r *LgpReader) SetOffset(offset int64) (int64, error) {

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

	currentOffset := r.offset

	if currentOffset > 0 {
		_, _ = r.SetOffset(0)
	}

	br := bufio.NewReader(r.stream)

	versionString, _ := br.ReadString('\n')
	uuidString, _ := br.ReadString('\n')

	r.offset, _ = r.stream.Seek(currentOffset, io.SeekStart)

	versionString = string(trapBOM([]byte(versionString)))
	r.version = strings.TrimSpace(versionString)
	r.uuid = strings.TrimSpace(uuidString)

}

func trapBOM(b []byte) []byte {
	trimmedBytes := bytes.Trim(b, "\xef\xbb\xbf")
	return trimmedBytes
}

func (r *LgpReader) Read() *Event {

	return r.readEvent()
}

func (r *LgpReader) readEvent() *Event {

	node, n := r.parser.NextNode()
	r.offset += int64(n)
	if node == nil {
		return nil
	}

	e := parseEventLogItemData(node, r.objects)
	return &e
}

func (r *LgpReader) streamReadEvent(stream EventsStream, ctx context.Context) {

	go r.readEvents(stream, ctx)

}

func (r *LgpReader) readEvents(stream EventsStream, ctx context.Context) {

	p := r.parser

	limiter := make(chan struct{}, cap(stream))
	pp.Println("limiter", cap(limiter))
	for {
		select {
		case <-r.stopChan:
			return
		case <-ctx.Done():
			return
		default:
			limiter <- struct{}{}
			node, n := p.NextNode()
			r.offset += int64(n)
			if node == nil {
				return
			}
			go func(n brackets.Node) {
				stream <- parseEventLogItemData(n, r.objects)
				<-limiter
			}(node)

		}
	}
}

func (r *LgpReader) StreamRead(ctx context.Context, bufSize ...int) EventsStream {

	size := defaultBufSize
	if len(bufSize) > 1 {
		size = bufSize[0]
	}

	events := make(EventsStream, size)

	r.streamReadEvent(events, ctx)

	return events
}

func (r *LgpReader) Stream(ctx context.Context, events EventsStream) {

	r.streamReadEvent(events, ctx)

}

func NewLgpReaderFromStream(lgpStream io.ReadSeekCloser, lgfStream io.ReadSeekCloser) (*LgpReader, error) {

	return &LgpReader{
		stream:   lgpStream,
		parser:   brackets.NewParser(lgpStream),
		objects:  NewLgfReader(lgfStream),
		stopChan: make(chan struct{}),
	}, nil

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

	if len(options.DirLgpFile) == 0 {
		options.DirLgpFile = filepath.Dir(path)
	}

	lgfStream, err := getLgfFile(options)

	if err != nil {
		return nil, err
	}

	reader, err := NewLgpReaderFromStream(lgpStream, lgfStream)
	if err != nil {
		return nil, err
	}

	reader.readMetadata()

	if _, err := reader.SetOffset(options.Offset); err != nil {
		return nil, err
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

		LgfFile := filepath.Join(opt.DirLgpFile, "1Cv8.lgf")

		lgfStream, err := os.OpenFile(LgfFile, os.O_RDONLY, 644)
		if err != nil {
			return nil, err
		}

		return lgfStream, nil
	}

}

func parseEventLogItemData(parsedData brackets.Node, objects Objects) Event {

	event := Event{}

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

	sessionDataSeparators := getSessionDataSeparators(parsedData.GetNode(18), objects)

	if len(sessionDataSeparators) > 0 {
		event.SessionDataSeparators = sessionDataSeparators
	}

	return event

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

		if d, ok := parser.(*ComplexDataMapParser); ok {
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

type ComplexDataType int

type ComplexData interface {
	Parse(node brackets.Node, eventType EventType)
}

const (
	UnknownComplexData      ComplexDataType = 0
	AuthenticationErrorData ComplexDataType = 1
	AuthenticationData      ComplexDataType = 6
	UpdateUserData          ComplexDataType = 30
)

type ComplexDataMapParser struct {
	fn   func(data map[string]interface{}, node brackets.Node, eventType EventType)
	Data map[string]interface{}
}

func (p *ComplexDataMapParser) Parse(node brackets.Node, eventType EventType) {

	p.fn(p.Data, node, eventType)
}

func NewComplexDataMapParser(fn func(data map[string]interface{}, node brackets.Node, eventType EventType)) *ComplexDataMapParser {
	return &ComplexDataMapParser{
		fn:   fn,
		Data: make(map[string]interface{}),
	}
}

func (c ComplexDataType) Parser() ComplexData {

	switch c {
	case UnknownComplexData:
		return nil
	case AuthenticationErrorData:
		return NewComplexDataMapParser(func(data map[string]interface{}, node brackets.Node, eventType EventType) {
			data["Пользователь ОС"] = getData(node.GetNode(1), eventType)
		})
	case AuthenticationData:
		return NewComplexDataMapParser(func(data map[string]interface{}, node brackets.Node, eventType EventType) {
			data["Имя"] = getData(node.GetNode(1), eventType)
			data["Текущий пользователь ОС"] = getData(node.GetNode(2), eventType)
		})
	case UpdateUserData:
		return NewComplexDataMapParser(func(data map[string]interface{}, node brackets.Node, eventType EventType) {
			data["Аутентификация ОС"] = getData(node.GetNode(1), eventType)
			data["Аутентификация 1С:Предприятия"] = getData(node.GetNode(2), eventType)
			data["Запрещено изменять пароль"] = getData(node.GetNode(3), eventType)
			data["Имя"] = getData(node.GetNode(4), eventType)
			data["Основной язык"] = getData(node.GetNode(5), eventType)
		})
	default:
		return nil
	}
}
