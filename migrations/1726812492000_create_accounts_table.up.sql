
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL,
    display_name VARCHAR(64) NOT NULL,
    email VARCHAR(256) UNIQUE NOT NULL,
    refresh_token_version INT NOT NULL DEFAULT 1,
    is_oauth BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),

    CONSTRAINT accounts_pkey PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS password_hashes (
    account_id INT UNIQUE NOT NULL,
    password_hash CHARACTER(70) NOT NULL,

    CONSTRAINT account_may_have_password_hash 
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,

    CONSTRAINT password_hashes_pkey PRIMARY KEY(account_id)
);
