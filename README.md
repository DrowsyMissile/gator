# gator

`gator` is a command-line RSS feed aggregator written in Go. It stores users, feeds, and follows in a PostgreSQL database.

## Requirements

To run gator you'll need:

- **[Go](https://go.dev/doc/install)** 1.27 or newer, to install the CLI
- **[PostgreSQL](https://www.postgresql.org/download/)** 15 or newer, to store data

Go programs compile to static binaries. Once `gator` is installed you can run it without the Go toolchain.

## Installation

Install the `gator` CLI with `go install`:

```sh
go install github.com/DrowsyMissile/gator@latest
```

This puts the `gator` binary in `$(go env GOPATH)/bin` (usually `~/go/bin`). Make sure that directory is on your `PATH`.

## Setup

### 1. Create the database

Start Postgres and create a database for gator:

```sh
psql -c 'CREATE DATABASE gator;'
```

### 2. Create the config file

gator reads its settings from `~/.gatorconfig.json`. Create that file and set `db_url` to your Postgres connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable"
}
```

Replace `username` and `password` with your Postgres credentials. gator adds `current_user_name` to this file itself when you register or log in.

### 3. Run the migrations

The schema lives in `sql/schema` as [goose](https://github.com/pressly/goose) migrations. Clone the repo, install goose, and apply the migrations:

```sh
git clone https://github.com/DrowsyMissile/gator.git
cd gator
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/gator" up
```

## Usage

```sh
gator <command> [args...]
```

Start by registering a user. This also logs you in:

```sh
gator register alice
```

Then add a feed and look at what you follow:

```sh
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
gator following
```

### Commands

| Command | Description |
| --- | --- |
| `register <name>` | Create a user and log in as them |
| `login <name>` | Switch to an existing user |
| `users` | List all users and mark the current one |
| `addfeed <name> <url>` | Add a feed and follow it as the current user |
| `feeds` | List every feed and the user who added it |
| `follow <url>` | Follow a feed someone else already added |
| `following` | List the feeds the current user follows |
| `unfollow <url>` | Stop following a feed |
| `agg` | Fetch an RSS feed and print its contents |
| `reset` | Delete all users, along with their feeds and follows |

`addfeed`, `follow`, `following`, and `unfollow` require a logged-in user.

## Development

During development you can run the code directly without installing it:

```sh
go run . <command> [args...]
```

Use `go run .` only while developing. To use gator normally, run the installed `gator` binary.
