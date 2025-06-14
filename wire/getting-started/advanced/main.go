package main

import (
	"fmt"

	argsanderror "github.com/LiangNing7/go-example/wire/getting-started/advanced/args_and_error"
	"github.com/LiangNing7/go-example/wire/getting-started/advanced/bindinginterfaces"
	"github.com/LiangNing7/go-example/wire/getting-started/advanced/bindingstruct"
	"github.com/LiangNing7/go-example/wire/getting-started/advanced/bindingvalues"
	"github.com/LiangNing7/go-example/wire/getting-started/advanced/providersets"
	"github.com/LiangNing7/go-example/wire/getting-started/advanced/structfields"
	"github.com/LiangNing7/go-example/wire/getting-started/advanced/structproviders"
)

func main() {
	// inject args and error.
	{
		e, err := argsanderror.InitializeEvent("Hello World!")
		if err != nil {
			fmt.Println(err)
		}
		e.Start()
	}

	// Provider Sets.
	{
		e, err := providersets.InitializeEvent("Provider Sets")
		if err != nil {
			fmt.Println(err)
		}
		e.Start()
	}

	// Struct Providers
	{
		e, err := structproviders.InitializeEvent("Struct Providers", 1)
		if err != nil {
			fmt.Println(err)
		}
		e.Start()
	}

	// Struct fields
	{
		m := structfields.InitializeMessage("Struct fields", 1)
		fmt.Println(m)
	}

	// Binding Values
	{
		m := bindingvalues.InitializeMessage()
		fmt.Printf("%+v\n", m)
	}

	// Binding Interfaces
	{
		w := bindinginterfaces.InitializeWriter()
		bindinginterfaces.Write(w, "Binding Interfaces")
	}

	// Binding struct to interface
	{
		msg := &bindingstruct.Message{
			Content: "content",
			Code:    1,
		}
		err := bindingstruct.RunStore(msg)
		if err != nil {
			fmt.Println(err)
		}
		err = bindingstruct.WireRunStore(msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}
