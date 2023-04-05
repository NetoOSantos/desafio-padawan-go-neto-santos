CREATE DATABASE CONVERSAO;

USE CONVERSAO;

CREATE TABLE IF NOT EXISTS conversions (
    id INTEGER PRIMARY KEY,
    amount REAL,
    from_currency TEXT,
    to_currency TEXT,
    rate REAL,
    result REAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
