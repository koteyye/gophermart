-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS enum ('REGISTERED', 'PROCESSED', 'PROCESSING', 'INVALID');

CREATE TYPE rewards AS enum ('PERCENT', 'NATURAL');

CREATE TABLE orders (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY UNIQUE,
    order_number VARCHAR NOT NULL,
    status order_status DEFAULT 'REGISTERED',
    accrual INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE matches (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY UNIQUE,
    match_name TEXT NOT NULL UNIQUE,
    reward INT NOT NULL,
    reward_type rewards NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE goods (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY UNIQUE,
    order_id uuid NOT NULL,
    match_id uuid NOT NULL,
    price INT NOT NULL,
    accrual INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (match_id) REFERENCES matches(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE goods;

DROP TABLE matches;

DROP TABLE orders;

DROP TYPE rewards;

DROP TYPE order_status;
-- +goose StatementEnd
