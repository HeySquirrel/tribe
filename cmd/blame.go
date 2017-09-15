package cmd

import (
	"fmt"
	"github.com/heysquirrel/tribe/blame"
	"github.com/heysquirrel/tribe/blame/model"
	"github.com/spf13/cobra"
	"os"
)

var endpoints []int

var blameCmd = &cobra.Command{
	Use:   "blame",
	Short: "Show a detailed history of each line of a file",
	Long:  `Better long description here`,
	Args:  cobra.MinimumNArgs(1),
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

		data, err := model.NewBlame(filename, start, end)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		blame := blame.NewBlameApp(data)
		defer blame.Close()

		blame.Loop()
	},
}

func init() {
	RootCmd.AddCommand(blameCmd)
	blameCmd.Flags().IntSliceVarP(&endpoints, "lines", "L", []int{1, 20}, "Line numbers to blame")
}
