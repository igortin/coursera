package person

import (
	"fmt"
)

func ShowName(p Person) {
	fmt.Println(p.Name)
}

func NewPerson(id int, name string) *Person {
	return &Person{
		Name: name,
		Id:   id,
	}
}
