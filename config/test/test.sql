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