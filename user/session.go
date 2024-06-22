package user

type Session struct {
	Id         int
	UserId     int
	Token      string
	ClientInfo string
	ExpireTime int64
}

//TODO: make sessions with cookies
