package main

import "fmt"

type Message string

func NewMessage() Message {
	return Message("Hi there!")
}

type Greeter struct {
	Message Message
}

// NewGreeter 把外部构造好的 Message 作为参数注入进来.
func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

type Event struct {
	Greeter Greeter
}

// NewEvent 把已经构造好的 Greeter 注入进来.
func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

// func main() {
// 	message := NewMessage()
// 	greeter := NewGreeter(message)
// 	event := NewEvent(greeter)
//
// 	event.Start()
// }

// func InitializeEvent() Event {
// 	message := NewMessage()
// 	greeter := NewGreeter(message)
// 	event := NewEvent(greeter)
// 	return event
// }

func main() {
	event := InitializeEvent()
	event.Start()
}
