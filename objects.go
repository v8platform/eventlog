package eventlog

import (
	"fmt"
	"github.com/v8platform/brackets"
	"io"
	"log"
)

type LgfReader struct {
	stream            io.Reader
	objects           map[string]string
	referencedObjects map[string][]string
	parser            *brackets.Parser
	curNode           brackets.Node
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
		stream:            r,
		parser:            brackets.NewParser(r),
		objects:           make(map[string]string),
		referencedObjects: make(map[string][]string),
	}

}

func getKeyValue(objectType int, id ...int) string {

	key := fmt.Sprintf("%d", objectType)

	for _, i := range id {
		key += fmt.Sprintf(".%d", i)
	}

	return key
}

func (r *LgfReader) ReferencedObjectValue(objectType int, id ...int) (value, uuid string) {

	if len(id) == 0 || (len(id) == 1 && id[0] == 0) {
		return "", ""
	}

	key := getKeyValue(objectType, id...)
	if valueUuid, ok := r.referencedObjects[key]; ok {
		return valueUuid[0], valueUuid[1]
	}

	r.readTill(objectType, id...)

	if valueUuid, ok := r.referencedObjects[key]; ok {
		return valueUuid[0], valueUuid[1]
	}

	log.Printf("error get referenced object value for type <%d> & id <%d>", objectType, id)

	return
}

func (r *LgfReader) ObjectValue(objectType int, id ...int) (value string) {

	if len(id) == 0 || (len(id) == 1 && id[0] == 0) {
		return ""
	}

	key := getKeyValue(objectType, id...)
	if value, ok := r.objects[key]; ok {
		return value
	}

	r.readTill(objectType, id...)

	if value, ok := r.objects[key]; ok {
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

	r.curNode = r.parser.NextNode()
	return r.curNode != nil
}

func (r *LgfReader) readTill(object int, id ...int) {

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

			r.referencedObjects[key] = value

			if len(id) > 0 && objectType == object && keyNode == id[0] {
				break
			}
		case ObjectTypeSessionDataSeparatorValue:

			keyNode := node.Int(2)
			key2Node := node.Int(3)
			var key = getKeyValue(objectType, keyNode, key2Node)

			valueNode := node.Get(1, 1)

			r.objects[key] = valueNode

			if len(id) > 1 &&
				objectType == object &&
				keyNode == id[0] &&
				key2Node == id[1] {
				break
			}

		default:

			keyNode := node.Int(2)

			var key = getKeyValue(objectType, keyNode)

			valueNode := node.Get(1)

			r.objects[key] = valueNode

			if len(id) > 0 && objectType == object && keyNode == id[0] {
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
