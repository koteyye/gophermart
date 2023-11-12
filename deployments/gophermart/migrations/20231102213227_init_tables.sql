-- +goose Up
-- +goose StatementBegin
create type order_status as enum ('new', 'processed', 'processing', 'invalid');

create type balance_operations_state as enum ('run', 'done', 'error');

create table if not exists users (
	id uuid default gen_random_uuid() not null primary key,
	user_name varchar(128) not null unique,
	hashed_password varchar(128) not null,
	created_at timestamp default now(),
	updated_at timestamp default now(),
	deleted_at timestamp
);

create table if not exists orders (
	id uuid default gen_random_uuid() not null primary key,
	order_number varchar not null unique,
	status order_status default 'new',
	accrual int not null default 0,
	created_at timestamp default now(),
	updated_at timestamp default now(),
	deleted_at timestamp,
	user_created uuid,
	foreign key (user_created) references users(id)
);

create table if not exists balance (
	id uuid default gen_random_uuid() not null primary key,
	user_id uuid not null,
	current_balance int not null default 0,
	withdrawn int not null default 0,
	created_at timestamp default now(),
	updated_at timestamp default now(),
	deleted_at timestamp,
	foreign key (user_id) references users(id)
);

create table if not exists balance_operations (
	id uuid default gen_random_uuid() not null primary key,
	order_id uuid,
	balance_id uuid,
	sum_operation int,
	operation_state balance_operations_state default 'run',
	created_at timestamp default now(),
	updated_at timestamp default now(),
	deleted_at timestamp,
	foreign key (order_id) references orders(id),
	foreign key (balance_id) references balance(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists balance_operations;

drop table if exists balance;

drop table if exists orders;

drop table if exists users;

drop type order_status;

drop type balance_operations_state;
-- +goose StatementEnd
