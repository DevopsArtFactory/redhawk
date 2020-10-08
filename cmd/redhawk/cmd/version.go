package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/redhawk/cmd/redhawk/cmd/builder"
	"github.com/DevopsArtFactory/redhawk/pkg/version"
)

// Get redhawk version
func NewCmdVersion() *cobra.Command {
	return builder.NewCmd("version").
		WithDescription("Print the version information").
		SetAliases([]string{"v"}).
		RunWithNoArgs(funcVersion)
}

// funcVersion
func funcVersion(_ context.Context, _ io.Writer) error {
	return version.Controller{}.Print(version.Get())
}
