package main

import (
	"fmt"

	"github.com/cavemanjay/sema/v5/pkg/labels"
	"github.com/charmbracelet/huh"
)

type CommitDetails struct {
	CommitLabel   string
	ChangeScope   string
	CommitMessage string
}

func GetCommitDetails() (CommitDetails, error) {
	var (
		commitLabelOptions = make([]huh.Option[string], 0, len(labels.Labels()))
		details            CommitDetails
	)

	for _, l := range labels.Labels() {
		commitLabelOptions = append(commitLabelOptions, huh.Option[string]{Key: l.String(), Value: l.Name})
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Commit label").
				Options(commitLabelOptions...).
				Value(&details.CommitLabel)),
		huh.NewGroup(
			huh.NewInput().
				Title("Change Scope").
				Description("readme, database, etc.").
				Value(&details.ChangeScope)),
		huh.NewGroup(
			huh.NewInput().
				Title("Commit Message").Value(&details.CommitMessage).
				Validate(
					func(s string) error {
						if len(s) == 0 {
							return fmt.Errorf("commit message must be > 0 chars")
						}
						return nil
					},
				),
		),
	)

	if err := form.Run(); err != nil {
		return CommitDetails{}, err
	}

	return details, nil
}
