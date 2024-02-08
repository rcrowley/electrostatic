package find

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var latest int
	cmd := &cobra.Command{
		Use:   "find [--html|--markdown] all|latest <n>",
		Short: "find HTML and/or Markdown files",
		Args: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) > 1 {
				return errors.New("`electrostatic find` requires at least one more argument")
			}
			switch args[0] {
			case "all":
				if len(args) > 1 {
					return errors.New("`electrostatic find all` doesn't accept any further arguments")
				}
			case "latest":
				if len(args) < 2 {
					return errors.New("`electrostatic find latest <n>` requires a numeric final argument")
				}
				if len(args) > 2 {
					return errors.New("`electrostatic find latest <n>` doesn't accept any further arguments")
				}
				if latest, err = strconv.Atoi(args[1]); err != nil {
					return errors.New("`electrostatic find latest <n>` requires a numeric final argument")
				}
			default:
				return fmt.Errorf("`electrostatic find %s` is not a valid command", args[0])
			}
			return nil
		},
		DisableFlagsInUseLine: true,
		ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) > 0 && toComplete == "" {
				toComplete = args[len(args)-1]
			}
			switch toComplete {
			case "all":
				return nil, cobra.ShellCompDirectiveNoFileComp // | cobra.ShellCompDirectiveNoSpace
			case "latest":
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return []string{
				"--html", "--markdown",
				"all", "latest",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	}
	findHTML := cmd.Flags().Bool("html", false, "find only HTML files")
	findMarkdown := cmd.Flags().Bool("markdown", false, "find only Markdown files")
	cmd.MarkFlagsMutuallyExclusive("html", "markdown")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		log.Printf("find.Command().Run findHTML: %#v", *findHTML)
		log.Printf("find.Command().Run findMarkdown: %#v", *findMarkdown)
		log.Printf("find.Command().Run latest: %#v", latest)
		log.Printf("find.Command().Run args: %#v", args)
	}
	return cmd
}
