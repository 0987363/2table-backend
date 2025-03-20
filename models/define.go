package models

const (
	AuthenticationHeader            = "X-Druid-Authentication"             // 账号密码token
	AuthenticationApplicationHeader = "X-Druid-Application-Authentication" // 应用token
)

const (
	GRPCSuccess        = 0
	GRPCDBError        = 1
	GRPCRequestFalt    = 2
	GRPCSQLError       = 3
	GRPCSQLNotFound    = 4
	GRPCDBUnknownError = 5
)

const (
	InvalidDataValue = -99999
)
