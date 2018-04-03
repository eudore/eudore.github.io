#!/usr/bin/env python  
# -*- coding: utf-8 -*-  

import sys
import web
import json
from web.session import Session
from web.wsgiserver import CherryPyWSGIServer
#from web.wsgiserver.ssl_builtin import BuiltinSSLAdapter

from python 	import MariaDB	as DB
from python 	import mail


reload(sys)
sys.setdefaultencoding("utf-8")

#CherryPyWSGIServer.ssl_certificate = "/etc/nginx/www.wejass.com/www.wejass.com.pem"
#CherryPyWSGIServer.ssl_private_key = "/etc/nginx/www.wejass.com/www.wejass.com.key"


urls = (
	'/auth','OAuth',
	'/note/(.*)','Note',
	'/file','File',
	'/py/(.*)','Tasks'
	)


web.config.debug = True
app = web.application(urls, globals())
application = app.wsgifunc()
render = web.template.render('templates/')

class OAuth:
	def GET(self):
		i = web.input();
		if i.type == 'LOGOUT':
			raise web.seeother('/')
	def POST(self):
		i = web.input();
		try:
			if i.get('type') == 'LOGIN':
				userType = DB.Login(i.get('name'), i.get('passwd'))
				if userType=='user':
					url ='/index.html?wejass=welcomeWeJass';
				else:
					url='/login.html'
				web.setcookie('token',userType, secure=True)
				raise web.seeother(url);
			elif i.type == 'CREATE':
				return DB.Sigin(i.get('name'),i.get('passwd'),i.get('code'))
			elif i.type == 'LOGOUT':
				web.setcookie('token',expires=-1, secure=True)
				return '0'
		except BaseException,e:
			return e.message;

class Note:
	def GET(self,path):
		i = web.input()
		#try:

		web.header('Content-Type','text/html;charset=utf-8', unique=True)
		return DB.GetContent(path)
		#return '{"result":true,"edittime": 1516542625.267140,"content": \'<div>NL Library<br>于2017年2月开始编写，其目标是为了实现一套高效的系统，研究各模块间相互协同、整合、优化，通常许多演示都难以加入到应用，各种相互冲突导致难以应用，通过统一定义处理，明确触发顺序避免可能存现的bug，不是复杂问题在环境下通过相互依赖可以简单实现，并拥有大量的新思路比一些通用的Ui都要高效。期间迭代过大量的版本，从哈希表存储转变纯数组存储就导致全部代码重构。在新的版本下，核心基于纯数组存储，，由于数组和哈希表5倍读写速度差异，使用定制的专用存储结构，封装的堆和封装链表进行存储，部分由哈希表实现低效复杂，但由数组存储可以轻松实现。使用的存储结构基本都是只能在特定的情况下使用，但是保证的高效，其他结构也只需要对模型进行小部分变动，目前一共设计出了两个存储模型。<br><br><a href="NL/NLTimer" target="_blank">Timer</a></div><div><a href="NL/NLEvent" target="_blank">Event</a><br><a href="NL/NLBuff" target="_blank">Buff</a><br><br>模型一基于Widget的自定义值存储索引<br>该模型只能用于Widget的子类单位、物品等对象使用，但单位、物品基本是最常用的，虽然该模型有限制但是非常实用。目前用于Event模块<br>在给对象注册时堆分配一个唯一索引，并讲索引保存但自定义值中，Widget的子类才有Get/Set****UserData这两个方法操作自定义值，其他对象只能使用无效属性位置去保存索引，读写属性时从自定义获得索引，用索引去操作对应的数组和数组偏移空间的数据，如果需要保存动态数量的数据，则保存数量链表的头，具体数据放到链表即可。<br><br>模型二使用Type与Index转化<br>该模型通常用于只读数据，对于可写数据不方便维护，不值得使用，目前用于Item模块存储使用数据<br>通过实现一个二分查找的数组的函数进行快速定位对于Type的索引位置，以200物品数据为例进行数据读取平均比较6.675次即可，使用遍历需要比较100次；获得索引为第一步，然后对索引加是偏移位获得对应类别数据保存是起始索引值，同另外一个Len数组同索引位置获取到数据数量，然后去对应的数组的位置读取数据即可。<br>Item模块的数据都是有lua处理好的，lua进行索引转换数据排序，然后输出对应的数据到指定位置，在Jass写入数据没有Lua中简单，而且大量数据需要Lua去读取处理，最后由Lua统一保存整理输出即可。<br><br>中心计时器分离时间滚动和任务处理成双计时器，完美结构瞬间计时器0.01秒等待问题。<br>Evnt模块其移除清理函数可以运行动态注册事件，优先级概念完美解决了触发顺序，基于数组得以高效实现。<br>Bonus使用Event模块而简单实现，也可以进行任何修改、扩展<br>Item防止切背包单位死亡加入主从背包，加入物编识别属性减少物编工作量，添加描述就自动化完成属性工作，特效部分同理，合成和套装使用lua使得数据直观易编辑。</div>\'}'
		if i.get('type') == 'Index':
			web.header('content-type','text/json')
			return DB.GetIndex(path)
		elif i.get('type') == 'Content':
			web.header('Content-Type','text/html;charset=utf-8', unique=True)
			return DB.GetContent(path)
		path=DB.GetHash(path)
		return render.note(json.loads(DB.GetContent(path))[path]["data"]);
		#except:
		#	return name;
	def POST(self,name):
		i = web.input()
		if i.get('type') == 'Index':
			pass
		elif i.get('type') == 'Content':
			return DB.SetContent(name,web.data())
		return "templates post:"+name

class File:
	def GET(self):
		return """<html><head><meta charset="UTF-8" /></head><body>
			<form method="POST" enctype="multipart/form-data" action="">
			<input type="file" name="myfile" />
			<br/>
			<input type="submit" />
			</form>
			</body></html>"""

	def POST(self):
		x = web.input(myfile={})
		filedir = '/data/web/download/test'
		if 'myfile' in x:
			filepath = x.myfile.filename.replace('\\', '/')
			filename = filepath.split('/')[-1]  # the filename with extension
			fout = open(filedir + '/' + filename, 'wb')
			fout.write(x.myfile.file.read())
			fout.close()
		web.debug(x['myfile'].filename) # 这里是文件名
		web.debug(x['myfile'].value) # 这里是文件内容
		web.debug(x['myfile'].file.read()) # 或者使用一个文件对象
		raise web.seeother('/File')

class Tasks:
	def GET(self,name):
		i = web.input()
		if name == 'test':
			return web.config.debug;
		elif name == "mail":
			mail.send(i)
			return name
		elif name == 'json':
			web.header('content-type','text/json')
			data={'b301ac313c5fb6ca614954045951efd5':{'GetUnitx':'0','GetUnitY':'0'},'21f7583c9a2d735cfa7a5d15fcd7c9d3':{'CJ':'1','BJ':'1'} }
			return json.dumps(data)
		else:
			return "py takes "+name
	def POST(self,name):
		i = web.input()
	#	if name == "mail":
		mail.send(i)
	#	return self.GET()


class RequestHandler():
	def POST():
		return web.data()
def notfound():
	return web.notfound("Sorry, the page you were looking for was not found.")

if __name__ == "__main__":
	app.run()
