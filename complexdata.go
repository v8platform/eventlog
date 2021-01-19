package eventlog

import "github.com/v8platform/brackets"

type ComplexDataParser interface {
	Parse(node brackets.Node, eventType EventType)
}

type ComplexDataType int

const (
	UnknownComplexData      ComplexDataType = 0
	AuthenticationErrorData ComplexDataType = 1
	AuthenticationData      ComplexDataType = 6
	UpdateUserData          ComplexDataType = 30
)

type complexDataMapParser struct {
	fn   func(data map[string]interface{}, node brackets.Node, eventType EventType)
	Data map[string]interface{}
}

func (p *complexDataMapParser) Parse(node brackets.Node, eventType EventType) {

	p.fn(p.Data, node, eventType)
}

func newComplexDataMapParser(fn func(data map[string]interface{}, node brackets.Node, eventType EventType)) *complexDataMapParser {
	return &complexDataMapParser{
		fn:   fn,
		Data: make(map[string]interface{}),
	}
}

func (c ComplexDataType) Parser() ComplexDataParser {

	switch c {
	case UnknownComplexData:
		return nil
	case AuthenticationErrorData:
		return newComplexDataMapParser(func(data map[string]interface{}, node brackets.Node, eventType EventType) {
			data["Пользователь ОС"] = getData(node.GetNode(1), eventType)
		})
	case AuthenticationData:
		return newComplexDataMapParser(func(data map[string]interface{}, node brackets.Node, eventType EventType) {
			data["Имя"] = getData(node.GetNode(1), eventType)
			data["Текущий пользователь ОС"] = getData(node.GetNode(2), eventType)
		})
	case UpdateUserData:
		return newComplexDataMapParser(func(data map[string]interface{}, node brackets.Node, eventType EventType) {
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
