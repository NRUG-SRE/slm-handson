FROM mcr.microsoft.com/devcontainers/base:ubuntu

ARG GO_VERSION=1.21.13
ARG NODE_VERSION=18

# Install Go
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:/home/vscode/go/bin:${PATH}"
ENV GOPATH="/home/vscode/go"

# Install Node.js
RUN curl -fsSL "https://deb.nodesource.com/setup_${NODE_VERSION}.x" | bash - \
    && apt-get install -y nodejs \
    && apt-get clean && rm -rf /var/lib/apt/lists/*
