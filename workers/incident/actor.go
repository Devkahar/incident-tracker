package incident

import (
	"incident-tracker/config"
	"incident-tracker/models"
	"time"

	"github.com/vladopajic/go-actor/actor"
)

const (
	numPoller    = 1
	numProcessor = 10
	pollInterval = 10 * time.Second
)

func NewActor(appCtx *config.ApplicationContext) []actor.Actor {
	mailbox := actor.NewMailbox[*models.Incident]()
	var a []actor.Actor
	a = append(a, mailbox)
	for i := 0; i < numProcessor; i++ {
		a = append(a, actor.New(newProcessor(mailbox, appCtx)))
	}

	a = append(a, actor.New(newPoller(appCtx, time.NewTicker(pollInterval), mailbox)))
	return a
}
