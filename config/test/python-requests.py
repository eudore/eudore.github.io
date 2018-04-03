# -*- coding: UTF-8 -*-

import json
import requests
import pandas as pd


url = 'http://18.16.202.149:9200'

def input():
	df=pd.DataFrame(pd.read_excel(io=u"D:\\Workspaces\\data\\部分线上服务器zabbix监控.xlsx"),columns=['hostname',u'公网IP',u'内网IP'])
	for i,line in df.iterrows():
		print(line.to_json(force_ascii=False))
		r = requests.post("%s/%s/%s" % (url, 'info/local',i),data=line.to_json(force_ascii=False))
		if r.status_code == requests.codes.ok:
			print(r.content)

def search(strs,path='info/local',action='_search'):
	my_paramm={"q":strs,"size":10}
	try:
		r=requests.get("%s/%s/%s" % (url, path, action), params=my_paramm,timeout=60)
		print(r.url)
		print(str(r.elapsed.microseconds)+'μs')
		print(r.content)
		print('')
		for i in json.loads(r.content)["hits"]["hits"]:
			print(json.dumps(i["_source"],ensure_ascii=False))
	except requests.Timeout:
		pass
	except:
		print('Error')

search("*gateway*")
# url headers content text status_code history requests.codes.ok
