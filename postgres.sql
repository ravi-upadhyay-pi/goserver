create table users(
  username varchar(32) not null primary key,
  password varchar(32) not null,
  first_name varchar(32) not null,
  last_name varchar(32) not null,
  email_id varchar(32) not null,
  phone_number varchar(15) not null,
  created_on timestamp not null,
  updated_on timestamp not null);

create table article(
  username varchar(32) not null,
	id varchar(128) not null,
	title varchar(128) not null,
	body text not null,
	tags varchar(128) not null,
	format varchar(10) not null,
	next bigint not null,
	previous bigint not null,
	private boolean not null,
	created_on timestamp not null,
	updated_on timestamp not null,
	primary key(username, id),
	foreign key(username) references users(username));
