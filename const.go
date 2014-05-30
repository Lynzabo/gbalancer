// Copyright 2014. All rights reserved.
// Use of this source code is governed by a GPLv3
// Author: Wenming Zhang <zhgwenming@gmail.com>

package main

const (
	MaxBackends             uint = 128
	MaxForwarders           uint = 8192
	MaxForwardersPerBackend uint = 4096
)

const (
	WsrepAddresses = "wsrep_incoming_addresses"
	WsrepConnected = "wsrep_connected"
	ListenAddr     = "127.0.0.1"
	ListenPort     = "3306"
	CheckInterval  = 60
	IPvsAddr       = "127.1.1.1"
)

const (
	FlagDown int = 0
	FlagUp   int = 1
)

const (
	ReqRefused int = 1
)
