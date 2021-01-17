package eventlog

import (
	"github.com/k0kubun/pp"
	"github.com/v8platform/brackets"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_parseEventLogItemData(t *testing.T) {
	type args struct {
		data string
	}

	file, _ := os.OpenFile("./tests/1Cv8.lgf", os.O_RDONLY, 046)

	objects := NewLgfReader(file)

	tests := []struct {
		name     string
		metadata Objects
		data     string
		want     Event
	}{
		{
			"simple",
			objects,
			"{2020041213435,N,\n{0,0},1,1,2,2,2,I,\"2: Произвольный текст 7\",1,\n{\"S\",\"2: Простой текст 7\"},\"\",1,1,0,2,0,\n{0}\n},",
			Event{},
		},
		{
			"transaction",
			objects,
			"{20201005114853,U,\n" +
				"{243b06bad83e0,7b3156},71,36,1,13732,11,I,\"\",55,\n" +
				"{\"R\",490:bace0cc47a56444311eaedd56d0dbdf8},\"Отчет производства за \"\",смену Уни00023710 от 03.09.2020 14:05:56\",1,1,0,1101,0,\n" +
				"{0}" +
				"}",
			Event{},
		},
		{
			"transaction",
			objects,
			"{20200919203835,N,\n{0,0},4,1,2,2,35,I,\"Добавление новых идентификаторов объектов метаданных:\nПодсистема.УИ_УниверсальныеИнструменты,\nПодсистема.УИ_УниверсальныеИнструменты.Подсистема.УИ_Поддержка,\nПодсистема.УИ_УниверсальныеИнструменты.Подсистема.УИ_Отладка,\nРоль.УИ_УниверсальныеИнструменты,\nСправочник.УИ_Алгоритмы,\nОтчет.УИ_КонсольОтчетов,\nОбработка.УИ_ГрупповаяОбработкаСправочниковИДокументов,\nОбработка.УИ_РедакторКонстант,\nОбработка.УИ_СтруктураХраненияБазыДанных,\nОбработка.УИ_УдалениеПомеченныхОбъектов,\nОбработка.УИ_ВыполнениеРегламентныхЗаданийНаКлиенте,\nОбработка.УИ_КонсольЗапросов,\nОбработка.УИ_КонсольЗаданий,\nОбработка.УИ_РегистрацияИзмененийДляОбменаДанными,\nОбработка.УИ_КонсольКода,\nОбработка.УИ_ПоискИУдалениеДублей,\nОбработка.УИ_ПоискСсылокНаОбъект,\nОбработка.УИ_ТехПоддержка,\nОбработка.УИ_РедакторРеквизитовОбъекта,\nОбработка.УИ_ДинамическийСписок,\nОбработка.УИ_КонсольHTTPЗапросов,\nОбработка.УИ_ВыгрузкаЗагрузкаДанныхXMLСФильтрами,\nОбработка.УИ_ПреобразованиеДанныхJSON,\nОбработка.УИ_КонструкторРегулярныхВыражений,\nОбработка.УИ_НавигаторПоКонфигурации,\nОбработка.УИ_ФайловыйМенеджер,\nОбработка.УИ_КонсольСравненияДанных,\nОбработка.УИ_ИнформацияОЛицензиях1С,\nОбработка.УИ_КонсольВебСервисов,\nОбработка.УИ_ЗагрузкаДанныхИзТабличногоДокумента,\nОбработка.УИ_ДанныеДляОтладки,\nОбработка.УИ_МенеджерХранилищНастроек,\nОбработка.УИ_РедакторJSON,\nОбработка.УИ_РедакторHTML,\nОбработка.УИ_УниверсальныйОбменДаннымиXML\n\",0,\n{\"U\"},\"\",0,0,0,4,0,\n{2,1,1,2,1}\n}",
			Event{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := strings.NewReader(tt.data)
			parser := brackets.NewParser(r)

			if got := parseEventLogItemData(parser.NextNode(), tt.metadata); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEventLogItemData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseEventLogFiles(t *testing.T) {

	tests := []struct {
		name  string
		file  string
		index int
		want  Event
	}{
		{
			"simple",
			"./tests/20200930210000.lgp",
			10,
			Event{},
		},
		{
			"simple",
			"./tests/20200930210000.lgp",
			29,
			Event{},
		},
	}

	file, _ := os.OpenFile("./tests/1Cv8.lgf", os.O_RDONLY, 046)

	objects := NewLgfReader(file)

	fileLog, _ := os.OpenFile("./tests/20200930210000.lgp", os.O_RDONLY, 046)

	parser := brackets.NewParser(fileLog)

	nodes := parser.ReadAllNodes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := parseEventLogItemData(nodes[tt.index], objects); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEventLogItemData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseEvents(t *testing.T) {

	file, _ := os.OpenFile("./tests/1Cv8.lgf", os.O_RDONLY, 046)

	objects := NewLgfReader(file)

	fileLog, _ := os.OpenFile("./tests/20200930210000.lgp", os.O_RDONLY, 046)

	parser := brackets.NewParser(fileLog)

	nodes := parser.ReadAllNodes()
	var events []Event
	for _, node := range nodes {
		events = append(events, parseEventLogItemData(node, objects))
	}

	pp.Println(events[38])
	pp.Println(nodes[38])
}

func Test_readTill(t *testing.T) {

	tests := []struct {
		name     string
		metadata string
	}{
		{

			"simple",
			"{2,\"Aleksej.local\",1},\n{3,\"Designer\",1},\n{4,\"_$InfoBase$_.RestoreFinish\",1},\n{4,\"_$Session$_.AuthenticationError\",2},\n{1,ae022e20-dbf2-11ea-599b-005056ae0f31,\"Антонина Парунина\",1},\n{4,\"_$Session$_.Authentication\",3},\n{13,1,1},\n{4,\"_$Session$_.Start\",4},\n{4,\"_$User$_.New\",5},\n{4,\"_$Session$_.Finish\",6},\n{3,\"1CV8C\",2},\n{9,530a3164-4ef1-4b3b-8269-13764ef4bf15,\"ОбластьДанныхВспомогательныеДанные\",1},\n{10,\n{\"N\",0},1,1},\n{9,6df2bb92-558c-4453-9de4-e4176e8f93dc,\"ОбластьДанныхОсновныеДанные\",2},\n{10,\n{\"N\",0},2,1},\n{11,\n{2,1,1,2,1},1},\n{12,\n{2,1,1,2,1},1},\n{4,\"_$Transaction$_.Begin\",7},\n{4,\"_$Transaction$_.Commit\",8},\n{4,\"_$Data$_.Update\",9},\n{5,4bb0f7c3-62f3-4352-9bc8-e243dd18fe4a,\"Справочник.ВерсииРасширений\",1},\n{5,71c3e5d3-7504-433f-9d48-7c287b2b863b,\"Константа.ПараметрыБлокировкиРаботыСВнешнимиРесурсами\",2},\n{4,\"Работа с внешними ресурсами заблокирована\",10},\n{5,fc67b510-3fb4-4305-92d2-c252bc718f03,\"РегистрСведений.СведенияОПользователях\",3},\n{5,e8c048fc-0f39-4778-aaf3-d6bff5acb06a,\"РегистрСведений.ЗамерыСтатистики\",4},\n{3,\"BackgroundJob\",3},\n{4,\"_$Job$_.Start\",11},\n{5,6809b99f-ad6b-493e-839c-14dbfc4faf93,\"РегистрСведений.ВсеОбновленияНовостей\",5},\n{4,\"БИП:Новости.Все обновления новостей\",12},\n{5,0fdbd1ab-716d-46b5-aa75-fb09812ced5b,\"РегистрСведений.ЗамерыВремени\",6},\n{5,0894db61-51f6-4068-9634-e55d95e2a8ac,\"Константа.ПараметрыИтоговИАгрегатов\",7},\n{4,\"_$Data$_.New\",13},\n{5,344cb0de-315d-4d05-93f8-8a1f53208470,\"Справочник.КлючевыеОперации\",8},\n{4,\"_$Job$_.Succeed\",14},\n{4,\"БИП:Новости.Изменение данных\",15},\n{5,4cd6c242-ca58-4a0d-915d-c3b96be0cd61,\"Константа.НастройкиНовостей\",9},\n{5,8c35433e-e0e0-4e81-abf1-91a4b275fd1f,\"РегистрСведений.РассчитанныеОтборыПоНовостям_РедкоМеняющиеся\",10},\n{4,\"БИП:Новости.Сервис и регламент\",16},\n{3,\"1CV8\",4},\n{1,071523a4-516f-4fce-ba4b-0d11ab7a1893,\"\",2},\n{11,\n{2,1,1,2,1},2},\n{13,1,2},",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := strings.NewReader(tt.metadata)
			objects := NewLgfReader(r)

			objects.readTill(0, 0)

		})
	}
}
