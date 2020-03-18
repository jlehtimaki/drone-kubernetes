# drone-eks-kubernetes
Drone plugin for AWS EKS. This will help to do Kubernetes deployments to your EKS cluster

## Build
Build the binary with the following commands:

```export CGO_ENABLED=0
go build
```

## Docker

Build the docker image with:
```
docker build --rm=true -t lehtux/drone-eks-kubernetes .
```

## Usage
```
docker run --rm -it -e AWS_ACCESS_KEY=.... -e PLUGIN_ASSUME_ROLE=.... -e AWS_SECRET_KEY=.... 
-e PLUGIN_ACTIONS="apply" -e PLUGIN_MANIFEST_DIR="manifests/" lehtux/drone-awscli
```

## Parameters
| Paramenter            | Description                   |Required       | Default Value | Allowed Values |
| -------------         |:-------------:                |:-------------:|:-------------:|:-------------: |
| AWS_ACCESS_KEY        | AWS Access key                | YES           | -             | -              |
| AWS_SECRET_KEY        | AWS Access key secret         | YES           | -             | -              |
| AWS_REGION            | AWS Region                    | NO            | eu-west-1     | -              |
| PLUGIN_ASSUME_ROLE    | AWS Assume role               | NO            | -             | Role ARN       |
| PLUGIN_ACTIONS        | AWS Client command to be run  | YES           | test          | test/apply/delete|
| PLUGIN_KUBECTL_VERSION| Kubectl version to be installed| NO           | v1.7.3        | -              |
| PLUGIN_NAMESPACE      | Kubernetes namespace          | NO            | default       | -              |
| PLUGIN_CLUSTER_NAME   | EKS Cluster name              | NO            | EKS-Cluster   | -              |
| PLUGIN_MANIFEST_DIR   | Directory holding the manifests| NO           | ./            | -              |
| PLUGIN_KUSTOMIZE      | Use Kustomize                 | NO            | false         | true / false   |
| PLUGIN_VERSION        | Version to deploy             | NO            | -             | -              |
| PLUGIN_IMAGE          | Image name of the deployment. Used with Kustomize | NO | -    | -              |

## Drone pipeline example
```yaml
kind: pipeline
type: kubernetes
name: Drone example pipeline

steps:
  - name: Deploy test app
    image: lehtux/drone-eks-kubernetes
    environment:
      AWS_REGION: "eu-west-1"
      AWS_ACCESS_KEY_ID:
        from_secret: access_key
      AWS_SECRET_ACCESS_KEY:
        from_secret: access_key_secret
    settings:
      assume_role: arn:aws:iam::xxxxxx:role/EKS
      actions: ["apply"]
      namespace: default
      kubectl_version: v1.16.6
      manifest_dir: deployments/deployment.yml

```