package client

import (
	"errors"
	"os"

	"github.com/AlecAivazis/survey/v2"

	"github.com/DevopsArtFactory/redhawk/pkg/color"
)

// SelectTarget select only one target from candidates
func SelectTarget(opt []string) (string, error) {
	if len(opt) == 1 {
		if len(opt[0]) == 0 {
			return opt[0], errors.New("endpoint is empty")
		}
		color.Blue.Fprintf(os.Stdout, "you have only one choice : "+opt[0])
		return opt[0], nil
	}

	var target string

	prompt := &survey.Select{
		Message: "Choose an instance:",
		Options: opt,
	}
	survey.AskOne(prompt, &target)

	if len(target) == 0 {
		return target, errors.New("choosing an instance has been canceled")
	}

	return target, nil
}
