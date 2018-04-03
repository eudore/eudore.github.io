#!/usr/bin/env python 
# -*- coding: utf-8 -*-
import sys
import json
import logging
import smtplib
from email.mime.text import MIMEText

def send(to,sub,mesg,type="plain"):
	_user = "wt@shitou.com"
	_pwd  = "3ROCKS.CN"
	_to   = to
	msg = MIMEText(mesg,type,'utf-8')
	msg["Subject"] = sub
	msg["From"]    = _user
	msg["To"]      = _to
	try:
		s = smtplib.SMTP_SSL("smtp.mxhichina.com", 465)
		s.login(_user, _pwd)
		s.sendmail(_user, _to, msg.as_string())
		s.quit()
		result = '邮件发送成功！'
	except smtplib.SMTPConnectError, e:
		result = '邮件发送失败，连接失败:', e.smtp_code, e.smtp_error
	except smtplib.SMTPAuthenticationError, e:
		result = '邮件发送失败，认证错误:', e.smtp_code, e.smtp_error
	except smtplib.SMTPSenderRefused, e:
		result = '邮件发送失败，发件人被拒绝:', e.smtp_code, e.smtp_error
	except smtplib.SMTPRecipientsRefused, e:
		result = '邮件发送失败，收件人被拒绝:', e.smtp_code, e.smtp_error
	except smtplib.SMTPDataError, e:
		result = '邮件发送失败，数据接收拒绝:', e.smtp_code, e.smtp_error
	except smtplib.SMTPException, e:
		result = '邮件发送失败, ', e.message
	except Exception, e:
		result = '邮件发送异常, ', str(e)
	
	print(result)
	r={"to": to,"Subject": sub,"Message": mesg,"Result": result}
	logging.basicConfig(filename='/var/log/zabbix/zabbix_alert.log',format='%(asctime)s - %(name)s - %(levelname)s - %(module)s:  %(message)s',datefmt='%Y-%m-%d %H:%M:%S %p',  level=10)    
	logging.info(json.dumps(r, encoding="UTF-8", ensure_ascii=False))

if __name__ == "__main__":
	send(sys.argv[1],sys.argv[2],sys.argv[3])