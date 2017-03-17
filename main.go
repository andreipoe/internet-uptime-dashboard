package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"net"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Dashboard
	DASHBOARD_PORT = 8080

	// Monitoring
	UPTIME_CHECK_INTERVAL = 30 * time.Second
	UPTIME_CHECK_URL      = "https://www.google.com"

	// Logging
	LOG_UP   = false
	LOG_DOWN = true
	LOG_DBG  = false

	// Databse
	DB_NAME              = "uptime.db"
	DB_TABLE_INSTANT     = "instant_feed"
	DB_TABLE_DAILY       = "daily_averages"
	DB_CLEANUP_INTERVAL  = 24 * time.Hour
	DB_CLEANUP_THRESHOLD = -48 * time.Hour
)

var db *UptimeDB

func checkDbError(err error) {
	if err != nil {
		log.Println("Database error:", err)
	}
}

type UptimeDB struct {
	*sql.DB
}

// Open the database and create tables if needed
func initDatabase() *UptimeDB {
	// var err error
	sqlDb, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		log.Println("Database initialisation error:", err)
		panic(err)
	}

	db := &UptimeDB {
		sqlDb,
	}

	if !db.tableExists(DB_TABLE_INSTANT) {
		_, err = db.Exec("CREATE TABLE " + DB_TABLE_INSTANT +
			"(time TIMESTAMP PRIMARY KEY, up BOOL)")
		checkDbError(err)
	}
	if !db.tableExists(DB_TABLE_DAILY) {
		_, err = db.Exec("CREATE TABLE " + DB_TABLE_DAILY +
			"(date TIMESTAMP PRIMARY KEY, up_count INTEGER, down_count INTEGER)")
		checkDbError(err)
	}

	return db
}

// Check whether a table already exists in the database
func (db *UptimeDB) tableExists(table string) bool {
	var name string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)

	if name == table {
		return true
	} else if err != sql.ErrNoRows {
		checkDbError(err)
	}

	return false
}

// Update the database following the latest check
func (db *UptimeDB) recordInternetState(state bool) {
	now     := time.Now()

	// Start a transaction for the two table updates
	txOk    := true
	tx, err := db.Begin()
	if err != nil {
		txOk = false
	}

	// Update the instant table
	_, err = db.Exec("INSERT INTO " + DB_TABLE_INSTANT + "(time, up) VALUES(?,?)", now,state)
	if err != nil {
		txOk = false
	}

	// Ensure the daily table has a row for today
    year, month, day := now.Date()
    today			 := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	var upCount, downCount int
	err = db.QueryRow("SELECT up_count,down_count FROM " + DB_TABLE_DAILY + " WHERE date=?", today).Scan(&upCount, &downCount)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err1 := db.Exec("INSERT INTO " + DB_TABLE_DAILY + "(date,up_count,down_count) VALUES(?,?,?)",
				today,0,0)
			if err1 != nil {
				txOk = false
			}

			upCount   = 0
			downCount = 0
		} else {
			txOk = false
		}
	}

	// Update the daily table
	if state {
		upCount++
	} else {
		downCount++
	}
	_, err = db.Exec("INSERT OR REPLACE INTO " + DB_TABLE_DAILY + "(date,up_count,down_count) VALUES(?,?,?)",
		today,upCount,downCount)
	if err != nil {
		txOk = false
	}

	// Commit only if all operations complted succesfully
	if txOk {
		tx.Commit()
		if LOG_DBG {
			log.Printf("Today's checks now: %d up, %d down\n", upCount, downCount)
		}
	} else {
		log.Println("Database error when recording state")
		tx.Rollback()
	}
}

// Periodically delete old records from the instant table after they have fallen outside the dashbaord's window
func (db *UptimeDB) cleanupOldInstantData() {
	for {
		threshold := time.Now().Add(DB_CLEANUP_THRESHOLD)

		result, err := db.Exec("DELETE FROM " + DB_TABLE_INSTANT + " WHERE time < datetime(?)", threshold)
		checkDbError(err)
		if LOG_DBG && err == nil {
			deleted, _ := result.RowsAffected()
			log.Printf("Deleted %d instant records before %v\n", deleted,threshold)
		}

		time.Sleep(DB_CLEANUP_INTERVAL)
	}
}

// Check if Google is reachable
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

// Periodically check if the connection is up
func monitorUptime() {
	for {
		state := checkInternetUp()
		if state {
			if LOG_UP {
				log.Println("Internet is up!")
			}
		} else {
			if LOG_DOWN {
				log.Println("Internet is DOWN!")
			}
		}
		db.recordInternetState(state)

		time.Sleep(UPTIME_CHECK_INTERVAL)
	}
}

func main() {
	db = initDatabase()
	defer db.Close()

	go monitorUptime()
	go db.cleanupOldInstantData()

	server := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", server)

	fmt.Printf("Listening on port %d...\n", DASHBOARD_PORT)
	http.ListenAndServe(":" + strconv.Itoa(DASHBOARD_PORT), nil)
}
