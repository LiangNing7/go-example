package main

import (
	"encoding/json"
	"fmt"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

type Foo struct {
	Array []*structpb.Value `protobuf:"bytes,1,rep,name=array,proto3" json:"array,omitempty"`
}

func main() {
	l1, _ := structpb.NewList([]any{"1", "2"})
	l2, _ := structpb.NewList([]any{"3", "4"})
	p := Foo{
		Array: []*structpb.Value{
			structpb.NewListValue(l1),
			structpb.NewListValue(l2),
		},
	}

	d, _ := json.Marshal(p)
	fmt.Println(string(d))
}
