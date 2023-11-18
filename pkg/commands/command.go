package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/caarlos0/env/v10"
)

// Config описывает конфигурацию приложения.
type Config interface {
	// SetFlags устанавливает флаги командной строки.
	//
	// NOTE: SetFlags предполагает передачу указателей в fs, для чего
	// необходима передача ресивера конфигурации по указателю.
	SetFlags(fs *flag.FlagSet)

	// Validate возвращает ошибку, если одно из полей конфигурации не валидно.
	Validate() error
}

// ExecFunc определяет функцию выполнения командой.
type ExecFunc[T Config] func(ctx context.Context, c T) error

// Command определяет команду выполнения.
type Command[T Config] struct {
	fs     *flag.FlagSet
	exec   ExecFunc[T]
	config T
}

// New возвращает новый экземпляр Command.
func New[T Config](name string, exec ExecFunc[T]) *Command[T] {
	cmd := &Command[T]{
		fs:     flag.NewFlagSet(name, flag.ContinueOnError),
		exec:   exec,
		config: newConfig[T](),
	}
	cmd.init()
	return cmd
}

// init инициализирует первичную конфигурацию.
func (cmd *Command[T]) init() {
	// Отключение автоматического срабатывания Usage
	// в случае возникновения ошибки при парсинге флагов.
	cmd.fs.Usage = func() {}

	out := cmd.fs.Output()
	cmd.config.SetFlags(cmd.fs)
	cmd.fs.SetOutput(out)
}

// usage печатает формат использования команды.
func (cmd *Command[T]) usage() {
	fmt.Fprintf(cmd.fs.Output(), "Usage of %s:\n", cmd.fs.Name())
	cmd.fs.PrintDefaults()
}

// parseFlags парсит флаги командной строки.
func (cmd *Command[T]) parseFlags() error {
	args := cleanArgs(os.Args[1:])
	return cmd.fs.Parse(args)
}

func cleanArgs(args []string) []string {
	clean := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-test.") || strings.HasPrefix(args[i], "--test.") {
			continue
		}
		clean = append(clean, args[i])
	}
	return clean
}

// parseEnv парсит переменные окружения.
func (cmd *Command[T]) parseEnv() error {
	if !isPtr(cmd.config) {
		return env.Parse(&cmd.config)
	}
	return env.Parse(cmd.config)
}

// Execute запускает команду и блокируется до её завершения.
func (cmd *Command[T]) Execute(ctx context.Context) error {
	err := cmd.parseFlags()
	if err != nil {
		cmd.usage()
		return nil
	}

	err = cmd.parseEnv()
	if err != nil {
		return fmt.Errorf("parsing env: %w", err)
	}

	err = cmd.config.Validate()
	if err != nil {
		return fmt.Errorf("config validation: %w", err)
	}

	err = cmd.exec(ctx, cmd.config)
	if err != nil {
		return fmt.Errorf("failed to execute func: %w", err)
	}

	return nil
}

func newConfig[T Config]() T {
	var zero T
	rt := reflect.TypeOf(zero)
	if rt.Kind() == reflect.Pointer {
		zero = reflect.New(rt.Elem()).Interface().(T)
	}
	return zero
}

func isPtr[T Config](c T) bool {
	rt := reflect.TypeOf(c)
	return rt.Kind() == reflect.Pointer
}