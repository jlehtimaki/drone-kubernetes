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
		Sensitive     bool
		RoleARN       string
		Region        string
		ServerAddress string
		K8SCert       string
		K8SToken      string
		K8SUser       string
	}

	Kube struct {
		Type           string
		Version        string
		Commands       []string
		ManifestDir    string
		ClusterName    string
		Namespace      string
		Kustomize      string
		AppVersion     string
		ImageName      string
		Rollout        string
		RolloutTimeout string
	}

	// Plugin represents the plugin instance to be executed
	Plugin struct {
		Config Config
		Kube   Kube
	}
)

var (
	allowedCommands = []string{"apply", "delete", "diff"}
)

func allowedCommand(command string) bool {
	for _, com := range allowedCommands {
		if com == command {
			return true
		}
	}
	return false
}

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

	if p.Kube.Type == "EKS" {
		fmt.Println("Using EKS type of Kubernetes settings")

		// Printing AWSCli client version
		commands = append(commands, exec.Command(awsCliExe, "--version"))

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
		if allowedCommand(action) {
			commands = append(commands, kubeCommand(p.Kube, action))
		} else {
			return fmt.Errorf("valid actions are: apply, destroy.  You provided %s", action)
		}
	}

	if p.Kube.Rollout == "true" {
		commands = append(commands, checkRolloutStatus(p.Kube)...)
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

		if p.Kube.Kustomize == "true" {
			// Pipeline the kustomize build command with kubectl command
			c1 := exec.Command(kustomizeExe, "build", p.Kube.ManifestDir)
			c2 := c

			// initialize error
			var err error

			// pipe the commands
			c2.Stdin, err = c1.StdoutPipe()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to pipeline commands")
			}
			c2.Stdout = os.Stdout

			// run the commands
			err = c2.Start()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to execute kubectl command")
			}
			err = c1.Run()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to execute kustomize command")
			}
			// wait for the first command to finish
			err = c2.Wait()
			if err != nil && !strings.Contains(c.String(), "diff") {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to wait kustomize command")
			}
		} else {
			err := c.Run()
			// If kubectl command is diff ignore exit code since diff returns exit 1 if the is changes
			if err != nil && !strings.Contains(c.String(), "diff") {
				logrus.Info(c.String())
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to execute a command")
			}
		}

		logrus.Debug("Command completed successfully")
	}

	return nil
}

func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
