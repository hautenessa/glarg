package glarg

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
)

type Subcommand interface {
	Description() string
	FlagSet() *flag.FlagSet
	SetupSubcommand() Subcommand
	HasInvalidFlags() bool
	Execute(ctx context.Context) int
}

type ArgumentConsumer interface {
	SetArgs(args []string)
}

type ArgumentUnpacker interface {
	UnpackArgs() error
}

type Subcommands struct {
	flagSet  *flag.FlagSet
	args     []string
	Name     string
	Children []Subcommand
}

func (self *Subcommands) Description() string {
	names := make([]string, len(self.Children))
	for i, v := range self.Children {
		names[i] = v.FlagSet().Name()
	}
	return fmt.Sprintf("Subcommands: %s", strings.Join(names, ", "))
}

func (self *Subcommands) Usage() {
	for _, v := range self.Children {
		log.Printf("  %-10s %-60s", v.FlagSet().Name(), v.Description())
	}
}

func (self *Subcommands) FlagSet() *flag.FlagSet {
	return self.flagSet
}

func (self *Subcommands) SetupSubcommand() Subcommand {
	self.flagSet = flag.NewFlagSet(self.Name, flag.ExitOnError)
	for i, v := range self.Children {
		self.Children[i] = v.SetupSubcommand()
	}
	return self
}

func (sef *Subcommands) HasInvalidFlags() bool {
	return false
}

func (self *Subcommands) Execute(ctx context.Context) int {
	if len(self.args) < 2 {
		log.Printf("Missing subcommand.")
		self.Usage()
		return 1
	}

	// strip "ourself" off the args and then find the
	// requested subcommand
	myArgs := self.args[1:]
	var subcmd Subcommand
	for _, v := range self.Children {
		if v.FlagSet().Name() == myArgs[0] {
			v.FlagSet().Parse(myArgs[1:])
			subcmd = v
			break
		}
	}

	// No subcommand, print the usage.
	if subcmd == nil {
		log.Printf("Unknown subcommand provided: %s.", myArgs[0])
		self.Usage()
		return 1
	}

	// If the subcommand needs to convert the flagset into
	// other data, it does it here.
	if au, ok := subcmd.(ArgumentUnpacker); ok {
		if err := au.UnpackArgs(); err != nil {
			log.Printf("Invalid arguments. %s", err)
			subcmd.FlagSet().PrintDefaults()
			return 1
		}
	}

	// The subcommand can communicate if it has invalid data here.
	// The subcommand is expected to output its own errors, but
	// not the defaults.
	if subcmd.HasInvalidFlags() {
		subcmd.FlagSet().PrintDefaults()
		return 1
	}

	if ac, ok := subcmd.(ArgumentConsumer); ok {
		ac.SetArgs(myArgs)
	}

	return subcmd.Execute(ctx)
}

func (self *Subcommands) SetArgs(args []string) {
	self.args = args
}

func handleInterupt(ctx context.Context, cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
			break
		}
	}
}

func Invoke(ctx context.Context, cmd Subcommand, args []string) int {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	go handleInterupt(ctx, cancel)

	setupCmd := cmd.SetupSubcommand()

	if nested, ok := setupCmd.(ArgumentConsumer); ok {
		nested.SetArgs(args)
	}

	return setupCmd.Execute(ctx)
}

type SubcommandNoOp struct {
	flagSet          *flag.FlagSet
	Name             string
	UnpackArgsError  error
	InvalidFlagsBool bool
	ExecuteInt       int
}

func (self *SubcommandNoOp) Description() string {
	return "Not actually implemented."
}
func (self *SubcommandNoOp) FlagSet() *flag.FlagSet {
	return self.flagSet
}
func (self *SubcommandNoOp) SetupSubcommand() Subcommand {
	self.flagSet = flag.NewFlagSet(self.Name, flag.ExitOnError)
	return self
}
func (self *SubcommandNoOp) UnpackArgs() error {
	return self.UnpackArgsError
}
func (self *SubcommandNoOp) HasInvalidFlags() bool {
	return self.InvalidFlagsBool
}
func (self *SubcommandNoOp) Execute(ctx context.Context) int {
	return self.ExecuteInt
}
