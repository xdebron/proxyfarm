package main

func main() {
	go fire_and_forget_querier()
	go start_server()
	start_checker()
	/*	var proxy proxy_object
		proxy._id = 0
		proxy._ip = "51.79.31.19"
		proxy._port = 8080
		proxy._type = 1
		checker_queue <- proxy*/
	//send()
}
