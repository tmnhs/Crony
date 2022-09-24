package main

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/server"
	"os"
)

const ServerName = "node"

func main() {
	err := server.InitNodeServer(ServerName)
	if err != nil {
		fmt.Println("init node server error:", err.Error())
		os.Exit(1)
	}
}
