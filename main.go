// Copyright 2014. All rights reserved.
// Use of this source code is governed by a GPLv3
// Author: Wenming Zhang <zhgwenming@gmail.com>

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zhgwenming/gbalancer/config"
	"github.com/zhgwenming/gbalancer/engine"
	logger "github.com/zhgwenming/gbalancer/log"
	"github.com/zhgwenming/gbalancer/wrangler"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

type Request struct {
	conn    net.Conn
	backend *Backend
	err     error
}

type Forwarder struct {
	backend *Backend
	request *Request
	bytes   uint
}

var (
	wgroup       = &sync.WaitGroup{}
	log          = logger.NewLogger()
	sigChan      = make(chan os.Signal, 1)
	configFile   = flag.String("config", "gbalancer.json", "Configuration file")
	unixSocket   = flag.Bool("unixsock", false, "listen to unix domain socket in addition - default path /var/lib/mysql/mysql.sock")
	failover     = flag.Bool("failover", false, "whether to enable failover mode for scheduling")
	daemonMode   = flag.Bool("daemon", false, "daemon mode")
	ipvsMode     = flag.Bool("ipvs", false, "to use lvs as loadbalancer")
	ipvsRemote   = flag.Bool("remote", false, "independent director")
	printVersion = flag.Bool("version", false, "print gbalancer version")
)

func init() {
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
}

func PrintVersion() {
	fmt.Printf("gbalancer version: %s\n", VERSION)
	os.Exit(0)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if *printVersion {
		PrintVersion()
	}

	file, _ := os.Open(*configFile)

	if *daemonMode {
		os.Chdir("/")
	}

	decoder := json.NewDecoder(file)
	config := config.Configuration{
		Service:    "galera",
		Addr:       "127.0.0.1",
		Port:       "3306",
		UnixSocket: DEFAULT_UNIX_SOCKET,
	}

	err := decoder.Decode(&config)
	if err != nil {
		log.Println("error:", err)
	}
	//log.Printf("%v", config)
	log.Printf("Listen on %s:%s, backend: %v", config.Addr, config.Port, config.Backend)

	tcpAddr := config.Addr + ":" + config.Port

	status := make(chan map[string]int, MaxBackends)
	//status := make(chan *BEStatus)

	// start the wrangler
	wgl := wrangler.NewWrangler(config, status)

	go wgl.Monitor()

	done := make(chan struct{})
	if *ipvsMode {
		wgroup.Add(1)
		if *ipvsRemote {
			ipvs := engine.NewIPvs(config.Addr, config.Port, "wlc", done, wgroup)
			go ipvs.RemoteSchedule(status)
		} else {
			//ipvs := NewIPvs(IPvsLocalAddr, config.Port, "sh", done)
			ipvs := engine.NewIPvs(IPvsLocalAddr, config.Port, "wlc", done, wgroup)
			go ipvs.LocalSchedule(status)
		}
	} else {
		listener, err := net.Listen("tcp", tcpAddr)

		if err != nil {
			log.Fatal(err)
		}

		job := make(chan *Request)

		// start the scheduler
		sch := NewScheduler(*failover)
		go sch.schedule(job, status)

		// tcp listener
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("%s\n", err)
				}
				//log.Println("main: got a connection")
				req := &Request{conn: conn}
				job <- req
			}
		}()

		// unix listener
		if *unixSocket {
			ul, err := net.Listen("unix", config.UnixSocket)

			if err != nil {
				fmt.Printf("%s\n", err)
				log.Fatal(err)
			}
			go func() {
				for {
					conn, err := ul.Accept()
					if err != nil {
						log.Printf("%s\n", err)
					}
					//log.Println("main: got a connection")
					req := &Request{conn: conn}
					job <- req
				}
			}()
		}
	}
	for sig := range sigChan {
		log.Printf("captured %v, exiting..", sig)
		close(done)
		wgroup.Wait()
		return
	}

}
