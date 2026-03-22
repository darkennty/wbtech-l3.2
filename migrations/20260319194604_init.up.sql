CREATE TABLE urls (
    short_url VARCHAR PRIMARY KEY,
    long_url VARCHAR NOT NULL
);

CREATE TABLE stats (
    id SERIAL PRIMARY KEY,
    url VARCHAR NOT NULL,
    time TIMESTAMP NOT NULL,
    user_agent VARCHAR NOT NULL
);