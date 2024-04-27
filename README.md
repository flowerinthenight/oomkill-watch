[![main](https://github.com/flowerinthenight/oomkill-watch/actions/workflows/main.yml/badge.svg)](https://github.com/flowerinthenight/oomkill-watch/actions/workflows/main.yml)

## Overview

`oomkill-watch` is a simple wrapper to the `kubectl get events -w` command, specifically filtering the `OOMKilling` events and optionally forwarding them to a Slack channel. It is designed to be long-running, ensuring restarts of the `kubectl` child process in situations where it terminates due to network or timeout errors.

## Installation

To run locally, you can either install it using [Homebrew](https://brew.sh/):

```sh
$ brew install flowerinthenight/tap/oomkill-watch
```
or using your Go environment:

```sh
$ go install github.com/flowerinthenight/oomkill-watch
```

When started, the resulting `kubectl` child process will use its current config for accessing a cluster. Therefore, make sure that your `kubectl` is pointing to the intended cluster first before running the tool.

## Deployment

The provided [Dockerfile](./Dockerfile) is for reference only as it's not configured to access any cluster.

I've only tried deploying this tool on a GKE cluster. You might want to use other commands to configure `kubectl` inside the pod for proper cluster access if you're not using GKE. Here's a snippet of the deployment file I'm using:

```yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  selector:
    matchLabels:
      app: oomkill-watch
  replicas: 1
  revisionHistoryLimit: 5
  template:
    metadata:
      labels:
        app: oomkill-watch
    spec:
      containers:
      - name: oomkill-watch
        image: "{your-image-here}"
        command:
        - '/bin/bash'
        - '-c'
        - |
          gcloud container clusters get-credentials {clustername} && \
          /app/oomkill-watch -slack {channel-webhook}
```

It overrides the command from the `Dockerfile` by running `gcloud container clusters get-credentials ...` first to configure `kubectl`'s cluster access before running the tool.
