package posgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const defaultDBTimeout = 3 * time.Second

type PostgresRepository struct {
	dbPool *pgxpool.Pool
	ctx    context.Context
	dsn    string
}

func NewPosgresRepository(ctx context.Context, dsn string) *PostgresRepository {
	return &PostgresRepository{
		ctx: ctx,
		dsn: dsn,
	}
}

func (p *PostgresRepository) Open() error {
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

func (p *PostgresRepository) Close() {
	p.dbPool.Close()
}

func (p *PostgresRepository) DBPool() *pgxpool.Pool {
	return p.dbPool
}
