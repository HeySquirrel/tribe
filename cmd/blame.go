package cmd

import (
	"github.com/heysquirrel/tribe/app"
	"github.com/spf13/cobra"
)

var lineNumbers string

var blameCmd = &cobra.Command{
	Use:   "blame",
	Short: "Show a detailed history of each line of a file",
	Long:  `Better long description here`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New(args[0])
		defer app.Close()

		app.Loop()
	},
}

func init() {
	RootCmd.AddCommand(blameCmd)
	blameCmd.Flags().StringVarP(&lineNumbers, "lines", "L", "", "Line numbers to blame")
}
