#!/bin/bash

# centos6: sudo yum install -y jwhois
# cenros7: sudo yum install -y whois

for domain in $*
do
	d1=`/usr/bin/whois $domain | grep Expiry | cut -b 26-35`
	while [ -z $d1 ];	# 万网注册域名格式不同、可能获取信息失败
	do
		d1=`/usr/bin/whois $domain | grep Expiration | cut -b 18-27`
	done
	t1=$(date +%s -d "$d1")
	t2=$(date +%s -d "$((date))")
	d2=`echo $(($t1-$t2))`
	#d3=$(awk "BEGIN{print $d2/86400 }")	#计算实数 
	d3=$(( $d2 / 86400 ))
	echo $d3
done
