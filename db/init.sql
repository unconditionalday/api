-- Install the extension we just compiled

CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;

/*
For simplicity, we are directly adding the content into this table as
a column containing text data. It could easily be a foreign key pointing to
another table instead that has the content you want to vectorize for
semantic search, just storing here the vectorized content in our "items" table.

"768" dimensions for our vector embedding is critical - that is the
number of dimensions our open source embeddings model output, for later in the
blog post.
*/

CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    language VARCHAR(10),
    image_url VARCHAR(255),
    image_width INT,
    image_height INT,
    summary TEXT,
    source VARCHAR(255),
    date TIMESTAMP,
    embedding vector(768)
);
