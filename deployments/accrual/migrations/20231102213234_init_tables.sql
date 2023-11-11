-- +goose Up
-- +goose StatementBegin
create type order_status as enum ('registered', 'invalid', 'processing', 'processed');
comment on type order_status is 'возможные статусы заказа';

create type rewards as enum ('percent', 'natural');
comment on type rewards is 'возможные типы вознагрождений';

create table orders (
                        id uuid not null default gen_random_uuid() unique,
                        order_number varchar not null,
                        status order_status default 'registered',
                        accrual int not null default 0,
                        created_at timestamp not null default now(),
                        updated_at timestamp not null default now(),
                        deleted_at timestamp
);
comment on table orders is 'зарегистрированные заказы на расчет вознагрождений';
comment on column orders.order_number is 'номер заказа';
comment on column orders.status is 'статус заказа';
comment on column orders.accrual is 'сумма вознагрождения по заказу';

create table matches (
                         id uuid not null default gen_random_uuid() unique,
                         match_name text not null unique,
                         reward int not null,
                         reward_type rewards not null,
                         created_at timestamp not null default now(),
                         updated_at timestamp not null default now(),
                         deleted_at timestamp
);
comment on table matches is 'товары с их механиками вознагрождения';
comment on column matches.match_name is 'название товара';
comment on column matches.reward is 'сумма вознагрождения';
comment on column matches.reward_type is 'механика вознагрождения';

create table goods (
                       id uuid not null default gen_random_uuid(),
                       order_id uuid not null,
                       match_id uuid not null,
                       price int not null,
                       accrual int,
                       foreign key (order_id) references orders(id),
                       foreign key (match_id) references matches(id)
);
comment on table goods is 'товары в заказе с их стоимостью и вознагрождением';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table goods

drop table matches

drop table orders

drop type rewards

drop type order_status
-- +goose StatementEnd
