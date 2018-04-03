#!/bin/bash

for domain in $*
do
	d1=$(echo | openssl s_client -servername $domain -connect $domain:443 2>/dev/null | openssl x509 -noout -dates | grep notAfter | cut -d= -f2)
	t1=$(date +%s -d "$d1")
	t2=$(date +%s -d "$((date))")
	d2=`echo $(($t1-$t2))`
	#d3=$(awk "BEGIN{print $d2/86400 }") 
	d3=$(( $d2 / 86400 ))
	echo $d3
done

