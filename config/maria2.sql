CREATE DATABASE IF NOT EXISTS Jass DEFAULT CHARACTER SET utf8 collate utf8_general_ci;


DROP TABLE IF EXISTS `tb_auth_oauth2_app`;
CREATE TABLE IF NOT EXISTS `tb_auth_oauth2_app`(
	`Name` VARCHAR(16) NOT NULL,
	`ClientID` VARCHAR(20) NOT NULL,
	`ClientSecret` VARCHAR(40) NOT NULL,
	`Callback` VARCHAR(100) NOT NULL,
	`Scopes` VARCHAR(200) NOT NULL
);

DROP TABLE IF EXISTS `tb_auth_oauth2_code`;
CREATE TABLE IF NOT EXISTS `tb_auth_oauth2_code`(
    Codo VARCHAR(32) PRIMARY KEY COMMENT '授权码',
    Redirect_uri VARCHAR(200) NOT NULL COMMENT '重定向URI',
    ClientID VARCHAR(20) NOT NULL COMMENT '客户端ID',
    Expires TIMESTAMP NOT NULL COMMENT '过期时间'
)

DROP TABLE IF EXISTS `tb_auth_oauth2_pass`;
CREATE TABLE IF NOT EXISTS `tb_auth_oauth2_pass`(
	`Login` VARCHAR(32) PRIMARY KEY COMMENT '登录用户',
	`Pass` VARCHAR(64) NOT NULL COMMENT '登录密码Hash',
	`UID` INT UNSIGNED NOT NULL COMMENT '认证ID'
)

DROP TABLE IF EXISTS `tb_auth_oauth2_source`;
CREATE TABLE IF NOT EXISTS `tb_auth_oauth2_source`(
	`Name` VARCHAR(16) NOT NULL UNIQUE COMMENT '认证方式',
	`App` VARCHAR(16) COMMENT '认证用用名称',
	`ClientID` VARCHAR(80) NOT NULL COMMENT 'ID',
	`ClientSecret` VARCHAR(80) NOT NULL COMMENT 'Secret'
);

DROP TABLE IF EXISTS `tb_auth_oauth2_login`;
CREATE TABLE IF NOT EXISTS `tb_auth_oauth2_login`(
	`Source` INT UNSIGNED NOT NULL COMMENT '认证方式',
	`OID` CHAR(21) NOT NULL COMMENT '认证ID',
	`UID` INT UNSIGNED NOT NULL COMMENT '用户ID',
	`Name` VARCHAR(100) COMMENT '认证名称',
	`Stats` INT COMMENT '认证状态 0正常 1注册'
);

DROP TABLE IF EXISTS `tb_auth_user_info`;
CREATE TABLE IF NOT EXISTS `tb_auth_user_info`(
	`UID` int UNSIGNED PRIMARY KEY Auto_Increment,
	`Name` VARCHAR(16) NOT NULL UNIQUE COMMENT '',
	`Status` INT UNSIGNED DEFAULT 0,
	`Level` INT UNSIGNED DEFAULT 0,
	`LoginIP` INT UNSIGNED DEFAULT 0,
	`LoginTime` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`SiginTime` DATE
);

DROP TABLE IF EXISTS `tb_auth_user_token`;
CREATE TABLE IF NOT EXISTS `tb_auth_user_token`(
	Token VARCHAR(32) PRIMARY KEY COMMENT 'Token',
	Refresh VARCHAR(32) COMMENT '更新令牌',
	Scope int COMMENT '权限范围',
    Expires TIMESTAMP NOT NULL COMMENT '过期时间'
)

DROP TABLE IF EXISTS `tb_auth_project`;
CREATE TABLE IF NOT EXISTS `tb_auth_project`(
	`PID` INT UNSIGNED PRIMARY KEY Auto_Increment,
	`Name` VARCHAR(16) NOT NULL,
	`UID` INT UNSIGNED NOT NULL,
	`Range` INT UNSIGNED DEFAULT 0,
	`Level` INT UNSIGNED DEFAULT 0 COMMENT 'public protected group private',
	`Info` VARCHAR(50),
	`EditTime` DATE,
	`CreateTime` DATE
);
# Guest、Reporter、Developer、Master、Owner、Admin
DROP TABLE IF EXISTS `tb_auth_authorized`;
CREATE TABLE IF NOT EXISTS `tb_auth_authorized`(
	`UID` INT UNSIGNED NOT NULL,
	`PID` INT UNSIGNED NOT NULL,
	`Level` INT UNSIGNED DEFAULT 0,
	PRIMARY KEY(`UID`,`PID`)
);

