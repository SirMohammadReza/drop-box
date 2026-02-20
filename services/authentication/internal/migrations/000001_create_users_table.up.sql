CREATE TABLE
    IF NOT EXISTS users (
        id BIGSERIAL PRIMARY KEY,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ,
        name VARCHAR(30) NOT NULL,
        phone_number VARCHAR(13) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL
    );