package eventlog

import "strings"

type TransactionStatusType string

func (t TransactionStatusType) String() string {
	switch t {
	case TransactionStatusCommitted:
		return "Зафиксирована"
	case TransactionStatusCanceled:
		return "Отменена"
	case TransactionStatusNotCompleted:
		return "Не завершена"
	case TransactionStatusNoTransaction:
		return "Нет транзакции"
	default:
		return ""
	}
}

const (
	TransactionStatusCommitted     TransactionStatusType = "U" // "U" => "Зафиксирована",
	TransactionStatusCanceled      TransactionStatusType = "C" //"C" => "Отменена",
	TransactionStatusNotCompleted  TransactionStatusType = "R" //"R" => "Не завершена",
	TransactionStatusNoTransaction TransactionStatusType = "N" //"N" => "Нет транзакции",
)

type SeverityType string

func (t SeverityType) String() string {
	switch t {
	case SeverityInfo:
		return "Информация"
	case SeverityError:
		return "Ошибка"
	case SeverityWarn:
		return "Предупреждение"
	case SeverityNote:
		return "Примечание"
	default:
		return ""
	}
}

const (
	SeverityInfo  SeverityType = "I" //"I" => "Информация",
	SeverityError SeverityType = "E" //"E" => "Ошибка",
	SeverityWarn  SeverityType = "W" //"W" => "Предупреждение",
	SeverityNote  SeverityType = "N" //"N" => "Примечание",
)

type ApplicationType string

func (t ApplicationType) String() string {
	switch t {
	case Application1CV8:
		return "Толстый клиент"
	case Application1CV8C:
		return "Тонкий клиент"
	case ApplicationWebClient:
		return "Веб-клиент"
	case ApplicationDesigner:
		return "Конфигуратор"
	case ApplicationCOMConnection:
		return "Внешнее соединение (COM, обычное)"
	case ApplicationWSConnection:
		return "Сессия web-сервиса"
	case ApplicationBackgroundJob:
		return "Фоновое задание"
	case ApplicationSystemBackgroundJob:
		return "Системное фоновое задание"
	case ApplicationSrvrConsole:
		return "Консоль кластера"
	case ApplicationCOMConsole:
		return "Внешнее соединение (COM, административное)"
	case ApplicationJobScheduler:
		return "Планировщик заданий"
	case ApplicationDebugger:
		return "Отладчик"
	case ApplicationRAS:
		return "Сервер администрирования"
	default:
		return "Неопределено"
	}
}

const (
	Application1CV8                ApplicationType = "1CV8"                //"1CV8" => "Толстый клиент",
	Application1CV8C               ApplicationType = "1CV8C"               //"1CV8C" => "Тонкий клиент",
	ApplicationWebClient           ApplicationType = "WebClient"           //"WebClient" => "Веб-клиент",
	ApplicationDesigner            ApplicationType = "Designer"            //"Designer" => "Конфигуратор",
	ApplicationCOMConnection       ApplicationType = "COMConnection"       //"COMConnection" => "Внешнее соединение (COM, обычное)",
	ApplicationWSConnection        ApplicationType = "WSConnection"        //"WSConnection" => "Сессия web-сервиса",
	ApplicationBackgroundJob       ApplicationType = "BackgroundJob"       //"BackgroundJob" => "Фоновое задание",
	ApplicationSystemBackgroundJob ApplicationType = "SystemBackgroundJob" //"SystemBackgroundJob" => "Системное фоновое задание",
	ApplicationSrvrConsole         ApplicationType = "SrvrConsole"         //"SrvrConsole" => "Консоль кластера",
	ApplicationCOMConsole          ApplicationType = "COMConsole"          //"COMConsole" => "Внешнее соединение (COM, административное)",
	ApplicationJobScheduler        ApplicationType = "JobScheduler"        //"JobScheduler" => "Планировщик заданий",
	ApplicationDebugger            ApplicationType = "Debugger"            //"Debugger" => "Отладчик",
	ApplicationRAS                 ApplicationType = "RAS"                 //"RAS" => "Сервер администрирования",
)

