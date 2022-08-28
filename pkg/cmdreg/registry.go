package cmdreg

import (
	"fmt"
	"github.com/posener/complete/v2"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/maps"
)

const defaultProgramName = "tb"

type registryConfig struct {
	program string
}
type Option func(c *registryConfig) error

func WithProgramName(program string) Option {
	return func(c *registryConfig) error {
		c.program = program
		return nil
	}
}

type commandConfig struct {
	completions complete.Completer
}
type CommandOption func(c *commandConfig) error

func WithCompletion(c complete.Completer) CommandOption {
	return func(cfg *commandConfig) error {
		cfg.completions = c
		return nil
	}
}

type Command interface {
	Exec(args []string)
}

type commandSet struct {
	command   Command
	completer complete.Completer
}

type CommandRegistry struct {
	program  string
	commands map[string]*commandSet
}

func New(options ...Option) *CommandRegistry {

	cfg := &registryConfig{program: defaultProgramName}
	for _, o := range options {
		if err := o(cfg); err != nil {
			log.Fatal(err)
		}
	}

	return &CommandRegistry{
		program:  cfg.program,
		commands: make(map[string]*commandSet),
	}
}

func (c *CommandRegistry) Register(cmd string, command Command, options ...CommandOption) {
	cfg := &commandConfig{}

	completer, ok := command.(complete.Completer)
	if ok && completer != nil {
		cfg.completions = completer
	}

	for _, o := range options {
		err := o(cfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	c.commands[cmd] = &commandSet{
		command:   command,
		completer: cfg.completions,
	}
}

func (c *CommandRegistry) RegisterFunc(cmd string, fn CommandFunc, options ...CommandOption) {
	c.Register(cmd, fn, options...)
}

func (c *CommandRegistry) List() []string {
	return maps.Keys(c.commands)
}

func (c *CommandRegistry) execCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments")
	}
	cmd := args[0]
	cmd = filepath.Base(cmd)

	runner, ok := c.commands[cmd]
	if !ok {
		return fmt.Errorf("unknown command '%s'", cmd)
	}
	runner.command.Exec(args)
	return nil
}

func (c *CommandRegistry) Exec(args []string) {
	if len(args) == 0 {
		log.Fatal("no program in arguments")
	}

	prog := filepath.Base(args[0])
	c.setupCompletions(prog)
	if prog == c.program {
		args = args[1:]
	}

	err := c.execCommand(args)
	if err != nil {
		log.Printf("%s", err)
		c.PrintHelp(os.Stderr)
		os.Exit(1)
	}
}

func (c *CommandRegistry) PrintHelp(w io.Writer) {
	commands := c.List()
	fmt.Fprintf(w, "available commands:\n")
	for _, c := range commands {
		fmt.Fprintf(w, "  %s\n", c)
	}
}

type CommandFunc func(args []string)

func (c CommandFunc) Exec(args []string) {
	c(args)
}
