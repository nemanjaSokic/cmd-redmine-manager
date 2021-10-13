package main

import (
	"cmd-redmine-manager/auth"
	"cmd-redmine-manager/redmine"
	"fmt"
	"os"
)

func main() {

	// const verzija = "v0.1"
	// const URL = "https://support.navigator.rs/"

	//Commandline utility za Redmine

	// issue_id := flag.Int("i", 0, "Issue id you are working on.Required.")
	// issue_url := flag.String("u", "", "Issue URL https://support.navigator.rs/issues/{issue_id}.")
	// uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")
	// flag.Parse()
	// if *issue_id==0 || strings.Compare(*issue_url,"")==0 {
	// 	panic(errors.New("Issue id (i) or issue url (u) must be passed. Look at redmine --help"))
	// }
	// fmt.Printf("textPtr: %s, metricPtr: %s, uniquePtr: %t\n", *textPtr, *metricPtr, *uniquePtr)

	username, password, error := auth.Credentials()
	if error != nil {
		panic(error)
	}
	fmt.Printf("You are logged in as %v\n", username)

	fmt.Println("Searching for task...")
	c := redmine.NewClient(fmt.Sprintf("%v:%v", username, password))
	issue, err := c.Issue(79738)
	if err != nil {
		fatal("Failed to show issue: %s\n", err)
	}
	fmt.Printf(`
Id: %d
Subject: %s
Project: %s
Tracker: %s
Status: %s
Priority: %s
Author: %s
CreatedOn: %s
UpdatedOn: %s
%s
`[1:],
		issue.Id,
		issue.Subject,
		issue.Project.Name,
		issue.Tracker.Name,
		issue.Status.Name,
		issue.Priority.Name,
		issue.Author.Name,
		issue.CreatedOn,
		issue.UpdatedOn,
		issue.Description)
}
func fatal(format string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, err)
	} else {
		fmt.Fprint(os.Stderr, format)
	}
	os.Exit(1)
}
