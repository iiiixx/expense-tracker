CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(40) NOT NULL,
    password VARCHAR(20) NOT NULL,
);

CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES user(id),
    amout DECIMAL(10,2) NOT NULL,
    category VARCHAR(40) NOT NULL,
    description TEXT,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)