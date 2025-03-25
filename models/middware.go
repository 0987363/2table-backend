package models

const (
	MiddwareKeyLogger              = "Logger"
	MiddwareKeyMongodb             = "Mongodb"
	MiddwareKeyGoogle              = "Google"
	MiddwareKeyMicro               = "Micro"
	MiddwareKeyKafkaMicro          = "MicroKafka"
	MiddwareKeyRedis               = "Redis"
	MiddwareKeyRequestID           = "RequestID"
	MiddwareKeyTimestamp           = "Timestamp"
	MiddwareKeyGeoInfo             = "UserGeoInfo"
	MiddwareKeyStack               = "StackClient"
	MiddwareKeySupportedEncryption = "SupportedEncryption"
	MiddwareFieldPrevUUID          = "PrevUUID"
	MiddwareFieldPrevMessageCount  = "PrevCount"
	MiddwareFieldSource            = "Source"
	MiddwareFieldSourceId          = "SourceId"
)
const (
	MiddwareFieldUUID            = "UUID"
	MiddwareFieldConnectionID    = "ConnectionID"
	MiddwareFieldBegin           = "Begin"
	MiddwareFieldClientString    = "ClientString"
	MiddwareFieldConnectionUUIDS = "Connection-UUIDS"
)

const (
	MicroKeyRequestID           = "Request-Id"
	MicroKeyConnectionID        = "Connection-Id"
	MicroKeySupportedEncryption = "Supported-Encryption"
	MicroKeyRequestUUID         = "Request-Uuid"
	MicroKeyPrevUUID            = "Prev-Uuid"
	MicroKeyMessageCount        = "Prev-Message-Count"
	MicroKeySource              = "Source"
	MicroKeySourceId            = "Source-Id"
)
