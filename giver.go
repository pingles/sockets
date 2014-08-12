package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		panic(err)
	}
	
	fmt.Println("giving")
	
	for i := 0;; i++ {
		_, err := fmt.Fprintf(conn, "giving %d\n", i)
		fmt.Printf("gave %d\n", i)
		if err != nil {
			fmt.Println("error giving, sleeping")
			d, _ := time.ParseDuration("10s")
			<- time.After(d)
			
			panic(err)
		}
	}
}