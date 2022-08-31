package command

import (
	"context"

	"github.com/spf13/cobra"

	toolservice "crawl-worker/internal/service/tool"
)

func NewGetChatByNameCommand(service toolservice.IService) *cobra.Command {
	command := &cobra.Command{
		Use:     "getchatbyname",
		Aliases: []string{"gcbn"},
		Run: func(cmd *cobra.Command, args []string) {
			chatName := cmd.Flag("name").Value.String()
			service.GetChatFromName(context.Background(), chatName)
		},
	}

	command.Flags().StringP("name", "n", "", "chat name")

	return command
}
