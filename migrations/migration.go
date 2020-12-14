package migrations

var Schema = `
create table if not exists news (
	id serial,
    title text,
    descriptions text,
    link integer,

	primary key (id)
)`