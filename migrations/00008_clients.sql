-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS clients (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        phone VARCHAR(15) NOT NULL,
        email TEXT NOT NULL,
        notifications BOOLEAN NOT NULL DEFAULT TRUE,
        studio_id VARCHAR(20) NOT NULL REFERENCES studios(id) ON DELETE DO NOTHING
      

    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;

-- +goose StatementEnd