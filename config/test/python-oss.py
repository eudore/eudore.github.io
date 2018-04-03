#!/usr/bin/env python  
# -*- coding: utf-8 -*- 

# AccessKeyID：LTAIoq1zEjIUpHUN
# AccessKeySecret：CZ8X8rq0s7p1qjFiDba5GTIeoQJ0vO


import oss2
from itertools import islice


def Init():
	auth = oss2.Auth('LTAIoq1zEjIUpHUN','CZ8X8rq0s7p1qjFiDba5GTIeoQJ0vO')
	bucket = oss2.Bucket(auth, 'http://oss-cn-hongkong-internal.aliyuncs.com', 'wejass')


	print([b.name for b in oss2.BucketIterator(oss2.Service(auth, 'http://oss-cn-hongkong-internal.aliyuncs.com'))])

	for b in oss2.ObjectIterator(bucket,prefix='css/'):
		print(b.key)

Init()
