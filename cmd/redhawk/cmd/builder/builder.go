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

package builder

import (
	"context"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Builder interface {
	WithDescription(description string) Builder
	WithLongDescription(description string) Builder
	SetAliases(alias []string) Builder
	AddCommands(children ...*cobra.Command) Builder
	SetFlags() Builder
	WithFlags(adder func(*pflag.FlagSet)) Builder
	RunWithNoArgs(action func(context.Context, io.Writer) error) *cobra.Command
	RunWithArgs(action func(context.Context, io.Writer, []string) error) *cobra.Command
	RunWithArgsAndCmd(action func(context.Context, io.Writer, *cobra.Command, []string) error) *cobra.Command
	RunWithCmdAndNoArgs(action func(context.Context, io.Writer, *cobra.Command) error) *cobra.Command
}

type builder struct {
	cmd cobra.Command
}

// NewCmd creates a new command builder.
func NewCmd(use string) Builder {
	return &builder{
		cmd: cobra.Command{
			Use: use,
		},
	}
}

// Write short description
func (b builder) WithDescription(description string) Builder {
	b.cmd.Short = description
	return b
}

// Write long description
func (b builder) WithLongDescription(description string) Builder {
	b.cmd.Long = description
	return b
}

// Set command alias
func (b builder) SetAliases(alias []string) Builder {
	b.cmd.Aliases = alias
	return b
}

//Run command without Argument
func (b builder) RunWithNoArgs(function func(context.Context, io.Writer) error) *cobra.Command {
	b.cmd.Args = cobra.NoArgs
	b.cmd.RunE = func(*cobra.Command, []string) error {
		return returnErrorFromFunction(function(b.cmd.Context(), b.cmd.OutOrStderr()))
	}
	return &b.cmd
}

// Run command with extra arguments
func (b builder) RunWithArgs(function func(context.Context, io.Writer, []string) error) *cobra.Command {
	b.cmd.RunE = func(_ *cobra.Command, args []string) error {
		return returnErrorFromFunction(function(b.cmd.Context(), b.cmd.OutOrStderr(), args))
	}
	return &b.cmd
}

// Run command with extra arguments
func (b builder) RunWithArgsAndCmd(function func(context.Context, io.Writer, *cobra.Command, []string) error) *cobra.Command {
	b.cmd.RunE = func(_ *cobra.Command, args []string) error {
		return returnErrorFromFunction(function(b.cmd.Context(), b.cmd.OutOrStderr(), &b.cmd, args))
	}
	return &b.cmd
}

// Run command with cmd but no args
func (b builder) RunWithCmdAndNoArgs(function func(context.Context, io.Writer, *cobra.Command) error) *cobra.Command {
	b.cmd.RunE = func(_ *cobra.Command, args []string) error {
		return returnErrorFromFunction(function(b.cmd.Context(), b.cmd.OutOrStderr(), &b.cmd))
	}
	return &b.cmd
}

// SetFlags attaches flags to commands
func (b builder) SetFlags() Builder {
	SetCommandFlags(&b.cmd)
	return b
}

// Set Child of command
func (b builder) AddCommands(children ...*cobra.Command) Builder {
	for _, child := range children {
		b.cmd.AddCommand(child)
	}
	return b
}

// Handle Error from real function
func returnErrorFromFunction(err error) error {
	return err
}

func (b builder) WithFlags(adder func(*pflag.FlagSet)) Builder {
	adder(b.cmd.Flags())
	return b
}
