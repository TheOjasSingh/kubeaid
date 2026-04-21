package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sredoctor",
	Short: "SRE Doctor - Diagnose Kubernetes issues quickly",
	Long:  "CLI tool for SREs and DevOps engineers to diagnose cluster issues",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
