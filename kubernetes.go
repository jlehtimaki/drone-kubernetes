package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const kubeExe = "kubectl"
const kustomizeExe = "kustomize"

var (
	path = "/bin/kubectl"
)

func kubeApply(kube Kube) *exec.Cmd {
	if kube.Kustomize == "true" {
		if kube.Namespace != "" {
			cmd := fmt.Sprintf("kustomize build %s | kubectl -n %s apply -f -", kube.ManifestDir, kube.Namespace)
			return exec.Command("bash", "-c", cmd)
		}
		cmd := fmt.Sprintf("kustomize build %s | kubectl apply -f -", kube.ManifestDir)
		return exec.Command("bash", "-c", cmd)
	}

	var args []string
	if kube.Namespace != "" {
		args = append(args, "--namespace", kube.Namespace)
	}
	args = append(args, "apply", "-f", kube.ManifestDir)
	return exec.Command(kubeExe, args...)
}

func kubeDelete(kube Kube) *exec.Cmd {
	if kube.Kustomize == "true" {
		if kube.Namespace != "" {
			cmd := fmt.Sprintf("kustomize build %s | kubectl -n %s delete -f -", kube.ManifestDir, kube.Namespace)
			return exec.Command("bash", "-c", cmd)
		}
		cmd := fmt.Sprintf("kustomize build %s | kubectl delete -f -", kube.ManifestDir)
		return exec.Command("bash", "-c", cmd)
	}

	var args []string
	if kube.Namespace != "" {
		args = append(args, "--namespace", kube.Namespace)
	}
	args = append(args, "delete", "-f", kube.ManifestDir)
	return exec.Command(kubeExe, args...)
}

func kustomizeSetVersion(kube Kube) *exec.Cmd {
	imageName := fmt.Sprintf("%s:%s", kube.ImageName, kube.AppVersion)
	return exec.Command(kustomizeExe, "edit", "set", "image", imageName)
}

func kubeTest() *exec.Cmd {
	return exec.Command(kubeExe, "get", "pods", "--all-namespaces")
}

func installKubectl(version string) error {
	downloadUrl := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/kubectl", version)
	err := downloadFile(path, downloadUrl)
	if err != nil {
		return err
	}
	err = addExecRights()
	if err != nil {
		return err
	}
	return nil
}

func addExecRights() error {
	cmd := exec.Command("chmod", "+x", "/bin/kubectl")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
