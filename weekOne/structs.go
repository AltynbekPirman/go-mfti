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
	res := p.GetName()
	fmt.Println(res, p.Name)
	p.SetName("myNewName")
	fmt.Println(p.Name)
	a.SetName("newPersonNameFromAccount") // methods of embedded struct are inherited

	a.SetID(45); a.Person.SetID(25) // if both structs have same methods, original method will be prioritized
	fmt.Printf("%#v \n", a)
	s := MySlice{1, 2, 3} // methods can be created for other types other than struct
	fmt.Println(s.GetLength())

}

// does not change original struct as copy of struct is passed
func (p Person) GetName() string {
	res := p.Name
	p.Name = ""
	return res
}

// changes original struct as pointer struct is passed
func (p *Person) SetName(newName string) string {
	p.Name = newName
	return p.Name
}

func (p *Person) SetID(id int) {
	p.Id = id
}


func (a *Account) SetID(id int) {
	a.Id = id
}

type MySlice []int

func (s MySlice) GetLength() int {
	return len(s)
}
