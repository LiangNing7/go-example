package main

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func main() {
	e, err := casbin.NewEnforcer("./model.pml", "./policy.csv")
	if err != nil {
		log.Fatalf("NewEnforcer failed:%v\n", err)
	}

	check(e, "zhangsan", "/index", "GET")
	check(e, "zhangsan", "/home", "GET")
	check(e, "zhangsan", "/users", "POST")
	check(e, "wangwu", "/users", "POST")
}
