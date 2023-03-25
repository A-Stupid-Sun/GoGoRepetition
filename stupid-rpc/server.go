package stupid_rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"stupid-rpc/codec"
	"sync"
)

const MagicNumber = 0x1234

type Option struct {
	MagicNumber int        // MagicNumber marks this's a geerpc request
	CodecType   codec.Type // client may choose different Codec to encode body
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

/*
	RPC报文结构如下所示：

| Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
| <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|

单个报文可能传输多段数据：
| Option | Header1 | Body1 | Header2 | Body2 | ...
*/

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}
func Accept(lis net.Listener) { DefaultServer.Accept(lis) } //用来服务器注册和启动
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer conn.Close()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc Server connetion error", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("Invalid MagicNumber %x shows it's not rpc connection", opt.MagicNumber)
		return
	}
	codecFunc := codec.NewCodecFuncMap[opt.CodecType]
	if codecFunc == nil {
		log.Printf("Lack of this codec type : %s for connection", opt.CodecType)
		return
	}
	server.serveCodec(codecFunc(conn))

}

type request struct {
	header       *codec.Header
	argv, replyv reflect.Value
}

func (server *Server) readRequestHandler(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header                        //代表header结构
	if err := cc.ReadHeader(&h); err != nil { //使用特定编码方式读取头文件，无法解码则返回错误。否则修改h内柔
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHandler(cc)
	if err != nil {
		return nil, err
	}
	req := &request{header: h}
	// TODO: now we don't know the type of request argv
	// day 1, just suppose it's string
	req.argv = reflect.New(reflect.TypeOf("")) //不懂

	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	} //读取参数

	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
	// 串行发送回复报文
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO, should call registered rpc methods to get the right replyv
	// day 1, just print argv and send a hello message
	defer wg.Done() //使用信号量 wg--
	log.Println(req.header, req.argv.Elem())

	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.header.Seqnum))

	server.sendResponse(cc, req.header, req.replyv.Interface(), sending)
}

func (server *Server) serveCodec(codec codec.Codec) {
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)  // wait until all request are handled
	for {
		req, err := server.readRequest(codec) // 多线程读取请求由sever实现，所以调用server.readRequest
		if err != nil {
			if req == nil {
				break //req==nil代表读取失败，读取过程中失败不可恢复，直接退出
			}
			req.header.Error = err.Error()
			server.sendResponse(codec, req.header, struct{}{}, sending)
			continue
		} // header读取失败，直接错误处理
		wg.Add(1)
		go server.handleRequest(codec, req, sending, wg)
	}
	wg.Wait() //等wg==0
	codec.Close()
}
