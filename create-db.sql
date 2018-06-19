# create database user
CREATE DATABASE social_app;
CREATE USER 'ubuntu'@'localhost' IDENTIFIED BY 'IS1501';
GRANT Alter, Alter routine, Create, Create routine, Delete, Drop, Execute, Index, Insert, References, Select, Update ON social_app.* TO 'ubuntu'@'localhost';

# create tables
use social_app;


drop table if exists FRIEND;
drop table if exists TOKEN;
drop table if exists INTEREST;
drop table if exists AREA;
drop table if exists MOMENT;
drop table if exists USER;

create table USER (
   user_id              varchar(20) not null,
   user_name            varchar(30) not null,
   password             varchar(30) not null,
   primary key (user_id)
)default charset=utf8;
insert into USER 
    values ('12345', 'abc', 'password');
insert into USER
    values ('17801055134', 'Shane', 'password');
insert into USER
	values ('billy.ustb@gmail.com', 'Billy', 'password');


create table FRIEND (
  fk_host_id     varchar(20) not null, 
  fk_follower_id varchar(20) not null,
  primary key (fk_host_id, fk_follower_id),
  foreign key (fk_host_id) references USER(user_id) on delete cascade,
  foreign key (fk_follower_id) references USER(user_id) on delete cascade
)default charset=utf8;
insert into FRIEND
  values ('12345', '17801055134');

create table MOMENT (
   moment_id            bigint not null,
   moment_time          varchar(128) not null,
   if_tag               bool not null default false,
   if_text              bool not null default false,
   if_image             bool not null default false,
   fk_user_id           varchar(12) not null,
   likes                int not null default 0,
   primary key (moment_id),
   foreign key (fk_user_id) references USER(user_id) on delete cascade
)default charset=utf8;
insert into MOMENT
    values (1, '2006-01-02 15:04:05', false, false, false, '12345', 0);
insert into MOMENT
	values (2, '2018-06-03 17:23:23', true, true, false, '17801055134', 1);

create table INTEREST (
    fk_user_id      varchar(20) not null,
    interest_tag    varchar(20) not null,
    foreign key (fk_user_id) references USER(user_id) on delete cascade,
    primary key (fk_user_id, interest_tag)
)default charset=utf8;
insert into INTEREST
    values ('12345', 'tag');
insert into INTEREST
    values ('17801055134', 'Golang');

create table AREA (
    interest_tag    varchar(20) not null,
    fk_moment_id       bigint not null,
    foreign key (fk_moment_id) references MOMENT(moment_id) on delete cascade,
    primary key (interest_tag, fk_moment_id)
) default charset=utf8;
insert into AREA
	values ('tag', 1);

create table TOKEN (
  token_id  varchar(80) not null,
  fk_user_id   varchar(20) not null,
  primary key (token_id, fk_user_id),
  foreign key (fk_user_id) references USER(user_id) on delete cascade
)default charset=utf8;
insert into TOKEN 
    values ('abcsuccessfullylogin', '12345');

# test data
insert into USER 
     values ('123456456', 'TT', 'passwd');
insert into TOKEN 
     values ('9085cb83222cfa67d4c07d0159cde207', '123456456');
insert into MOMENT
     values (1528817152888728750, '2018-06-12 23:29:12', true, true, true, '123456456', 0);

