package main

import (
  "os"

  "github.com/sirupsen/logrus"
  "github.com/urfave/cli"
)

var revision string // build number set at compile-time

func main() {
  app := cli.NewApp()
  app.Name = "kubernetes plugin"
  app.Usage = "kubernetes plugin"
  app.Action = run
  app.Version = revision
  app.Flags = []cli.Flag{

    //
    // plugin args
    //

    cli.StringSliceFlag{
      Name:   "actions",
      Usage:  "a list of actions to have kubectl perform",
      EnvVar: "PLUGIN_ACTIONS",
      Value:  &cli.StringSlice{"test"},
    },
    cli.StringFlag{
      Name:   "assume_role",
      Usage:  "A role to assume before running the awscli commands",
      EnvVar: "PLUGIN_ASSUME_ROLE",
    },
    cli.StringFlag{
      Name:   "kubectl_version",
      Usage:  "kubectl version number",
      EnvVar: "PLUGIN_KUBECTL_VERSION",
    },
    cli.StringFlag{
      Name:   "cluster_name",
      Usage:  "EKS Cluster Name",
      EnvVar: "PLUGIN_CLUSTER_NAME",
      Value:  "EKS-Cluster",
    },
    cli.StringFlag{
      Name:   "manifest_dir",
      Usage:  "Directory that holds manifests",
      EnvVar: "PLUGIN_MANIFEST_DIR",
      Value:  "./",
    },
    cli.StringFlag{
      Name:   "kubernetes_namespace",
      Usage:  "Namespace for Kubernetes",
      EnvVar: "PLUGIN_NAMESPACE",
      Value:  "default",
    },
    cli.StringFlag{
      Name:   "aws_region",
      Usage:  "AWS Region to use",
      EnvVar: "AWS_REGION",
      Value:  "eu-west-1",
    },
  }

  if err := app.Run(os.Args); err != nil {
    logrus.Fatal(err)
  }
}

func run(c *cli.Context) error {
  plugin := Plugin{
    Config: Config{
      RoleARN:          c.String("assume_role"),
      Region:           c.String("aws_region"),
    },
    Kube: Kube{
      Version:          c.String("kubectl_version"),
      Commands:         c.StringSlice("actions"),
      ClusterName:      c.String("cluster_name"),
      ManifestDir:      c.String("manifest_dir"),
      Namespace:        c.String("kubernetes_namespace"),
    },
  }

  return plugin.Exec()
}