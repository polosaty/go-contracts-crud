package migrations

import (
	"context"
)

func migration01(ctx context.Context, db DBInterface) error {
	_, err := db.Exec(
		ctx,
		`
create table if not exists company (
    id   bigserial
        constraint company_id_pk
            primary key,
    name varchar(512) not null,
    code varchar(255),
    constraint company_name_code_pk
        unique (name, code)
);

create unique index if not exists company_id_uindex
    on company(id);

create table if not exists contract (
    id              bigserial
        constraint contract_id_pk
            primary key,
    trader_id       bigint
        constraint contract_trader_id_fk
            references company
            on delete restrict,
    buyer_id        bigint
        constraint contract_buyer_id_fk
            references company
            on delete restrict,
    number          varchar(255),
    sign_date       timestamp with time zone,
    expiration_date timestamp with time zone,
    sum             numeric(10, 2) not null
);

create unique index contract_trader_id_buyer_id_number_uindex
    on contract(trader_id, buyer_id, number);

create table if not exists buy (
    id          bigserial
        constraint buy_id_pk
            primary key,
    contract_id bigint                   not null
        constraint buy_contract_id_fk
            references contract
            on delete cascade,
    timestamp   timestamp with time zone not null,
    sum         numeric(10, 2)           not null
);

create unique index if not exists buy_contract_id_timestamp_uindex
    on buy(contract_id asc, timestamp desc);


INSERT INTO revision VALUES(1);  
`)
	return err
}
