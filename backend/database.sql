create table article
(
    id       integer
        constraint article_pk
            primary key autoincrement,
    title    text    not null,
    summary  text,
    image    text,
    created  integer not null,
    modified integer default 0,
    link     text    not null,
    markdown text    not null,
    html     text    not null
);

create unique index article_link_uindex
    on article (link);

create unique index article_title_uindex
    on article (title);

create table article_tag
(
    article_id integer not null
        references article
            on delete cascade,
    tag        text
);

create table assets
(
    id   integer
        constraint assets_pk
            primary key autoincrement,
    name text   not null,
    data blob   not null,
    link string not null
);

create unique index assets_link_uindex
    on assets (link);

create unique index assets_name_uindex
    on assets (name);

create table author
(
    id          integer
        constraint author_pk
            primary key autoincrement,
    name        text not null,
    password    text not null,
    information text
);

create table article_author
(
    article_id integer not null
        references article
            on delete cascade,
    author_id  integer not null
        references author
            on delete cascade
);

create unique index author_name_uindex
    on author (name);

