package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "oscarsthegrouch",
	Short: "oscarsthegrouch is a CLI tool for managing your trash",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("oscarsthegrouch is a CLI tool for managing your trash")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
