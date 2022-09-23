//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package tcp

import (
	"fmt"
	"goim-client/api/grpc"
	"goim-client/internal"
	"io"
	"log"
	"net"
	"time"
)

func HandleTCPMsg(c *net.TCPConn, p *grpc.Proto) {
	switch p.Op {
	case 8:
		if !internal.HbOpen {
			go DocSyncTaskCronJobByTCP(c)
			internal.HbOpen = true
		}
	case 3:
		fmt.Println("receive: heartbeat")
	case 9:
		body := p.GetBody()

		if len(body) > 16 {
			fmt.Printf("messageReceived： ver=%v,body=%s\n", p.Ver, string(body[16:]))
		}
	case 5:
		body := p.GetBody()

		if len(body) > 16 {
			fmt.Printf("messageReceivedByWS： ver=%v,body=%s\n", p.Ver, string(body[16:]))
		}
	case 1000:
		fmt.Printf("messageReceivedByWS： ver=%v,body=%s\n", p.Ver, string(p.Body))
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}

func DocSyncTaskCronJobByTCP(c *net.TCPConn) {
	for {
		if internal.HeartbeatInterval.After(time.Now()) {
			HeartbeatTCP(c)
			internal.HeartbeatInterval.Add(internal.HearBeatSpec)
			time.Sleep(time.Second * 10)
		}
	}
}

func HeartbeatTCP(c *net.TCPConn) {
	msgBuf := internal.Heartbeat()
	c.Write(msgBuf)
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
	conn.Write(internal.Auth())
	go SendMsg(conn)
	// 接收来自服务器端的广播消息
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		realBuf := buf[:length]
		p, err := internal.ParseMsg(realBuf)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("读取错误:", err)
			return
		}
		HandleTCPMsg(conn, p)

	}
}

// 向服务器端发消息
func SendMsg(conn net.Conn) {
	for {
		time.Sleep(time.Second * 10)
		conn.Write(internal.PackageMsg([]byte(fmt.Sprintf("hello,time is %s", time.Now().Format("15:04:05")))))
	}
}
