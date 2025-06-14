package structfields

// NOTE: 结构体字段作为 Provier.

// type Message struct {
// 	Content string
// 	Code    int
// }
//
// // NewMessage 注意，这里返回的是指针类型.
// func NewMessage(content string, code int) *Message {
// 	return &Message{
// 		Content: content,
// 		Code:    code,
// 	}
// }

type Content string

type Message struct {
	Content Content
	Code    int
}

// NewMessage 注意，这里返回的是指针类型
func NewMessage(content string, code int) *Message {
	return &Message{
		Content: Content(content),
		Code:    code,
	}
}
