package network

import (
	"fmt"
	"github.com/prometheus/common/log"
	"net"
	"time"
)

func InitTcp(address string, headLen, tagLen int32) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Error(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error(err)
	}
	go acceptTcp(listener, headLen, tagLen)
}

func acceptTcp(listener *net.TCPListener, headLen, tagLen int32) {
	for {
		var (
			conn *net.TCPConn
			err  error
		)
		if conn, err = listener.AcceptTCP(); err != nil {
			log.Error(err)
			return
		}
		if err = conn.SetKeepAlive(false); err != nil {
			log.Error(err)
			return
		}
		if err = conn.SetReadBuffer(1024); err != nil {
			log.Error(err)
			return
		}
		if err = conn.SetWriteBuffer(1024); err != nil {
			log.Error(err)
			return
		}
		go serveTCP(conn, headLen, tagLen)
	}
}

func serveTCP(conn *net.TCPConn, headLen, tagLen int32) {
	client := NewTCPClient(conn, headLen, tagLen)
	_ = client.conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
	go func() {
		for {
			tag, data, err := client.Read()
			if err != nil {
				log.Error(err)
				_ = client.Close()
				return
			}
			_ = client.conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
			fmt.Println("tag:", tag)
			fmt.Println(string(data))
		}
	}()
}
