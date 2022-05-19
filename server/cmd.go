package main

import "flag"

type Cmd struct {
	Port string
}

func NewCmd() *Cmd {
	cmd := &Cmd{}
	flag.StringVar(&cmd.Port, "port", "9090", "-port 9090")
	flag.Parse()
	return cmd
}
