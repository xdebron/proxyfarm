package main

import (
  "math/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
  "crypto/tls"
  "h12.io/socks"
  "database/sql"
  "time"
  "strings"

  _ "github.com/go-sql-driver/mysql"
)

/*
0-undefined
1-HTTP(S)
2-HTTPS
2-SOCKS4
3-SOCKS4A
4-SOCKS5
*/

type proxy_object struct {
	_id   int
	_ip   string
	_port int
	_type int
}

var checker_db, _ = sql.Open("mysql", DB_CONN_STR)
var checker_queue = make(chan proxy_object, 5000)

func start_checker(){
  for i:=0;i<CHECKER_THREADS;i++{
    go checker_worker()
  }
  populate_checker_queue()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func populate_checker_queue(){
  for {
    if len(checker_queue) < 100{
      rows, err := checker_db.Query("SELECT id,ip,port,type FROM proxylist where last_crawled<UNIX_TIMESTAMP() - 300 ORDER BY last_crawled asc LIMIT 1000")
      if err ==nil {
        for rows.Next() {
          var proxy proxy_object
          rows.Scan(&proxy._id, &proxy._ip, &proxy._port, &proxy._type)
          checker_queue <- proxy
        }
      }
    }
    time.Sleep(time.Second * 1)
  }

}

func make_check_req(client *http.Client) bool {
  req, _ := http.NewRequest("GET", SERVER_URL, nil)
  //req.Close = true
  req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")

  resp, err := client.Do(req)
  if err == nil {
    //return true
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return strings.Contains(string(body), SERVER_URL_CONTAINS)
  } else {
    //fmt.Println(err)
  }
  return false
}

func create_client(proxy proxy_object) *http.Client{
  tr := &http.Transport{}
  client := &http.Client{}
  client.Timeout = 5 * time.Second
  tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

  if proxy._type == 1 {
    u, _ := url.Parse(fmt.Sprintf("http://%s:%d", proxy._ip, proxy._port))
	  tr.Proxy = http.ProxyURL(u)
  }else if proxy._type == 2 {
    u, _ := url.Parse(fmt.Sprintf("https://%s:%d", proxy._ip, proxy._port))
	  tr.Proxy = http.ProxyURL(u)
  }else if proxy._type == 3{
    var dialer = socks.DialSocksProxy(socks.SOCKS4, fmt.Sprintf("%s:%d", proxy._ip, proxy._port))
    tr.Dial = dialer
  }else if proxy._type == 4{
    var dialer = socks.DialSocksProxy(socks.SOCKS4A, fmt.Sprintf("%s:%d", proxy._ip, proxy._port))
    tr.Dial = dialer
  }else if proxy._type == 5{
    var dialer = socks.DialSocksProxy(socks.SOCKS5, fmt.Sprintf("%s:%d", proxy._ip, proxy._port))
    tr.Dial = dialer
  }
  client.Transport = tr
  return client
}

func checker_worker() {
	for {
		check := <-checker_queue
		if check._type != 0 {
        client := create_client(check)
        if make_check_req(client){
          mysql_queue <- fmt.Sprintf("UPDATE proxylist set type=%d, last_crawled=UNIX_TIMESTAMP(), working=1 where id=%d", check._type, check._id)
        } else {
          mysql_queue <- fmt.Sprintf("UPDATE proxylist set last_crawled=UNIX_TIMESTAMP(), working=0 where id=%d", check._id)
        }
		}else{
      for i:=1; i<6;i++{
        try_proxy := check
        try_proxy._type = i
        client := create_client(try_proxy)
  			if make_check_req(client) {//OKİ
  				mysql_queue <- fmt.Sprintf("UPDATE proxylist set type=%d, last_crawled=UNIX_TIMESTAMP(), working=1 where id=%d", i, check._id)
          break
  			} else {
          if i == 5{
            mysql_queue <- fmt.Sprintf("UPDATE proxylist set last_crawled=UNIX_TIMESTAMP(), working=0 where id=%d", check._id)
          }

        }

      }
    }
	}
}
