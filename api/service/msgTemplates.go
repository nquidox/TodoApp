package service

const (
	/* Errors */

	ParseErr          = "ID parse error"
	ValidationErr     = "Validation error"
	BodyReadErr       = "Body read error"
	ServerResponseErr = "Server response error"
	WriteBytesErr     = "Error sending data"

	/* User Errors */

	UserCreateErr = "User create error"
	UserReadErr   = "User read error"
	UserUpdateErr = "User update error"
	UserDeleteErr = "User delete error"
	EmailErr      = "Email is incorrect"
	PasswordErr   = "Password is incorrect"

	/* Session Errors */

	SessionCreateErr   = "Session create error"
	SessionCloseErr    = "Session close error"
	CookieReadErr      = "Cookie read error"
	InvalidTokenErr    = "Invalid token"
	TokenReadErr       = "Error getting token"
	TokenValidationErr = "Error validating token"
	UUIDParseErr       = "Error parsing uuid"

	/* TODO Lists Errors */

	ListCreateErr = "List create error"
	ListReadErr   = "List read error"
	ListUpdateErr = "List update error"
	ListDeleteErr = "List delete error"

	/* TODO Tasks Errors */

	TaskCreateErr = "Task create error"
	TaskReadErr   = "Task read error"
	TaskUpdateErr = "Task update error"
	TaskDeleteErr = "Task delete error"

	/* JSON Errors */

	JSONReadErr          = "JSON read error"
	JSONSerializingErr   = "Serializing error"
	JSONDeserializingErr = "Deserializing error"

	/* DB Errors */

	TableInitErr = "Table init error"
	DBReadErr    = "Error reading from DB"
	DBNotFound   = "record not found"

	/* OK Messages */

	LoginSuccess  = "Login success"
	LogoutSuccess = "Logout success"

	UserCreateSuccess = "User created successfully"
	UserReadSuccess   = "User read successfully"
	UserUpdateSuccess = "User updated successfully"
	UserDeleteSuccess = "User deleted successfully"

	TodoListCreateSuccess = "Todo list created successfully"
	TodoListReadSuccess   = "Todo list read successfully"
	TodoListUpdateSuccess = "Todo list updated successfully"
	TodoListDeleteSuccess = "Todo list deleted successfully"

	TaskCreateSuccess = "Task created successfully"
	TaskReadSuccess   = "Task read successfully"
	TaskUpdateSuccess = "Task updated successfully"
	TaskDeleteSuccess = "Task deleted successfully"
)
