package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)
var cmdName = &cobra.Command{
	Use: "name ",
	Short:"print app name",
	Long:"print app name on screen",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("toycli")
	},
}