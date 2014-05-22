package main

import (
	"fmt"
	"net/http"
	"os"
	"jc/snowflakeweb/snowflake"
)


func main() {

	http.HandleFunc("/", hello)
	fmt.Println("listening...")
	port := os.Getenv("PORT")
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	s := snowflake.NewSnowflake(1, 1)
	//fmt.Println(s.NextId())
	fmt.Fprintln(res, s.NextId())
 }
