CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    created_on TIMESTAMP,
    last_login TIMESTAMP
);

INSERT INTO users (
    username,
    hashed_password
)
VALUES (
    'admin',
    'fortestingonlyTODOremove'
);