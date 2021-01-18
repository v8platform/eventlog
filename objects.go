package eventlog

import (
	"fmt"
	"github.com/v8platform/brackets"
	"io"
	"log"
	"sync"
)

type LgfReader struct {
	stream  io.Reader
	objects *sync.Map
	parser  *brackets.Parser
	curNode brackets.Node
	muRead  *sync.RWMutex

	needRead bool

	waitRead chan struct{}
}

const (
	ObjectTypeNone                      = 0
	ObjectTypeUsers                     = 1
	ObjectTypeComputers                 = 2
	ObjectTypeApplications              = 3
	ObjectTypeEvents                    = 4
	ObjectTypeMetadata                  = 5
	ObjectTypeServers                   = 6
	ObjectTypeMainPorts                 = 7
	ObjectTypeAddPorts                  = 8
	ObjectTypeSessionDataSeparator      = 9
	ObjectTypeSessionDataSeparatorValue = 10
	ObjectTypeUnknown                   = 11
)

func NewLgfReader(r io.Reader) *LgfReader {

	return &LgfReader{
		stream:   r,
		parser:   brackets.NewParser(r),
		objects:  &sync.Map{},
		waitRead: make(chan struct{}, 1),
		muRead:   &sync.RWMutex{},
	}

}

func getKeyValue(objectType int, id ...int) string {

	key := fmt.Sprintf("%d", objectType)

	for _, i := range id {
		key += fmt.Sprintf(".%d", i)
	}

	return key
}
func (r *LgfReader) getReferencedObjectValue(key string) (string, string, bool) {

	r.muRead.RLock()
	defer r.muRead.RUnlock()

	val, ok := r.objects.Load(key)

	if ok {
		valueUuid := val.([]string)
		return valueUuid[0], valueUuid[1], true
	}

	return "", "", false
}

func (r *LgfReader) getObjectValue(key string) (string, bool) {

	r.muRead.RLock()
	defer r.muRead.RUnlock()

	val, ok := r.objects.Load(key)

	if ok {

		return val.(string), true
	}

	return "", false
}

func (r *LgfReader) ReferencedObjectValue(objectType int, id ...int) (value, uuid string) {

	if len(id) == 0 || (len(id) == 1 && id[0] == 0) {
		return "", ""
	}

	key := getKeyValue(objectType, id...)
	if value, uuid, ok := r.getReferencedObjectValue(key); ok {
		return value, uuid
	}

	r.readTill(objectType, key)

	if value, uuid, ok := r.getReferencedObjectValue(key); ok {
		return value, uuid
	}

	log.Printf("error get referenced object value for type <%d> & id <%d>", objectType, id)

	return
}

func (r *LgfReader) ObjectValue(objectType int, id ...int) (value string) {

	if len(id) == 0 || (len(id) == 1 && id[0] == 0) {
		return ""
	}

	key := getKeyValue(objectType, id...)
	if value, ok := r.getObjectValue(key); ok {
		return value
	}

	r.readTill(objectType, key)

	if value, ok := r.getObjectValue(key); ok {
		return value
	}

	log.Printf("error get object value for type <%d> & id <%d>", objectType, id)

	return
}

func (r *LgfReader) initParser() {

	if r.parser == nil {
		r.parser = brackets.NewParser(r.stream)
	}

}

func (r *LgfReader) Read() bool {

	r.initParser()

	r.curNode, _ = r.parser.NextNode()
	return r.curNode != nil
}

func (r *LgfReader) readTill(object int, needKey ...string) {

	r.muRead.Lock()
	defer r.muRead.Unlock()

	for r.Read() {

		node := r.curNode

		objectType := node.Int(0)

		// Skip unknown object types
		if objectType >= ObjectTypeUnknown {
			continue
		}

		switch objectType {
		case ObjectTypeUsers, ObjectTypeMetadata, ObjectTypeSessionDataSeparator:

			keyNode := node.Int(3)

			var key = getKeyValue(objectType, keyNode)

			valueNode1 := node.Get(2)
			valueNode2 := node.Get(1)

			value := []string{valueNode1, valueNode2}

			r.objects.Store(key, value)

			if len(needKey) > 0 && objectType == object && key == needKey[0] {
				break
			}
		case ObjectTypeSessionDataSeparatorValue:

			keyNode := node.Int(2)
			key2Node := node.Int(3)
			var key = getKeyValue(objectType, keyNode, key2Node)

			valueNode := node.Get(1, 1)

			r.objects.Store(key, valueNode)

			if len(needKey) > 1 &&
				objectType == object &&
				key == needKey[0] {
				break
			}

		default:

			keyNode := node.Int(2)

			var key = getKeyValue(objectType, keyNode)

			valueNode := node.Get(1)

			r.objects.Store(key, valueNode)

			if len(needKey) > 0 && objectType == object && key == needKey[0] {
				break
			}

		}
	}

}

type RefObject struct {
	Name  string
	Uuid  string
	Value string
}
