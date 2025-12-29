package main

import (
	"flag"
	"fmt"
	"taoey/memos-utils/memos-sync/master"
	"taoey/memos-utils/memos-sync/slave"
)

func main() {
	process := flag.String("process", "", "Name of the process")
	flag.Parse()
	switch *process {
	case "master":
		master.Run()
	case "slave":
		slave.Run()
	default:
		fmt.Println("Unknown process")
	}
}
