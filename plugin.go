package main

import (
  "fmt"
  "github.com/sirupsen/logrus"
  "os"
  "os/exec"
  "strings"
)

type (
  // Config holds input parameters for the plugin
  Config struct {
    Sensitive       bool
    RoleARN         string
    Region          string
    ServerAddress   string
    K8SCert         string
    K8SToken        string
    K8SUser         string
  }

  Kube struct {
    Type            string
    Version         string
    Commands        []string
    ManifestDir     string
    ClusterName     string
    Namespace       string
    Kustomize       string
    AppVersion      string
    ImageName       string
  }

  // Plugin represents the plugin instance to be executed
  Plugin struct {
    Config      Config
    Kube        Kube
  }
)

// Exec executes the plugin
func (p Plugin) Exec() error {
  // Install specified version of kubectl
  if p.Kube.Version != "" {
    err := installKubectl(p.Kube.Version)
    if err != nil {
      return err
    }
  }

  // Initialize commands
  var commands []*exec.Cmd

  // Print Kubectl version
  commands = append(commands, exec.Command(kubeExe, "version", "--client=true"))
  commands = append(commands, exec.Command(awsCliExe, "--version"))

  if p.Kube.Type == "EKS" {
    fmt.Println("Using EKS type of Kubernetes settings")
    // Assume AWS Role
    if p.Config.RoleARN != "" {
      assumeRole(p.Config.RoleARN)
    }
    // Get kubeconfig config
    commands = append(commands, awsGetKubeConfig(p.Kube.ClusterName, p.Config.Region))
  }

  if p.Kube.Type == "Baremetal" {
    fmt.Println("Using Baremetal type of Kubernetes settings")
    commands = append(commands, bareMetalSetKubeConfig(p.Config.K8SToken, p.Config.K8SCert, p.Config.ServerAddress, p.Config.K8SUser)...)
  }

  // Set version with Kustomize
  if p.Kube.AppVersion != "" {
    commands = append(commands, kustomizeSetVersion(p.Kube))
  }
  // Add commands listed in actions
  for _, action := range p.Kube.Commands {
    switch action {
    case "apply":
      commands = append(commands, kubeApply(p.Kube))
    case "delete":
      commands = append(commands, kubeDelete(p.Kube))
    case "test":
      commands = append(commands, kubeTest())
    default:
      return fmt.Errorf("valid actions are: apply, destroy.  You provided %s", action)
    }
  }

  // Run commands
  for _, c := range commands {
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr
    if !p.Config.Sensitive {
      trace(c)
    }

    if strings.Contains(c.String(), "edit") {
      c.Dir = p.Kube.ManifestDir
    }

    err := c.Run()
    if err != nil {
      logrus.WithFields(logrus.Fields{
        "error": err,
      }).Fatal("Failed to execute a command")
    }
    logrus.Debug("Command completed successfully")
  }

  return nil
}

func trace(cmd *exec.Cmd) {
  fmt.Println("$", strings.Join(cmd.Args, " "))
}