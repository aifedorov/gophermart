package domain

import (
	"context"
	"time"

	"github.com/aifedorov/gophermart/internal/client/accrual"
	repository "github.com/aifedorov/gophermart/internal/order/repository/db"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"go.uber.org/zap"
)

type Poller interface {
	StartPollingWithOrderNumber(number string) error
}

type poller struct {
	ctx           context.Context
	repo          repository.Repository
	accrualClient accrual.HTTPClient
}

func NewPoller(ctx context.Context, repo repository.Repository, accrualClient accrual.HTTPClient) Poller {
	return &poller{
		ctx:           ctx,
		repo:          repo,
		accrualClient: accrualClient,
	}
}

var pollingTimeouts = []time.Duration{
	200 * time.Millisecond,
	400 * time.Millisecond,
	600 * time.Millisecond,
	1000 * time.Millisecond,
	1600 * time.Millisecond,
	2600 * time.Millisecond,
	4200 * time.Millisecond,
}

func (p *poller) StartPollingWithOrderNumber(number string) error {
	for _, timeout := range pollingTimeouts {
		select {
		case <-p.ctx.Done():
			logger.Log.Debug("poller: context was cancelled")
			return p.ctx.Err()
		default:
		}

		logger.Log.Debug("poller: star polling", zap.String("orderNumber", number))
		res, ok, err := p.accrualClient.GetAccrualByOrderNumber(number)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		switch res.Status {
		case accrual.StatusRegistered, accrual.StatusProcessing:
			logger.Log.Debug("poller: status isn't terminated", zap.Any("order status", res.Status))
		case accrual.StatusInvalid:
			err := p.repo.UpdateOrderStatusByNumber(number, repository.OrderstatusINVALID, nil)
			if err != nil {
				return err
			}
			logger.Log.Debug("poller: finish polling", zap.String("orderNumber", number), zap.Any("order status", res.Status))
			return nil
		case accrual.StatusProcessed:
			err := p.repo.UpdateOrderStatusByNumber(number, repository.OrderstatusPROCESSED, res.Amount)
			if err != nil {
				return err
			}
			logger.Log.Debug("poller: finish polling", zap.String("orderNumber", number), zap.Any("order status", res.Status))
			return nil
		}
		time.Sleep(timeout)
	}
	logger.Log.Debug("poller: finish unsuccessful polling", zap.String("orderNumber", number))
	return nil
}
