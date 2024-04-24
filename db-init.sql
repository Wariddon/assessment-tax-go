DROP TABLE IF EXISTS test CASCADE;

CREATE TABLE test
(
    id          SERIAL PRIMARY KEY,
    test VARCHAR(255) UNIQUE NOT NULL
);

DROP TABLE IF EXISTS allowances CASCADE;

CREATE TABLE allowances
(
    id  SERIAL PRIMARY KEY,
    allowance VARCHAR(255) UNIQUE NOT NULL,
    amount FLOAT
);

INSERT INTO allowances 
    (allowance, amount)
VALUES 
    ('personal', 60000),
    ('donation', 100000),
    ('k-receipt', 50000)
