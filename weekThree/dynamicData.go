package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)


type User3 struct {
	ID 			int
	RealName	string `unpack:"-"`
	Login		string
	Flags		int
}


func PrintReflect(u interface{}) error {
	val := reflect.ValueOf(u).Elem()
	fmt.Printf("%T have %d fields:\n", u, val.NumField())	// numField returns number of fields for struct
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i) // returns i-th field of struct
		typeField := val.Type().Field(i)

		fmt.Printf("\tname=%v, type=%v, value=%v, tag=%v\n", typeField.Name, typeField.Type.Kind(), valueField,
			typeField.Tag)
	}
	return nil
}


func dynReflect(){
	u := User3{
		ID:			42,
		RealName:	"Socrates",
		Flags:		32,
	}
	err := PrintReflect(&u)
	if err != nil {
		panic(err)
	}
}


func dyn() {
	data := []byte{
		128, 36, 17, 0,
		9, 0, 0, 0,
		118, 46, 114, 111, 109, 97, 110, 111, 118,
		16, 0, 0, 0,
	}

	u := User3{}

	err := unpackRef(&u, data)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)
}


// lets unpack byte slice from perl's pack function to struct
func unpackRef(u interface{}, data []byte) error {

	r := bytes.NewReader(data) // returns reader for reading and seeking in byte slice

	val := reflect.ValueOf(u).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if typeField.Tag.Get("unpack") == "-" {
			continue
		}

		switch typeField.Type.Kind() {
		case reflect.Int:
			var value uint32
			err := binary.Read(r, binary.LittleEndian, &value)
			if err != nil {
				return err
			}
		case reflect.String:
			var lenRaw uint32
			err := binary.Read(r, binary.LittleEndian, &lenRaw)
			if err != nil {
				return err
			}
			dataRaw := make([]byte, lenRaw)
			err = binary.Read(r, binary.LittleEndian, &dataRaw)
			if err != nil {
				return err
			}
			valueField.SetString(string(dataRaw))
		default:
			return fmt.Errorf("invalid type: %v for field %v", typeField.Type.Kind(), typeField.Name)
		}
	}
	return nil
}
