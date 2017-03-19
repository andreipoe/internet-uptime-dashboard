# Internet Uptime Dashboard

A simple, self-contained solution for home Internet connection uptime monitoring written in Go.

## Features

This tool trades customisability for simplicity. The frontend consists of a single web-page showing a combination of uptime data charts. All the data is collected, stored and served using a single Go applicaton and a standalone sqlite3 database.

## Configuration

The server's settings can be tweaked in the `const` block of `main.go`. The web dashboard settings are the top of the `dashboard.js` file. 

## Running instructions

### Prerequisites

* A [Go](https://golang.org/doc/install) installation.
* The [go-sqlite3](https://github.com/mattn/go-sqlite3) package.

### Setting up

1. Download or clone this repository.
2. Build the server, e.g. using `go build`.
3. Run the server (and optionally set it to auto-start on boot).
4. Be amazed by how bad your ISP's service is.

## TODOs

* Some kind of export function, e.g. to csv.

----------

## Credits

This project is made possible by free and open source software: [Go](https://golang.org/), [Chart.js](http://www.chartjs.org/), [Moment.js](https://momentjs.com/) and [SQLite](http://www.sqlite.org/).
