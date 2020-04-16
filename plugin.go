package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  "time"

  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sts"
  "github.com/sirupsen/logrus"
)

type (
  // Config holds input parameters for the plugin
  Config struct {
    Sensitive        bool
    RoleARN          string
    Region           string
  }

  Kube struct {
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

  if p.Config.RoleARN != "" {
    assumeRole(p.Config.RoleARN)
  }

  // Initialize commands
  var commands []*exec.Cmd

  // Print Kubectl version
  commands = append(commands, exec.Command(kubeExe, "version", "--client=true"))
  commands = append(commands, exec.Command(awsCliExe, "--version"))

  // Get kubeconfig config
  commands = append(commands, awsGetKubeConfig(p.Kube.ClusterName, p.Config.Region))

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

func assumeRole(roleArn string) {
  client := sts.New(session.New())
  duration := time.Hour * 1
  stsProvider := &stscreds.AssumeRoleProvider{
    Client:          client,
    Duration:        duration,
    RoleARN:         roleArn,
    RoleSessionName: "drone",
  }

  value, err := credentials.NewCredentials(stsProvider).Get()
  if err != nil {
    logrus.WithFields(logrus.Fields{
      "error": err,
    }).Fatal("Error assuming role!")
  }
  os.Setenv("AWS_ACCESS_KEY_ID", value.AccessKeyID)
  os.Setenv("AWS_SECRET_ACCESS_KEY", value.SecretAccessKey)
  os.Setenv("AWS_SESSION_TOKEN", value.SessionToken)
}



func trace(cmd *exec.Cmd) {
  fmt.Println("$", strings.Join(cmd.Args, " "))
}