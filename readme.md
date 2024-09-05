# Minicloud

> Minicloud is a small Minecraft cloud system that allows you to run Minecraft servers easily.

## Features

- Easy to use **(WIP)**
- Fast: We use Golang for the backend with Gate as the proxy.
- Secure: All servers are isolated with Docker.
- Scalable: You can run as many servers as you want.

## Installation

> TODO

## Development

### Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [Go](https://go.dev/doc/install)

### Setting up the project

1. Clone the repository
2. Setup the database with `docker compose -f dev.compose.yml up -d`
3. Copy the example.dev.config.json to config.json
4. Run the project with `go run cmd/plugin.go`
