package main

import (
	"bytes"
	"encoding/xml"
	"io"
)

type User struct {
	ID		int `xml:"id,attr"`
	Login	string`xml:"login"`
	Name	string`xml:"name"`
	Browser	string`xml:"browser"`
}

type Users struct {
	Version string `xml:"version,attr"`
	List	[]User `xml:"user"`
}

func CountStruct() {
	logins := make([]string, 0)
	v := Users{}

	err := xml.Unmarshal(xmlData, &v)
	if err != nil {
		panic(err)
	}
	for _, u := range v.List {
		logins = append(logins, u.Login)
	}
}


func CountDecoder() {
	input := bytes.NewReader(xmlData)
	decoder := xml.NewDecoder(input)
	logins := make([]string, 0)
	var login string
	for {
		tok, tokErr := decoder.Token()
		if tokErr == io.EOF || tok == nil {
			break
		} else if tokErr != nil {
			panic(tokErr)
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "login" {
				err := decoder.DecodeElement(&login, &tok)
				if err != nil {
					panic(err)
				}
				logins = append(logins, login)
			}
		}
	}
}



var xmlData = []byte(`<?xml version="1.0" encoding="utf-8"?>
	<users>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
	</browser>
		</user>
	</users>`)
