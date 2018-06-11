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
insert into USER
    values ('17801055134', 'Shane', 'password');

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
  fk_host_id     varchar(12) not null, 
  fk_follower_id varchar(12) not null,
  primary key (fk_host_id, fk_follower_id),
  foreign key (fk_host_id) references USER(user_id) on delete cascade,
  foreign key (fk_follower_id) references USER(user_id) on delete cascade
);
insert into FRIEND
  values ('12345', '17801055134');

drop table if exists TOKEN;
create table TOKEN (
  token_id  varchar(80) not null,
  fk_user_id   varchar(12) not null,
  primary key (token_id, fk_user_id),
  foreign key (fk_user_id) references USER(user_id) on delete cascade
);
insert into TOKEN 
    values ('abcsuccessfullylogin', '12345');

drop table if exists INTEREST;
create table INTEREST (
    fk_user_id      varchar(12) not null,
    interest_tag    varchar(20) not null,
    foreign key (fk_user_id) references USER(user_id) on delete cascade,
    primary key (fk_user_id, interest_tag)
);
insert into INTEREST
    values ('12345', 'tag');
insert into INTEREST
    values ('17801055134', 'Golang');

drop table if exists AREA;
create table AREA (
    interest_tag    varchar(20) not null,
    fk_moment_id       bigint not null,
    foreign key (fk_moment_id) references MOMENT(moment_id) on delete cascade,
    primary key (interest_tag, fk_moment_id)
);

