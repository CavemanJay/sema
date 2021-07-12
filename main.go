package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	. "github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
)

const help = `Labels explained:

    Feat:     new feature for the user
    Fix:      bug fix for the user
    Docs:     changes to the documentation
    Style:    formatting with no production code change
    Refactor: refactoring production code
    Test:     adding missing tests, refactoring tests
    Chore:    updating grunt tasks
`

func init() {
	h := flag.Bool("help", false, "Display help message")
	flag.Parse()
	if *h {
		fmt.Println(help)
		os.Exit(0)
	}
}

func main() {
	label, err := label()
	if err != nil {
		log.Fatal(err)
	}

	scope, err := scope()
	if err != nil {
		log.Fatal(err)
	}

	message, err := message()
	if err != nil {
		log.Fatal(err)
	}

	command := fmt.Sprintf("%s(%s): %s", label, scope, message)
	fmt.Println("Commit: ", Green(command), "\n")

	commit := exec.Command("git", "commit", "-m", command)

	var out bytes.Buffer
	commit.Stdout = &out

	if err = commit.Run(); err != nil {
		fmt.Println(Red(out.String()))
		fmt.Println(Red(err.Error()))
		return
	}
	fmt.Println(out.String())
}

func label() (choice string, err error) {
	prompt := promptui.Select{
		Label: "Select commit label",
		Items: []string{
			"Feat",
			"Fix",
			"Docs",
			"Style",
			"Refactor",
			"Test",
			"Chore",
		},
	}
	_, choice, err = prompt.Run()
	return
}

func scope() (string, error) {
	valiadtor := func(input string) (err error) {
		if len(input) > 10 {
			return errors.New("input too long")
		}
		return
	}
	prompt := promptui.Prompt{
		Label:    "Change scope",
		Validate: valiadtor,
	}
	return prompt.Run()
}

func message() (msg string, err error) {
	prompt := promptui.Prompt{Label: "Commit message"}
	msg, err = prompt.Run()
	return
}
