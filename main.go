package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rcrowley/electrostatic/find"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFlags(0)
}

func main() {

	rootCmd := &cobra.Command{
		Use:                   "electrostatic",
		Short:                 "",
		Long:                  ``,
		Args:                  cobra.NoArgs,
		CompletionOptions:     cobra.CompletionOptions{DisableDefaultCmd: true},
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:    "shell-completion",
		Hidden: true,
		Short:  "print shell completion program for the current shell",
		Args:   cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			shell := os.Getenv("SHELL")
			if shell != "" {
				shell = filepath.Base(shell)
			}
			switch shell {
			case "", "bash":
				cmd.Root().GenBashCompletionV2(os.Stdout, true /* includeDesc */)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true /* includeDesc */)
			default:
				log.Fatalf("unsupported SHELL=%q", shell)
			}
		},
		DisableFlagsInUseLine: true,
	})

	rootCmd.AddCommand(find.Command())

	rootCmd.Execute()
}
