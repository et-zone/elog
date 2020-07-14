package glogs

import (
	cron "gopkg.in/robfig/cron.v2"
)

type Cron struct {
	cron *cron.Cron
}

var cro *Cron

func newCron() *Cron {
	if cro != nil {
		return cro
	}
	cr := &Cron{cron: cron.New()}
	cro = cr
	return cro
}

func (this *Cron) addFunc(spec string, cmd func()) (cron.EntryID, error) {
	return this.cron.AddFunc(spec, cmd)
}

func (this *Cron) start() {
	this.cron.Start()

}
func (this *Cron) stop() {
	this.cron.Stop()
}
