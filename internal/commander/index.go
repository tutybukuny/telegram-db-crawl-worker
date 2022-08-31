package commander

import (
	"github.com/spf13/cobra"

	"crawl-worker/internal/commander/command"
	toolservice "crawl-worker/internal/service/tool"
	"crawl-worker/pkg/container"
	"crawl-worker/pkg/l"
)

type Commander struct {
	ll          l.Logger             `container:"name"`
	toolService toolservice.IService `container:"name"`

	rootCommand *cobra.Command
}

func New() *Commander {
	cmd := &Commander{rootCommand: command.NewRootCmd()}
	container.Fill(cmd)

	initCommands(cmd.toolService, cmd.rootCommand)

	return cmd
}

func initCommands(service toolservice.IService, rootCommand *cobra.Command) {
	getChatByName := command.NewGetChatByNameCommand(service)
	rootCommand.AddCommand(getChatByName)
}

func (c *Commander) Execute() {
	if err := c.rootCommand.Execute(); err != nil {
		c.ll.Error("error when execute command", l.Error(err))
	}
}
