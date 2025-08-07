package incident

import (
	"incident-tracker/config"
	"incident-tracker/models"
	"incident-tracker/repository"
	"time"

	"github.com/vladopajic/go-actor/actor"
)

type poller struct {
	appCtx *config.ApplicationContext
	ticker *time.Ticker
	outC   actor.MailboxSender[*models.Incident]
}

func newPoller(appCtx *config.ApplicationContext, ticker *time.Ticker, outC actor.MailboxSender[*models.Incident]) *poller {
	return &poller{
		appCtx: appCtx,
		ticker: ticker,
		outC:   outC,
	}
}

func (p *poller) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		p.appCtx.Logger.Info("[Incident Poller] is shutting down")
		return actor.WorkerEnd
	case <-p.ticker.C:
		repo := repository.NewIncidentRepository(p.appCtx.DB)
		incidents, err := repo.GetIncidentByRequestStatusAndAIModel(models.StatusPending, "", 100, 0)
		if err != nil {
			p.appCtx.Logger.Sugar().Errorf("Error getting pending incidents: %v", err)
			return actor.WorkerContinue
		}

		for _, incident := range incidents {
			p.outC.Send(ctx, &incident)
		}

		return actor.WorkerContinue
	}
}
