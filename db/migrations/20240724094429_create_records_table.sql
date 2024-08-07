-- migrate:up
CREATE TABLE IF NOT EXISTS records (
    id SERIAL NOT NULL,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    userId INT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(id),
    FOREIGN KEY (userId) REFERENCES users(id)

);

-- migrate:down
