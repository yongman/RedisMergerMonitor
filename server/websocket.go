package server

import (
	c "../conf"
	"../inspector"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func StatusServer(ws *websocket.Conn) {
	fmt.Println("ws:new client")
	for {
		<-inspector.ChanDone
		inspector.MapMutex.Lock()
		//choose the info we care
		data, err := json.Marshal(inspector.SlaveInfoFilter(inspector.ServerInfoSnap))
		inspector.MapMutex.Unlock()
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = ws.Write([]byte(data))
		if err != nil {
			break
		}
	}
	fmt.Println("ws:client closed")
}

func RunWebsocketServer(meta *c.MonitorConf) {

	http.Handle("/state", websocket.Handler(StatusServer))

	err := http.ListenAndServe(meta.WsListen, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
