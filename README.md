### oomkill-watch

[![main](https://github.com/flowerinthenight/oomkill-watch/actions/workflows/main.yml/badge.svg)](https://github.com/flowerinthenight/oomkill-watch/actions/workflows/main.yml)

`oomkill-watch` is a simple wrapper to the `kubectl get events -w` command, filtering the `OOMKilling` events and optionally forwarding them to a Slack channel. It is designed to be long-running, handling restarts of the `kubectl` child process on events where it terminates due to network or timeout errors.

```sh
# Build:
$ docker build --rm -t oomkill-watch .

# Run using docker:
$ docker run -it --rm oomkill-watch
```
