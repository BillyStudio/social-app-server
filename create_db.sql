# create database user
CREATE DATABASE social_app;
CREATE USER 'ubuntu'@'localhost' IDENTIFIED BY 'IS1501';
GRANT Alter, Alter routine, Create, Create routine, Delete, Drop, Execute, Index, Insert, References, Select, Update ON social_app.* TO 'ubuntu'@'localhost';

# create tables
use social_app;

drop table if exists USER;
create table USER (
   user_id              varchar(12) not null,
   user_name            varchar(20) not null,
   password             varchar(30) not null,
   primary key (user_id)
);
insert into USER 
    values ('12345', 'abc', 'password');


drop table if exists MOMENT;
create table MOMENT (
   moment_id            bigint not null,
   moment_time          varchar(128) not null,
   if_tag               bool not null,
   if_text              bool not null,
   if_image             bool not null,
   fk_user_id           varchar(12) not null,
   primary key (moment_id),
   foreign key (fk_user_id) references USER(user_id) on delete cascade
);
insert into MOMENT
    values (1, '2006-01-02 15:04:05', false, false, false, '12345');

drop table if exists FRIEND;
create table FRIEND (
  host_id     varchar(12) not null,
  follower_id varchar(12) not null,
  primary key (host_id),
  foreign key (fk_host_id) references USER(user_id) on delete cascade
);
insert into FRIEND
  values ('12345', '17801055134');

drop table if exists TOKEN;
create table TOKEN (
  token_id  varchar(80) not null,
  user_id   varchar(12) not null
);
