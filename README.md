# Internet Uptime Dashboard

A simple, self-contained solution for Internet services uptime monitoring writting in go.

## Features

This project trades customisability for simplicity. It consists of a single web-page showing a combination of uptime data charts. All the data is collected using, stored and served using go. The few user settings, such as the port of the web server, can be tweaked in the `const` block of `main.go`. The top of `dashboard.js` contains constants for the graphs' colours and the API server's URL.

## Running instructions

0. Make sure you have `go` installed.
1. Download or clone this repository.
2. Run `go build` to build the server.
3. Run the server, and optionally set it to auto-start on boot.
4. Set the URL of your server in `API_URL` at the top of `dashboard.js`.
5. Be amazed by how bad your ISP's service is.

## TODOs

* Some kind of export to csv.
