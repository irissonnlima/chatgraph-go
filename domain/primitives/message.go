package domain_primitives

type Message struct {
	TypeMessage    string
	ContentMessage string
}

func (m Message) Stringfy() string {
	return "TypeMessage: " + m.TypeMessage + " ContentMessage: " + m.ContentMessage + "\n"
}
