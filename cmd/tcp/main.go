//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package main

import (
	"fmt"
	"goim-client/internal"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	signalChan = make(chan os.Signal, 1)
)

const (
	tcpUrl = `192.168.32.124:3101`
)

func main() {
	Start(tcpUrl)

	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-signalChan
		fmt.Printf("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//svc.Close()
			fmt.Println(`service exit`)
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func Start(tcpAddrStr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddrStr)
	if err != nil {
		log.Printf("Resolve tcp addr failed: %v\n", err)
		return
	}

	// 向服务器拨号
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Dial to server failed: %v\n", err)
		return
	}

	// 向服务器发消息
	//go SendMsg(conn)
	conn.Write(internal.Auth())
	go SendMsg(conn)
	// 接收来自服务器端的广播消息
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		realBuf := buf[:length]
		p, err := internal.ParseMsg(realBuf)
		if err != nil {

			log.Println("读取错误:", err)
			return
		}
		internal.HandleTCPMsg(conn, p)

	}
}

// 向服务器端发消息
func SendMsg(conn net.Conn) {
	for {
		time.Sleep(time.Second * 10)
		conn.Write(internal.PackageMsg([]byte("hello")))
	}
}
