package posgre

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultDBTimeout = 3 * time.Second

type PostgreRepository struct {
	dbPool *pgxpool.Pool
	ctx    context.Context
	dsn    string
}

func NewPosgresRepository(ctx context.Context, dsn string) *PostgreRepository {
	return &PostgreRepository{
		ctx: ctx,
		dsn: dsn,
	}
}

func (p *PostgreRepository) Open() error {
	ctx, cancel := context.WithTimeout(p.ctx, defaultDBTimeout)
	defer cancel()

	dbpool, err := pgxpool.New(context.Background(), p.dsn)
	if err != nil {
		return err
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return err
	}

	p.dbPool = dbpool
	return nil
}

func (p *PostgreRepository) Close() {
	p.dbPool.Close()
}

func (p *PostgreRepository) DBPool() *pgxpool.Pool {
	return p.dbPool
}
