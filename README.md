
# proxyfarm
Farm proxies all around the world

**proxyfarm** is a proxy farming tool written with GoLang for high performance and concurrency. Currently can handle 130k proxies with 1 gb ram and 1 core processor.

# TODO

 - Add syn scan before proxy checking.
 - Add syn scan for subnets of existing proxies.

# Dependencies

 - github.com/go-sql-driver/mysql
 - h12.io/socks
 - And of course working mysql server.
# Usage
First of all you need to create new database for proxyfarm using provided sql file. Install dependencies given above then compile and run. It will immediately create https server for you.
 - https://localhost/addproxy for adding proxies to database manually.
 - https://localhost/listproxy for listing working proxies.
