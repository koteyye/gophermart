package gophermart

import (
	"context"
	"flag"
	"fmt"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/config"
	"github.com/sergeizaitcev/gophermart/internal/gophermart/server"
)

// Command определяет команду с набором флагов для gophermart.
type Command struct {
	fs     *flag.FlagSet
	actual *config.Config
}

// NewCommand возвращает новый экземпляр Command.
func NewCommand() *Command {
	fs := flag.NewFlagSet("gophermart", flag.ContinueOnError)

	// Отключение автоматического срабатывания при парсинге флагов.
	fs.Usage = func() {}

	return &Command{
		fs: fs,
		actual: &config.Config{
			Level: slog.LevelInfo,
		},
	}
}

// Usage печает формат использования команды.
func (cmd *Command) Usage() {
	fmt.Fprintf(cmd.fs.Output(), "Usage of %s:\n", cmd.fs.Name())
	cmd.fs.PrintDefaults()
}

// Parse парсит флаги командной строки.
func (cmd *Command) Parse(args []string) error {
	c := cmd.actual.Clone()

	cmd.fs.TextVar(&c.Level, "v", cmd.actual.Level, "logging level")
	cmd.fs.StringVar(&c.RunAddress, "a", cmd.actual.RunAddress, "run address")
	cmd.fs.StringVar(&c.DatabaseURI, "d", cmd.actual.DatabaseURI, "database uri")
	cmd.fs.StringVar(
		&c.AccrualSystemAddress,
		"r",
		cmd.actual.AccrualSystemAddress,
		"accrual system address",
	)

	err := cmd.fs.Parse(args)
	if err != nil {
		return fmt.Errorf("parsing flags: %w", err)
	}

	cmd.actual = c

	return nil
}

// Run запускает gophermart и блокируется до тех пор, пока не сработает
// контекст или функция не вернёт ошибку.
func (cmd *Command) Run(ctx context.Context) error {
	cfg, err := cmd.initConfig()
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	srv := server.New(cfg)

	return srv.Run(ctx)
}

func (cmd *Command) initConfig() (*config.Config, error) {
	cfg := cmd.actual.Clone()

	err := config.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	return cfg, nil
}
