package inspector

import (
	"../conf"
	"../redis"
	"fmt"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

type InfoMap map[string]*redis.RedisInfo

var ServerInfoSnap InfoMap
var MapMutex sync.RWMutex
var ChanDone chan string

var addrFilterMap map[string]bool

type Client struct {
	C  chan string
	Ws *websocket.Conn
}

var wsClients []*Client

func Init() {
	ServerInfoSnap = make(map[string]*redis.RedisInfo)
	addrFilterMap = make(map[string]bool)
	MapMutex = sync.RWMutex{}
	ChanDone = make(chan string)
	wsClients = make([]*Client, 0)
}

func Run(meta *conf.MonitorConf) {
	for {
		inner := func(meta *conf.MonitorConf) {
			MapMutex.Lock()
			defer MapMutex.Unlock()
			for _, s := range meta.Servers {
				if meta.ServerType == "redis" {
					info := fetchRedisServerInfo(s, meta.ServerType)
					if info == nil {
						continue
					}
					infoFilter(info)
					ServerInfoSnap[s] = info

					addrFilterMap[s] = true
					//check the master is in servers to monitor
					if info.Get("role") == "slave" {
						ip := info.Get("master_host")
						port := info.Get("master_port")
						addr := fmt.Sprintf("%s:%s", ip, port)
						_, ok := addrFilterMap[addr]
						if !ok {
							//add the master to servers to monitor
							meta.Servers = append(meta.Servers, addr)
							addrFilterMap[addr] = true
						}
					}
				}
			}
		}
		if len(wsClients) > 0 {
			inner(meta)
			//notify all the clients
			for _, c := range wsClients {
				c.C <- "go"
			}
		}
		fmt.Printf("total: %d clients\n", len(wsClients))
		time.Sleep(1000 * time.Millisecond)
	}
}

func fetchRedisServerInfo(addr string, serverType string) *redis.RedisInfo {
	info, err := redis.FetchInfo(addr, "all")
	if err != nil {
		return nil
	}
	return info
}

func infoFilter(info *redis.RedisInfo) {
	filterMap := map[string]bool{}
	filterMap["uptime_in_seconds"] = true
	filterMap["user_memory_human"] = true
	filterMap["loading"] = true
	filterMap["instantaneous_ops_per_sec"] = true
	filterMap["role"] = true
	filterMap["aof_rewrite_in_progress"] = true

	//used by slave
	filterMap["master_repl_offset"] = true

	if info.Get("role") == "slave" && info.Get("master_host") != "127.0.0.1" {
		filterMap["master_host"] = true
		filterMap["master_port"] = true
		filterMap["master_link_status"] = true
		filterMap["master_sync_in_progress"] = true
		filterMap["slave_repl_offset"] = true
		filterMap["m_repl_offset"] = true
		//get its master info
		ip := info.Get("master_host")
		port := info.Get("master_port")
		addr := fmt.Sprintf("%s:%s", ip, port)
		minfo, ok := ServerInfoSnap[addr]
		if ok {
			(*info)["m_repl_offset"] = minfo.Get("master_repl_offset")
		}
	}
	for k, _ := range *info {
		_, ok := filterMap[k]
		if !ok {
			delete(*info, k)
		}
	}
}

func SlaveInfoFilter(snapshot InfoMap) InfoMap {
	slaveMap := map[string]*redis.RedisInfo{}
	for k, v := range snapshot {
		if v.Get("role") == "slave" {
			slaveMap[k] = v
		}
	}
	return slaveMap
}

func ClientRegiste(ws *websocket.Conn) *Client {
	client := &Client{
		C:  make(chan string),
		Ws: ws,
	}
	wsClients = append(wsClients, client)
	return client
}

func ClientUnreg(cli *Client) {
	for i, c := range wsClients {
		if c == cli {
			wsClients = append(wsClients[:i], wsClients[i+1:]...)
			break
		}
	}
}
