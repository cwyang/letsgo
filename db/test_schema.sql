create database test_notes character set utf8mb4 collate utf8mb4_unicode_ci;

create user 'test_user'@'localhost';
grant create, drop, alter, index, select, insert, update, delete on test_notes.* to 'test_user'@'localhost';
alter user 'test_user'@'localhost' identified by 'pass';
