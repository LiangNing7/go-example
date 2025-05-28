package main

import (
	"fmt"

	"github.com/jinzhu/copier"
)

// User 定义源结构体.
type User struct {
	Name string
	Age  int32
	Role string
}

func (u *User) DoubleAge() int32 {
	return u.Age * 2
}

// Employee 定义目标结构体.
type Employee struct {
	Name      string
	Age       int32
	DoubleAge int32
	SuperRole string
}

func (e *Employee) Role(role string) {
	e.SuperRole = "Super " + role
}

func main() {
	user := User{
		Name: "LiangNing7",
		Age:  18,
		Role: "Admin",
	}
	emp := Employee{}

	copier.Copy(&emp, &user)
	fmt.Printf("%+v", emp)
}
