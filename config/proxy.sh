#!/bin/bash
BasePath=$(dirname `cd $(dirname $0);pwd`)
_install_ssr(){
	tar xzvf $BasePath/download/shadowsocksr.tar.gz -C /usr/local/
	echo "/usr/local/shadowsocksr/logrun.sh" >> /etc/rc.local 
	chmod a+x /etc/rc.d/rc.local
	cd /usr/local/shadowsocksr
	bash /usr/local/shadowsocksr/initcfg.sh
	sed -i 's/^\(API_INTERFACE =\).*\(#.*\)$/\1 \"mudbjson\" \2/g' /usr/local/shadowsocksr/userapiconfig.py
	python /usr/local/shadowsocksr/mujson_mgr.py -a -u 8080 -p 8080 -k c152 -m aes-256-cfb -O auth_chain_a -o tls1.2_ticket_auth
	python /usr/local/shadowsocksr/mujson_mgr.py -a -u 8081 -p 8081 -k 63e6 -m aes-256-cfb -O auth_chain_a -o tls1.2_ticket_auth
	python /usr/local/shadowsocksr/mujson_mgr.py -a -u 8082 -p 8082 -k b2c1 -m aes-256-cfb -O auth_chain_a -o tls1.2_ticket_auth
	python /usr/local/shadowsocksr/mujson_mgr.py -a -u 8083 -p 8083 -k 1344 -m aes-256-cfb -O auth_chain_a -o tls1.2_ticket_auth
	python /usr/local/shadowsocksr/mujson_mgr.py -a -u 8084 -p 8084 -k 1033 -m aes-256-cfb -O auth_chain_a -o tls1.2_ticket_auth
	python /usr/local/shadowsocksr/mujson_mgr.py -a -u 8085 -p 8085 -k r3k8d -m aes-256-cfb -O auth_chain_a -o tls1.2_ticket_auth
	python /usr/local/shadowsocksr/mujson_mgr.py -a -u 8086 -p 8086 -k rgy48 -m aes-256-cfb -O auth_chain_a -o tls1.2_ticket_auth
	#/usr/local/shadowsocksr/stop.sh
	/usr/local/shadowsocksr/logrun.sh
	netstat -atunlp | grep 'python2.7'
}
_config_v2(){
	systemctl stop firewalld
    systemctl disable firewalld
    setenforce 0
	cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
	cp /usr/local/v2ray-v3.6-linux-64/systemd/v2ray.service /usr/lib/systemd/system
	systemctl enable v2ray
	systemctl restart v2ray
}
_config_prot(){
	for i in `seq 8080 8089`
	do
		iptables -A INPUT -p tcp --dport $i
		iptables -A OUTPUT -p tcp --sport $i
	done
	iptables-save > /etc/sysconfig/iptables
	iptables -n -v -L -t filter -x
}
_config_vultr(){
	yum -y install net-tools
	systemctl stop firewalld
	systemctl disable firewalld
	setenforce 0
	cp /usr/local/shadowsocksr/ssr.service /usr/lib/systemd/system
	systemctl start ssr
	systemctl enable ssr
	rpm -ivh http://soft.91yun.org/ISO/Linux/CentOS/kernel/kernel-3.10.0-229.1.2.el7.x86_64.rpm --force
	wget -N --no-check-certificate https://github.com/91yun/serverspeeder/raw/master/serverspeeder.sh && bash serverspeeder.sh
}
#_install_ssr
