package agent

import (
	"os"
	"os/exec"
)

const defaultGitEditor = "vi"

func try(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func BracketedOrEmpty(label string) string {
	if label == "" {
		return ""
	}
	return "(" + label + ")"
}

func (a *Agent) MaybeBreakingExclam() string {
	if a.Config.Commit.Breaking {
		return "!"
	}
	return ""
}

func (a *Agent) maybeBreakingSuffix() string {
	if a.Config.Commit.Breaking {
		return "BREAKING CHANGE: \n"
	}
	return ""
}
