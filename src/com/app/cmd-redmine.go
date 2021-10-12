package main

import (
	"cmd-redmine-manager/auth"
	"fmt"
)

func main() {

	// const verzija = "v0.1"
	// const URL = "https://support.navigator.rs/"

	//Commandline utility za Redmine

	username, password, error := auth.Credentials()
    if error != nil{
        panic(error)
    }
    fmt.Printf("You are logged in as %v/%v\n",username,password)
}
