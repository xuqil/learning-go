package demo

type User struct {
	Name string
	age  int
}

func NewUser() *User {
	return &User{
		Name: "Tom",
		age:  18,
	}
}
