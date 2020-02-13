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

var (
  path        = "/bin/kubectl"
)

func kubeApply(kube Kube) *exec.Cmd {
  var args []string
  if kube.Namespace != "" {
    args = append(args, "-n", kube.Namespace)
  }
  args = append(args, "apply", "-f", kube.ManifestDir)
  return exec.Command(kubeExe,args...)
}

func kubeDelete(kube Kube) *exec.Cmd {
  var args []string
  if kube.Namespace != "" {
    args = append(args, "-n", kube.Namespace)
  }
  args = append(args, "delete", "-f", kube.ManifestDir)
  return exec.Command(kubeExe, args...)
}

func kubeTest() *exec.Cmd {
  return exec.Command(kubeExe, "get","pods","--all-namespaces")
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

func addExecRights() error{
  cmd := exec.Command("chmod","+x","/bin/kubectl")
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
