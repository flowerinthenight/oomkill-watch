[![main](https://github.com/flowerinthenight/oomkill-watch/actions/workflows/main.yml/badge.svg)](https://github.com/flowerinthenight/oomkill-watch/actions/workflows/main.yml)

## Overview

`oomkill-watch` is a simple wrapper to the `kubectl get events -w` command, specifically filtering the `OOMKilling` events and optionally forwarding them to a Slack channel. It is designed to be long-running, handling restarts of the `kubectl` child process in situations where it terminates due to network or timeout errors.

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

```sh
# Build:
$ docker build --rm -t oomkill-watch .

# Run using docker:
$ docker run -it --rm oomkill-watch
```
