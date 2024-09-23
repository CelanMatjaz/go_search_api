
CREATE TABLE IF NOT EXISTS resume_presets (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    label VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),

    CONSTRAINT account_has_many_resume_presets
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS resume_sections (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    label VARCHAR(64) NOT NULL,
    text VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),

    CONSTRAINT account_has_many_resume_sections
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mtm_tags_resume_presets (
    id SERIAL PRIMARY KEY,
    tag_id INT NOT NULL,
    record_id INT NOT NULL,
    UNIQUE (tag_id, record_id),

    CONSTRAINT tag_has_many_resume_presets
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,

    CONSTRAINT resume_preset_has_many_tags
    FOREIGN KEY (record_id) REFERENCES resume_presets(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mtm_tags_resume_sections (
    id SERIAL PRIMARY KEY,
    tag_id INT NOT NULL,
    record_id INT NOT NULL,
    UNIQUE (tag_id, record_id),

    CONSTRAINT tag_has_many_resume_sections
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,

    CONSTRAINT resume_section_has_many_tags
    FOREIGN KEY (record_id) REFERENCES resume_sections(id) ON DELETE CASCADE
);
