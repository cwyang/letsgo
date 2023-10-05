create table notes (
       id integer not null primary key auto_increment,
       title varchar(100) not null,
       content text not null,
       created datetime not null,
       expires datetime not null
);

create index idx_notes_created on notes(created);

create table users (
       id integer not null primary key auto_increment,
       name varchar(255) not null,
       email varchar(255) not null,
       hashed_password char(60) not null,
       created datetime not null,
       active boolean not null default true
);

alter table users add constraint users_uc_email unique (email);

insert into users (name, email, hashed_password, created) values (
       'Alice Kim',
       'alice@mail.com',
       '12345678900987654321123456789012',
       '2023-10-05 18:31:22'
);
