package main

import (
	"hy.juck.com/go-publisher-server/cron"
	"hy.juck.com/go-publisher-server/router"
)

func main() {
	go func() {
		c := cron.InitCron()
		c.Run()
	}()
	// 开启路由监听
	router.NewHttpRouter()
}
