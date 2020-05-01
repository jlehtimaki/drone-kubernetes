package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func bareMetalSetKubeConfig(token string, cert string, server string, user string) []*exec.Cmd {
	fmt.Println("Setting up Baremetal Kubernetes configuration")
	if cert != "" {
		// Write certificate file
		writeCertToFile(cert)
	}

	// Assign all needed kubernetes commands and return them
	var commands []*exec.Cmd
	tokenString := fmt.Sprintf("--token=%s", token)
	serverString := fmt.Sprintf("--server=%s", server)
	userString := fmt.Sprintf("--user=%s", user)
	commands = append(commands, exec.Command(kubeExe, "config", "set-credentials", "default", tokenString))
	if cert != "" {
		commands = append(commands, exec.Command(kubeExe, "config", "set-cluster", "default", serverString, "--certificate-authority=ca.crt"))
	} else {
		commands = append(commands, exec.Command(kubeExe, "config", "set-cluster", "default", serverString, "--insecure-skip-tls-verify=true"))
	}
	commands = append(commands, exec.Command(kubeExe, "config", "set-context", "default", "--cluster=default", userString))
	commands = append(commands, exec.Command(kubeExe, "config", "use-context", "default"))
	return commands
}

func writeCertToFile(cert string) {
	err := ioutil.WriteFile("ca.crt", []byte(cert), 0644)
	if err != nil {
		fmt.Printf("Could not write certificate file: %s", err.Error())
	}
}
