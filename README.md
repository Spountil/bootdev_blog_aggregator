# Gator - Blog Aggregator CLI

A command-line RSS feed aggregator and reader built in Go and backed by PostgreSQL. This project is based on the "Build a Blog Aggregator" project from [Boot.dev](https://boot.dev).

`gator` allows multiple users to manage and follow RSS feeds, runs a continuous background scraping worker to ingest posts into a PostgreSQL database, and provides tools to browse stored posts directly from the terminal.

---

## Features

- **User Management**: Register and switch between users with active session persistence.
- **Feed Management**: Add RSS feeds, list all available feeds, and follow/unfollow feeds per user.
- **Automated Aggregator**: Continuously scrape and ingest posts across registered feeds using configurable interval tickers.
- **Post Reader**: Browse collected posts from feeds you follow directly in the CLI.

---

## Prerequisites

Before running the program, ensure you have the following installed:

1. **[Go](https://go.dev/doc/install)** (version 1.22 or higher)
2. **[PostgreSQL](https://www.postgresql.org/download/)** database server running locally or accessible remotely
3. **[Goose](https://github.com/pressly/goose)** (recommended for database schema migrations)
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

Make sure your Go binaries directory (`$(go env GOPATH)/bin`) is included in your system `PATH`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Installation

You can compile and install the CLI using [`go install`](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies).

### From Local Source

Clone the repository and install from the project root:

```bash
git clone https://github.com/Spountil/bootdev_blog_aggregator.git
cd bootdev_blog_aggregator
go install .
```

*Note: By default, `go install .` compiles the executable using the module name (`bootdev_blog_aggregator`). If you prefer the binary named `gator`, you can build it with:*

```bash
go build -o gator .
```

### From Remote Package

You can also install directly from GitHub:

```bash
go install github.com/Spountil/bootdev_blog_aggregator@latest
```

---

## Setup & Configuration

### 1. Database Setup

Ensure PostgreSQL is running, then create a new database for the aggregator (e.g., `gator`):

```bash
createdb gator
```

### 2. Database Schema Migrations

Apply the database schema migrations using Goose from the `sql/schema` directory:

```bash
cd sql/schema
goose postgres "postgres://<username>:<password>@localhost:5432/gator?sslmode=disable" up
cd ../..
```

*(Replace `<username>` and `<password>` with your PostgreSQL credentials.)*

### 3. Config File Setup

`gator` requires a JSON configuration file named `.gatorconfig.json` located in your user's home directory (`~/.gatorconfig.json`).

Create `~/.gatorconfig.json` with the following content:

```json
{
  "connection": "postgres://<username>:<password>@localhost:5432/gator?sslmode=disable",
  "db_url": "postgres://<username>:<password>@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

- **`connection` / `db_url`**: Your PostgreSQL connection string.
- **`current_user_name`**: The currently logged-in user (updated automatically by `register` and `login` commands).

---

## Available Commands

Run commands using `gator <command> [arguments]` (or `bootdev_blog_aggregator <command> [arguments]` depending on how you built the binary).

| Command | Arguments | Description |
| :--- | :--- | :--- |
| `register` | `<username>` | Registers a new user and sets them as the current active user. |
| `login` | `<username>` | Switches the current active user to an existing user. |
| `users` | *(none)* | Lists all registered users and indicates the current active user. |
| `addfeed` | `<name> <url>` | Adds a new RSS feed and automatically marks it as followed by the current user. |
| `feeds` | *(none)* | Lists all feeds along with their URL and the user who added them. |
| `follow` | `<url>` | Follows an existing feed by its URL for the current user. |
| `following` | *(none)* | Displays all feeds followed by the current user. |
| `unfollow` | `<url>` | Unfollows a feed for the current user. |
| `agg` | `<time_between_reqs>` | Starts the aggregator worker to scrape feeds continuously at the specified interval (e.g., `1m`, `30s`, `1h`). |
| `browse` | `[limit]` | Displays stored posts from the database (optional limit defaults to 2). |
| `reset` | *(none)* | Resets/deletes all users in the database (warning: cascades to feeds, follows, and posts). |

---

## Example Workflow

Here is a quick walkthrough to get started:

1. **Register a new user:**
   ```bash
   gator register lane
   ```

2. **Add an RSS feed:**
   ```bash
   gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
   ```

3. **Check following list:**
   ```bash
   gator following
   ```

4. **Start the background aggregator worker:**
   ```bash
   gator agg 1m
   ```
   *(Press `Ctrl+C` to stop the aggregator once feeds have been fetched).*

5. **Browse the collected posts:**
   ```bash
   gator browse 5
   ```

---

## Acknowledgements

- Built as part of the backend curriculum on [Boot.dev](https://boot.dev).
- Uses [sqlc](https://sqlc.dev/) for generating type-safe Go code from SQL queries.
- Uses [goose](https://github.com/pressly/goose) for database migrations.
