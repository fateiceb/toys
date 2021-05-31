package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var cmdType = &cobra.Command{
	Use: "type",
	Short: "print this app type",
	Long: "print this app type on screen",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("this is a cli app")
	},
}