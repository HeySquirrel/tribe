package cmd

import (
	"fmt"
	"github.com/heysquirrel/tribe/apis"
	"github.com/heysquirrel/tribe/apis/rally"
	"github.com/heysquirrel/tribe/blame"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/heysquirrel/tribe/config"
	"github.com/spf13/cobra"
	"os"
)

var endpoints []int

var blameCmd = &cobra.Command{
	Use:   "blame",
	Short: "Show a detailed history of each line of a file",
	Long:  `Better long description here`,
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

		apikey := config.RallyApiKey()
		api := apis.NewCachingServer(rally.New(apikey))

		annotate := model.NewCachingAnnotate(model.NewAnnotate(api))

		blame := blame.NewBlameApp(file, annotate)
		defer blame.Close()

		blame.Loop()
	},
}

func init() {
	RootCmd.AddCommand(blameCmd)
	blameCmd.Flags().IntSliceVarP(&endpoints, "lines", "L", []int{1, 20}, "Line numbers to blame")
}
