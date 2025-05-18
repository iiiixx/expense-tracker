CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(40) NOT NULL UNIQUE,  
    password CHAR(60) NOT NULL        
);

CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,  
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),        
    category VARCHAR(40) NOT NULL,
    description TEXT,
    DATE NOT NULL DEFAULT CURRENT_DATE
);