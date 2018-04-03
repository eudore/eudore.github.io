CREATE DATABASE IF NOT EXISTS Jass DEFAULT CHARACTER SET utf8 collate utf8_general_ci;

USE Jass;
DROP TABLE IF EXISTS tb_User;

CREATE TABLE IF NOT EXISTS tb_User(
	ID int UNSIGNED PRIMARY KEY Auto_Increment,
	User VARCHAR(50) UNIQUE,
	Pass text NOT NULL,
	Name text NOT NULL,
	TYPE INT UNSIGNED DEFAULT 0,
	LoginIP INT UNSIGNED,
	LoginTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	SiginTime DATE
)Engine InnoDB;

CREATE TABLE IF NOT EXISTS tb_Index(
	ID int UNSIGNED PRIMARY KEY Auto_Increment,
	Name text,
	PID int UNSIGNED,
	MD5 char(32),
	Active boolean default False
)Engine InnoDB;

CREATE TABLE IF NOT EXISTS tb_Content(
	MD5 char(32) PRIMARY KEY,
	EditTime TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	Content text
)Engine InnoDB;


CREATE TABLE IF NOT EXISTS tb_Session(
	ID CHAR(32) PRIMARY KEY,
	TYPE INT UNSIGNED DEFAULT 0,
	Timer TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)Engine InnoDB;


SHOW PROCEDURE STATUS;


USE Jass;
DROP PROCEDURE IF EXISTS pro_Login;
DELIMITER //
CREATE PROCEDURE pro_Login(in user VARCHAR(50),in pass text,in ip int,out out_name text,out out_type INT UNSIGNED)
BEGIN
	declare pro_id int default 0;
	SELECT ID,Name,TYPE INTO pro_id,out_name,out_type FROM Jass.tb_User WHERE tb_User.User=user AND tb_User.Pass=pass;
	IF(pro_id!=0)THEN
		UPDATE tb_User Set LoginIP=ip WHERE ID=pro_id;
	ELSE
		SET out_name = '';
		SET out_type = 0;
	END IF;
		SELECT out_name,out_type;
END//
DELIMITER ;

USE Jass;
DROP PROCEDURE IF EXISTS pro_Sigin;
DELIMITER //
CREATE PROCEDURE pro_Sigin(in user VARCHAR(50),in pass text,in name text,out flag boolean)
BEGIN
	IF NOT EXISTS(SELECT ID FROM tb_User WHERE tb_User.User=user) THEN 
		INSERT INTO tb_User(User,Pass,Name,SiginTime) VALUES(user,pass,name,curdate());
		set flag = true;
	ELSE
		set flag = false;
	END IF;
END//
DELIMITER ;


USE Jass;
DROP PROCEDURE IF EXISTS pro_Oauth;
DELIMITER //
CREATE PROCEDURE pro_Oauth(in token CHAR(32),in t INT)
BEGIN
	declare pro_TYPE int UNSIGNED;
	SELECT TYPE INTO pro_TYPE FROM tb_Session WHERE ID=token;
	SELECT((t&pro_TYPE)=t);
END//
DELIMITER ;




USE Jass;
DROP PROCEDURE IF EXISTS pro_AddIndex;
DELIMITER //
CREATE PROCEDURE pro_AddIndex(in path text)
BEGIN
	declare pro_PID int UNSIGNED DEFAULT 0;
	declare pro_path text;
	declare pro_name text;
	SET pro_path=REPLACE(path,'/',CHAR(0));
    SET pro_name=substring_index(pro_path,CHAR(0),-1);
    SET pro_path=LEFT(pro_path,LENGTH(pro_path)-LENGTH(pro_name)-1);
	IF NOT EXISTS(SELECT ID FROM tb_Index WHERE MD5=MD5(path)) and path!=pro_path THEN
		SELECT ID into pro_PID FROM tb_Index WHERE MD5=MD5(pro_path) limit 1;
		IF pro_PID=0 THEN
			CALL pro_AddIndex(pro_path);
			SELECT ID into pro_PID FROM tb_Index WHERE MD5=MD5(pro_path) limit 1;
		END IF;
		INSERT INTO tb_Index(Name,PID,MD5) VALUES(pro_name,pro_PID,MD5(path));
		UPDATE tb_Index SET Active=True WHERE ID=pro_PID;
	END IF;
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS pro_GetIndex;
DELIMITER //
CREATE PROCEDURE pro_GetIndex(in path text)
BEGIN
	declare pro_PID int UNSIGNED;
	SELECT ID into pro_PID FROM tb_Index WHERE MD5=path limit 1;
	SELECT Name,Active FROM tb_Index WHERE PID=pro_PID;
END//
DELIMITER ;


DROP PROCEDURE IF EXISTS pro_Sigin;
DELIMITER //
CREATE PROCEDURE pro_LoadIndex(in path text)
	SET @path=REPLACE(path,'/',CHAR(0))
	SELECT @PID:=ID FROM tb_Index WHERE MD5=MD5(@path)
	SELECT Name FROM tb_Index WHERE PID=@PID
BEGIN
END//
DELIMITER ;
