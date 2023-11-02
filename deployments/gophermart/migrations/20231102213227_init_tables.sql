-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create type order_status as enum ('new', 'processed', 'processing', 'invalid');

create table users (id uuid default gen_random_uuid() not null primary key,
user_name varchar(128) not null unique,
user_password varchar(128) not null,
created_at timestamp default now(),
updated_at timestamp default now(),
deleted_at timestamp);

create table orders (
id uuid default gen_random_uuid() not null primary key,
order_number int not null,
status order_status,
accrual int not null,
created_at timestamp default now(),
updated_at timestamp default now(),
deleted_at timestamp,
user_created uuid,
user_deleted uuid,
foreign key (user_created) references users(id),
foreign key (user_deleted) references users(id)
);

create table balance (
id uuid default gen_random_uuid() not null primary key,
user_id uuid not null,
current_balance int not null,
created_at timestamp default now(),
updated_at timestamp default now(),
deleted_at timestamp,
foreign key (user_id) references users(id)
);

create table balance_operations (
id uuid default gen_random_uuid() not null primary key,
order_id uuid,
balance_id uuid,
sum_operation int,
created_at timestamp default now(),
updated_at timestamp default now(),
deleted_at timestamp,
foreign key (order_id) references orders(id),
foreign key (balance_id) references balance(id)
);

-- +goose Down
-- +goose StatementBegin
drop table balance_operations;

drop table balance;

drop table orders;

drop table users

drop type order_status
-- +goose StatementEnd
