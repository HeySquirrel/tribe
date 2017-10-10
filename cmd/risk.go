package cmd

import (
	"fmt"
	"os"

	"github.com/HeySquirrel/tribe/git"
	"github.com/HeySquirrel/tribe/risk"
	"github.com/spf13/cobra"
)

var riskCmd = &cobra.Command{
	Use:   "risk",
	Short: "Attempt to estimate the risk of changing the given file.",
	Long:  `Combine measurements around the number of work items, contributors and frequency of edits to come up with a file risk score.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("%s does not exist.\n", file)
			os.Exit(1)
		}

		commits, err := git.Log("--follow", file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		risk := risk.Calculate(file, commits)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		risk.Write(os.Stdout)

	},
}

func init() {
	RootCmd.AddCommand(riskCmd)
}
