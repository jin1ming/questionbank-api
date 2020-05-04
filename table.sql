
CREATE TABLE `paper` (
    id integer(64) NOT NULL ,
    title varchar(255) NOT NULL, -- 试卷标题
    owner varchar(255) NOT NULL,
    primary key (`id`)
) ;


CREATE TABLE `paper_question` (
    question_id varchar(255) NOT NULL,
    paper_id integer(64) NOT NULL,
    primary key (`id`)
) ;


create table `user` (
    username varchar(255) NOT NULL,
    pwd_md5 varchar(32) ,
    role varchar(32) ,
    primary key (`username`)
) ;

