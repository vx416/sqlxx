	CREATE TABLE IF NOT EXISTS users(
		id INTEGER primary key AUTOINCREMENT,
		first_name TEXT,
    last_name TEXT,
    password_hash BLOB,
    created_at TEXT
	);

  CREATE TABLE IF NOT EXISTS habits(
    id INTEGER primary key AUTOINCREMENT,
    user_id INTEGER,
    name TEXT,
    type TEXT,
    FOREIGN KEY(user_id) REFERENCES user(id)
  );