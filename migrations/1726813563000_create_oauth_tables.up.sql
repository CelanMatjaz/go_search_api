
CREATE TABLE IF NOT EXISTS oauth_clients (
    id SERIAL,
    name BPCHAR NOT NULL,
    client_id BPCHAR NOT NULL UNIQUE,
    client_secret BPCHAR NOT NULL UNIQUE,
    scopes BPCHAR NOT NULL,
    code_endpoint BPCHAR NOT NULL,
    token_endpoint BPCHAR NOT NULL,
    account_data_endpoint BPCHAR NOT NULL,

    CONSTRAINT oauth_clients_pkey PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS account_oauth_data (
    id SERIAL,
    account_id INT,
    oauth_client_id INT,
    access_token VARCHAR(2048) NOT NULL,
    refresh_token VARCHAR(512) NOT NULL,

    CONSTRAINT account_oauth_data_has_account
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,

    CONSTRAINT account_oauth_data_has_oauth_client
    FOREIGN KEY (oauth_client_id) REFERENCES oauth_clients(id) ON DELETE CASCADE,

    CONSTRAINT oauth_account_data_pkey PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS oauth_client_name_index ON oauth_clients(name);
