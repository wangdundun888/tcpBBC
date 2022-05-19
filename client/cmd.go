package main

import "flag"

type Cmd struct {
	Ip string
}

func NewCmd() *Cmd {
	cmd := &Cmd{}
	flag.StringVar(&cmd.Ip, "ip", "127.0.0.1:9090", "-ip 127.0.0.1:9090")
	flag.Parse()
	return cmd
}
