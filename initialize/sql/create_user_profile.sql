CREATE TABLE IF NOT EXISTS UserProfile (
    id VARCHAR(255) PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth TEXT,
    address TEXT,
    user_id VARCHAR(255) references User(id)
) ENGINE=INNODB;