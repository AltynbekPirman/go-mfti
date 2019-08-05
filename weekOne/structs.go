package main

import "fmt"

type Account struct {
	Id int
	Name string
	Owner Person
	Person	// Embed person to account, Now account will have fields Age, Person.Name and Person.Id
}

type Person struct {
	Id int
	Name string
	Age int
}

var p Person


var a Account


func structs() {
	p = Person{
		Name: "myName",
		Age: 100,
	}
	a = Account{1, "myAccount", p, p}
	fmt.Printf("%#v \n", a)
	fmt.Println(a.Age)
	fmt.Println(a.Name, a.Person.Name)
}