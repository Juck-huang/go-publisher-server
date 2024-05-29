package cron

import (
	"github.com/robfig/cron/v3"
	"hy.juck.com/go-publisher-server/config"
)

var (
	G = config.G
)

type RunTask struct {
	C *cron.Cron
}

func InitCron() *RunTask {
	c := cron.New()
	return &RunTask{
		C: c,
	}
}

func (o *RunTask) addFunc(t string, c func()) error {
	_, err := o.C.AddFunc(t, c)
	if err != nil {
		return err
	}
	return nil
}

// 开始执行定时任务
func (o *RunTask) start() {
	o.C.Start()
}

func (o *RunTask) Run() {
	o.addFunc("* * * * *", DbExportFileClean)
	o.start()
}
