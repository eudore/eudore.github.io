#!/usr/bin/env python  
# -*- coding: utf-8 -*- 

import sys
import json
import datetime
import logging
import hashlib
import mysql.connector

reload(sys)
sys.setdefaultencoding('utf-8')
logging.basicConfig(filename='/data/web/logs/py-db.log',format='%(asctime)s - %(name)s - %(levelname)s - %(module)s: %(message)s',datefmt='%Y-%m-%d %H:%M:%S %p', level=20)

def getDB():
	config = {
		'user': 'root',
		'password': '',
		'host': 'localhost',
		'database': 'Jass',
		'charset':'utf8'
	}
	return mysql.connector.connect(**config)

def AddIndex(path):
	db=getDB()
	cursor=db.cursor()
	try:
		cursor.callproc('pro_AddIndex',[path])
		db.commit()
	except:
		db.rollback()
		print('rollback')
	finally:
		cursor.close()
		db.close()

def GetIndex(path):
	results = {}
	try:
		db=getDB()
		cursor=db.cursor()
		cursor.callproc('pro_GetIndex',[path])
		db.commit()
		for result in cursor.stored_results():
			for i in result.fetchall():
				results[i[0]]=i[1]
	except:
		db.rollback()
		print('rollback')
	finally:
		cursor.close()
		db.close()
	return json.dumps({path:results}, ensure_ascii=False)

def GetContent(path):
	result = {}
	try:
		db=getDB()
		cursor=db.cursor(buffered=True)
		sql="select EditTime,Content from tb_Content WHERE MD5='%s';" % (path)
		cursor.execute(sql)
		result=cursor.fetchone()
	except:
		db.rollback()
		print('rollback')
	finally:
		cursor.close()
		db.close()
	return json.dumps({'time':result[0].strftime('%Y-%m-%d'),'content':result[1]}, ensure_ascii=False) 
	#return json.dumps({path:{'time':result[0].strftime('%Y-%m-%d'),'data':result[1]}}, ensure_ascii=False) 

def SetContent(path,data):
	logging.info("SetContent-"+path)
	try:
		db=getDB()
		cursor=db.cursor(buffered=True)
		cursor.execute("UPDATE tb_Content SET Content=%s WHERE MD5=%s;",(data,path))
#		cursor.execute("UPDATE tb_Content SET Content='{}' WHERE MD5='{}';".format(data,path))
		db.commit()
	except:
		db.rollback()
		return json.dumps({"result": 1}, ensure_ascii=False) 
	finally:
		cursor.close()
		db.close()
	return json.dumps({"result": 0}, ensure_ascii=False) 
	

def GetHash(path):
	m = hashlib.md5()   
	m.update(path.replace('/',chr(0)))
	return m.hexdigest()


def Test():
	db=getDB()
	cursor=db.cursor()
	try:
		cursor.callproc('pro_Sigin',('web','py','name',''))
		db.commit()
		cursor.close()
		cursor=db.cursor()
		cursor.execute('select @_pro_Sigin_3')
		data = cursor.fetchone()
		print "Sigin_3 : %s"%data
	except:
		db.rollback()
	finally:
		cursor.close()
	
	cursor=db.cursor()
	cursor.execute("SELECT VERSION()")
	data = cursor.fetchone()
	print "Database version : %s"%data
	cursor.close()
	db.close()




def InitCJ():
	file = open("/data/webpy/data/InitCJ.txt")
	for line in file.readlines():
		AddIndex(line.strip('\n'))
	file.close()

def Init():
#	InitCJ()
	AddIndex('Jass/BJ/Timer');
	AddIndex('Jass/NL/NLPool');
	AddIndex('Jass/NL/NLLink');
	AddIndex('Jass/NL/NLTimer');
	AddIndex('Jass/NL/NLEvent');
	AddIndex('Jass/NL/NLBonus');
	AddIndex('Jass/NL/NLItem');

#print(GetIndex('Jass/CJ/Unit'))
#Init()
print(GetContent('6853fcdd2bb428a3160c59079175bf42'))
