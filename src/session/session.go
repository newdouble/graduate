package session

type Session struct {
	TraceID string
	HintCode string
	RealUID string
	Token string
	Ip string
}

func New() *Session {
	return &Session{}
}
