# Beego social app server

## Requirements
* Development Language: Go 1.10 or higher
* Server Framework: Beego 1.9
* Database: MySQL (5.6+)
* API documentation: Swagger UI
* ODBC: go-sql-driver

## Installations
1. Install Go and configurate environment variables.
2. Create three files named `src/`, `bin/`, and `pkg/` in your $GOPATH. Then install Beego framework in `$GOPATH/src/`
```
go get github.com/astaxie/beego
go get -u github.com/beego/bee
```
3. Install MySQL
4. Install go-sql-driver
```
go get -u github.com/go-sql-driver/mysql
```

## Updates
### V. 0.1.5
**New features:**
1. Support Chinese characters by implementing utf8

**Details:**

Controllers中param参数类型，可以有的值是 formData、query、path、body、header，formData 表示是 post 请求的数据，query 表示带在 url 之后的参数，path 表示请求路径上得参数，例如上面例子里面的 key，body 表示是一个 raw 数据请求，header 表示带在 header 信息中得参数。


**LEAVE:**
1. add USER_PROFILE table
2. Role based access control(RBAC)
3. Change password storage method
4. add interest API profile(including `Follow` API)
 
### V. 0.1.4
**New features:**
1. Modify FRIEND table
```sql
drop table if exists FRIEND;
create table FRIEND (
  host_id     varchar(12) not null, 
  follower_id varchar(12) not null,
  primary key (host_id, follower_id),
  foreign key (fk_host_id) references USER(user_id) on delete cascade,
  foreign key (fk_follwer_id) references USER(user_id) on delete cascade
);
insert into FRIEND
  values ('12345', '17801055134');
```
2. Create TOKEN table
```sql
drop table if exists TOKEN;
create table TOKEN (
  token_id  varchar(80) not null,
  fk_user_id   varchar(12) not null,
  primary key (token_id, fk_user_id),
  foreign key (fk_user_id) references USER(user_id) on delete cascade
);
insert into TOKEN 
    values ('abcsuccessfullylogin', '12345');
```
3. Create INTEREST table
```sql
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
```
4. Create AREA table
```sql
drop table if exists AREA;
create table AREA (
    interest_tag    varchar(20) not null,
    fk_moment_id       bigint not null,
    foreign key (fk_moment_id) references MOMENT(moment_id) on delete cascade,
    primary key (interest_tag, fk_moment_id)
);
```
5. To post a moment requires user's token, no longer user's id.

**Details**
The output JSON of `GetAll` API was encapsulated by a `map` data structure. But now, I change `map` to `slice` for simplifying the structure of the JSON.

**LEAVE:**
1. add USER_PROFILE table
2. support Chinese characters in moment text and tags
4. Role based access control(RBAC)
5. Change password storage method
6. add interest API profile(including `Follow` API)


### V. 0.1.3
**New features:**
1. add friend list table but haven't test it
```sql
drop table if exists FRIEND;
create table FRIEND (
  host_id     varchar(12) not null,
  follower_id varchar(12) not null,
  primary key (host_id),
  foreign key (fk_host_id) references USER(user_id) on delete cascade
);
insert into FRIEND
  values ('12345', '17801055134');
```
2. promote the login process by returning tokens rather than user id.

**LEAVE:**
1. add USER_PROFILE table
2. support Chinese characters in moment text and tags
3. test FRIEND_LIST table and add `FindFriend` API
4. Role based access control(RBAC)
5. Change password storage method

### V. 0.1.2
The first thing is to replace the string type location to a bool type showing if it exists. Another thing is to add foreign key into the table `MOMENT`, so that every moment has its owner. Thus the SQLs for building tables has changed as follows:
```sql
drop table if exists USER;
drop table if exists MOMENT;

create table USER (
   user_id              varchar(12) not null,
   user_name            varchar(20) not null,
   password             varchar(30) not null,
   primary key (user_id)
);
insert into USER 
    values ('12345', 'abc', 'password');

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
```

LEAVE:
1. add USER_PROFILE table
2. support Chinese characters in moment text and tags
3. add FRIEND_LIST table
4. Role based access control(RBAC)
5. Change password storage method

### V. 0.1.1
APIs about moment:
1. Get all moments
2. Post a moment
3. Get a moment by moment id
4. Delete a moment by moment id

APIs about user:
1. Get all users
2. Post a new user
3. Login user
4. Logout user
5. Get a user by user's phone

### V. 0.1.0
1. create two tables named `USER` and `COMMENT`
2. generate an api example from Beego

## SQLs used for building tables
see file *create_db.sql*
