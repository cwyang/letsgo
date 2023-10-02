insert into notes (title, content, created, expires) values (
       'Hello world',
       'Hello world...\nThis is a sample text\n<END>',
       utc_timestamp(),
       date_add(utc_timestamp(), interval 365 day)
);

insert into notes (title, content, created, expires) values (
       '한글',
       '안녕하세요...\n이것은 샘플 텍스트입니다.\n반갑습니다.',
       utc_timestamp(),
       date_add(utc_timestamp(), interval 365 day)
);

insert into notes (title, content, created, expires) values (
       '한글 2',
       '안녕하세요...\n이것은 샘플 텍스트입니다.\n반갑습니다.\n정말이에요',
       utc_timestamp(),
       date_add(utc_timestamp(), interval 7 day)
);
