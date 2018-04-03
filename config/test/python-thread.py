#!/usr/bin/python
# -*- coding: UTF-8 -*-
 
import threading
import time
 
class myThread (threading.Thread):
    def __init__(self, name, delay):
        threading.Thread.__init__(self)
        self.threadName = name
        self.delay = delay
    def run(self):
        count = 0
        while count < 5:
            time.sleep(self.delay)
            count += 1
            print "%s: %s" % ( self.threadName, time.ctime(time.time()) )


thread1 = myThread("Thread-1", 1)
thread2 = myThread("Thread-2", 2)
 
# 开启线程
thread1.start()
thread2.start()