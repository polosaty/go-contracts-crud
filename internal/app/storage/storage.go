package storage

import (
	"context"
	"errors"
	"time"
)

var ErrInsufficientContractSum = errors.New("insufficient balance for buy")
var ErrDuplicateCompany = errors.New("duplicate company")
var ErrCompanyNotFound = errors.New("company not found")
var ErrCompanyCantBeDeleted = errors.New("company cannot be deleted")
var ErrDuplicateContract = errors.New("duplicate contract")
var ErrContractNotFound = errors.New("contract not found")
var ErrContractCantBeDeleted = errors.New("contract cannot be deleted")
var ErrDuplicateBuy = errors.New("duplicate buy")

type Company struct {
	ID   *int64 `json:"id,omitempty"`
	Name string `json:"name"`
	Code string `json:"code,omitempty"`
}

type Contract struct {
	ID             *int64    `json:"id,omitempty"`
	TraderID       int64     `json:"trader_id"`
	BuyerID        int64     `json:"buyer_id"`
	Number         string    `json:"number"`
	SignDate       time.Time `json:"sign_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	Sum            float64   `json:"sum"`
}

type Buy struct {
	ID         *int64    `json:"id,omitempty"`
	ContractID int64     `json:"contract_id"`
	Timestamp  time.Time `json:"timestamp"`
	Sum        float64   `json:"sum"`
}

type Repository interface {
	CreateCompany(ctx context.Context, company *Company) (int64, error)
	ReadCompany(ctx context.Context, id int64) (*Company, error)
	ReadCompanyList(ctx context.Context, limit uint, offset uint64) ([]Company, error)
	UpdateCompany(ctx context.Context, id int64, company *Company) error
	DeleteCompany(ctx context.Context, id int64) error

	CreateContract(ctx context.Context, contract *Contract) (int64, error)
	ReadContract(ctx context.Context, id int64) (*Contract, error)
	ReadContractList(ctx context.Context, limit uint, offset uint64) ([]Contract, error)
	UpdateContract(ctx context.Context, id int64, contract *Contract) error
	DeleteContract(ctx context.Context, id int64) error

	CreateBuy(ctx context.Context, company *Buy) (int64, error)
}
