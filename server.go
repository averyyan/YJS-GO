package YJS_GO

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		// if r.Method != "GET" {
		//	fmt.Println("method is not GET")
		//	return false
		// }
		// if r.URL.Path != "/ws" {
		//	fmt.Println("path error")
		//	return false
		// }
		return true
	},
}

func serve() {
	go StartHttpServer()
	c := make(chan os.Signal, 1)
	ticker := time.NewTicker(50 * time.Second)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		select {
		case s := <-c:
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
				time.Sleep(time.Second)
				return
			case syscall.SIGHUP:
			// TODO reload
			default:
				return
			}
		case <-ticker.C:
			// service.NewOrGetRedisService().RegistHost(time.Minute)
		}

	}
}
func StartHttpServer() {
	err := http.ListenAndServe("9001", nil)
	if err != nil {
		panic("error|ListenAndServe:%v ")
	}
	fmt.Printf("start http server succeed")
}
