package service

const (
	/* Errors */

	ParseErr      = "ID parse error"
	ValidationErr = "Validation error"
	BodyReadErr   = "Body read error"

	/* User Errors */

	UserCreateErr = "User create error"
	UserReadErr   = "User read error"
	UserUpdateErr = "User update error"
	UserDeleteErr = "User delete error"
	EmailErr      = "Email is incorrect"
	PasswordErr   = "Password is incorrect"

	/* Session Errors */

	SessionCreateErr = "Session create error"
	SessionCloseErr  = "Session close error"
	CookieReadErr    = "Cookie read error"

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
	JSONDeserializingErr = "Deserializing error"

	/* DB Errors */

	DBReadErr = "Error reading from DB"

	/* OK Messages */

	UpdateOk = "User updated successfully"
	DeleteOk = "User deleted successfully"
)
