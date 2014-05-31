// Copyright 2014. All rights reserved.
// Use of this source code is governed by a GPLv3
// Author: Wenming Zhang <zhgwenming@gmail.com>

package main

import (
	"github.com/zhgwenming/gbalancer/utils"
)

type Sandbox struct {
	Name    string
	Network string
	Command string
	Args    []string
	Addr    string
}

func NewSandbox(name, network, cmd string) *Sandbox {
	args := make([]string, 0)
	return &Sandbox{name, network, cmd, args, ""}
}

func (s *Sandbox) Run() error {
	cmdline := "virt-sandbox --network address=" + s.Addr + "/24,source=" + s.Network + "  -n " + s.Name + " " + s.Command
	return utils.RunCommand(cmdline)
}

// virt-sandbox /bin/bash --network address=172.16.154.199/24,source=lan  -n test-virt-sandbox
func main() {
	CreateRequiredNetwork()
	sandbox := NewSandbox("test-sandbox", "vnet-eno16777736", "/bin/bash")
	sandbox.Addr = "172.16.154.199"
	sandbox.Run()
}
