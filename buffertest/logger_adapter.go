package one

import (
	"context"
	"go.uber.org/zap"
)

// LoggerAdapter wraps the usecase interface
// with a logging adapter which can be swapped out
type LoggerAdapter struct {
	Logger  *zap.Logger
	Usecase OneService
}

func (a *LoggerAdapter) logErr(err error) {
	if err != nil {
		a.Logger.Error(err.Error())
	}
}

// Processes an event to detemine if action is required
func (a *LoggerAdapter) CreateMetric(ctx context.Context, data string) error {
	defer a.Logger.Sync()
	
	err := a.Usecase.CreateCPUMetrics(ctx, data)
	a.logErr(err)
	return err
}

func (a *LoggerAdapter) CloseRepository() error {
	defer a.Logger.Sync()
	a.Logger.Info("closing Repository")
	err := a.Usecase.CloseRepository()
	a.logErr(err)
	return err
}
