BEGIN;
CREATE TABLE IF NOT EXISTS users(
	id bigserial,
	first_name varchar(50),
  last_name varchar(50),
  password_hash bytea,
  created_at timestamp,
  primary key (id)
);

CREATE TABLE IF NOT EXISTS habits(
  id bigserial,
  user_id int,
  name varchar(50),
  type varchar(30),
  primary key (id),
  foreign key (user_id) references users(id) ON DELETE SET NULL
);

DELETE FROM habits;
ALTER SEQUENCE habits_id_seq RESTART WITH 1;
DELETE FROM users;
ALTER SEQUENCE users_id_seq RESTART WITH 1;

INSERT INTO users (first_name, last_name, created_at) VALUES ('test_1', 'test_2', 'now()'), ('test_3', 'test4', 'now()');
INSERT INTO habits (user_id, name, type) VALUES (1, 'habit1', 'test'), (2, 'habit2', 'test');
COMMIT;