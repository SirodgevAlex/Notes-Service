CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    author_id INT,
    title TEXT
    text TEXT,
    FOREIGN KEY (author_id) REFERENCES users(id)
);