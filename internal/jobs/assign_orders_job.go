package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"

	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrdersJob{}

type AssignOrdersJob struct {
	assignOrderHandler commands.AssignOrderHandler
}

func NewAssignOrdersJob(assignOrderHandler commands.AssignOrderHandler) (cron.Job, error) {
	if assignOrderHandler == nil {
		return nil, errs.NewValueIsInvalidError("assignOrderHandler")
	}
	return &AssignOrdersJob{assignOrderHandler: assignOrderHandler}, nil
}

func (j *AssignOrdersJob) Run() {
	ctx := context.Background()
	command := commands.NewAssignOrderCommand()
	err := j.assignOrderHandler.Handle(ctx, command)
	if err != nil {
		log.Error(err)
	}
}