type EventScopeType string

func (t EventScopeType) String() string {
	switch t {
	case EventScopeAccess:
		return "Доступ"
	case EventScopeData:
		return "Данные"
	case EventScopeInfobase:
		return "Информационная база"
	case EventScopeJob:
		return "Фоновое задание"
	case EventScopeOpenIDProvider:
		return "Провайдер OpenID"
	case EventScopePerformError:
		return "Ошибка выполнения"
	case EventScopeSession:
		return "Сеанс"
	case EventScopeTransaction:
		return "Транзакция"
	case EventScopeUser:
		return "Пользователи"
	default:
		return string(t)
	}
}

const (
	EventScopeUndefined      EventScopeType = "_$Undefined$_"      //  => Неопределено
	EventScopeAccess         EventScopeType = "_$Access$_"         //  => Доступ
	EventScopeData           EventScopeType = "_$Data$_"           //  => Данные
	EventScopeInfobase       EventScopeType = "_$InfoBase$_"       //  => Информационная база
	EventScopeJob            EventScopeType = "_$Job$_"            //  => Фоновое задание
	EventScopeOpenIDProvider EventScopeType = "_$OpenIDProvider$_" //  => Провайдер OpenID
	EventScopePerformError   EventScopeType = "_$PerformError$_"   //  => Ошибка выполнения
	EventScopeSession        EventScopeType = "_$Session$_"        //  => Сеанс
	EventScopeTransaction    EventScopeType = "_$Transaction$_"    //  => Транзакция
	EventScopeUser           EventScopeType = "_$User$_"           //  => Пользователи
)

type EventCauseType string

