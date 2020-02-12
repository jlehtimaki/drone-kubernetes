package main

import (
  "os/exec"
)

const awsCliExe = "aws"

func awsGetKubeConfig(clusterName string, region string) *exec.Cmd {
  return exec.Command(awsCliExe, "eks","--region",region,"update-kubeconfig","--name",clusterName)
}
