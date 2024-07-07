CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    ticker TEXT NOT NULL,
    decimals smallint NOT NULL
);