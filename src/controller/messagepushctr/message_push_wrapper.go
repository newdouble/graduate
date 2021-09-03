package messagepushctr

const(
	_    EventType = iota
	EventIDCongestionPopup
)

type EventType int

func (e EventType) String() string  {
	switch e {
	case EventIDCongestionPopup:
		return ""
	}
}

type MessagePushReq struct {
	UID string
	Map map[string]interface{}
}
