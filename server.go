package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"LiveGoLib/httpflv"
	"LiveGoLib/proxy"
	"LiveGoLib/rtmp/core"
)

type ClientNetConn struct {
	conn net.Conn
}

func (cnc *ClientNetConn) GetNetConn() net.Conn {
	return cnc.conn
}

func NewClientNetConn(c net.Conn) *ClientNetConn {
	ret := &ClientNetConn{
		conn: c,
	}
	return ret
}

func StartHttpFlvServer(httpCli httpflv.ClientInfo, wg *sync.WaitGroup) {
	http.HandleFunc("/video/drop.flv", httpCli.Handleflv)
	go func() {
		defer wg.Done()
		defer httpCli.Prox.SetClientCnt(false)
		log.Fatal(http.ListenAndServe(":8000", nil))
	}()
}

func main() {
	//http.HandleFunc("/users", acceptClientInfo)
	prox := proxy.NewProxyPush()
	var httpCli httpflv.ClientInfo
	httpCli.Prox = prox
	httpWaitGroup := &sync.WaitGroup{}
	httpWaitGroup.Add(1)
	StartHttpFlvServer(httpCli, httpWaitGroup)
	listener, err := net.Listen("tcp", "127.0.0.1:1935")
	if err != nil {
		return
	}
	fmt.Println("toooooo")
	defer listener.Close()
	conn, err := listener.Accept()
	clientConn := NewClientNetConn(conn)
	if err != nil {
		return
	}
	core.HandShakeClientPkg(clientConn)
	core.StartHandleConn(clientConn, prox)
	httpWaitGroup.Wait()
}
