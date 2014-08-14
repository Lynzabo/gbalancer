// Copyright 2014. All rights reserved.
// Use of this source code is governed by a GPLv3
// Author: Wenming Zhang <zhgwenming@gmail.com>

package main

import (
	"flag"
	"fmt"
	"github.com/docker/spdystream"
	logger "github.com/zhgwenming/gbalancer/log"
	"github.com/zhgwenming/gbalancer/utils"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var (
	wgroup  = &sync.WaitGroup{}
	log     = logger.NewLogger()
	sigChan = make(chan os.Signal, 1)
	pidFile = flag.String("pidfile", "", "pid file")
)

func init() {
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *pidFile != "" {
		if err := utils.WritePid(*pidFile); err != nil {
			fmt.Printf("error: %s\n", err)
			log.Fatal("error:", err)
		}
		defer func() {
			if err := os.Remove(*pidFile); err != nil {
				log.Printf("error to remove pidfile %s:", err)
			}
		}()
	}

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			spdyConn, err := spdystream.NewConnection(conn, true)
			if err != nil {
				panic(err)
			}
			go spdyConn.Serve(spdystream.MirrorStreamHandler)
		}
	}()

	// waiting for exit signals
	for sig := range sigChan {
		log.Printf("captured %v, exiting..", sig)

		return
	}
}
