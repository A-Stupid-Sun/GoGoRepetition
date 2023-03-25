package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"stupid-rpc"
	"stupid-rpc/codec"
	"time"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("start server fail!", err)
		return
	}
	log.Printf("start server on : %s", l.Addr())
	addr <- l.Addr().String()
	stupid_rpc.Accept(l)
}
func main() {
	addr := make(chan string)
	go startServer(addr)

	//simple rpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer conn.Close()
	time.Sleep(time.Second)                                // server需要accpet启动
	json.NewEncoder(conn).Encode(stupid_rpc.DefaultOption) //设置
	cc := codec.NewGobCodec(conn)
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Stupid-Service",
			Seqnum:        uint64(i),
		}
		data := fmt.Sprintf("Stupid Service Call %x", i)
		cc.Write(h, data)
		cc.ReadHeader(h)
		var reply string
		cc.ReadBody(&reply)
		log.Println("reply data:", reply)
	}
}
