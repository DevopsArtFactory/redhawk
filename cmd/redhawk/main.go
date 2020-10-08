package main

import (
	"context"
	"errors"
	"os"

	Logger "github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/redhawk/cmd/redhawk/app"
	"github.com/DevopsArtFactory/redhawk/pkg/color"
)

func main() {
	if err := app.Run(os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, context.Canceled) {
			Logger.Debugln("ignore error since context is cancelled:", err)
		} else {
			color.Red.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
