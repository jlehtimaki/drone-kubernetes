FROM ubuntu:18.04 AS builder

RUN apt update && apt install -y curl unzip bash

WORKDIR /build

# Fetch Kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.17.3/bin/linux/amd64/kubectl
RUN chmod +x kubectl

# Fetch awscli installation
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip

# Fetch kustomize
RUN curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash

FROM ubuntu:18.04

# Install AWSCli needed binaries & certs
RUN apt update && apt install -y less ca-certificates

# Copy binaries
COPY --from=builder /build/kubectl /bin/
COPY --from=builder /build/aws .
COPY --from=builder /build/kustomize /bin/

# Install AWSCli
RUN ./install -b /bin/

# Copy drone-kubernetes binary
COPY drone-kubernetes /bin/

ENTRYPOINT ["/bin/drone-kubernetes"]