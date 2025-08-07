FROM mcr.microsoft.com/devcontainers/go:latest

# Install build dependencies
RUN apt-get update \    
    && apt-get install -y --no-install-recommends \    
        git \    
        curl \    
        unzip \    
        make \    
        bash \    
    && apt-get clean \    
    && rm -rf /var/lib/apt/lists/*

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest \    
    && go install golang.org/x/lint/golint@latest

WORKDIR /workspace

# Set default shell
SHELL ["/bin/bash", "-c"]