CREATE OR REPLACE VIEW v_auth_authorized
AS
SELECT U.`UID`,U.`Name` UName,P.`PID`,P.`Name` PName,A.`Level`,P.`Info`,P.`EditTime` FROM tb_auth_user U INNER JOIN tb_auth_authorized A On U.UID=A.UID INNER JOIN tb_auth_project P ON A.PID=P.PID;


DROP TABLE IF EXISTS `tb_file_project`
CREATE TABLE IF NOT EXISTS `tb_file_project`(
	`FID` INT UNSIGNED PRIMARY KEY Auto_Increment COMMENT 'file project ID',
	`UID` INT UNSIGNED,
	`Path` VARCHAR(32) NOT NULL,
	`Store` INT UNSIGNED DEFAULT 0 COMMENT 'store ID',
	`Info` VARCHAR(50),
	`Size` INT,
	`MaxSize` INT DEFAULT 0
)

DROP TABLE IF EXISTS `tb_file_store`;
CREATE TABLE IF NOT EXISTS `tb_file_store`(
	`ID` INT UNSIGNED DEFAULT 0 COMMENT 'store ID',
	`Name` VARCHAR(32) COMMENT 'store name',
	`Host` VARCHAR(64) NOT NULL COMMENT 'Host',
	`Config` VARCHAR(256) COMMENT 'config'
)

DROP TABLE IF EXISTS `tb_file_save`;
CREATE TABLE IF NOT EXISTS `tb_file_save`(
	`ID` INT UNSIGNED PRIMARY KEY Auto_Increment COMMENT 'file save ID',
	`FID` INT UNSIGNED DEFAULT 0 COMMENT 'file project ID',
	`Name` VARCHAR(100) COMMENT 'file name',
	`Path` VARCHAR(100) COMMENT 'file path',
	`Hash` char(32) UNIQUE COMMENT 'file path hash',
	`PHash` char(32) COMMENT 'file path parent hash',
#	`UID` int UNSIGNED DEFAULT 0,
#	`Store` INT UNSIGNED DEFAULT 0,
	`Size` INT UNSIGNED DEFAULT 0,
	`ModTime` TIMESTAMP(6)
);


DROP TABLE IF EXISTS `tb_note_save`;
CREATE TABLE IF NOT EXISTS `tb_note_save`(
	`Hash` char(32) PRIMARY KEY,
	`PHash` char(32),
	`Status` INT UNSIGNED DEFAULT 0,
	`UID` int UNSIGNED DEFAULT 0,
	`Uri` VARCHAR(100),
	`Format` CHAR(6),
	`Name` text,
	`Title` VARCHAR(50),
	`Content` TEXT,
	`EditTime` TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP
);


DROP TABLE IF EXISTS `tb_tools_proxy`
CREATE TABLE IF NOT EXISTS `tb_tools_proxy`(
	`ID` int UNSIGNED PRIMARY KEY Auto_Increment,
	`Server` char(15) NOT NULL,
	`Port` int NOT NULL,
	`Protocol` VARCHAR(32),
	`Method` VARCHAR(32),
	`Obfs` VARCHAR(32),
	`Password` VARCHAR(32),
	`Obfsparam` VARCHAR(32) DEFAULT '',
	`Protoparam` VARCHAR(32) DEFAULT '',
	`Remarks` VARCHAR(32),
	`Group` VARCHAR(32)
);


USE Jass;
DROP PROCEDURE IF EXISTS pro_auth_sigin;
DELIMITER //
CREATE PROCEDURE pro_auth_sigin(in user VARCHAR(50),in pass CHAR(32),in name VARCHAR(16),out flag boolean)
BEGIN
	IF EXISTS(SELECT `UID` FROM tb_auth_user WHERE tb_auth_user.User=user) THEN 
		set flag = false;
	ELSE
		INSERT INTO tb_auth_user(`User`,`Pass`,`Name`,`SiginTime`) VALUES(user,pass,name,curdate());
		set flag = true;
	END IF;
END//
DELIMITER ;

USE Jass;
DROP PROCEDURE IF EXISTS pro_auth_login;
DELIMITER //
CREATE PROCEDURE pro_auth_login(in in_user VARCHAR(50),in in_pass CHAR(32),in in_ip int,out out_uid INT UNSIGNED,out out_name VARCHAR(16),out out_level INT UNSIGNED)
BEGIN
	SELECT `UID`,`Name`,`Level` INTO out_uid,out_name,out_level FROM `tb_auth_user` WHERE `User`=in_user AND `Pass`=in_pass;
	IF(out_uid!=0)THEN
		UPDATE `tb_auth_user` Set `LoginIP`=in_ip WHERE `UID`=out_uid;
	ELSE
		SET out_uid = 0;
		SET out_name = '';
		SET out_level = 0;
	END IF;
	SELECT out_uid,out_name,out_level;
END//
DELIMITER ;
