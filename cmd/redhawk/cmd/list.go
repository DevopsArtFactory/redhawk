/*
Copyright 2020 The redhawk Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/DevopsArtFactory/redhawk/cmd/redhawk/cmd/builder"
	"github.com/DevopsArtFactory/redhawk/pkg/executor"
)

// List resources in detail
func NewCmdList() *cobra.Command {
	return builder.NewCmd("list").
		WithDescription("list infrastructure resources in AWS").
		SetFlags().
		RunWithCmdAndNoArgs(funcList)
}

// funcList
func funcList(ctx context.Context, out io.Writer, cmd *cobra.Command) error {
	return executor.RunExecutor(ctx, func(executor executor.Executor) error {
		return executor.Runner.List(out)
	})
}
