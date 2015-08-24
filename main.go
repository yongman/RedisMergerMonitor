package main

import (
	"./conf"
	"./inspector"
	"./server"
	"fmt"
)

func main() {
	//init
	inspector.Init()

	//read meta config
	meta, err := conf.LoadConf("redis-monitor.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	go server.RunWebsocketServer(meta)
	go inspector.Run(meta)
	server.RunHttpServer(meta)
}
