package domain_response

type TextMessage struct {
	Type    string
	Title   string
	Detail  string
	Caption string
}

type ButtonType string

const (
	Postback ButtonType = "postback"
	URL      ButtonType = "url"
)

type Button struct {
	Type   ButtonType
	Title  string
	Detail string
}

type ResponseMessage struct {
	TextMessage  TextMessage
	Buttons      []Button
	DiplayButton Button
}
