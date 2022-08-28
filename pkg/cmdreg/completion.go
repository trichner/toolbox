package cmdreg

import (
	"github.com/posener/complete/v2"
)

type completionCommand struct {
	complete.Command
	SubCompleter map[string]complete.Completer
}

func (c *completionCommand) SubCmdList() []string {
	subs := make([]string, 0, len(c.SubCompleter))
	for sub := range c.SubCompleter {
		subs = append(subs, sub)
	}
	return subs
}

func (c *completionCommand) SubCmdGet(cmd string) complete.Completer {
	if c.SubCompleter[cmd] == nil {
		return nil
	}
	return c.SubCompleter[cmd]
}

func (c *CommandRegistry) setupCompletions(program string) {

	if c.program != program {
		cmd, ok := c.commands[program]
		if !ok || cmd.completer == nil {
			//no completions, nothing to set up
			return
		}
		complete.Complete(program, cmd.completer)
		return
	}

	commands := map[string]complete.Completer{}
	for k, cmd := range c.commands {
		commands[k] = cmd.completer
	}
	cmd := &completionCommand{
		Command:      complete.Command{},
		SubCompleter: commands,
	}
	complete.Complete(program, cmd)
}
