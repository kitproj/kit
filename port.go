package main

import (
	"fmt"
	"log"
	"net"
)

func isPortFree(port uint16) error {
	log.Printf("checking port %d is free\n", port)
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		return err
	}
	_ = listen.Close()
	return nil
}
