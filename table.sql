
set character_set_database=utf8;
set character_set_server=utf8;

drop database questionbank;
create database questionbank;
use questionbank;

CREATE TABLE `paper` (
    id integer(64) NOT NULL AUTO_INCREMENT,
    title varchar(255) NOT NULL, -- 试卷标题
    owner varchar(255) NOT NULL,
    primary key (`id`)
) AUTO_INCREMENT = 1 ;


CREATE TABLE `paper_question` (
    question_id varchar(255) NOT NULL,
    paper_id integer(64) NOT NULL
) ;


create table `user` (
    username varchar(255) NOT NULL,
    pwd_md5 varchar(32) ,
    role varchar(32) ,
    primary key (`username`)
) ;

