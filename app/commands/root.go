package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golang-skeleton",
	Short: "Golang Skeleton CLI",
	Long:  `A CLI tool for Golang Skeleton project management`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add all commands here
	rootCmd.AddCommand(GenerateJWTSecretCmd)
}
