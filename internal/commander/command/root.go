package command

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tdtool",
		Short: "tdtool - a tool for telegram",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}
