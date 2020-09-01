package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
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

func kubeCommand(kube Kube, command string) *exec.Cmd {
	var args []string
	if kube.Namespace != "" {
		args = append(args, "--namespace", kube.Namespace)
	}
	if kube.Kustomize == "true" {
		args = append(args, command, "-f", "-")
	} else {
		args = append(args, command, "-f", kube.ManifestDir)
	}
	return exec.Command(kubeExe, args...)
}

func kustomizeSetVersion(kube Kube) *exec.Cmd {
	imageName := fmt.Sprintf("%s:%s", kube.ImageName, kube.AppVersion)
	return exec.Command(kustomizeExe, "edit", "set", "image", imageName)
}

func installKubectl(version string) error {
	arch := os.Getenv("GOARCH")
	if arch == "" {
		arch = "amd64"
	}
	downloadUrl := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/%s/kubectl", version, arch)
	logrus.Info("Installing Kubectl version ", version)
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
	err := os.Chmod("/bin/kubectl", 0777)
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(filepath string, url string) error {
	//Get the response bytes from the url
	logrus.Info("Downloading file ", url)
	response, err := http.Get(url)
	if err != nil {
	}
	defer response.Body.Close()

	//Create a empty file
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	// Check that file exists
	if !checkFileExists(filepath) {
		return fmt.Errorf("kubectl file not found")
	}
	return nil
}

func checkFileExists(filepath string) bool {
	// Returns true if file exists
	_, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return true
}
