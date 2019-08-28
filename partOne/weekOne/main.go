package main

import (
	"fmt"
	"os"
)

func main() {
	// Note: len returns byte size not count of symbols
	fmt.Println(len("♥")) // returns 3 NOT 1

	//	Anonymous function, function without name, can be used to assign to variable, to return, just run func once, etc
	func (a string) {
		fmt.Println("Passed string: ", a)
	}("myString")

	// Functions as first class objects. means that you can:
	//	1) assign func to variable
	//	2) pass func as argument to another func
	//	3) return func from another func
	//	4) have func as field of struct

	// Type of function??? Example below defines type strFuncType as function that takes in string and returns nothing
	type strFuncType func(string)

	// This func takes in callback
	worker := func(callb strFuncType) {
		callb("callback")
	}

	worker(printer)
	//worker(simplePrinter) // compile time error as simplePrinter is not of type strFuncType

	// Closure - outer returns func which has pointer to varInClosure,
	// Thus as long as f exists varInClosure will exist in scope of f
	outer := func(a string) strFuncType {
		varInClosure := a
		return func(b string) {
			fmt.Println(&varInClosure == &a)
			fmt.Printf("outer a: %s, inner b: %s \n", varInClosure, b)
		}
	}

	f := outer("Prefix")
	f("B")

	structs()
	interfaces()
	err := practice(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
}


func printer(str string) {
	fmt.Println("Passed: ", str)
}


func simplePrinter() {
	fmt.Println("my text is permanent")
}
