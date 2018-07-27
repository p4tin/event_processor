package main

import "fmt"

type messageProc string

func (g messageProc) Register() map[string]string {
	return map[string]string{
		"Type":      "QueueTask",
		"Provider":  "AWS",
		"QueueName": "A15-DEV-PAUL-TEST-1",
	}
}

func (g messageProc) Process(data interface{}) error {
	fmt.Println("Hello From a plugin!!!")
	return nil
}

// exported as symbol named "Greeter"
var MessageProc messageProc
