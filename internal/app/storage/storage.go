package storage

import (
	"context"
	"errors"
	"time"
)

var ErrInsufficientBalance = errors.New("insufficient balance for buy")
var ErrDuplicateCompany = errors.New("duplicate company")
var ErrCompanyNotFound = errors.New("company not found")

type Company struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name"`
	Code string `json:"code,omitempty"`
}

type Contract struct {
	Trader         *Company  `json:"trader"`
	Buyer          *Company  `json:"buyer"`
	Number         string    `json:"number"`
	SignDate       time.Time `json:"sign_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	Sum            float64   `json:"sum"`
}

type Buy struct {
	Contract  *Contract `json:"contract"`
	Timestamp time.Time `json:"timestamp"`
	Sum       float64   `json:"sum"`
}

type Repository interface {
	CreateCompany(ctx context.Context, company *Company) (int64, error)
	ReadCompany(ctx context.Context, id int64) (*Company, error)
	ReadCompanyList(ctx context.Context) ([]Company, error)
	UpdateCompany(ctx context.Context, id int64, company *Company) error
	DeleteCompany(ctx context.Context, id int64) error

	CreateContract(ctx context.Context, contract *Contract) (int64, error)
	ReadContract(ctx context.Context, id int64) (*Contract, error)
	ReadContractList(ctx context.Context) ([]Contract, error)
	UpdateContract(ctx context.Context, id int64, contract *Contract) error
	DeleteContract(ctx context.Context, id int64) error

	CreateBuy(ctx context.Context, company *Buy) (int64, error)
}
