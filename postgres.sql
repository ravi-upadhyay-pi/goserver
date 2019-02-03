create table users(
  username varchar(32) not null primary key,
  password varchar(32) not null,
  first_name varchar(32) not null,
  last_name varchar(32) not null,
  email_id varchar(32) not null,
  phone_number varchar(15) not null,
  created_on timestamp not null,
  updated_on timestamp not null);
