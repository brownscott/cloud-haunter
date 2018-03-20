package types

type ActionType string

func (ot ActionType) String() string {
	return string(ot)
}

const (
	LOG_ACTION          = ActionType("log")
	NOTIFICATION_ACTION = ActionType("notification")
)

type Action interface {
	Execute(string, []*Instance)
}

type Message interface {
	Message() string
}

type Dispatcher interface {
	Send(message Message) error
}
