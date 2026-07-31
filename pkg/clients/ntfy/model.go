package ntfy

type MessagePriority int

const (
	MinPriority     MessagePriority = 1
	LowPriority     MessagePriority = 2 // No vibration or sound. Notification will not visibly show up until notification drawer is pulled down.
	DefaultPriority MessagePriority = 3 // Short default vibration and sound. Default notification behavior
	HighPriority    MessagePriority = 4 // Long vibration burst, default notification sound with a pop-over notification.
	MaxPriority     MessagePriority = 5 // Really long vibration bursts, default notification sound with a pop-over notification.
)

type Message struct {
	Topic    string           `json:"topic"`
	Text     string           `json:"text"`
	Title    *string          `json:"title"`
	Priority *MessagePriority `json:"priority"`
	Tags     *string          `json:"tags"`
	Click    *string          `json:"click"`
}

func (m *Message) SetTitle(t string) *Message {
	m.Title = &t
	return m
}

func (m *Message) SetPriority(p MessagePriority) *Message {
	m.Priority = &p
	return m
}

func (m *Message) SetTags(t string) *Message {
	m.Tags = &t
	return m
}

func (m *Message) SetClick(url string) *Message {
	m.Click = &url
	return m
}
