package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"jc/snowflakeweb/snowflake"
)

const workers uint64 = 5
var requestChannel [workers]chan (chan uint64)
var responseChannel chan uint64
var currentWorker int

func startSnowflake(datacenterId uint64, workerId uint64) {
	s := snowflake.NewSnowflake(datacenterId, workerId)
	requestChannel[workerId] = make(chan (chan uint64))

	for {
		responseChan := <- requestChannel[workerId]
		fmt.Println(workerId, "handling request")
		id := s.NextId()
		fmt.Println("id:", id)
		responseChan <- id
	}
	
}

func main() {

	currentWorker = 0
	responseChannel = make(chan uint64)
	datacenterId, _ := strconv.ParseUint(os.Getenv("DATACENTER_ID"), 0, 64)
	var i uint64

	for i = 0; i < workers; i++ {
		go startSnowflake(datacenterId, i)
	}

	http.HandleFunc("/", handler)
	fmt.Println("listening...")
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}

func handler(res http.ResponseWriter, req *http.Request) {

	requestChannel[currentWorker] <- responseChannel
	currentWorker = (currentWorker + 1) % int(workers)

	id := <- responseChannel
	fmt.Fprintln(res, id)

}
