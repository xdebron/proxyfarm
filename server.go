package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var server_db, _ = sql.Open("mysql", DB_CONN_STR)

func AddProxy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Write([]byte("<form action='' method='post'><textarea name='proxies'></textarea><input type='submit'></form>"))
	if strings.Compare(req.Method, "POST") == 0 {
		proxies := req.FormValue("proxies")
		proxies_arr := strings.Split(proxies, "\n")
		valid := 0
		nonvalid := 0
		exists := 0
		for _, proxy := range proxies_arr {
			proxy = strings.Trim(proxy, "\r \t")
			proxy_split := strings.Split(proxy, ":")
			if len(proxy_split) == 2 {
				ip := net.ParseIP(proxy_split[0])
				port, err := strconv.ParseInt(proxy_split[1], 10, 16)
				if ip != nil && err == nil && port < 65536 && port > 0 {
					query := fmt.Sprintf("INSERT INTO proxylist VALUES (NULL, '%s', '%d', 0, 0, 0)", ip, port)
					_, err := server_db.Exec(query)
					if err == nil {
						valid += 1
						continue
					} else {
						//fmt.Println(err)
						exists += 1
						continue
					}
				}
			}
			nonvalid += 1
		}
		w.Write([]byte(fmt.Sprintf("<br>Added: %d<br>Not-Valid: %d<br>Already Exists: %d", valid, nonvalid, exists)))
	}
}

func CheckProxy(w http.ResponseWriter, req *http.Request) {
	if strings.Compare(req.Method, "POST") == 0 {
		fmt.Println(req.RemoteAddr)
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		data := req.FormValue("data")
		w.Write([]byte(data))
	} else {
		w.Write([]byte("NOPE."))
	}
}

func ListProxy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	rows, err := server_db.Query("SELECT ip,port,type FROM proxylist where working=1 GROUP BY ip ORDER BY type desc, last_crawled asc")
	if err == nil {
		for rows.Next() {
			var proxy proxy_object
			rows.Scan(&proxy._ip, &proxy._port, &proxy._type)
			w.Write([]byte(fmt.Sprintf("%s:%d %d\n", proxy._ip, proxy._port, proxy._type)))
		}
	}
}

func start_server() {
	http.HandleFunc("/addproxy", AddProxy)
	http.HandleFunc("/listproxy", ListProxy)
	http.HandleFunc("/", CheckProxy)
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err == nil {
		fmt.Println("Http Server Started")
	} else {
		fmt.Println(err)
	}
}
