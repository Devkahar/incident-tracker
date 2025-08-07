package workers

import (
	"incident-tracker/config"
	"incident-tracker/workers/incident"

	"github.com/vladopajic/go-actor/actor"
)

func CreateActors(applicationContext *config.ApplicationContext) actor.Actor {
	actors := incident.NewActor(applicationContext)
	return actor.Combine(actors...).Build()
}
