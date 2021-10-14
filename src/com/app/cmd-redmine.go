package main

import (
	"cmd-redmine-manager/auth"
	"cmd-redmine-manager/redmine"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	printVersion = flag.Bool("version", false, "print version")
)
var Client *redmine.Client

func initClient(key string) {
	Client = redmine.NewClient(key)
}

func main() {

	const version = "v0.1"
	const name = "RedmineCLI"

	username, password, error := auth.Credentials()
	if error != nil {
		panic(error)
	}
	fmt.Printf("You are logged in as %v\n", username)

	initClient(fmt.Sprintf("%v:%v", username, password))
	fmt.Printf("Client %v is initialized.\n", Client)

	flag.Parse()

	if *printVersion {
		fmt.Printf("%s %s\n", name, version)
		return
	}
	if flag.NArg() < 1 {
		usage()
	}

	rand.Seed(time.Now().UnixNano())

	switch flag.Arg(0) {
	case "t", "test":
		SelectProject()
	case "i", "issue":
		switch flag.Arg(1) {
		case "a", "add":
			createIssue()
		case "c", "create":
			if flag.NArg() == 4 {
				addIssue(flag.Arg(2), flag.Arg(3))
			} else {
				usage()
			}
		case "u", "update":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				updateIssue(id)
			} else {
				usage()
			}
		case "n", "notes":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				notesIssue(id)
			} else {
				usage()
			}
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				showIssue(id)
			} else {
				usage()
			}
		case "x", "close":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid issue id: %s\n", err)
				}
				closeIssue(id)
			} else {
				usage()
			}
		case "l", "list":
			listIssues(nil)
		case "p", "project":
			filter := &redmine.IssueFilter{
				ProjectId: fmt.Sprint(SelectProject()),
			}
			listIssues(filter)
		case "m", "mine":
			filter := &redmine.IssueFilter{
				AssignedToId: "me",
			}
			listIssues(filter)
		default:
			usage()
		}
	case "p", "project":
		switch flag.Arg(1) {
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid project id: %s\n", err)
				}
				showProject(id)
			} else {
				usage()
			}
		case "l", "list":
			listProjects()
		default:
			usage()
		}
	case "m", "membership":
		switch flag.Arg(1) {
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid membership id: %s\n", err)
				}
				showMembership(id)
			} else {
				usage()
			}
		case "l", "list":
			if flag.NArg() == 3 {
				projectId, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid project id: %s\n", err)
				}
				listMemberships(projectId)
			} else {
				usage()
			}
		default:
			usage()
		}
	case "u", "user":
		switch flag.Arg(1) {
		case "s", "show":
			if flag.NArg() == 3 {
				id, err := strconv.Atoi(flag.Arg(2))
				if err != nil {
					fatal("Invalid user id: %s\n", err)
				}
				showUser(id)
			} else {
				usage()
			}
		case "l", "list":
			listUsers()
		default:
			usage()
		}
	default:
		usage()
	}
}
func fatal(format string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, err)
	} else {
		fmt.Fprint(os.Stderr, format)
	}
	os.Exit(1)
}
func usage() {
	fmt.Println(`godmine <command> <subcommand> [arguments]
Project Commands:
  show     s show given project.
             $ godmine p s 1
  list     l listing projects.
             $ godmine p l
Issue Commands:
  add      a create issue with text editor.
             $ godmine i a
  create   c create issue from given arguments.
             $ godmine i c subject description
  update   u update given issue.
             $ godmine i u 1
  show     s show given issue.
             $ godmine i s 1
  close    x close given issue.
             $ godmine i x 1
  notes    n add notes to given issue.
             $ godmine i n 1
  list     l listing issues.
             $ godmine i l
Membership Commands:
  show     s show given membership.
             $ godmine m s 1
  list     l listing memberships of given project.
             $ godmine m l 1
User Commands:
  show     s show given user.
             $ godmine u s 1
  list     l listing users.
             $ godmine u l
 `)
	os.Exit(1)
}

func addIssue(subject, description string) {
	var issue redmine.Issue
	issue.ProjectId = SelectProject()
	issue.Subject = subject
	issue.Description = description
	_, err := Client.CreateIssue(issue)
	if err != nil {
		fatal("Failed to create issue: %s\n", err)
	}
}

func createIssue() {
	issue, err := issueFromEditor("")
	if err != nil {
		fatal("%s\n", err)
	}
	issue.ProjectId = SelectProject()
	_, err = Client.CreateIssue(*issue)
	if err != nil {
		fatal("Failed to create issue: %s\n", err)
	}
}

func updateIssue(id int) {
	issue, err := Client.Issue(id)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
	issueNew, err := issueFromEditor(fmt.Sprintf("%s\n%s\n", issue.Subject, issue.Description))
	if err != nil {
		fatal("%s\n", err)
	}
	issue.Subject = issueNew.Subject
	issue.Description = issueNew.Description
	issue.ProjectId = SelectProject()
	err = Client.UpdateIssue(*issue)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
}
func closeIssue(id int) {
	issue, err := Client.Issue(id)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
	is, err := Client.IssueStatuses()
	if err != nil {
		fatal("Failed to get issue statuses: %s\n", err)
	}
	for _, s := range is {
		if s.IsClosed {
			issue.StatusId = s.Id
			err = Client.UpdateIssue(*issue)
			if err != nil {
				fatal("Failed to update issue: %s\n", err)
			}
			break
		}
	}
}

func notesIssue(id int) {
	issue, err := Client.Issue(id)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}

	content, err := notesFromEditor(issue)
	if err != nil {
		fatal("%s\n", err)
	}
	issue.Notes = content
	issue.ProjectId = SelectProject()
	err = Client.UpdateIssue(*issue)
	if err != nil {
		fatal("Failed to update issue: %s\n", err)
	}
}

