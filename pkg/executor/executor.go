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

package executor

import (
	"context"

	"github.com/DevopsArtFactory/redhawk/pkg/builder"
	"github.com/DevopsArtFactory/redhawk/pkg/runner"
)

type Executor struct {
	Context context.Context
	Runner  *runner.Runner
}

var NewExecutor = createNewExecutor

// RunExecutor for command line
func RunExecutor(ctx context.Context, action func(Executor) error) error {
	flags, err := builder.GetFlags()
	if err != nil {
		return err
	}

	if err := builder.ValidateFlags(flags); err != nil {
		return err
	}

	b, err := builder.CreateNewBuilder(flags)
	if err != nil {
		return err
	}

	executor, err := NewExecutor(ctx, b)
	if err != nil {
		return err
	}

	//Run function with executor
	err = action(executor)

	return alwaysSucceedWhenCancelled(ctx, err)
}

// createNewExecutor creates new executor
func createNewExecutor(ctx context.Context, b *builder.Builder) (Executor, error) {
	executor := Executor{
		Context: ctx,
		Runner:  runner.New(),
	}

	executor.Runner.Builder = b

	return executor, nil
}

// alwaysSucceedWhenCancelled makes response true if user canceled
func alwaysSucceedWhenCancelled(ctx context.Context, err error) error {
	// if the context was cancelled act as if all is well
	if err != nil && ctx.Err() == context.Canceled {
		return nil
	}
	return err
}
