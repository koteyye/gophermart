-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('new', 'processed', 'processing', 'invalid');

CREATE TYPE operation_status AS ENUM ('run', 'done', 'error');

CREATE TABLE IF NOT EXISTS users (
	id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
	login varchar(128) NOT NULL UNIQUE,
	hashed_password varchar(128) NOT NULL,
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now(),
	deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS balance (
	id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
	user_id uuid NOT NULL,
	amount int NOT NULL DEFAULT 0,
	withdrawn int NOT NULL DEFAULT 0,
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now(),
	deleted_at timestamp,
	FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS orders (
	number varchar NOT NULL PRIMARY KEY,
	status order_status DEFAULT 'new',
	accrual int NOT NULL DEFAULT 0,
	user_created uuid,
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now(),
	deleted_at timestamp,
	FOREIGN KEY (user_created) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS operations (
	id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
	order_number varchar,
	amount int,
	status operation_status DEFAULT 'run',
	balance_id uuid,
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now(),
	deleted_at timestamp,
	FOREIGN KEY (balance_id) REFERENCES balance(id),
	FOREIGN KEY (order_number) REFERENCES orders(number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS operations;

DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS balance;

DROP TABLE IF EXISTS users;

DROP TYPE operation_status;

DROP TYPE order_status;
-- +goose StatementEnd
