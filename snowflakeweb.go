package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"github.com/jazzychad/gosnowflake/snowflake"
)

const workers uint64 = 5
var requestChannel [workers]chan (chan uint64)
var responseChannel chan uint64
var currentWorker int

func startSnowflake(datacenterId uint64, workerId uint64) {
	// create a snowflake generator
	s := snowflake.NewSnowflake(datacenterId, workerId)

	// make the request channel to read
	requestChannel[workerId] = make(chan (chan uint64))

	for {
		// get channel to reply on
		responseChan := <- requestChannel[workerId]
		fmt.Println(workerId, "handling request")

		// generate id
		id := s.NextID()
		fmt.Println("id:", id)

		// send the id back on the response channel
		responseChan <- id
	}
	
}

func main() {

	// setup
	currentWorker = 0
	responseChannel = make(chan uint64)
	datacenterId, _ := strconv.ParseUint(os.Getenv("DATACENTER_ID"), 0, 64)

	// start some workers
	var i uint64
	for i = 0; i < workers; i++ {
		go startSnowflake(datacenterId, i)
	}

	// start the http server
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	fmt.Println("listening...")
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}

func handler(res http.ResponseWriter, req *http.Request) {
	// handle request
	requestChannel[currentWorker] <- responseChannel
	currentWorker = (currentWorker + 1) % int(workers)

	id := <- responseChannel
	fmt.Fprintln(res, id)
}

func faviconHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "")
}
