package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"net"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Dashboard and API Server
	SERVER_PORT = 18080

	// Monitoring
	UPTIME_CHECK_INTERVAL = 30 * time.Second
	UPTIME_CHECK_URL      = "https://www.google.com"

	// Logging
	LOG_UP   = false
	LOG_DOWN = true
	LOG_DBG  = false
	LOG_API  = true

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

func apiInstant(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]

	// Ignore all methods apart from GET
	if r.Method != "GET" {
		w.WriteHeader(http.StatusTeapot)
		if LOG_API {
			log.Println("GET/daily request from", remoteIP)
		}
		return
	}

	if LOG_API {
		log.Println("GET/instant request from", remoteIP)
	}

	// Read the data from the database
	rows, err := db.Query("SELECT * FROM " + DB_TABLE_INSTANT)
	if err != nil {
		log.Printf("Failed to read from table %s: %v\n", DB_TABLE_INSTANT, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// JSON object template struct
	type instantInfo struct {
		Timestamp  string `json:"timestamp"`
		State      bool   `json:"up"`
	}

	// Put the data into a slice of structs
	var data []instantInfo
	for rows.Next() {
		var t time.Time
		var state bool

		err = rows.Scan(&t, &state)
		checkDbError(err)

		data = append(data, instantInfo{t.Format("2006-01-02 15:04"), state})
	}

	// Encode the data as JSON
	jdata, err := json.Marshal(map[string]interface{} {"data": data})
	if err != nil {
		log.Println("Failed to encode instant data into JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the JSON data to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jdata)
}

func apiDaily(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]

	// Ignore all methods apart from GET
	if r.Method != "GET" {
		w.WriteHeader(http.StatusTeapot)
		if LOG_API {
			log.Println("Not responing to method", r.Method, "from", remoteIP)
		}
		return
	}

	if LOG_API {
		log.Println("GET/daily request from", remoteIP)
	}

	// Read the data from the database
	rows, err := db.Query("SELECT * FROM " + DB_TABLE_DAILY)
	if err != nil {
		log.Printf("Failed to read from table %s: %v\n", DB_TABLE_DAILY, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// JSON object template struct
	type dailyInfo struct {
		Day  string `json:"day"`
		Up   int    `json:"up"`
		Down int    `json:"down"`
	}

	// Put the data into a slice of structs
	var data []dailyInfo
	for rows.Next() {
		var day time.Time
		var up, down int

		err = rows.Scan(&day, &up, &down)
		checkDbError(err)

		data = append(data, dailyInfo{day.Format("2006-01-02"), up, down})
	}

	// Encode the data as JSON
	jdata, err := json.Marshal(map[string]interface{} {"data": data})
	if err != nil {
		log.Println("Failed to encode daily data into JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the JSON data to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jdata)
}

func main() {
	db = initDatabase()
	defer db.Close()

	go monitorUptime()
	go db.cleanupOldInstantData()

	// Dashboard server
	server := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", server)

	// Data API
	http.HandleFunc("/api/instant", apiInstant)
	http.HandleFunc("/api/daily", apiDaily)

	fmt.Printf("Listening on port %d...\n", SERVER_PORT)
	http.ListenAndServe(":" + strconv.Itoa(SERVER_PORT), nil)
}
