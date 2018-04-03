#!/usr/bin/env python  
# -*- coding: utf-8 -*- 

# AccessKeyID：LTAIoq1zEjIUpHUN
# AccessKeySecret：CZ8X8rq0s7p1qjFiDba5GTIeoQJ0vO

import os
import sys
import time
import multiprocessing,re
import logging
import base64,hashlib,oss2
import platform,signal,pyinotify

# const config
AccessKeyID = 'LTAIoq1zEjIUpHUN'
AccessKeySecret = 'CZ8X8rq0s7p1qjFiDba5GTIeoQJ0vO'
endpoint = 'http://oss-cn-hongkong-internal.aliyuncs.com'
bucket_name = 'wejass'


logging.basicConfig(filename='/data/web/logs/sync-oss.log',format='%(asctime)s - %(name)s - %(levelname)s - %(module)s: %(message)s',datefmt='%Y-%m-%d %H:%M:%S %p', level=20)


# 计算文件MD5
def calculate_file_md5(file_name, block_size=64 * 1024):
    with open(file_name, 'rb') as f:
        md5 = hashlib.md5()
        while True:
            data = f.read(block_size)
            if not data:
                break
            md5.update(data)
    return base64.b64encode(md5.digest())

# 计算文件crc64
def calculate_file_crc64(file_name, block_size=64 * 1024, init_crc=0):
    with open(file_name, 'rb') as f:
        crc64 = oss2.utils.Crc64(init_crc)
        while True:
            data = f.read(block_size)
            if not data:
                break
            crc64.update(data)
    return crc64.crc
		

class SyncOss():
	__count__ = [0,0,0,0]
	INSERT,UPDATE,DELETE,FORMER=0,1,2,3
	def __init__(self,root='',ross='',sync='',nosync=''):
		# authorization
		self.AccessKeyID = 'LTAIoq1zEjIUpHUN'
		self.AccessKeySecret = 'CZ8X8rq0s7p1qjFiDba5GTIeoQJ0vO'
		self.endpoint = 'http://oss-cn-hongkong-internal.aliyuncs.com'
		self.bucket_name = 'wejass'
		self.getbucket()
		# directory
		self.root = os.path.abspath(root)
		self.ross = ross
		self.sync = sync
		if not self.root:
			self.root = os.getcwd()
		else:
			os.chdir(self.root)
		if not self.sync:
			self.sync = os.listdir(self.root)
		self.sync = list(set(self.sync)-set(nosync))


	# Count the number of operations objects 
	def count(self,type,c=1):
		self.__count__[type]=self.__count__[type]+c

	# request resulte into logging
	def log(self,result,info):
		msg='{0} - http status: {1} - data: {2}'.format(info,result.status,result.headers)
		if result.status == 200:
			logging.info(msg)
		elif result.status >= 400:
			logging.error(msg)
		else:
			logging.warning(msg)

	def getbucket(self):
		self.auth = oss2.Auth(self.AccessKeyID,self.AccessKeySecret)
		self.bucket = oss2.Bucket(self.auth, self.endpoint, self.bucket_name)

	# Delete the local root does not exist path oss object
	def delete(self,path):
		localpath = os.path.join(self.root,path)
		osspath = os.path.join(self.ross,path)
		if (not os.path.exists(localpath)) and self.bucket.object_exists(osspath):
			self.bucket.delete_object(osspath)
			#self.log(self.bucket.delete_object(osspath),'delete the local {0}/{1} not exist {2}'.format(root,oss,oss))
			self.count(SyncOss.DELETE)
			return 0
		dellist = []
		for obj in oss2.ObjectIterator(self.bucket,prefix='{0}/'.format(osspath)):
			if not os.path.exists(os.path.join(self.root,obj.key) ):
				dellist.append(obj.key)
		if dellist:
			self.bucket.batch_delete_objects(dellist)
			#self.log(self.bucket.batch_delete_objects(dellist),'delete the local {0} not exist {1}'.format(root,dellist))
			self.count(SyncOss.DELETE,c=len(dellist))

	# Upload file from local to oss
	def upload(self,path):
		localpath = os.path.join(self.root,path)
		osspath = os.path.join(self.ross,path)
		result = None
		encode_md5 = calculate_file_md5(localpath)
		if self.bucket.object_exists(osspath):
			crc64 = calculate_file_crc64(localpath)
			meta = self.bucket.get_object_meta(osspath)
			if str(crc64) != meta.headers['x-oss-hash-crc64ecma']:
				logging.info(osspath)
				logging.info(localpath)
				result = self.bucket.put_object_from_file(osspath, localpath, headers={'Content-MD5': encode_md5})
				#self.log(result,'put object update from file: {0} to OSS:{1}'.format(localpath,osspath))
				self.count(SyncOss.UPDATE)
			else:
				self.count(SyncOss.FORMER)
		else:
			result = self.bucket.put_object_from_file(osspath, localpath, headers={'Content-MD5': encode_md5})
			#self.log(result,'put object inert from file: {0} to OSS:{1}'.format(localpath,osspath))
			self.count(SyncOss.INSERT)
		return result

	# Upload dir from local to oss
	def insert(self,path):
		if os.path.isdir(os.path.join(self.root,path)):
			for obj in os.listdir(os.path.join(self.root,path)):
				self.insert(os.path.join(path,obj))
		else:
			self.upload(path)

	# Sync local to Oss
	def syncall(self,sync=[],nosync=[]):
		if not sync:
			sync = self.sync
		for obj in sync:
			if obj not in nosync:
				self.delete(obj)
				self.insert(obj)
		print('Sync {0}/{1} to {2}/{3}/{4}'.format(self.root,self.sync,self.endpoint,self.bucket_name,self.ross))
		for i,name in enumerate(['INSERT','UPDATE','DELETE','FORMER']):
			print('{0} file: {1}'.format(name,self.__count__[i]))



