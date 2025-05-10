package cmd

import (
	"fmt"
	"os"

	"github.com/astr0n8t/k8s-portmapper/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-portmapper",
	Short: "SHORT DESCRIPTION",
	Long: `LONG DESCRIPTION 
	WITH MULTIPLE LINES`,
	// Simply call the internal run command
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
