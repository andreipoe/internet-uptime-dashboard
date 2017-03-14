package main

import (
	"fmt"
	"strconv"
	"net/http"
)

const DASHBOARD_PORT = 8080

func main () {
	server := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", server)

	fmt.Printf("Listening on port %d...\n", DASHBOARD_PORT)
	http.ListenAndServe(":" + strconv.Itoa(DASHBOARD_PORT), nil)
}
