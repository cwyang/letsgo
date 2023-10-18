alter table movies add constraint movies_runtime_check CHECK (runtime >= 0);
alter table movies add constraint movies_year_check CHECK (year between 1888 and date_part('year', now()));
alter table movies add constraint grenres_length_check CHECK (array_length(genres, 1) between 1 and 5);
