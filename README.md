# Blog aggregator

Never miss a post from your favorite blogs. `gator` is a command-line tool that lets you subscribe to RSS feeds and automatically pulls new posts into a local database, so you can browse everything in one place, from the terminal, without opening a browser.

Built in Go with a PostgreSQL backend.

---

## What it does

You register an account, add RSS feed URLs, and run the aggregator in the background. It fetches each feed on a schedule you control and stores new posts locally. When you want to read, you browse your latest posts directly from the CLI.

Multiple users can share the same instance — each with their own followed feeds.

---

## Prerequisites

**Go 1.21+**
Download from [go.dev/dl](https://go.dev/dl/) or install via your package manager.

**PostgreSQL**

macOS:
```bash
brew install postgresql@16
brew services start postgresql@16
```

Ubuntu/Debian:
```bash
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

Windows: use the [official installer](https://www.postgresql.org/download/windows/).

---

## Installation

**1. Install the `gator` CLI**

```bash
go install github.com/CamilleOnoda/gator@latest
```

This compiles the program and places the `gator` binary in `$GOPATH/bin`. Make sure that directory is on your `$PATH` — if `gator` isn't found after installing, add this to your shell config (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Once installed, you run the program as `gator` — no Go toolchain needed.

**2. Set up the database**

```bash
createdb gator
```

**3. Create a config file**

The app looks for a config file at `~/.gatorconfig.json`. Create it with your database connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

---

## Quick start

```bash
# Create an account and log in
gator register alice
gator login alice

# Add an RSS feed (creates it and follows it automatically)
gator addfeed "Hacker News" https://news.ycombinator.com/rss

# Start the aggregator — it fetches new posts every 30 seconds
gator agg 30s

# In another terminal, browse your latest posts
gator browse 10
```

---

## All commands

### Account

| Command | Description |
|---|---|
| `register <username>` | Create a new user account |
| `login <username>` | Switch to an existing account |
| `users` | List all registered users |
| `reset` | Delete all users (useful for development) |

### Feeds

| Command | Description |
|---|---|
| `addfeed <name> <url>` | Add a new RSS feed and follow it |
| `feeds` | List all feeds in the system |
| `follow <url>` | Follow an existing feed by URL |
| `following` | List feeds you're currently following |
| `unfollow <url>` | Unfollow a feed |

### Reading

| Command | Description |
|---|---|
| `agg <interval>` | Start fetching feeds on a schedule (e.g. `30s`, `5m`) |
| `browse [limit]` | Show your latest posts (default: 2) |

`agg` runs continuously until you stop it with `Ctrl+C`. Run it in a separate terminal while you use `browse` in another.

---

## Project structure

```
gator/
├── main.go                 # CLI entrypoint: wiring, command registration
├── internal/
│   ├── app/
│   │   └── types.go        # Core application types (State, Command, CLIcommands)
│   │
│   ├── config/
│   │   ├── config.go       # Config file handling (~/.gatorconfig.json)
│   │   └── rss_feed.go     # RSS feed data structures
│   │
│   ├── database/
│   │   ├── db.go           # Database connection setup
│   │   ├── feeds.sql.go    # Generated queries (feeds)
│   │   ├── users.sql.go    # Generated queries (users)
│   │   └── models.go       # Generated models (sqlc)
│   │
│   ├── handlers/
│   │   └── handlers.go     # CLI command handlers (login, feeds, browse, etc.)
│   │
│   └── services/
│       └── services.go     # Core logic (feed fetching, scraping, processing)
│
├── sql/
│   ├── queries/            # Raw SQL queries used by sqlc
│   │   ├── feeds.sql
│   │   └── users.sql
│   │
│   └── schema/             # Database migrations (goose)
│       ├── 001_users.sql
│       ├── 002_feeds.sql
│       ├── 003_feed_follows.sql
│       ├── 004_feed_lastfetched.sql
│       └── 005_posts.sql
│
├── sqlc.yaml               # sqlc configuration
├── go.mod / go.sum         # Go module dependencies
└── README.md
```

---

## Tech stack

- **Go** — standard library HTTP client, XML parsing, context-based timeouts
- **PostgreSQL** — stores users, feeds, follows, and posts
- **sqlc** — generates type-safe Go code from raw SQL
- **goose** — database schema migrations
- **google/uuid** — unique IDs for all records
- **lib/pq** — PostgreSQL driver for Go
