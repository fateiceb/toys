package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "toycli", // appname
}

func init() {
	rootCmd.AddCommand(cmdName)
	rootCmd.AddCommand(cmdType)
}
func Execute(){
	rootCmd.Execute()
}