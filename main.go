package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"net"
	"net/http"
)

const (
	// Dashboard
	DASHBOARD_PORT        = 8080

	// Monitoring
	UPTIME_CHECK_INTERVAL = 30 * time.Second
	UPTIME_CHECK_URL      = "https://www.google.com"

	// Logging
	LOG_UP   = false
	LOG_DOWN = true
	LOG_DBG  = false
)

func checkInternetUp() bool {
	var timeouts = [...]int32 { 2, 4, 8 }

	for _, timeout := range timeouts {
		timeoutDuration := time.Duration(timeout) * time.Second

		// Check that we can open a connection
		conn, err := net.DialTimeout("tcp", "www.google.com:https", timeoutDuration)
        if conn != nil {
            defer conn.Close()
        }

		if err == nil {
			remoteAddr := conn.RemoteAddr().String()
			if LOG_DBG {
				log.Printf("Reachable %s\n", remoteAddr)
			}

			// Check that the intended other side responsed, i.e. not someone on the local network
			if strings.HasPrefix(remoteAddr, "192.168.") {
				if LOG_DBG {
					log.Printf("Response came from inside the local network.")
				}
			} else {
				// Check that we can load a page
				client := http.Client {
					Timeout: timeoutDuration,
				}
				response, err := client.Get(UPTIME_CHECK_URL)
				if err != nil {
					if LOG_DBG {
						log.Printf("Failed to GET: %v\n", err)
					}
				} else {
					defer response.Body.Close()
					if LOG_DBG {
						log.Println("Response status", response.StatusCode)
					}

					if response.StatusCode < 400 {
						return true
					}
				}
			}
		} else {
			if LOG_DBG {
				log.Printf("Dial error: %v\n", err)
			}
		}

		if LOG_DBG {
			log.Printf("Unreachable %d\n", timeout)
		}
	}

	return false
}

func monitorUptime() {
	for {
		if checkInternetUp() {
			if LOG_UP {
				log.Println("Internet is up!")
			}
		} else {
			if LOG_DOWN {
				log.Println("Internet is DOWN!")
			}
		}

		time.Sleep(UPTIME_CHECK_INTERVAL)
	}
}

func main() {
	go monitorUptime()

	server := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", server)

	fmt.Printf("Listening on port %d...\n", DASHBOARD_PORT)
	http.ListenAndServe(":" + strconv.Itoa(DASHBOARD_PORT), nil)
}
