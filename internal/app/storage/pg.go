package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/polosaty/go-contracts-crud/internal/app/storage/migrations"
	"log"
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

func NewStorageFromPool(pool dbInterface) *PG {
	return &PG{db: pool}
}

func (s *PG) CreateCompany(ctx context.Context, company *Company) (id int64, err error) {
	err = s.db.QueryRow(ctx,
		`INSERT INTO "company" (name, code) VALUES($1, $2) `+
			` RETURNING id`, company.Name, company.Code).
		Scan(&id)

	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pgerrcode.IsIntegrityConstraintViolation(pge.SQLState()) {
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

func (s *PG) ReadCompanyList(ctx context.Context, limit uint, offset uint64) ([]Company, error) {
	if limit > 100 || limit == 0 {
		limit = 100
	}
	companies := make([]Company, 0, limit)
	rows, err := s.db.Query(ctx,
		`SELECT id, name, code FROM company LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}

	for rows.Next() {
		var v Company
		err = rows.Scan(&v.ID, &v.Name, &v.Code)
		if err != nil {
			return nil, fmt.Errorf("cant parse row from select companies: %w", err)
		}
		companies = append(companies, v)
	}

	return companies, nil
}

func (s *PG) UpdateCompany(ctx context.Context, id int64, company *Company) error {
	tag, err := s.db.Exec(ctx,
		`UPDATE company SET name = $1, code = $2 WHERE id = $3`, company.Name, company.Code, id)
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pgerrcode.IsIntegrityConstraintViolation(pge.SQLState()) {
			return ErrDuplicateCompany
		}
		return fmt.Errorf("update company error: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrCompanyNotFound
	}

	return nil
}

func (s *PG) DeleteCompany(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx,
		`DELETE FROM company WHERE id = $1`, id)
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pgerrcode.IsIntegrityConstraintViolation(pge.SQLState()) {
			return ErrCompanyCantBeDeleted
		}
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrCompanyNotFound
	}
	return nil
}

// CreateContract - создает контракт с 0-ой суммой
// возвращает id созданного контракта, либо ошибку если таковой создать не удалось
func (s *PG) CreateContract(ctx context.Context, contract *Contract) (id int64, err error) {
	if contract.Sum <= 0 {
		return 0, ErrInsufficientContractSum
	}
	err = s.db.QueryRow(ctx,
		`INSERT INTO contract (trader_id, buyer_id, number, sign_date, expiration_date, "sum") `+
			` VALUES($1, $2, $3, $4, $5, $6) RETURNING id`,
		contract.TraderID,
		contract.BuyerID,
		contract.Number,
		contract.SignDate,
		contract.ExpirationDate,
		contract.Sum,
	).
		Scan(&id)

	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pgerrcode.IsIntegrityConstraintViolation(pge.SQLState()) {
			// contract already exists
			// Handle  duplicate key value violates
			return 0, ErrDuplicateContract
		}
		return 0, fmt.Errorf("create contract error: %w", err)
	}

	return
}

func (s *PG) ReadContract(ctx context.Context, id int64) (*Contract, error) {
	var contract Contract
	err := s.db.QueryRow(ctx,
		`SELECT trader_id, buyer_id, number, sign_date, expiration_date, "sum" `+
			`FROM  contract WHERE id = $1`, id).
		Scan(
			&contract.TraderID,
			&contract.BuyerID,
			&contract.Number,
			&contract.SignDate,
			&contract.ExpirationDate,
			&contract.Sum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrContractNotFound
		}
		return nil, err
	}
	return &contract, nil
}

func (s *PG) ReadContractList(ctx context.Context, limit uint, offset uint64) ([]Contract, error) {
	if limit > 100 || limit == 0 {
		limit = 100
	}
	contracts := make([]Contract, 0, limit)
	rows, err := s.db.Query(ctx,
		`SELECT id, trader_id, buyer_id, number, sign_date, expiration_date, "sum" `+
			`FROM contract LIMIT $1 OFFSET $2`, limit, offset)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrContractNotFound
		}
		return nil, err
	}

	for rows.Next() {
		var v Contract
		err = rows.Scan(&v.ID, &v.TraderID, &v.BuyerID, &v.Number, &v.SignDate, &v.ExpirationDate, &v.Sum)
		if err != nil {
			return nil, fmt.Errorf("cant parse row from select contracts: %w", err)
		}
		contracts = append(contracts, v)
	}

	return contracts, nil
}

// UpdateContract
//  обновить контракт по запросу извне
func (s *PG) UpdateContract(ctx context.Context, id int64, contract *Contract) error {
	if contract.Sum <= 0 {
		return ErrInsufficientContractSum
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `SELECT id FROM contract WHERE id = $1 FOR UPDATE`, id); err != nil {
		return fmt.Errorf("lock contract for update: %w", err)
	}
	var buysSum float64
	err = tx.QueryRow(ctx, `SELECT coalesce(sum("sum"), 0) FROM buy WHERE contract_id = $1`, id).Scan(&buysSum)
	if err != nil {
		return fmt.Errorf("read amount of purchases: %w", err)
	}

	if buysSum > contract.Sum {
		return ErrInsufficientContractSum
	}

	tag, err := tx.Exec(ctx,
		`UPDATE contract `+
			` SET buyer_id = $1, trader_id = $2, number = $3, sign_date = $4, expiration_date = $5, sum = $6 `+
			` WHERE id = $7`,
		contract.BuyerID, contract.TraderID, contract.Number, contract.SignDate, contract.ExpirationDate, contract.Sum, id)

	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pgerrcode.IsIntegrityConstraintViolation(pge.SQLState()) {
			return ErrDuplicateContract
		}
		return fmt.Errorf("update contract error: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrContractNotFound
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("cant commit tx %w", err)
	}

	return nil
}

func (s *PG) DeleteContract(ctx context.Context, id int64) error {
	tag, err := s.db.Exec(ctx,
		`DELETE FROM contract WHERE id = $1`, id)
	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pgerrcode.IsIntegrityConstraintViolation(pge.SQLState()) {
			return ErrContractCantBeDeleted
		}
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrContractNotFound
	}
	return nil
}

// CreateBuy
//  создать покупку по запросу извне
//  для того чтобы избежать превышения суммы контракта,
//  на момент создания покупки контракт блокируется
func (s *PG) CreateBuy(ctx context.Context, buy *Buy) (int64, error) {
	tx, err := s.db.Begin(ctx)
	var id int64
	if err != nil {
		return 0, fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var contractSum float64
	err = tx.QueryRow(ctx, `SELECT "sum" FROM contract WHERE id = $1 FOR UPDATE`, buy.ContractID).
		Scan(&contractSum)
	if err != nil {
		return 0, fmt.Errorf("lock contract for create buy: %w", err)
	}
	var buysSum float64
	err = tx.QueryRow(ctx, `SELECT coalesce(sum("sum"), 0) FROM buy WHERE contract_id = $1`, buy.ContractID).
		Scan(&buysSum)
	if err != nil {
		return 0, fmt.Errorf("read amount of purchases: %w", err)
	}

	if buysSum+buy.Sum > contractSum {
		return 0, ErrInsufficientContractSum
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO buy(contract_id, "timestamp", "sum")  VALUES($1, $2, $3) RETURNING id`,
		buy.ContractID, buy.Timestamp, buy.Sum).Scan(&id)

	var pge *pgconn.PgError
	if errors.As(err, &pge) {
		if pgerrcode.IsIntegrityConstraintViolation(pge.SQLState()) {
			log.Println("create buy error: ", err)
			return 0, ErrDuplicateBuy
		}
		return 0, fmt.Errorf("create buy error: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("cant commit tx %w", err)
	}

	return id, nil
}
