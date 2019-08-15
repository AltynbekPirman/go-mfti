package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	//jsonIntro()
	//jsonAdv()
	//unknownJson()
	//dynReflect()
	//dyn()
	//cgenn()
	//CountStruct()
	CountDecoder()
}

func jsonIntro(){
	type User3 struct {
		ID 			int
		Username 	string
		phone		string // field that starts with lower case is PRIVATE field
	}

	var jsonStr = `{"id": 24, "username": "Sartr", "phone": "123"}`

	data := []byte(jsonStr)
	u := &User3{}
	// Unmarshal takes in byte slice(our json string converted to []byte) and interface to which json will be written
	err := json.Unmarshal(data, u)
	if err != nil {
		panic(err)
	}
	fmt.Println(u) // -> &{24 Sartr} - Note phone will be empty as it is private and not available from encoding/json package

	u.phone = "44564"
	u.ID = 456
	// Marshal returns marshalled json as byte slice and error
	res, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res)) // {"ID":456,"Username":"Sartr"} - Note: NO phone in converted json as phone is private

}


type UserAdv struct {
	ID 			int `json:"user_id,string"` // user_id - json key, type of value is string
	Username	string
	Address		string `json:",omitempty"` // json key remains as field name and omitempty excludes from json if field is empty
	Company		string `json:"-"` // do not include to json
}

func jsonAdv() {
	user1 := UserAdv{
		ID: 45,	// json key must be user_id and value will be "45" NOT 45
		Address: "", // will be omitted as it is empty(default value)
		Username: "Kirkegaard",
		Company: "Danish church", // not included in json
	}

	res, err := json.Marshal(user1)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res)) // -> {"user_id":"45","Username":"Kirkegaard"}

	jsonUser := `{"user_id":"45","Username":"Kirkegaard", "Company": "someComp", "Address": "Kopenhagen"}`
	jsonUserByte := []byte(jsonUser)
	var user2 UserAdv
	err = json.Unmarshal(jsonUserByte, &user2)
	if err != nil {
		panic(err)
	}
	fmt.Println(user2) // {45 Kirkegaard Kopenhagen }
}

// if json structure is unknown or changes constantly, You can use empty interface to serialize and deserialize json
func unknownJson(){
	var jsonStr = `[{"id": 17, "username": "iivan", "phone": 0}, {"id": "17", "address": "none", "company": "Mail.ru"}]`
	data := []byte(jsonStr)

	var users []interface{}
	err := json.Unmarshal(data, &users) // users will be slice of empty interfaces consisting of map[string]interface {}
	if err != nil {
		panic(err)
	}
	for _, i := range users {
		fmt.Printf("%v\n", i)
	}
}
