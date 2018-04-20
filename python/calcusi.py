#!/usr/bin/env python  
# -*- coding: utf-8 -*-  

# Calculate Subresource Integrity

import os
import sys
import re
import base64
import hashlib
import requests


root='/data/web/static'
def checksum(data):
	return checksha512(data)

def checksha256(data):
		digest_sha256 = hashlib.sha256(data).digest()
		hash_base64 = base64.b64encode(digest_sha256).decode()
		return 'sha256-{}'.format(hash_base64)

def checksha512(data):
		digest_sha512 = hashlib.sha512(data).digest()
		hash_base64 = base64.b64encode(digest_sha512).decode()
		return 'sha512-{}'.format(hash_base64)

def getdata(path):
	if re.match('https?://.*', path):
		r = requests.get("%s" % (path))
		if r.status_code == requests.codes.ok:
			return r.text.encode('utf-8')
	else:
		with open(path, 'rb') as f:
			return f.read()
	return None

def writeSI(path,data):
	i=0
	data.append(['',''])
	
	fileread = open(path,'r');
	lines = fileread.readlines();
	fileread.close();

	filewrite = open(path,'w');
	for line in lines:
		if line == data[i][0]:
			line = data[i][1]
			i=i+1
		filewrite.write(line);
	filewrite.close();

def readSI(file):
	base=os.path.dirname(file)
	fileread = open(file,'r');
	lines = fileread.readlines();
	fileread.close();

	diff=[]
	for line in lines:
		r = re.match('\s*<script.*src=[\"\'](\S*\.js)[\"\'].*></script>', line) or re.match('\s*<link.*href=[\"\'](\S*\.css)[\"\'].*>', line)
		if r :
			if re.match('https?://.*', r.group(1)):
				filepath = r.group(1)
			elif re.match('^/.*', r.group(1)):
				filepath = root + r.group(1)
			else:
				filepath = os.path.join(base,r.group(1))
			data = getdata(filepath)
			if data==None:
				continue
			data = checksum(data)
			r2 = re.match('.*\s+integrity=[\"\'](\S*)[\"\'].*', line)
			if r2:
				if data == r2.group(1):
					continue 
				data = re.sub('\s+integrity=[\"\']\S*[\"\']', ' integrity=\''+data+'\'', line)
			else:
				data = re.sub('=[\"\']+{}[\"\']+'.format(r.group(1)),'=\'{}\' integrity=\'{}\''.format(r.group(1),data), line)
			diff.append([line,data])
	if diff:
		for i in diff:
			print("{}  ->  {}".format(i[0].strip(),i[1].strip()) )
			pass
	return diff


def integrity(path):
	if os.path.isdir(path):
		for obj in os.listdir(path):
			integrity(os.path.join(path,obj))
	else:
		if os.path.basename(path).endswith(".html"):
			readSI(path)

def main():
	if len(sys.argv)>1 and os.path.exists(sys.argv[1]):
		p=sys.argv[1]
		d=readSI(p)
		writeSI(p, d)

if __name__ == "__main__":
	main()
