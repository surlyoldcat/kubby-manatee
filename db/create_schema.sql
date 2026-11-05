DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS movies;

CREATE TABLE movies (
    movie_id INTEGER PRIMARY KEY,
    title    TEXT NOT NULL,
    genres   TEXT NOT NULL  -- pipe-delimited, e.g. 'Adventure|Animation'
);

CREATE TABLE links (
    movie_id INTEGER PRIMARY KEY REFERENCES movies (movie_id),
    imdb_id  TEXT NOT NULL,  -- text to preserve leading zeros
    tmdb_id  INTEGER
);

CREATE TABLE ratings (
    user_id  INTEGER NOT NULL,
    movie_id INTEGER NOT NULL REFERENCES movies (movie_id),
    rating   NUMERIC(2,1) NOT NULL CHECK (rating BETWEEN 0.5 AND 5.0),
    rated_at BIGINT NOT NULL,  -- Unix epoch seconds
    PRIMARY KEY (user_id, movie_id)
);

CREATE TABLE tags (
    tag_id    SERIAL PRIMARY KEY,
    user_id   INTEGER NOT NULL,
    movie_id  INTEGER NOT NULL REFERENCES movies (movie_id),
    tag       TEXT NOT NULL,
    tagged_at BIGINT NOT NULL  -- Unix epoch seconds
);

CREATE INDEX ratings_movie_id_idx ON ratings (movie_id);
CREATE INDEX tags_movie_id_idx ON tags (movie_id);
