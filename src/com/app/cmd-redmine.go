package main

import (
	"cmd-redmine-manager/auth"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	// const verzija = "v0.1"
	// const URL = "https://support.navigator.rs/"

	//Commandline utility za Redmine

	username, password, error := auth.Credentials()
	if error != nil {
		panic(error)
	}
	fmt.Printf("You are logged in as %v\n", username)

	fmt.Println("Searching for task...")
	response, err := http.Get("https://" + username + ":" + password + "@support.navigator.rs/issues.json?issue_id=85523")
	// fmt.Println(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(responseData))
}
