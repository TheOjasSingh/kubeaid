package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeClient() (*kubernetes.Clientset, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func ListPods(namespace string) {
	clientset, err := getKubeClient()
	if err != nil {
		fmt.Println("Error connecting to cluster:", err)
		return
	}

	ns := namespace
	if ns == "" {
		ns = ""
	}

	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error fetching pods:", err)
		return
	}

	fmt.Println("\nNAME\t\tSTATUS\t\t\tNAMESPACE")
	fmt.Println("-------------------------------------------------------")

	var issues []string

	for _, pod := range pods.Items {
		status := string(pod.Status.Phase)
		displayStatus := status

		// Detect deeper issues
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Waiting != nil {
				reason := cs.State.Waiting.Reason

				if reason == "CrashLoopBackOff" {
					displayStatus = "CrashLoopBackOff"
					issues = append(issues,
						fmt.Sprintf("%s → CrashLoopBackOff (container restarting)", pod.Name))
				}
			}

			if cs.State.Terminated != nil {
				if cs.State.Terminated.Reason == "OOMKilled" {
					displayStatus = "OOMKilled"
					issues = append(issues,
						fmt.Sprintf("%s → OOMKilled (out of memory)", pod.Name))
				}
			}
		}

		// Pending pods
		if pod.Status.Phase == v1.PodPending {
			displayStatus = "Pending"
			issues = append(issues,
				fmt.Sprintf("%s → Pending (possible scheduling issue)", pod.Name))
		}

		// Color formatting
		switch displayStatus {
		case "Running":
			color.Green("%s\t%s\t%s", pod.Name, displayStatus, pod.Namespace)
		case "Pending":
			color.Yellow("%s\t%s\t%s", pod.Name, displayStatus, pod.Namespace)
		case "CrashLoopBackOff", "OOMKilled", "Failed":
			color.Red("%s\t%s\t%s", pod.Name, displayStatus, pod.Namespace)
		default:
			fmt.Printf("%s\t%s\t%s\n", pod.Name, displayStatus, pod.Namespace)
		}
	}

	// Print Issues Summary
	if len(issues) > 0 {
		fmt.Println("\n⚠ Issues detected:")
		for _, issue := range issues {
			color.Red("- %s", issue)
		}
	} else {
		color.Green("\n✔ No major issues detected")
	}
}
