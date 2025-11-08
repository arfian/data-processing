BEGIN;

CREATE TABLE IF NOT EXISTS products (
    id int PRIMARY KEY NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    brand VARCHAR(100) NOT NULL,
    category VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) DEFAULT 0.00,
    currency VARCHAR(20) NULL,
    stock int DEFAULT 0,
    ean VARCHAR(50) NULL,
    color VARCHAR(50) NULL,
    size VARCHAR(50) NULL,
    availability VARCHAR(50) NULL,
    internal_id int NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL
);

COMMIT;