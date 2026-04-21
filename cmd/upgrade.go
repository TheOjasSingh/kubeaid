package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade kubeaid to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		err := upgradeBinary()
		if err != nil {
			fmt.Println("❌ Upgrade failed:", err)
			return
		}
		fmt.Println("✅ kubeaid upgraded successfully!")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func upgradeBinary() error {
	repo := "TheOjasSingh/kubeaid"

	osType := runtime.GOOS
	arch := runtime.GOARCH

	if osType != "linux" {
		return fmt.Errorf("only linux supported")
	}

	var file string
	switch arch {
	case "amd64":
		file = "kubeaid-linux-amd64"
	case "arm64":
		file = "kubeaid-linux-arm64"
	default:
		return fmt.Errorf("unsupported architecture: %s", arch)
	}

	url := fmt.Sprintf(
		"https://github.com/%s/releases/latest/download/%s",
		repo,
		file,
	)

	fmt.Println("⬇️ Downloading latest version...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpFile := "/tmp/kubeaid-new"

	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Chmod(tmpFile, 0755)
	if err != nil {
		return err
	}

	fmt.Println("⚠️ Run this to complete upgrade:")
	fmt.Println("sudo mv /tmp/kubeaid-new /usr/local/bin/kubeaid")
	if err != nil {
		return fmt.Errorf("need sudo permissions: %v", err)
	}

	return nil
}
