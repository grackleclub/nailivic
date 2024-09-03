CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL,
    username VARCHAR ( 255 ) UNIQUE NOT NULL,
    hashed_password VARCHAR ( 255 ) NOT NULL,
    created_on TIMESTAMP,
    last_login TIMESTAMP
)
