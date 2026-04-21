package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

func getStatusIcon(status string) string {
	switch status {
	case "Running":
		return "🟢 Running"
	case "Pending":
		return "🟡 Pending"
	case "CrashLoopBackOff":
		return "🔴 CrashLoopBackOff"
	case "OOMKilled":
		return "🔴 OOMKilled"
	default:
		return status
	}
}

func ListPods(namespace string) {
	clientset, err := getKubeClient()
	if err != nil {
		fmt.Println("Error connecting to cluster:", err)
		return
	}

	ns := namespace

	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error fetching pods:", err)
		return
	}

	// Counters
	healthy := 0
	warning := 0
	critical := 0

	var issues []string

	// First pass: detect status + count
	type PodInfo struct {
		Name      string
		Namespace string
		Status    string
	}

	var podList []PodInfo

	for _, pod := range pods.Items {
		displayStatus := string(pod.Status.Phase)

		// Detect deeper issues FIRST
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Waiting != nil {
				if cs.State.Waiting.Reason == "CrashLoopBackOff" {
					displayStatus = "CrashLoopBackOff"
					issues = append(issues,
						fmt.Sprintf("🔴 %s → CrashLoopBackOff (restarting)", pod.Name))
				}
			}

			if cs.State.Terminated != nil {
				if cs.State.Terminated.Reason == "OOMKilled" {
					displayStatus = "OOMKilled"
					issues = append(issues,
						fmt.Sprintf("🔴 %s → OOMKilled (out of memory)", pod.Name))
				}
			}
		}

		// Pending check
		if pod.Status.Phase == v1.PodPending {
			displayStatus = "Pending"
			issues = append(issues,
				fmt.Sprintf("🟡 %s → Pending (scheduling issue)", pod.Name))
		}

		// Count AFTER final status decided
		switch displayStatus {
		case "Running":
			healthy++
		case "Pending":
			warning++
		case "CrashLoopBackOff", "OOMKilled", "Failed":
			critical++
		}

		podList = append(podList, PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    displayStatus,
		})
	}

	// ✅ NOW print summary (after counting)
	fmt.Println("\nCLUSTER OVERVIEW")
	fmt.Println("────────────────────────────────────────────")
	fmt.Printf("🟢 Healthy Pods: %d\n", healthy)
	fmt.Printf("🟡 Warnings:     %d\n", warning)
	fmt.Printf("🔴 Critical:     %d\n", critical)

	// ✅ Print table
	fmt.Println("\nPODS")
	fmt.Println("────────────────────────────────────────────")
	fmt.Printf("%-30s %-20s %-20s\n", "NAME", "STATUS", "NAMESPACE")

	for _, p := range podList {
		statusWithIcon := getStatusIcon(p.Status)

		fmt.Printf("%-30s %-20s %-20s\n",
			p.Name,
			statusWithIcon,
			p.Namespace,
		)
	}

	// ✅ Issues section
	if len(issues) > 0 {
		fmt.Println("\nISSUES")
		fmt.Println("────────────────────────────────────────────")
		for _, issue := range issues {
			fmt.Println(issue)
		}
	} else {
		fmt.Println("\n✔ No major issues detected")
	}
}
