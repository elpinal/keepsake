# ARCHITECTURE

## Directory Structure

### `entry/`

In Keepsake, each web page is recorded as an "entry" with a time stamp which
denotes when the entry is added.

### `gettitle/`

To get the title of a web page, we use a simple heuristic algorithm in this package.

### `log/`

We define a package called `log`, other than the `log` package in the Go
standard library, because we want to use log levels and emit logs as JSON
objects.

### `resources/`

HTML and CSS.

### `server/`

Handlers for the web server. This package does not depend on any specific
database. Instead, it defines and depends on a storage interface.

### `storage/`

This package implements the storage interface. Currently we use SQLite3 as a
backend.

### `main.go`

The entry point of the executable, `keepsake`.
