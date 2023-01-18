package types

import "fmt"

type User struct {
	Name string
	Age  int
	// 如果同属一个包，phone 可以被测试访问到，如果是不同包，就访问不到了
	phone string
}

func NewUser(name string, age int, phone string) User {
	return User{
		Name:  name,
		Age:   age,
		phone: phone,
	}
}

func NewUserPtr(name string, age int, phone string) *User {
	return &User{
		Name:  name,
		Age:   age,
		phone: phone,
	}
}

func (u User) GetAge() int {
	return u.Age
}

func (u User) GetPhone() string {
	return u.phone
}

func (u *User) ChangeName(newName string) {
	u.Name = newName
}

func (u User) private() {
	fmt.Println("private")
}

type Message string
