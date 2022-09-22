package main

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/config"
	"os"
)

func main() {
	err := config.LoadConfig("admin")
	if err != nil {
		fmt.Println(err)
		os.Exit(111)
	}
}
