ALTER TABLE users
    ADD COLUMN app_publish_amount BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN app_amount BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN address VARCHAR(255) DEFAULT '',
    ADD COLUMN sign VARCHAR(255) DEFAULT '',
    ADD COLUMN status INT NOT NULL DEFAULT 1,
    ADD COLUMN code VARCHAR(255) DEFAULT '',
    ADD COLUMN twitter VARCHAR(255) DEFAULT '',
    ADD COLUMN telegram VARCHAR(255) DEFAULT '',
    ADD COLUMN discord VARCHAR(255) DEFAULT '';

-- add dataset table
CREATE TABLE datasets (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    created_at BIGINT DEFAULT (strftime('%s', 'now')),
    created_by VARCHAR(255),
    data_source_type VARCHAR(255),
    indexing_technique VARCHAR(255),
    permission VARCHAR(255),
    owner VARCHAR(255)
);

-- add user_assets table
CREATE TABLE user_assets (
    id SERIAL PRIMARY KEY,
    uid VARCHAR(255) NOT NULL,
    mainchain VARCHAR(255) NOT NULL DEFAULT 'SOLANA',
    token VARCHAR(255) NOT NULL DEFAULT 'SOL',
    amount DECIMAL(20,8) NOT NULL DEFAULT 0.00000000,
    frozen DECIMAL(20,8) NOT NULL DEFAULT 0.00000000,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- add user_withdraws table
CREATE TABLE user_withdraws (
    id SERIAL PRIMARY KEY,
    uid VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    mainchain VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    amount DECIMAL(20,8),
    description VARCHAR(255)
);