class OssInotify(pyinotify.ProcessEvent):
	def __init__(self, syncoss,pidfile):
		self.sync = syncoss
		self.pidfile = pidfile

	def process_IN_CREATE(self, event):
		logging.info('event: create- path: {0}/{1}'.format(event.path,event.name))
		self.sync.insert(os.path.relpath(event.pathname,self.sync.root))
		logging.info('over')

	def process_IN_DELETE(self, event):
		logging.info('event: delete - path: {0}/{1}'.format(event.path,event.name))
		self.sync.delete(os.path.relpath(event.pathname,self.sync.root))
		logging.info('over')

	def process_IN_MODIFY(self, event):
		logging.info('event: modify - path: {0}/{1}'.format(event.path,event.name))
		self.sync.insert(os.path.relpath(pyinotify,self.sync.root))
		logging.info(event.pathname)
		logging.info('over')

	def process_IN_MOVED_FROM(self, event):
		logging.info('event: move from - path: {0}/{1}'.format(event.path,event.name))

	def process_IN_MOVED_TO(self, event):
		logging.info('event: move to - path: {0}/{1}'.format(event.path,event.name))

	def __del__( self ):
		os.remove(self.pidfile)


def status(pidfile):
	if os.path.exists(pidfile):
		with open(pidfile) as pid:
			pid = pid.read().strip('\n')
		try:
			os.kill(int(pid),0)
		except OSError:
			os.remove(pidfile)
			return 0
		else:
			return pid
	return 0

def stop(pidfile):
	pid = status(pidfile)
	if pid:
		os.kill(int(pid), signal.SIGKILL)
		os.remove(pidfile)

def start(path,sync,pidfile): 
	if not status(pidfile):
		flags = pyinotify.IN_CREATE | pyinotify.IN_DELETE | pyinotify.IN_MODIFY | pyinotify.IN_MOVED_FROM | pyinotify.IN_MOVED_TO
		wm = pyinotify.WatchManager()
		for i in sync:
			wm.add_watch(os.path.join(path,i), flags, rec=True)
		notifier = pyinotify.Notifier(wm, OssInotify(SyncOss(root=path, sync=sync),pidfile))
		notifier.loop(daemonize=True,pid_file=pidfile)


def isInotify():
	with open('/boot/config-{0}'.format(platform.release())) as f:
		for line in f:
			if line.startswith('CONFIG_INOTIFY_USER='):
				return (line.split('=')[1]=='y\n')
	return False


cmds = []
envs = {'root':'', 'ross':'', 'sync':'', 'nosync':'', 'pidfile':'/run/sync-oss.pid'}
envs['root'] = '/data/web/static/'
envs['sync'] = ['js','css']
#envs['sync'] = ['css','js','images','favicon.ico']

def main():
	if not isInotify():
		logging.error('The system does not support Inotify')
		return 1
	for i in sys.argv[1:]:
		kv = re.findall('--(.*)=(.*)',i)
		if kv:
			kv = kv[0]
			if kv[0] in ['sync','nosync']:
				envs[kv[0]] = kv[1].split(',')
			else:
				envs[kv[0]] = kv[1]
		elif i in ['start', 'stop', 'status', 'sync', 'args']:
			cmds.append(i)
	if not cmds:
		print('command: start/stop/status/sync')
		print('args: root default current directory')
		print('args: ross default null')
		print('args: sync default current directory all list')
		print('args: nosync default null')
		print('args: pidfile default /run/sync-oss.pid')
		cmds.append('status')
	if cmds[-1] != 'status':
		cmds.append('status')
	for cmd in cmds:
		if cmd == 'start':
			root='/data/web/html'
			sync=['css','js','images','favicon.ico']
			multiprocessing.Process(target = start, args = (envs['root'], envs['sync'],envs['pidfile'],)).start()
		elif cmd == 'stop':
			stop(envs['pidfile'])
		elif cmd == 'status':
			pid = status(envs['pidfile'])
			if pid:
				print('Sync Oss is running pid = {0}'.format(pid))
			else:
				print('Sync Oss is stop')
		elif cmd == 'sync':
			SyncOss(root=envs['root'], ross=envs['ross'], sync=envs['sync']).syncall()
		elif cmd == 'args':
			print('cmds = {0}'.format(cmds))
			print('envs = {0}'.format(envs))
		time.sleep(0.2)

if __name__ == "__main__":
	main()