func showIssue(id int) {
	issue, err := Client.Issue(id)
	if err != nil {
		fatal("Failed to show issue: %s\n", err)
	}
	assigned := ""
	if issue.AssignedTo != nil {
		assigned = issue.AssignedTo.Name
	}

	fmt.Printf(`
Id: %d
Subject: %s
Project: %s
Tracker: %s
Status: %s
Priority: %s
Author: %s
Assigned: %s
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
		assigned,
		issue.CreatedOn,
		issue.UpdatedOn,
		issue.Description)
}

func listIssues(filter *redmine.IssueFilter) {
	issues, err := Client.IssuesByFilter(filter)
	if err != nil {
		fatal("Failed to list issues: %s\n", err)
	}
	for _, i := range issues {
		fmt.Printf("%4d: %s\n", i.Id, i.Subject)
	}
}

func showProject(id int) {
	project, err := Client.Project(id)
	if err != nil {
		fatal("Failed to show project: %s\n", err)
	}

	fmt.Printf(`
Id: %d
Name: %s
Identifier: %s
CreatedOn: %s
UpdatedOn: %s
%s
`[1:],
		project.Id,
		project.Name,
		project.Identifier,
		project.CreatedOn,
		project.UpdatedOn,
		project.Description)
}

func listProjects() {
	issues, err := Client.Projects()
	if err != nil {
		fatal("Failed to list projects: %s\n", err)
	}
	for _, i := range issues {
		fmt.Printf("%4d: %s\n", i.Id, i.Name)
	}
}

func showMembership(id int) {
	membership, err := Client.Membership(id)
	if err != nil {
		fatal("Failed to show membership: %s\n", err)
	}

	fmt.Printf(`
Id: %d
Project: %s
User: %s
Role: `[1:],
		membership.Id,
		membership.Project.Name,
		membership.User.Name)
	for i, role := range membership.Roles {
		if i != 0 {
			fmt.Print(", ")
		}
		fmt.Printf(role.Name)
	}
	fmt.Println()
}

func listMemberships(projectId int) {
	memberships, err := Client.Memberships(projectId)
	if err != nil {
		fatal("Failed to list memberships: %s\n", err)
	}
	for _, i := range memberships {
		fmt.Printf("%4d: %s\n", i.Id, i.User.Name)
	}
}

func showUser(id int) {
	user, err := Client.User(id)
	if err != nil {
		fatal("Failed to show user: %s\n", err)
	}

	fmt.Printf(`
Id: %d
Login: %s
Firstname: %s
Lastname: %s
Mail: %s
CreatedOn: %s
`[1:],
		user.Id,
		user.Login,
		user.Firstname,
		user.Lastname,
		user.Mail,
		user.CreatedOn)
}

func listUsers() {
	users, err := Client.Users()
	if err != nil {
		fatal("Failed to list users: %s\n", err)
	}
	for _, i := range users {
		fmt.Printf("%4d: %s\n", i.Id, i.Login)
	}
}

func issueFromEditor(contents string) (*redmine.Issue, error) {
	file := ""
	newf := fmt.Sprintf("%d.txt", rand.Int())
	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "godmine", newf)
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "godmine", newf)
	}
	defer os.Remove(file)
	editor := getEditor()

	if contents == "" {
		contents = "### Subject Here ###\n### Description Here ###\n"
	}

	ioutil.WriteFile(file, []byte(contents), 0600)

	if err := run([]string{editor, file}); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	text := string(b)

	if text == contents {
		return nil, errors.New("Canceled")
	}
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return nil, errors.New("Canceled")
	}
	var subject, description string
	if len(lines) == 1 {
		subject = lines[0]
	} else {
		subject, description = lines[0], strings.Join(lines[1:], "\n")
	}
	var issue redmine.Issue
	issue.Subject = subject
	issue.Description = description
	return &issue, nil
}

func getEditor() string {
	editor := ""
	if runtime.GOOS == "windows" {
		editor = "notepad"
	} else {
		editor = "vim"
	}
	return editor
}
func notesFromEditor(issue *redmine.Issue) (string, error) {
	file := ""
	newf := fmt.Sprintf("%d.txt", rand.Int())
	if runtime.GOOS == "windows" {
		file = filepath.Join(os.Getenv("APPDATA"), "godmine", newf)
	} else {
		file = filepath.Join(os.Getenv("HOME"), ".config", "godmine", newf)
	}
	defer os.Remove(file)
	editor := getEditor()

	body := "### Notes Here ###\n"
	contents := issue.GetTitle() + "\n" + body

	ioutil.WriteFile(file, []byte(contents), 0600)

	if err := run([]string{editor, file}); err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	text := strings.Join(strings.SplitN(string(b), "\n", 2)[1:], "\n")

	if text == body {
		return "", errors.New("Canceled")
	}
	return text, nil
}
func SelectProject() int {
	projects, err := Client.Projects()
	if err != nil {
		fatal("Failed to list projects: %s\n", err)
	}
	fmt.Println(projects)
	return 1
}
func run(argv []string) error {
	cmd, err := exec.LookPath(argv[0])
	if err != nil {
		return err
	}
	var stdin *os.File
	if runtime.GOOS == "windows" {
		stdin, _ = os.Open("CONIN$")
	} else {
		stdin = os.Stdin
	}
	p, err := os.StartProcess(cmd, argv, &os.ProcAttr{Files: []*os.File{stdin, os.Stdout, os.Stderr}})
	if err != nil {
		return err
	}
	defer p.Release()
	w, err := p.Wait()
	if err != nil {
		return err
	}
	if !w.Exited() || !w.Success() {
		return errors.New("failed to execute text editor")
	}
	return nil
}
