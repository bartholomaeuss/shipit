# Shipit

ShipIt is a Go-powered companion CLI for the classic “clone → tweak → ship it to a test box” workflow. It wraps a few high-friction steps—cloning repositories into isolated scratch space, syncing those clones to a remote host, cleaning them up, and triggering remote Docker-based deploy scripts—into a single tool.

## Features

- **Ephemeral cloning** – `shipit repo clone` creates a uniquely named temp directory, clones the requested GitHub repository, and copies it to a target host via `scp`, keeping your laptop and test box tidy.
- **One-command cleanup** – `shipit repo clean` deletes local scratch directories and can prune the matching directories on your remote host over SSH.
- **Remote deploy script runner** – `shipit deploy run` connects to your host and executes `scripts/docker/prerun.sh`, `run.sh`, and `postrun.sh` so you can rehearse release automation without leaving your terminal.

## Command map

| Command | Purpose |
| --- | --- |
| `shipit repo clone` | Clone a GitHub repository into a temp dir and copy it to the remote host’s home directory. |
| `shipit repo clean` | Delete Shipit temp directories locally and on the remote host (all or a specific directory). |
| `shipit deploy run` | SSH into the host and execute `scripts/docker/prerun.sh`, `run.sh`, and `postrun.sh` in order. |

## Requirements

- Go **1.22** or newer (see `go.mod`)
- Git, SSH, and SCP available in your `PATH`
- Access to a remote host configured in `~/.ssh/config` (defaults to `test`)
- Docker scripts (`scripts/docker/{prerun,run,postrun}.sh`) checked into the repositories you want to deploy

_ShipIt relies on your existing SSH tooling—`ssh`, `scp`, Host entries, `~/.ssh/config`, agent forwarding, etc.—and does not attempt to configure or override any of it._

## Build

This repo is meant to be opened inside the provided dev container (`.devcontainer/devcontainer.json`). The container image already supplies Go 1.22, Cobra CLI, Delve, and the rest of the tooling, so you only need to run the build command and download the resulting binary from the container to your host OS.

```bash
GOOS=windows GOARCH=amd64 go build -buildvcs=false -o shipit .
```

## Usage

ShipIt is a Cobra CLI. Run `shipit --help` for the top-level summary and `shipit <command> --help` for command-specific flags and defaults.

### Repository workflow

#### Clone into scratch space

```bash
shipit repo clone \
  --url https://github.com/example/project.git \
  --user jane \
  --host staging
```

- `--url` (required): GitHub repository to clone.
- `--user` (required): SSH username; used when copying the clone to the remote host.
- `--host`: SSH host alias from your config; defaults to `test`.

The command prints the temporary directory and the `cd` instruction so you can immediately jump into the clone locally, while also copying it to `~` on the remote host.

#### Clean up scratch directories

Delete every ShipIt scratch directory locally and remotely:

```bash
shipit repo clean --all --user jane --host staging
```

Delete just one directory (locally and remotely) by passing `--specific-dir /tmp/shipit-repo-12345`. All cleanup commands are safety-checked to only touch paths prefixed with `shipit-repo-`.

### Deployment workflow

```bash
shipit deploy run \
  --dir ~/shipit-repo-12345 \
  --user jane \
  --host staging
```

`shipit deploy run` SSHes to the host, ensures the user context is correct, and executes the repository’s Docker helper scripts in order: `prerun.sh`, `run.sh`, then `postrun.sh`. Use `--dir` to tell ShipIt which remote repository directory owns those scripts. Flags mirror the repo commands: `--host` defaults to `test` and `--user` is optional if your SSH config already defines one.

## Example end-to-end flow

1. Clone and copy a repo: `shipit repo clone --url … --user jane --host staging`
2. Make changes locally or remotely; test on the remote host with `shipit deploy run --dir ~/shipit-repo-… --user jane`
3. Clean out temp directories when you are done: `shipit repo clean --all --user jane --host staging`

## Development

- Format and vet: `go fmt ./... && go vet ./...`
- Run any future tests with `go test ./...`
- Generate new Cobra commands with `cobra-cli add <name>`

## License

Distributed under the MIT License. See `LICENSE` for details.
