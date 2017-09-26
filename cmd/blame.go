package cmd

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/blame"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/spf13/cobra"
	"os"
)

var endpoints []int

var blameCmd = &cobra.Command{
	Use:   "blame",
	Short: "Why the @*$% does this code exist?",
	Long:  `Access historical work items or issues, frequent contributors and your entire git history with one simple command so that you quickly determine why a line of code exists.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		var start, end int

		switch len(endpoints) {
		case 1:
			start = endpoints[0]
			end = start + 20
		case 2:
			start = endpoints[0]
			end = endpoints[1]
		default:
			cmd.Help()
			os.Exit(1)
		}

		file, err := model.NewFile(filename, start, end)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		server, err := apis.NewRallyFromConfig("rally1")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		api := apis.NewCachingServer(server)

		annotate := model.NewCachingAnnotate(model.NewAnnotate(api))

		blame := blame.NewBlameApp(file, annotate)
		defer blame.Close()

		blame.Loop()
	},
}

func init() {
	RootCmd.AddCommand(blameCmd)
	blameCmd.Flags().IntSliceVarP(&endpoints, "lines", "L", []int{1, 20}, "line numbers to blame")
}
