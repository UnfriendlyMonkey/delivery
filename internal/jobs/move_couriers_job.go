// Package jobs
package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"

	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &MoveCouriersJob{}

type MoveCouriersJob struct {
	moveCouriersHandler commands.MoveCouriersHandler
}

func NewMoveCouriersJob(moveCouriersHandler commands.MoveCouriersHandler) (cron.Job, error) {
	if moveCouriersHandler == nil {
		return nil, errs.NewValueIsInvalidError("moveCouriersHandler")
	}

	return &MoveCouriersJob{moveCouriersHandler: moveCouriersHandler}, nil
}

func (j *MoveCouriersJob) Run() {
	ctx := context.Background()

	command := commands.NewMoveCouriersCommand()
	err := j.moveCouriersHandler.Handle(ctx, command)
	if err != nil {
		log.Error(err)
	}
}
