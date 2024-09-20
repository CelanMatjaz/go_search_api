
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL,
    account_id INT NOT NULL,
    label VARCHAR(64) NOT NULL,
    color VARCHAR(7) UNIQUE NOT NULL,

    CONSTRAINT account_has_tags 
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,

    CONSTRAINT tags_pkey PRIMARY KEY(id)
)
