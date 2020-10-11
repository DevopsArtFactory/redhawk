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

	executor, err := createNewExecutor()
	if err != nil {
		return err
	}

	executor.Runner.Builder = b

	//Run function with executor
	err = action(executor)

	return alwaysSucceedWhenCancelled(ctx, err)
}

// createNewExecutor creates new executor
func createNewExecutor() (Executor, error) {
	executor := Executor{
		Context: context.Background(),
		Runner:  runner.New(),
	}

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
