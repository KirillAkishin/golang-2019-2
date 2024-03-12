package main

import "fmt"

type Person struct {
	Id      int
	Name    string
	Address string
}

type Account struct {
	Id      int
	Name    string
	Cleaner func(string) string
	Owner   Person
	*Person
}

func main() {

	bb := Person{
		Name:    "Василий",
		Address: "Москва",
	}

	// полное объявление структуры
	var acc Account = Account{
		Id:     1,
		Name:   "rvasily",
		Person: &bb,
	}
	fmt.Printf("%#v\n", acc)

	// короткое объявление структуры
	acc.Owner = Person{2, "Romanov Vasily", "Moscow"}

	// fmt.Printf("%#v\n", acc)

	fmt.Println(acc.Name)
	fmt.Println(acc.Person.Name)

	bb.Name = "New Name"

	fmt.Println(acc.Name)
	fmt.Println(acc.Person.Name)

}