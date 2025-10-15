-- +goose Up
CREATE TABLE accesses (
    id serial PRIMARY KEY,
    endpoint text NOT NULL,
    role INT NOT NULL,
    UNIQUE (endpoint, role)
);

-- +goose Down
DROP TABLE accesses;
