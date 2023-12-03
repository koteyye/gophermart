package commands_test

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sergeizaitcev/gophermart/pkg/commands"
)

var _ commands.Config = (*testConfig)(nil)

type testConfig struct {
	String string `env:"STRING"`
	Int64  int64  `env:"INT64"`
}

func (m *testConfig) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&m.String, "string", "string", "")
	fs.Int64Var(&m.Int64, "int64", 1, "")
}

func (*testConfig) Validate() error {
	return nil
}

func ExampleCommand() {
	os.Setenv("INT64", "2")

	exec := func(_ context.Context, got *testConfig) error {
		fmt.Printf("String=%s\nInt64=%d\n", got.String, got.Int64)
		return nil
	}

	cmd := commands.New("test", exec)
	cmd.Execute(context.Background())

	// Output:
	//
	// String=string
	// Int64=2
}