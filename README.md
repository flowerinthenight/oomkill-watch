## oomkill-watch

[![main](https://github.com/flowerinthenight/oomkill-watch/actions/workflows/main.yml/badge.svg)](https://github.com/flowerinthenight/oomkill-watch/actions/workflows/main.yml)

`oomkill-watch` is a simple wrapper to the `kubectl get events -w` command, filtering the `OOMKilling` events then forwarding them to a Slack channel.

```sh
# Build:
$ docker build --rm -t oomkill-watch .

# Run using docker:
$ docker run -it --rm oomkill-watch
```
