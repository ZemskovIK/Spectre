CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    login VARCHAR(200) UNIQUE NOT NULL,
    phash BLOB NOT NULL,
    access_level INTEGER DEFAULT 1
);


INSERT INTO users (login, phash, access_level) 
VALUES ('admin', '$2a$10$sj7WSFEVuSDHPr8.FywNZuCbjpBl40ai2xm6Bd089POhWFeTJW57e', 6);