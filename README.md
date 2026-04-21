# 🚑 kubeaid

> Diagnose Kubernetes issues in seconds — not minutes.

kubeaid is a lightweight CLI tool for SREs and DevOps engineers that scans your Kubernetes cluster and highlights issues like CrashLoopBackOff, Pending pods, and OOMKilled containers — with clear, actionable output.

---

## ⚡ Why kubeaid?

Debugging Kubernetes usually means:

* jumping between `kubectl describe`, logs, and dashboards
* wasting time finding the root cause

kubeaid simplifies that into one command:

```bash
kubeaid k8s pods
```

---

## ✨ Features

* 🔍 List all pods across namespaces
* ⚠ Detect common issues:

  * CrashLoopBackOff
  * Pending pods
  * OOMKilled containers
* 🎨 Clean, color-coded output
* 🚀 Zero setup (uses your existing kubeconfig)

---

## 📦 Installation

### One-line install (Linux)

```bash
curl -sL https://raw.githubusercontent.com/TheOjasSingh/kubeaid/main/install.sh | bash
```

---

### Manual install

Download the binary from Releases and move it to your PATH:

```bash
chmod +x kubeaid
sudo mv kubeaid /usr/local/bin/
```

---

## 🚀 Usage

### List all pods

```bash
kubeaid k8s pods
```

### Filter by namespace

```bash
kubeaid k8s pods -n default
```

---

## 🧪 Example Output

```
NAME                STATUS              NAMESPACE
-------------------------------------------------------
nginx               Running             default
redis               Pending             prod
payment-service     CrashLoopBackOff    prod

⚠ Issues detected:
- redis → Pending (possible scheduling issue)
- payment-service → CrashLoopBackOff (container restarting)
```

---

## ⚙️ Requirements

* Kubernetes cluster access
* Valid kubeconfig (`~/.kube/config`)
* `kubectl` working locally

---

## 🛠 Tech Stack

* Go
* Cobra CLI
* Kubernetes client-go

---

## 🧭 Roadmap

* [ ] Log analyzer (`kubeaid logs`)
* [ ] Root cause suggestions engine
* [ ] Pre-deploy checks
* [ ] kubectl plugin support
* [ ] AI-assisted diagnosis

---

## 🤝 Contributing

Contributions are welcome!
Feel free to open issues or submit PRs.

---

## ⭐ Support

If you find kubeaid useful, give it a star ⭐ on GitHub — it helps a lot!

---

## 📜 License

MIT License
