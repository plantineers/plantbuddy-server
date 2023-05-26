# PlantBuddy Server

[![Lint commit message](https://github.com/plantineers/plantbuddy-server/actions/workflows/commit-lint.yml/badge.svg)](https://github.com/plantineers/plantbuddy-server/actions/workflows/commit-lint.yml)
[![Go Workflow](https://github.com/plantineers/plantbuddy-server/actions/workflows/go.yml/badge.svg)](https://github.com/plantineers/plantbuddy-server/actions/workflows/go.yml)
[![GitHub release](https://img.shields.io/github/release/plantineers/plantbuddy-server.svg)](https://github.com/plantineers/plantbuddy-server/releases/)

This server has been developed by:

- [Maximilian Floto](https://github.com/mfloto)
- [Yannick Kirschen](https://github.com/yannickkirschen)

## Usage

- **Start the server:** `go run cmd/main.go` (running on port `3333`)
- **Format the code:** `go fmt ./...`
- **Build everything:** `go build ./...`
- **Run tests:** `go test ./...`

## Contributing

Here are recommendations and some rules for contributing to this project.

### IDE

We suggest using Visual Studio Code with the extensions you find in `.vscode/extensions.json`. You should get
a recommendation to install them when opening the project.

### Save actions

It is mandatory to enable save actions in your IDE! You will need to configure save actions both for Go (`go fmt`) and the `.editorconfig`
file. If you use Visual Studio Code, you can use the extension [Save Actions](https://marketplace.visualstudio.com/items?itemName=emeraldwalk.RunOnSave)
which is already configured for this project.

### Commit messages

This project uses [ConventionalCommits](https://www.conventionalcommits.org/en/v1.0.0/) as style guide for commit messages.
Please note that every commit message will be scanned on push and fails in case the rules are not met!

### Pull requests

Pushing on `main` is forbidden (exceptions apply to [@yannickkirschen](https://github.com/yannickkirschen), as he is silly enough
to deal with Git hell. If he breaks the git history with a force push, he buys the team a shot in a bar of their choice).
Please create a pull request and wait for a review. When merging, squashing is enabled. That means, all commits
will result in a single commit having the title you provide when merging. Please double check this title, as it will be
linted as well!

## Basic directory structure

### Go packages

- `cmd/`: main applications for this project, executable via the command line.
- `config/`: configuration model (read from `buddy.json`, see [Configuration](#configuration))
- `db`: database session model (see [Access the database](#access-the-database))
- `model`: data models for the API (see [api.yaml](./api.yaml))
- `plant`: business logic for working with plants
- `utils`: utility functions (see [utils](#utils))

### Other directories or files

- `/docs/sql`: SQL scripts for development
- `auth-proposal.md`: proposal for authentication and authorization
- `plantbuddy.sqlite`: the database used for development
- `server-requests.http`: example requests for the API

## Configuration

The configuration file `buddy.json` is located in the root directory of the project. Properties
are mapped to the structure `config.Config`. If you want to add a new property, you have to
add it to the structure and to the `buddy.json` file.

Accessing the configuration is done via the global variable `config.PlantBuddyConfig`.

## Access the database

For accessing the database, we use a wrapping session to handle the connection. Our goal is to
have one session for each HTTP request. That means you have to create a session at the beginning of the
logic inside a request's handler. The session is closed at the end of the handler.

```go
var session = db.NewSession()
defer session.Close()

err := session.Open()
if err != nil {
    return nil, err
}

// Do something with session
```

See [Code structure](#code-structure) on how to use session in a repository pattern.

## Code structure

We want to have a dedicated package for every business domain. I.e. the package `plant` contains
all logic for working with plants and the following files:

- `plant-endpoint.go`: contains the HTTP handlers for the plant domain
- `plant-database.go`: contains the repository interface to access the database
- `plant-sqlite.go`: contains the implementation of the repository interface for SQLite

All repository implementations provide a function to create a new repository. This function
always looks roughly like this:

```go
func NewRepository(session *db.Session) (*PlantSqliteRepository, error) {
 if !session.IsOpen() {
  return nil, errors.New("session is not open")
 }

 return &PlantSqliteRepository{db: session.DB}, nil
}
```

**Please take a look on existing structures and functions to learn more on naming stuff.**

## Utils

Here are the utilities we developed to make life easier.

### Get a path parameter

As we use basic HTTP routing, it is difficult to get a path parameter. That's why we
created the function `utils.PathParameterFilter(path string, prefix string) (int64, error)` to solve
this issue.

*Example*: assume your path is `/v1/hello/123` and your prefix is `/v1/hello/`. Calling
`PathParameterFilter` will result in `123` as a path parameter (`int64`).

## Fun

**Get number of lines of code:** `git ls-files | grep '\.go' | xargs wc -l`
