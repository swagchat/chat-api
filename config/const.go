package config

type ctxKey int

const (
	// AppName is Application name
	AppName = "chat-api"
	// APIVersion is API version
	APIVersion = "0"
	// BuildVersion is API build version
	BuildVersion = "0.9.1"

	// KeyLength is key length
	KeyLength = 32
	// TokenLength is token length
	TokenLength = 32

	// HeaderUserID is http header for userID
	HeaderUserID = "X-Sub"
	// HeaderUsername is http header for username
	HeaderUsername = "X-Preferred-Username"
	// HeaderWorkspace is http header for workspace
	HeaderWorkspace = "X-Realm"
	// HeaderClientID is http header for clientID
	HeaderClientID = "X-ClientId"
	// HeaderRealmRoles is http header for roles
	HeaderRealmRoles = "X-Realm-Roles"
	// HeaderAccountRoles is http header for account roles
	HeaderAccountRoles = "X-Account-Roles"

	CtxDsCfg ctxKey = iota
	CtxClientID
	CtxUserID
	CtxWorkspace
	CtxRoomUser
	CtxSubscription

	RoleGeneral int32 = 1

	RetrieveRoomMessagesDefaultLimit = 50
)
