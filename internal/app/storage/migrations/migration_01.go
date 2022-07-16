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

create unique index if not exists contract_id_uindex
    on contract(id);

create table if not exists buy (
    id          bigserial
        constraint buy_id_pk
            primary key,
    contract_id bigint                   not null
        constraint buy_contract_id_fk
            references contract
            on delete cascade,
    timespamp   timestamp with time zone not null,
    sum         numeric(10, 2)           not null
);

create table if not exists revision (
    version bigserial
        constraint revision_version_pk
            primary key
);


INSERT INTO revision VALUES(1);  
`)
	return err
}
