// Copyright 2014. All rights reserved.
// Use of this source code is governed by a GPLv3
// Author: Wenming Zhang <zhgwenming@gmail.com>

package main

import (
	"github.com/docker/spdystream"
	"net"
	"time"
)

const (
	STREAMPORT = "6900"
)

type Backend struct {
	spdyconn *spdystream.Connection
	address  string
	flags    int
	index    int
	ongoing  uint
	RxBytes  uint64
	TxBytes  uint64
}

func NewBackend(addr string, useTunnel bool) *Backend {
	b := &Backend{address: addr}
	if useTunnel {
		if conn, err := NewStreamConn(addr, STREAMPORT); err == nil {
			b.spdyconn = conn
		}
	}
	return b
}

func NewStreamConn(addr, port string) (*spdystream.Connection, error) {
	conn, err := net.DialTimeout("tcp", addr+":"+port, time.Second)
	if err != nil {
		log.Printf("dail spdy error: %s", err)
		return nil, err
	}

	spdyConn, err := spdystream.NewConnection(conn, false)
	if err != nil {
		log.Printf("spdystream create connection error: %s", err)
		return nil, err
	}

	go spdyConn.Serve(spdystream.NoOpStreamHandler)
	if _, err = spdyConn.Ping(); err != nil {
		return nil, err
	} else {
		return spdyConn, nil
	}
}
