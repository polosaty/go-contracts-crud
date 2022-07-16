package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/polosaty/go-contracts-crud/internal/app/storage/migrations"
)

type PG struct {
	db dbInterface
}

var _ Repository = (*PG)(nil)

type dbInterface interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Ping(context.Context) error
	Close()
}

func NewStoragePG(uri string) (*PG, error) {
	ctx := context.Background()
	conf, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: parse dsn problem (dsn=%v): %w", uri, err)
	}

	conn, err := pgxpool.ConnectConfig(ctx, conf)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database(uri=%v): %w", uri, err)
	}

	repo := &PG{
		db: conn,
	}

	err = migrations.Migrate(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("can't apply migrations: %w", err)
	}

	return repo, nil
}

func (s *PG) CreateCompany(ctx context.Context, company *Company) (id int64, err error) {
	err = s.db.QueryRow(ctx,
		`INSERT INTO "company" (name, code) VALUES($1, $2)
			RETURNING id`, company.Name, company.Code).
		Scan(&id)

	//https://github.com/jackc/pgconn/issues/15#issuecomment-867082415
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pge.SQLState() == "23505" {
			// company already exists
			// Handle  duplicate key value violates
			return 0, ErrDuplicateCompany
		}
		return 0, fmt.Errorf("create company error: %w", err)
	}

	return
}

func (s *PG) ReadCompany(ctx context.Context, id int64) (*Company, error) {
	var company Company
	err := s.db.QueryRow(ctx,
		`SELECT name, code FROM  company WHERE id = $1`, id).
		Scan(&company.Name, &company.Code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}
	return &company, nil
}

func (s *PG) ReadCompanyList(ctx context.Context) ([]Company, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PG) UpdateCompany(ctx context.Context, id int64, company *Company) error {
	//TODO implement me
	panic("implement me")
}

func (s *PG) DeleteCompany(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *PG) CreateContract(ctx context.Context, contract *Contract) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PG) ReadContract(ctx context.Context, id int64) (*Contract, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PG) ReadContractList(ctx context.Context) ([]Contract, error) {
	//TODO implement me
	panic("implement me")
}

func (s *PG) UpdateContract(ctx context.Context, id int64, contract *Contract) error {
	//TODO implement me
	panic("implement me")
}

func (s *PG) DeleteContract(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *PG) CreateBuy(ctx context.Context, company *Buy) (int64, error) {
	//TODO implement me
	panic("implement me")
}
