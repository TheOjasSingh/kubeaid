package cmd

import (
	"github.com/TheOjasSingh/kubeaid/internal/k8s"
	"github.com/spf13/cobra"
)

var namespace string

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes related commands",
}

var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "List and analyze pods",
	Run: func(cmd *cobra.Command, args []string) {
		k8s.ListPods(namespace)
	},
}

func init() {
	rootCmd.AddCommand(k8sCmd)

	podsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace filter")
	k8sCmd.AddCommand(podsCmd)
}
