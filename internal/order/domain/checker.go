package domain

import (
	"context"
	"database/sql"
	"errors"
	"time"

	repository "github.com/aifedorov/gophermart/internal/order/repository/db"
	"github.com/aifedorov/gophermart/internal/pkg/logger"
	"go.uber.org/zap"
)

type Checker interface {
	Run() error
}

type checker struct {
	ctx    context.Context
	repo   repository.Repository
	poller Poller
}

func NewChecker(ctx context.Context, repo repository.Repository, poller Poller) Checker {
	return &checker{
		ctx:    ctx,
		repo:   repo,
		poller: poller,
	}
}

func (c *checker) Run() error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-c.ctx.Done():
			logger.Log.Debug("poller: context was cancelled")
			return c.ctx.Err()
		case <-ticker.C:
			err := c.processNewOrders()
			if err != nil {
				logger.Log.Error("app: error processing new orders", zap.Error(err))
			}
		}
	}
}

func (c *checker) processNewOrders() error {
	order, err := c.repo.GetNewTopUpOrder()
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	go func() {
		err := c.poller.StartPollingWithOrderNumber(order.Number)
		if err != nil {
			logger.Log.Error("poller failed to start orders", zap.Error(err))
		}
	}()
	return nil
}
