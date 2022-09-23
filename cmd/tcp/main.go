//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package main

import (
	"fmt"
	"goim-client/internal/tcp"
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
	tcp.Start(tcpUrl)

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
