package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var clean bool

var rootCmd = &cobra.Command{
	Use:   "go-telnet [ADDRESS] [PORT]",
	Short: "telnet replica utility written in Go",
	Long:  "Runs any programme passed to goenvdir with a specified environment variables stored in a specified folder.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("invalid number of arguments")
		}
		return nil
	},
	Run: rootCommand,
}

// Execute is a method that runs the root command of the programme
func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&clean, "clean", "c", false, "run programme with the empty environment variables")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCommand(cmd *cobra.Command, args []string) {
	pr := internal.NewProgRunnerImpl()
	pr.ClearEnv = clean
	pr.EnvPath = args[0]
	pr.ChildProg = args[1]
	if err := pr.Execute(); err != nil {
		os.Exit(1)
	}
}
