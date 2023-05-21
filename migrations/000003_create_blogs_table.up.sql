CREATE TABLE IF NOT EXISTS blogs(
    ID bigserial primary key,
    title varchar(200),
    content text,
    slug varchar(200),
    created_at timestamp(0) with time zone NOT NULL,
    user_id int not null REFERENCES users ON DELETE CASCADE
);