func (t EventCauseType) String() string {
	switch t {
	case EventCauseAccess:
		return "Доступ"
	case EventCauseAccessDenied:
		return "Отказ в доступе"
	case EventCauseDelete:
		return "Удаление"
	case EventCauseDeletePredefinedData:
		return "Удаление предопределенных данных"
	case EventCauseDeleteVersions:
		return "Удаление версий"
	case EventCauseNew:
		return "Добавление"
	case EventCauseNewPredefinedData:
		return "Добавление предопределенных данных"
	case EventCauseNewVersion:
		return "Добавление версии"
	case EventCausePos:
		return "Проведение"
	case EventCausePredefinedDataInitialization:
		return "Инициализация предопределенных данных"
	case EventCausePredefinedDataInitializationDataNotFound:
		return "Инициализация предопределенных данных. Данные не найдены"
	case EventCauseSetPredefinedDataInitialization:
		return "Установка инициализации предопределенных данных"
	case EventCauseSetStandardODataInterfaceContent:
		return "Изменение состава стандартного интерфейса OData"
	case EventCauseTotalsMaxPeriodUpdate:
		return "Изменение максимального периода рассчитанных итогов"
	case EventCauseTotalsMinPeriodUpdate:
		return "Изменение минимального периода рассчитанных итогов"
	case EventCausePost:
		return "Проведение"
	case EventCauseUnpost:
		return "Отмена проведения"
	case EventCauseUpdate:
		return "Изменение"
	case EventCauseUpdatePredefinedData:
		return "Изменение предопределенных данных"
	case EventCauseVersionCommentUpdate:
		return "Изменение комментария версии"
	case EventCauseConfigExtensionUpdate:
		return "Изменение расширения конфигурации"
	case EventCauseConfigUpdate:
		return "Изменение конфигурации"
	case EventCauseDBConfigBackgroundUpdateCancel:
		return "Отмена фонового обновления"
	case EventCauseDBConfigBackgroundUpdateFinish:
		return "Завершение фонового обновления"
	case EventCauseDBConfigBackgroundUpdateResume:
		return "Продолжение (после приостановки) процесса фонового обновления"
	case EventCauseDBConfigBackgroundUpdateStart:
		return "Запуск фонового обновления"
	case EventCauseDBConfigBackgroundUpdateSuspend:
		return "Приостановка (пауза) процесса фонового обновления"
	case EventCauseDBConfigExtensionUpdate:
		return "Изменение расширения конфигурации"
	case EventCauseDBConfigExtensionUpdateError:
		return "Ошибка изменения расширения конфигурации"
	case EventCauseDBConfigUpdate:
		return "Изменение конфигурации базы данных"
	case EventCauseDBConfigUpdateStart:
		return "Запуск обновления конфигурации базы данных"
	case EventCauseDumpError:
		return "Ошибка выгрузки в файл"
	case EventCauseDumpFinish:
		return "Окончание выгрузки в файл"
	case EventCauseDumpStart:
		return "Начало выгрузки в файл"
	case EventCauseEraseData:
		return "Удаление данных информационной баз"
	case EventCauseEventLogReduce:
		return "Сокращение журнала регистрации"
	case EventCauseEventLogReduceError:
		return "Ошибка сокращения журнала регистрации"
	case EventCauseEventLogSettingsUpdate:
		return "Изменение параметров журнала регистрации"
	case EventCauseEventLogSettingsUpdateError:
		return "Ошибка при изменение настроек журнала регистрации"
	case EventCauseInfoBaseAdmParamsUpdate:
		return "Изменение параметров информационной базы"
	case EventCauseInfoBaseAdmParamsUpdateError:
		return "Ошибка изменения параметров информационной базы"
	case EventCauseIntegrationServiceActiveUpdate:
		return "Изменение активности сервиса интеграции"
	case EventCauseIntegrationServiceSettingsUpdate:
		return "Изменение настроек сервиса интеграции"
	case EventCauseMasterNodeUpdate:
		return "Изменение главного узла"
	case EventCausePredefinedDataUpdate:
		return "Обновление предопределенных данных"
	case EventCauseRegionalSettingsUpdate:
		return "Изменение региональных установок"
	case EventCauseRestoreError:
		return "Ошибка загрузки из файла"
	case EventCauseRestoreFinish:
		return "Окончание загрузки из файла"
	case EventCauseRestoreStart:
		return "Начало загрузки из файла"
	case EventCauseSecondFactorAuthTemplateDelete:
		return "Удаление шаблона вторго фактора аутентификации"
	case EventCauseSecondFactorAuthTemplateNew:
		return "Добавление шаблона вторго фактора аутентификации"
	case EventCauseSecondFactorAuthTemplateUpdate:
		return "Изменение шаблона вторго фактора аутентификации"
	case EventCauseSetPredefinedDataUpdate:
		return "Установить обновление предопределенных данных"

	case EventCauseTARImportant:
		return "Ошибка"
	case EventCauseTARInfo:
		return "Сообщение"
	case EventCauseTARMess:
		return "Предупреждение"

	case EventCauseCancel:
		return "Отмена"
	case EventCauseFail:
		return "Ошибка выполнения"
	case EventCauseStart:
		return "Запуск"
	case EventCauseSucceed:
		return "Успешное завершение"
	case EventCauseTerminate:
		return "Принудительное завершение"

	case EventCauseNegativeAssertion:
		return "Отклонено"
	case EventCausePositiveAssertion:
		return "Подтверждено"

	case EventCauseAuthentication:
		return "Аутентификация"
	case EventCauseAuthenticationError:
		return "Ошибка аутентификации"
	case EventCauseAuthenticationFirstFactor:
		return "Аутентификация первый фактор"
	case EventCauseConfigExtensionApplyError:
		return "Ошибка применения расширения конфигурации"
	case EventCauseFinish:
		return "Завершение"
	case EventCauseBegin:
		return "Начало"

	case EventCauseCommit:
		return "Фиксация"
	case EventCauseRollback:
		return "Отмена"

	case EventCauseAuthenticationLock:
		return "Блокировка аутентификации"
	case EventCauseAuthenticationUnlock:
		return "Разблокировка аутентификации"
	case EventCauseAuthenticationUnlockError:
		return "Ошибка разблокировки аутентификации"
	case EventCauseDeleteError:
		return "Ошибка удаления"
	case EventCauseUpdateError:
		return "Ошибка изменения"
	case EventCauseNewError:
		return "Ошибка добавления"

	default:
		return string(t)
	}
}

