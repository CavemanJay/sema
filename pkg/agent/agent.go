package agent

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
)

type (
	Agent struct {
		Config      *Config
		repo        *git.Repository
		workTree    *git.Worktree
		commitTitle string
	}

	Config struct {
		Commit Commit
		Push   Push
	}

	Commit struct {
		Long     bool
		Breaking bool
	}

	Push struct {
		Force bool
		Tags  bool
	}
)

func New(config *Config) *Agent {
	return &Agent{
		Config: config,
	}
}

func (a *Agent) Init() (err error) {
	a.repo, err = git.PlainOpen(".")
	if err != nil {
		return
	}
	a.workTree, err = a.repo.Worktree()
	return
}

func (a *Agent) Title(label, scope, message string) {
	a.commitTitle = fmt.Sprintf("%s%s: %s", label, scope, message)
}

func (a *Agent) Add() error {
	return try(exec.Command("git", "add", "."))
}

func (a *Agent) Commit() (string, error) {
	if a.Config.Commit.Long {
		return a.longCommit()
	} else {
		return a.shortCommit()
	}
}

func (a *Agent) Push() error {
	args := []string{"push"}
	if a.Config.Push.Force {
		args = append(args, "--force")
	}
	if a.Config.Push.Tags {
		args = append(args, "--tags")
	}
	return try(exec.Command("git", args...))
}

func (a *Agent) longCommit() (string, error) {
	path, err := a.createCommitTemplate()
	if err != nil {
		return "", fmt.Errorf("failed to create commit template file: %s", err)
	}
	msg, err := editCommitTemplate(path)
	if err != nil {
		return "", fmt.Errorf("failed to edit template: %s", err)
	}
	_, err = a.workTree.Commit(msg, &git.CommitOptions{})
	return msg, err
}

func (a *Agent) createCommitTemplate() (string, error) {
	file, err := os.CreateTemp("", "sema-commit-template-")
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.WriteString(a.commitTitle + "\n\n" + a.maybeBreakingSuffix())
	return file.Name(), err
}

func editCommitTemplate(path string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		if err := try(exec.Command("cmd", "/C", editor(), path)); err != nil {
			return "", err
		}
	case "linux":
		cmd := exec.Command("bash", "-c", fmt.Sprintf(`%s %s`, editor(), path))
		if err := try(cmd); err != nil {
			return "", err
		}
	default:
		log.Fatalf("Running on an unsupported OS: %s\n", runtime.GOOS)
	}
	return readCommitMessageFromTemplate(path)
}

func editor() string {
	output, err := exec.Command("git", "var", "GIT_EDITOR").Output()
	if err != nil {
		return defaultGitEditor
	}
	return strings.TrimSpace(string(output))
}

func readCommitMessageFromTemplate(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	return string(contents), err
}

func (a *Agent) shortCommit() (string, error) {
	// display(a.commitTitle)
	_, err := a.workTree.Commit(a.commitTitle, &git.CommitOptions{})
	return a.commitTitle, err
}
