package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cavemanjay/sema/v5/pkg/agent"
	"github.com/charmbracelet/huh"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

const gitHubURL = "https://github.com/cavemanjay/sema"

const (
	add      = "add"
	push     = "push"
	force    = "force"
	long     = "long"
	breaking = "breaking"
	tags     = "tags"
)

var version = "unknown"

var app = &cli.App{
	Name:        "sema",
	Usage:       "Semantic commits made simple",
	Description: gitHubURL,
	Version:     version,
	Authors: []*cli.Author{{
		Name:  "Viktor A. Rozenko Voitenko",
		Email: "sharp.vik@gmail.com",
	}, {
		Name:  "Antoine Langlois",
		Email: "antoine.l@antoine-langlois.net",
	}, {
		Name:  "Jaydlc",
		Email: "slash_opt@proton.me",
	}},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    add,
			Aliases: []string{"a"},
			Usage:   "begin by running 'git add'",
		},
		&cli.BoolFlag{
			Name:    push,
			Aliases: []string{"p"},
			Usage:   "run 'git push' on successful commit",
		},
		&cli.BoolFlag{
			Name:    force,
			Aliases: []string{"f"},
			Usage:   "force push changes with 'git push -f'",
		},
		&cli.BoolFlag{
			Name:    long,
			Aliases: []string{"l"},
			Usage:   "open editor to elaborate commit message",
		},
		&cli.BoolFlag{
			Name:    breaking,
			Aliases: []string{"b"},
			Usage:   "mark commit as introducing breaking changes",
		},
		&cli.BoolFlag{
			Name:    tags,
			Aliases: []string{"t"},
			Usage:   "push tags along with commits",
		},
	},
	UseShortOptionHandling: true,
	Action:                 run,
	Commands: []*cli.Command{{
		Name:   "github",
		Usage:  "Open sema GitHub repository in browser",
		Action: github,
	}},
}

func run(c *cli.Context) error {
	a := agent.New(config(c))
	if err := a.Init(); err != nil {
		return err
	}

	details, err := GetCommitDetails(a)
	if err != nil {
		return err
	}

	fmt.Println(huh.NewInput().Title("Commit Label ").Inline(true).Value(&details.CommitLabel).View())
	fmt.Println(huh.NewInput().Title("Change Scope ").Inline(true).Value(&details.ChangeScope).View())
	fmt.Println(huh.NewInput().Title("Commit Message ").Inline(true).Value(&details.CommitMessage).View())

	formattedDetails := details
	formattedDetails.ChangeScope = agent.BracketedOrEmpty(details.ChangeScope)
	formattedDetails.CommitLabel += a.MaybeBreakingExclam()

	a.Title(formattedDetails.CommitLabel, formattedDetails.ChangeScope, formattedDetails.CommitMessage)

	if c.Bool(add) {
		if err := a.Add(); err != nil {
			return err
		}
	}

	commitBody, err := a.Commit()
	if err != nil {
		return err
	}

	fmt.Println(huh.NewText().Title("Commit").Lines(lines(commitBody)).ShowLineNumbers(true).Value(&commitBody).View())

	if c.Bool(push) {
		if err := a.Push(); err != nil {
			return err
		}
	}

	return nil
	// return pipe(do.Init, do.Title).
	// 	thenIf(c.Bool(add), do.Add).
	// 	then(do.Commit).
	// 	thenIf(c.Bool(push), do.Push).
	// 	run()
}

func config(c *cli.Context) *agent.Config {
	return &agent.Config{
		Commit: agent.Commit{
			Long:     c.Bool(long),
			Breaking: c.Bool(breaking),
		},
		Push: agent.Push{
			Force: c.Bool(force),
			Tags:  c.Bool(tags),
		},
	}
}

func github(_ *cli.Context) error {
	browser.Stdout = io.Discard
	browser.Stderr = io.Discard
	return browser.OpenURL(gitHubURL)
}

func main() {
	if err := app.Run(os.Args); err != nil && !errors.Is(err, huh.ErrUserAborted) {
		fmt.Fprintln(os.Stderr, err)
	}
}

func lines(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}
