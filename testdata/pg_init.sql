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
  constraint fk_user foreign key (user_id) references users(id)
);