package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/redhawk/cmd/redhawk/cmd/builder"
	"github.com/DevopsArtFactory/redhawk/pkg/executor"
)

// Scan resources in detail
func NewCmdScanResources() *cobra.Command {
	return builder.NewCmd("scan-resources").
		WithDescription("Scan infrastructure resources in AWS").
		SetFlags().
		RunWithNoArgs(funcScanResources)
}

// funcScanResources
func funcScanResources(ctx context.Context, out io.Writer) error {
	return executor.RunExecutor(ctx, func(executor executor.Executor) error {
		return executor.Runner.ScanResources(out)
	})
}
