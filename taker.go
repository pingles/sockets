package main

import (
	"fmt"
	"bufio"
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"
)

func handleConn(controlChan chan bool, conn net.Conn) {
	fmt.Println("connected to giver, waiting 10s")
	// sleep for 10s
	d, _ := time.ParseDuration("10s")
	<- time.After(d)
	fmt.Println("taking from giver")
	
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// reached EoF
	controlChan <- true 

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading from connection:", err)
	}
}

func main() {
	connections := []net.Conn{}
	
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	
	controlChan := make(chan bool, 1)
	
	go func() {
		for sig := range signals {
			switch sig {
			case syscall.SIGINT:
				for _, conn := range connections {
					fmt.Println("CloseWrite()")
					conn.(*net.TCPConn).CloseRead()
				}
				<- controlChan
				fmt.Println("control channel triggered, exiting")
				os.Exit(0)
				break
			}
		}
	}()

	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		connections = append(connections, conn)
		
		if err != nil {
			panic(err)
		}
		go handleConn(controlChan, conn)
	}
}