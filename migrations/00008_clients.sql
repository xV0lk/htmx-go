-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    IF NOT EXISTS clients (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        phone VARCHAR(15) NOT NULL,
        email TEXT NOT NULL,
        notifications BOOLEAN NOT NULL DEFAULT TRUE,
        studio_id INT NOT NULL REFERENCES studios(id) ON DELETE NO ACTION
      

    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;

-- +goose StatementEnd