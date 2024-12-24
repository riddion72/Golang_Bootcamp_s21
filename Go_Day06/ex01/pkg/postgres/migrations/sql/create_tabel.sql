CREATE TABLE IF NOT EXISTS articles
(
    id int NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    article_name text UNIQUE NOT NULL,
    article_content text NOT NULL
);