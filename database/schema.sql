CREATE TABLE IF NOT EXISTS todos (
  id INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  done BOOLEAN
);

INSERT INTO todos (name, done) VALUES ("test", "0");
