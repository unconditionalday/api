-- +goose Up
CREATE TABLE feeds (
    title VARCHAR(255) NOT NULL,
    link VARCHAR(255) PRIMARY KEY,
    language VARCHAR(10),
    image_title VARCHAR(255),	
    image_url VARCHAR(255),
    summary TEXT,
    source VARCHAR(255),
    date timestamptz,
    embedding vector(384)
);

-- +goose Down
DROP TABLE feeds;
