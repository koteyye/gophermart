-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSED', 'PROCESSING', 'INVALID');

CREATE TABLE IF NOT EXISTS users (
	id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
	login varchar NOT NULL UNIQUE,
	hashed_password varchar NOT NULL,
	current_balance int NOT NULL DEFAULT 0,
	withdrawn_balance int NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS orders (
	id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
	user_created uuid NOT NULL,
	order_number varchar NOT NULL,
	status order_status DEFAULT 'NEW',
	accrual int NOT NULL DEFAULT 0,
	created_at timestamp DEFAULT now(),
	updated_at timestamp DEFAULT now(),
	FOREIGN KEY (user_created) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS operations (
	id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
	user_created uuid NOT NULL,
	order_number varchar NOT NULL,
	amount int NOT NULL DEFAULT 0,
	created_at timestamp DEFAULT now(),
	FOREIGN KEY (user_created) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS operations;

DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS users;

DROP TYPE order_status;
-- +goose StatementEnd
