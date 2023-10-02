create database notes character set utf8mb4 collate utf8mb4_unicode_ci;

use notes

create table notes (
       id integer not null primary key auto_increment,
       title varchar(100) not null,
       content text not null,
       created datetime not null,
       expires datetime not null
);

create index idx_notes_created on notes(created);

create user 'user'@'localhost';
grant select, insert, update on notes.* to 'user'@'localhost';

-- alter user 'user'@'localhost' identified by 'password';