const (
	EventCauseAccess       EventCauseType = "Access"       //"_$Access$_.Access" => "Доступ.Доступ",
	EventCauseAccessDenied EventCauseType = "AccessDenied" //"_$Access$_.AccessDenied" => "Доступ.Отказ в доступе",

	EventCauseDelete                                   EventCauseType = "Delete"                                   //"_$Data$_.Delete" => "Данные.Удаление",
	EventCauseDeletePredefinedData                     EventCauseType = "DeletePredefinedData"                     //"_$Data$_.DeletePredefinedData" => " Данные.Удаление предопределенных данных",
	EventCauseDeleteVersions                           EventCauseType = "DeleteVersions"                           //"_$Data$_.DeleteVersions" => "Данные.Удаление версий",
	EventCauseNew                                      EventCauseType = "New"                                      //"_$Data$_.New" => "Данные.Добавление",
	EventCauseNewPredefinedData                        EventCauseType = "NewPredefinedData"                        //"_$Data$_.NewPredefinedData" => "Данные.Добавление предопределенных данных",
	EventCauseNewVersion                               EventCauseType = "NewVersion"                               //"_$Data$_.NewVersion" => "Данные.Добавление версии",
	EventCausePos                                      EventCauseType = "Pos"                                      //"_$Data$_.Pos" => "Данные.Проведение",
	EventCausePredefinedDataInitialization             EventCauseType = "PredefinedDataInitialization"             //"_$Data$_.PredefinedDataInitialization" => "Данные.Инициализация предопределенных данных",
	EventCausePredefinedDataInitializationDataNotFound EventCauseType = "PredefinedDataInitializationDataNotFound" //"_$Data$_.PredefinedDataInitializationDataNotFound" => "Данные.Инициализация предопределенных данных.Данные не найдены",
	EventCauseSetPredefinedDataInitialization          EventCauseType = "SetPredefinedDataInitialization"          //"_$Data$_.SetPredefinedDataInitialization" => "Данные.Установка инициализации предопределенных данных",
	EventCauseSetStandardODataInterfaceContent         EventCauseType = "SetStandardODataInterfaceContent"         //"_$Data$_.SetStandardODataInterfaceContent" => "Данные.Изменение состава стандартного интерфейса OData",
	EventCauseTotalsMaxPeriodUpdate                    EventCauseType = "TotalsMaxPeriodUpdate"                    //"_$Data$_.TotalsMaxPeriodUpdate" => "Данные.Изменение максимального периода рассчитанных итогов",
	EventCauseTotalsMinPeriodUpdate                    EventCauseType = "TotalsMinPeriodUpdate"                    //"_$Data$_.TotalsMinPeriodUpdate" => "Данные.Изменение минимального периода рассчитанных итогов",
	EventCausePost                                     EventCauseType = "Post"                                     //"_$Data$_.Post" => "Данные.Проведение",
	EventCauseUnpost                                   EventCauseType = "Unpost"                                   //"_$Data$_.Unpost" => "Данные.Отмена проведения",
	EventCauseUpdate                                   EventCauseType = "Update"                                   //"_$Data$_.Update" => "Данные.Изменение",
	EventCauseUpdatePredefinedData                     EventCauseType = "UpdatePredefinedData"                     //"_$Data$_.UpdatePredefinedData" => "Данные.Изменение предопределенных данных",
	EventCauseVersionCommentUpdate                     EventCauseType = "VersionCommentUpdate"                     //"_$Data$_.VersionCommentUpdate" => "Данные.Изменение комментария версии",

	EventCauseConfigExtensionUpdate            EventCauseType = "ConfigExtensionUpdate"            //"_$InfoBase$_.ConfigExtensionUpdate" => "Информационная база.Изменение расширения конфигурации",
	EventCauseConfigUpdate                     EventCauseType = "ConfigUpdate"                     //"_$InfoBase$_.ConfigUpdate" => "Информационная база.Изменение конфигурации",
	EventCauseDBConfigBackgroundUpdateCancel   EventCauseType = "DBConfigBackgroundUpdateCancel"   //"_$InfoBase$_.DBConfigBackgroundUpdateCancel" => "Информационная база.Отмена фонового обновления",
	EventCauseDBConfigBackgroundUpdateFinish   EventCauseType = "DBConfigBackgroundUpdateFinish"   //"_$InfoBase$_.DBConfigBackgroundUpdateFinish" => "Информационная база.Завершение фонового обновления",
	EventCauseDBConfigBackgroundUpdateResume   EventCauseType = "DBConfigBackgroundUpdateResume"   //"_$InfoBase$_.DBConfigBackgroundUpdateResume" => "Информационная база.Продолжение (после приостановки) процесса фонового обновления",
	EventCauseDBConfigBackgroundUpdateStart    EventCauseType = "DBConfigBackgroundUpdateStart"    //"_$InfoBase$_.DBConfigBackgroundUpdateStart" => "Информационная база.Запуск фонового обновления",
	EventCauseDBConfigBackgroundUpdateSuspend  EventCauseType = "DBConfigBackgroundUpdateSuspend"  //"_$InfoBase$_.DBConfigBackgroundUpdateSuspend" => "Информационная база.Приостановка (пауза) процесса фонового обновления",
	EventCauseDBConfigExtensionUpdate          EventCauseType = "DBConfigExtensionUpdate"          //"_$InfoBase$_.DBConfigExtensionUpdate" => "Информационная база.Изменение расширения конфигурации",
	EventCauseDBConfigExtensionUpdateError     EventCauseType = "DBConfigExtensionUpdateError"     //"_$InfoBase$_.DBConfigExtensionUpdateError" => "Информационная база.Ошибка изменения расширения конфигурации",
	EventCauseDBConfigUpdate                   EventCauseType = "DBConfigUpdate"                   //"_$InfoBase$_.DBConfigUpdate" => "Информационная база.Изменение конфигурации базы данных",
	EventCauseDBConfigUpdateStart              EventCauseType = "DBConfigUpdateStart"              //"_$InfoBase$_.DBConfigUpdateStart" => "Информационная база.Запуск обновления конфигурации базы данных",
	EventCauseDumpError                        EventCauseType = "DumpError"                        //"_$InfoBase$_.DumpError" => "Информационная база.Ошибка выгрузки в файл",
	EventCauseDumpFinish                       EventCauseType = "DumpFinish"                       //"_$InfoBase$_.DumpFinish" => "Информационная база.Окончание выгрузки в файл",
	EventCauseDumpStart                        EventCauseType = "DumpStart"                        //"_$InfoBase$_.DumpStart" => "Информационная база.Начало выгрузки в файл",
	EventCauseEraseData                        EventCauseType = "EraseData"                        //"_$InfoBase$_.EraseData" => " Информационная база.Удаление данных информационной баз",
	EventCauseEventLogReduce                   EventCauseType = "EventLogReduce"                   //"_$InfoBase$_.EventLogReduce" => "Информационная база.Сокращение журнала регистрации",
	EventCauseEventLogReduceError              EventCauseType = "EventLogReduceError"              //"_$InfoBase$_.EventLogReduceError" => "Информационная база.Ошибка сокращения журнала регистрации",
	EventCauseEventLogSettingsUpdate           EventCauseType = "EventLogSettingsUpdate"           //"_$InfoBase$_.EventLogSettingsUpdate" => "Информационная база.Изменение параметров журнала регистрации",
	EventCauseEventLogSettingsUpdateError      EventCauseType = "EventLogSettingsUpdateError"      //"_$InfoBase$_.EventLogSettingsUpdateError" => "Информационная база.Ошибка при изменение настроек журнала регистрации",
	EventCauseInfoBaseAdmParamsUpdate          EventCauseType = "InfoBaseAdmParamsUpdate"          //"_$InfoBase$_.InfoBaseAdmParamsUpdate" => "Информационная база.Изменение параметров информационной базы",
	EventCauseInfoBaseAdmParamsUpdateError     EventCauseType = "InfoBaseAdmParamsUpdateError"     //"_$InfoBase$_.InfoBaseAdmParamsUpdateError" => "Информационная база.Ошибка изменения параметров информационной базы",
	EventCauseIntegrationServiceActiveUpdate   EventCauseType = "IntegrationServiceActiveUpdate"   //"_$InfoBase$_.IntegrationServiceActiveUpdate" => "Информационная база.Изменение активности сервиса интеграции",
	EventCauseIntegrationServiceSettingsUpdate EventCauseType = "IntegrationServiceSettingsUpdate" //"_$InfoBase$_.IntegrationServiceSettingsUpdate" => "Информационная база.Изменение настроек сервиса интеграции",
	EventCauseMasterNodeUpdate                 EventCauseType = "MasterNodeUpdate"                 //"_$InfoBase$_.MasterNodeUpdate" => "Информационная база.Изменение главного узла",
	EventCausePredefinedDataUpdate             EventCauseType = "PredefinedDataUpdate"             //"_$InfoBase$_.PredefinedDataUpdate" => "Информационная база.Обновление предопределенных данных",
	EventCauseRegionalSettingsUpdate           EventCauseType = "RegionalSettingsUpdate"           //"_$InfoBase$_.RegionalSettingsUpdate" => "Информационная база.Изменение региональных установок",
	EventCauseRestoreError                     EventCauseType = "RestoreError"                     //"_$InfoBase$_.RestoreError" => "Информационная база.Ошибка загрузки из файла",
	EventCauseRestoreFinish                    EventCauseType = "RestoreFinish"                    //"_$InfoBase$_.RestoreFinish" => "Информационная база.Окончание загрузки из файла",
	EventCauseRestoreStart                     EventCauseType = "RestoreStart"                     //"_$InfoBase$_.RestoreStart" => "Информационная база.Начало загрузки из файла",
	EventCauseSecondFactorAuthTemplateDelete   EventCauseType = "SecondFactorAuthTemplateDelete"   //"_$InfoBase$_.SecondFactorAuthTemplateDelete" => "Информационная база.Удаление шаблона вторго фактора аутентификации",
	EventCauseSecondFactorAuthTemplateNew      EventCauseType = "SecondFactorAuthTemplateNew"      //"_$InfoBase$_.SecondFactorAuthTemplateNew" => "Информационная база.Добавление шаблона вторго фактора аутентификации",
	EventCauseSecondFactorAuthTemplateUpdate   EventCauseType = "SecondFactorAuthTemplateUpdate"   //"_$InfoBase$_.SecondFactorAuthTemplateUpdate" => "Информационная база.Изменение шаблона вторго фактора аутентификации",
	EventCauseSetPredefinedDataUpdate          EventCauseType = "SetPredefinedDataUpdate"          //"_$InfoBase$_.SetPredefinedDataUpdate" => "Информационная база.Установить обновление предопределенных данных",

	EventCauseTARImportant EventCauseType = "TARImportant" //"_$InfoBase$_.TARImportant" => "Тестирование и исправление.Ошибка",
	EventCauseTARInfo      EventCauseType = "TARInfo"      //"_$InfoBase$_.TARInfo" => "Тестирование и исправление.Сообщение",
	EventCauseTARMess      EventCauseType = "TARMess"      //"_$InfoBase$_.TARMess" => "Тестирование и исправление.Предупреждение",

	EventCauseCancel    EventCauseType = "Cancel"    //"_$Job$_.Cancel" => "Фоновое задание.Отмена",
	EventCauseFail      EventCauseType = "Fail"      //"_$Job$_.Fail" => "Фоновое задание.Ошибка выполнения",
	EventCauseStart     EventCauseType = "Start"     //"_$Job$_.Start" => "Фоновое задание.Запуск",
	EventCauseSucceed   EventCauseType = "Succeed"   //"_$Job$_.Succeed" => "Фоновое задание.Успешное завершение",
	EventCauseTerminate EventCauseType = "Terminate" //"_$Job$_.Terminate" => "Фоновое задание.Принудительное завершение",

	EventCauseNegativeAssertion EventCauseType = "NegativeAssertion" //"_$OpenIDProvider$_.NegativeAssertion" => "Провайдер OpenID.Отклонено",
	EventCausePositiveAssertion EventCauseType = "PositiveAssertion" //"_$OpenIDProvider$_.PositiveAssertion" => "Провайдер OpenID.Подтверждено",

	EventCauseAuthentication            EventCauseType = "Authentication"            //"_$Session$_.Authentication" => "Сеанс.Аутентификация",
	EventCauseAuthenticationError       EventCauseType = "AuthenticationError"       //"_$Session$_.AuthenticationError" => "Сеанс.Ошибка аутентификации",
	EventCauseAuthenticationFirstFactor EventCauseType = "AuthenticationFirstFactor" //"_$Session$_.AuthenticationFirstFactor" => "Сеанс.Аутентификация первый фактор",
	EventCauseConfigExtensionApplyError EventCauseType = "ConfigExtensionApplyError" //"_$Session$_.ConfigExtensionApplyError" => "Сеанс.Ошибка применения расширения конфигурации",
	EventCauseFinish                    EventCauseType = "Finish"                    //"_$Session$_.Finish" => "Сеанс.Завершение",
	EventCauseBegin                     EventCauseType = "Begin"                     //"_$Session$_.Start" => "Сеанс.Начало", "_$Transaction$_.Begin" => "Транзакция.Начало",

	EventCauseCommit   EventCauseType = "Commit"   //"_$Transaction$_.Commit" => "Транзакция.Фиксация",
	EventCauseRollback EventCauseType = "Rollback" //"_$Transaction$_.Rollback" => "Транзакция.Отмена",

	EventCauseAuthenticationLock        EventCauseType = "AuthenticationLock"        //"_$User$_.AuthenticationLock" => "Пользователи.Блокировка аутентификации",
	EventCauseAuthenticationUnlock      EventCauseType = "AuthenticationUnlock"      //"_$User$_.AuthenticationUnlock" => "Пользователи.Разблокировка аутентификации",
	EventCauseAuthenticationUnlockError EventCauseType = "AuthenticationUnlockError" //"_$User$_.AuthenticationUnlockError " => "Пользователи.Ошибка разблокировки аутентификации",
	EventCauseDeleteError               EventCauseType = "DeleteError"               //"_$User$_.DeleteError" => "Пользователи.Ошибка удаления",
	EventCauseUpdateError               EventCauseType = "UpdateError"               //"_$User$_.UpdateError" => "Пользователи. Ошибка изменения",
	EventCauseNewError                  EventCauseType = "NewError"                  //"_$User$_.NewError" => "Пользователи.Ошибка добавления",
	//EventCauseNew = "New" //"_$User$_.New" => "Пользователи.Добавление",
	//EventCauseUpdate= "Update" //"_$User$_.Update" => "Пользователи.Изменение",
	//EventCauseDelete= "Delete" //"_$User$_.Delete" => "Пользователи.Удаление",
)

type EventType string

func (t EventType) Scope() EventScopeType {
	scope, _ := t.getScopeCause()
	return scope
}

func (t EventType) Cause() EventCauseType {
	_, cause := t.getScopeCause()
	return cause
}

func (t EventType) getScopeCause() (scope EventScopeType, cause EventCauseType) {

	scopeCause := strings.Split(string(t), ".")

	if len(scopeCause) == 1 {
		return EventScopeType(scopeCause[0]), ""
	}

	if len(scopeCause) == 2 {
		return EventScopeType(scopeCause[0]), EventCauseType(scopeCause[1])
	}

	return EventScopeUndefined, ""
}

func (t EventType) String() string {

	scope, cause := t.getScopeCause()

	if len(cause) == 0 {
		return scope.String()
	}

	return scope.String() + "." + cause.String()
}
