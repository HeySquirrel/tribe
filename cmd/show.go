package cmd

import (
	"fmt"
	"github.com/heysquirrel/tribe/work"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show relevant information about a given work item or issue",
	Long:  `See the name, description and owner for relevant work items or issues in CA Agile Central or Jira.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workitemid := args[0]

		api, err := work.NewItemServer()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		workitem, err := api.GetItem(workitemid)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		Display(os.Stdout, workitem)
	},
}

func Display(writer io.Writer, workitem work.Item) {
	fmt.Fprintln(writer)

	fmt.Fprintf(writer, "%s - %s\n\n", workitem.GetId(), workitem.GetName())
	fmt.Fprintln(writer, sanitize.HTML(workitem.GetDescription()))
}

func init() {
	RootCmd.AddCommand(ShowCmd)
}
