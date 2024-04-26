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

```sh
# Build:
$ docker build --rm -t oomkill-watch .

# Run using docker:
$ docker run -it --rm oomkill-watch
```
