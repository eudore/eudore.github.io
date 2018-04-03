#!/usr/bin/env python  
# -*- coding: utf-8 -*- 

import os
import sys
import pyinotify
import platform
import threading
import logging

# config
pidfile = '/tmp/pyinotify.pid'


logging.basicConfig(filename='/data/web/logs/oss-sync.log',format='%(asctime)s - %(name)s - %(levelname)s - %(module)s: %(message)s',datefmt='%Y-%m-%d %H:%M:%S %p', level=10)

def isInotify():
	with open('/boot/config-{0}'.format(platform.release())) as f:
		for line in f:
			if line.startswith('CONFIG_INOTIFY_USER='):
				return (line.split('=')[1]=='y\n')
	return False

class Cdn(pyinotify.ProcessEvent):
	def process_IN_CREATE(self, event):
		logging.debug('event: create - path: {0}/{1}'.format(event.path,event.name))
	def process_IN_DELETE(self, event):
		logging.debug('event: detele - path: {0}/{1}'.format(event.path,event.name))
	def process_IN_MODIFY(self, event):
		logging.debug('event: modify - path: {0}/{1}'.format(event.path,event.name))
	def process_IN_MOVED_FROM(self, event):
		logging.debug('event: move from - path: {0}/{1}'.format(event.path,event.name))
	def process_IN_MOVED_TO(self, event):
		logging.debug('event: move to - path: {0}/{1}'.format(event.path,event.name))

def main():
	if not isInotify():
		logging.error('The system does not support Inotify')
		return 1
	flags = pyinotify.IN_CREATE | pyinotify.IN_DELETE | pyinotify.IN_MODIFY | pyinotify.IN_MOVED_FROM | pyinotify.IN_MOVED_TO
	wm = pyinotify.WatchManager()
	wm.add_watch('/data/web/html/js', flags, rec=True)
	wm.add_watch('/data/web/html/css', flags, rec=True)
	wm.add_watch('/data/web/html/images', flags, rec=True)
	notifier = pyinotify.Notifier(wm, Cdn())
	notifier.loop(daemonize=True,pid_file=pidfile)


if __name__ == "__main__":
	main()