package model

type Event struct {
	Action  string
	Actor   string
	Payload map[string]any
}

func NewEvent(action, actor string, payload map[string]any) Event {
	return Event{Action: action, Actor: actor, Payload: payload}
}
