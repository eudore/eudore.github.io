#!/bin/bash
_run_server(){
	docker run -d --restart=always -p 2379:2379 -p 2380:2380 -v /data/docker/etcd/node1:/node1.etcd -e ETCD_NAME=node1 -e ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379 -e ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 etcd
	docker run -d --restart=always -p 3306:3306 -v /data/docker/mariadb:/var/lib/mysql/ mariadb
	docker run -d --restart=always -p 5432:5432 -v /data/docker/postgresql:/var/lib/postgresql/data postgres
	docker run -d --restart=always -p 12001:11211 memcached memcached -m 64 -U 0
	docker run -d --restart=always -p 12002:11211 memcached memcached -m 64 -U 0
	docker run -d --restart=always -p 12003:11211 memcached memcached -m 64 -U 0
	# docker run -d --restart=always -p 8020:80 -e GOROOT=$GOROOT -e GOPATH=$GOPATH -v $GOROOT:$GOROOT godoc:go1.10
	# gogs	
	docker run -d --restart=unless-stopped -p 3010:8080 -v /data/docker/ranchermy/mysql:/var/lib/mysql rancher/server

	docker run --restart=always --name gogsdb -v /data/docker/gogsdb:/var/lib/postgresql/data -e POSTGRES_USER=gogs -e POSTGRES_PASSWORD=gogs -e POSTGRES_DB=gogs -d postgres
	docker run -d --restart=always --link gogsdb -p 1022:22 -p 3000:3000 -v /data/docker/gogs:/data gogs/gogs
	# sonar
	# docker run --name postgresqldb --restart=always -v /data/docker/postgresql:/var/lib/postgresql/data -e POSTGRES_USER=sonar -e POSTGRES_PASSWORD=sonar -e POSTGRES_DB=sonar -d postgres
	# docker run --name sonarqube --restart=always --link postgresqldb -e SONARQUBE_JDBC_USERNAME=sonar -e SONARQUBE_JDBC_PASSWORD=sonar -e SONARQUBE_JDBC_URL=jdbc:postgresql://postgresqldb:5432/sonar -v /data/docker/sonar/data:/opt/sonarqube/data -v /data/docker/sonar/extensions:/opt/sonarqube/extensions -p 9000:9000 -d sonarqube
}
_install_nginx(){
	id nginx || useradd nginx
	groups nginx || groupadd -g nginx nginx
	mkdir -pv /var/lib/nginx/tmp
	yum -y install gcc gcc-c++ git wget automake pcre pcre-devel zlib-devel openssl openssl-devel
	cd /tmp
	wget -O nginx-ct.zip -c https://github.com/grahamedgecombe/nginx-ct/archive/v1.3.2.zip
	unzip nginx-ct.zip
	git clone -b tls1.3-draft-18 --single-branch https://github.com/openssl/openssl.git openssl
	wget -P /tmp/ http://nginx.org/download/nginx-1.13.0.tar.gz
	tar axf /tmp/nginx-1.13.0.tar.gz
	#sed -i '/^#define NGINX_VERSION /s/1.8.0/1.0.0/g;/^#define NGINX_VER /s/nginx/wejass/g' /tmp/nginx-1.8.0/src/core/nginx.h
	#sed -i '/^tatic char ngx_http_server_string/s/nginx/wejass/g' /tmp/nginx-1.8.0/src/http/ngx_http_header_filter_module.c
	#sed -i '/^"<hr><center>nginx/s/nginx/wejass/g' /tmp/nginx-1.8.0/src/http/ngx_http_special_response.c
	cd /tmp/nginx-1.13.0
#	./configure --prefix=/usr/share/nginx --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --http-client-body-temp-path=/var/lib/nginx/tmp/client_body --http-proxy-temp-path=/var/lib/nginx/tmp/proxy --http-fastcgi-temp-path=/var/lib/nginx/tmp/fastcgi --http-uwsgi-temp-path=/var/lib/nginx/tmp/uwsgi --http-scgi-temp-path=/var/lib/nginx/tmp/scgi --pid-path=/run/nginx.pid --lock-path=/run/lock/subsys/nginx --user=nginx --group=nginx --with-file-aio --with-http_ssl_module --with-http_realip_module --with-http_addition_module --with-http_auth_request_module --with-http_sub_module --with-http_dav_module --with-http_flv_module --with-http_mp4_module --with-http_gunzip_module --with-http_gzip_static_module --with-http_random_index_module --with-http_secure_link_module --with-http_degradation_module --with-http_stub_status_module --with-mail_ssl_module --with-pcre --with-pcre-jit
	./configure --prefix=/usr/share/nginx --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --http-client-body-temp-path=/var/lib/nginx/tmp/client_body --http-proxy-temp-path=/var/lib/nginx/tmp/proxy --http-fastcgi-temp-path=/var/lib/nginx/tmp/fastcgi --http-uwsgi-temp-path=/var/lib/nginx/tmp/uwsgi --http-scgi-temp-path=/var/lib/nginx/tmp/scgi --pid-path=/run/nginx.pid --lock-path=/run/lock/subsys/nginx --user=nginx --group=nginx --with-file-aio --with-openssl=../openssl --with-openssl-opt='enable-tls1_3' --with-http_ssl_module --with-http_v2_module --with-http_realip_module --with-http_addition_module --with-http_auth_request_module --with-http_sub_module --with-http_dav_module --with-http_flv_module --with-http_mp4_module --with-http_gunzip_module --with-http_gzip_static_module --with-http_random_index_module --with-http_secure_link_module --with-http_degradation_module --with-http_stub_status_module --with-mail_ssl_module --with-pcre --with-pcre-jit --add-module=../nginx-ct-1.3.2/ --add-module=../oauth-mem
	make && make install
}
_install_nginx1.5() {
	id nginx || useradd nginx
	groups nginx || groupadd -g nginx nginx
	mkdir -pv /var/lib/nginx/tmp
	yum -y install gcc gcc-c++ git wget automake pcre pcre-devel zlib-devel openssl openssl-devel
	cd /tmp
	wget -O nginx-ct.zip -c https://github.com/grahamedgecombe/nginx-ct/archive/v1.3.2.zip && unzip nginx-ct.zip
	git clone -b OpenSSL_1_1_1-pre7 https://github.com/openssl/openssl.git openssl
	wget -P /tmp/ http://nginx.org/download/nginx-1.15.3.tar.gz && tar axf /tmp/nginx-1.15.3.tar.gz
	sed -i '/^#define NGINX_VERSION /s/1.15.3/1.0.0/g;/^#define NGINX_VER /s/nginx/wejass/g' /tmp/nginx-1.15.3/src/core/nginx.h
	sed -i '/^static u_char ngx_http_server_string/s/nginx/wejass/g' /tmp/nginx-1.15.3/src/http/ngx_http_header_filter_module.c
	sed -i '/^"<hr><center>nginx/s/nginx/wejass/g' /tmp/nginx-1.15.3/src/http/ngx_http_special_response.c
	cd /tmp/nginx-1.15.3/
	./configure --prefix=/usr/share/nginx --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --http-client-body-temp-path=/var/lib/nginx/tmp/client_body --http-proxy-temp-path=/var/lib/nginx/tmp/proxy --http-fastcgi-temp-path=/var/lib/nginx/tmp/fastcgi --http-uwsgi-temp-path=/var/lib/nginx/tmp/uwsgi --http-scgi-temp-path=/var/lib/nginx/tmp/scgi --pid-path=/run/nginx.pid --lock-path=/run/lock/subsys/nginx --user=nginx --group=nginx --with-file-aio --with-openssl=../openssl --with-openssl-opt='enable-tls1_3' --with-http_ssl_module --with-http_v2_module --with-http_realip_module --with-http_addition_module --with-http_auth_request_module --with-http_sub_module --with-http_dav_module --with-http_flv_module --with-http_mp4_module --with-http_gunzip_module --with-http_gzip_static_module --with-http_random_index_module --with-http_secure_link_module --with-http_degradation_module --with-http_stub_status_module --with-mail_ssl_module --with-pcre --with-pcre-jit --add-module=../nginx-ct-1.3.2/
	make && make install
}
_install_zip(){
	cd /tmp
	wget http://www.nih.at/libzip/libzip-1.1.2.tar.xz
	tar Jxvf libzip-1.1.2.tar.xz
	cd libzip-1.1.2
	mkdir -p /usr/local/libzip
	./configure prefix=/usr/local/libzip/
	make -j$CPU_NUM  && make install
}
_install_mutil(){
	mkdir {data,run,logs,binlogs,relaylogs,dump}
	chown -R mysql.mysql *
	mysql_install_db --datadir=/data/war3xydb/data/ --user=mysql
	ExecStart=/usr/bin/mysqld_safe --defaults-file=/data/war3xydb/my.cnf
}
_install_php(){
	yum -y install libxml2 libxml2-devel openssl openssl-devel curl-devel libjpeg-devel libpng-devel freetype-devel libmcrypt-devel bison libzip libzip-devel
	git clone https://github.com/php/php-src.git
	./buildconf
	./configure --prefix=/usr/local/php7 --exec-prefix=/usr/local/php7 --bindir=/usr/local/php7/bin --sbindir=/usr/local/php7/sbin --includedir=/usr/local/php7/include --libdir=/usr/local/php7/lib/php --mandir=/usr/local/php7/php/man --with-config-file-path=/usr/local/php7/etc --with-config-file-scan-dir=/usr/local/php7/etc/php.d --enable-mysqlnd --with-mysqli --with-pdo-mysql --enable-fpm --with-fpm-user=nginx --with-fpm-group=nginx --with-gd --with-iconv --with-zlib --enable-xml --enable-shmop --enable-sysvsem --enable-inline-optimization --enable-mbregex --enable-mbstring --enable-ftp --with-openssl --enable-pcntl --enable-sockets --with-xmlrpc --enable-zip --enable-soap --without-pear --with-gettext --enable-session --with-curl --with-jpeg-dir --with-freetype-dir --enable-opcache
	make && make install
}
_install_mem(){
	yum install -y gcc make
	wget -P /tmp/ https://github.com/libevent/libevent/releases/download/release-2.1.8-stable/libevent-2.1.8-stable.tar.gz
	tar axf /tmp/libevent-2.1.8-stable.tar.gz -C /tmp
	cd /tmp/libevent-2.1.8-stable
	./configure --prefix=/usr/local/libevent
	make && make install
	wget -P /tmp/ http://www.memcached.org/files/memcached-1.5.3.tar.gz
	tar axf /tmp/memcached-1.5.3.tar.gz  -C /tmp/
	cd /tmp/memcached-1.5.3
	./configure --prefix=/usr/local/memcached --with-libevent=/usr/local/libevent/
	make && make install
	/usr/local/memcached/bin/memcached -u root -m 256 -U 0 -p 12001 -l 0.0.0.0 -d
	/usr/local/memcached/bin/memcached -u root -m 256 -U 0 -p 12002 -l 0.0.0.0 -d
}
_install_ssl(){
	# openssl req -new -newkey rsa:2048 -sha256 -nodes -out example_com.csr -keyout example_com.key -subj "/C=CN/ST=ShenZhen/L=ShenZhen/O=Example Inc./OU=Web Security/CN=example.com"
	# openssl x509 -req -days 365 -in example_com.csr -signkey example_com.key -out example_com.crt
	if [ ! -f ~/.acme.sh/acme.sh ];then
		yum install -y socat
		cd /tmp
		git clone https://github.com/Neilpang/acme.sh.git
		cd acme.sh
		./acme.sh --install
	fi
	if [ ! -f /tmp/ct-submit-1.1.2/ct-submit-1.1.2 ];then
		cd /tmp
		wget -O ct-submit.zip -c https://github.com/grahamedgecombe/ct-submit/archive/v1.1.2.zip
		unzip ct-submit.zip
		cd ct-submit-1.1.2
		go build
	fi
	# first
	# ~/.acme.sh/acme.sh --issue -w /data/web -d www.wejass.com -d wejass.com -d cdn.wejass.com -d git.wejass.com
	# ~/.acme.sh/acme.sh --issue -w /data/web -d www.wejass.com -d wejass.com -d cdn.wejass.com -d git.wejass.com --keylength ec-256
	# renew
	~/.acme.sh/acme.sh --renew -d www.wejass.com -d wejass.com -d cdn.wejass.com -d git.wejass.com --force
	~/.acme.sh/acme.sh --renew -d www.wejass.com -d wejass.com -d cdn.wejass.com -d git.wejass.com --force --ecc

	cd ~/.acme.sh/www.wejass.com
	cat fullchain.cer www.wejass.com.key > www.wejass.com.pem
	/tmp/ct-submit-1.1.2/ct-submit-1.1.2 ct.googleapis.com/icarus <www.wejass.com.pem >www.wejass.com.sct
	openssl x509 -in www.wejass.com.pem -noout -subject
	openssl x509 -in www.wejass.com.pem -noout -pubkey | openssl asn1parse -noout -inform pem -out public.key
	PKP1=`openssl dgst -sha256 -binary public.key | openssl enc -base64`
	
	cd ~/.acme.sh/www.wejass.com_ecc
	cat fullchain.cer www.wejass.com.key > www.wejass.com.pem
	/tmp/ct-submit-1.1.2/ct-submit-1.1.2 ct.googleapis.com/icarus <www.wejass.com.pem >www.wejass.com.sct
	openssl x509 -in www.wejass.com.pem -noout -subject
	openssl x509 -in www.wejass.com.pem -noout -pubkey | openssl asn1parse -noout -inform pem -out public.key
	PKP2=`openssl dgst -sha256 -binary public.key | openssl enc -base64`
 
	cp -rf ~/.acme.sh/www.wejass.com* /etc/nginx/openssl/
	cd /etc/nginx/openssl/
	cp -rf www.wejass.com/ca.cer public/
	cp -rf www.wejass.com_ecc/www.wejass.com.sct ct/
	cp -rf www.wejass.com_ecc/www.wejass.com.sct ct/www.wejass.com_ecc.sct
	cd public
	openssl dhparam -out dhparam.pem 4096
	openssl rand 48 > tls_session_ticket.key
 	systemctl force-reload nginx
#	openssl ciphers -V 'TLSv1.2:TLSv1.3'  | column -t
#	curl -s https://www.wejass.com/css/slick-login.css | openssl dgst -sha256 -binary | openssl enc -base64 -A
}
_install_go(){
	yum install -y golang
	go get -u github.com/pangudashu/memcache
	go get -u github.com/go-sql-driver/mysql
	sed -i '/clientLocalFiles/a\\t\tclientMultiStatements |' /root/go/src/github.com/go-sql-driver/mysql/packets.go
	sed -i '/clientMultiStatements/a\\t\tclientMultiResults |' /root/go/src/github.com/go-sql-driver/mysql/packets.go
}
_install_python(){
	yum install -y gcc-c++
	cd /tmp
	wget https://www.python.org/ftp/python/2.7.5/Python-2.7.5.tar.bz2
	tar axf Python-2.7.5.tar.bz2
	./configure 
	make && make install
}
_install_devel(){
	yum install -y git python-devel hsh
	cd /tmp
	git clone git://github.com/webpy/webpy.git
	git clone https://github.com/aliyun/aliyun-oss-python-sdk.git 
	cd /tmp/aliyun-oss-python-sdk/
	python  setup.py install
	cd /tmp/webpy/
	python setup.py install
	pip install -U crcmod
}
_install_server(){
	echo "set ts=4">> /etc/vimrc
	yum install -y gcc golang mariadb-devel libmemcached libmemcached-devel
	yum install -y mariadb mariadb-server  
	yum install -y uwsgi uwsgi-plugin-python mysql-connector-python pyOpenSSL
	cp -rf /data/web/config/nginx.conf /etc/nginx/
	cp -rf /data/web/config/uwsgi.ini /etc/
	tar xzvf /data/web/config/wejass-SSL.tar.gz -C /etc/nginx
	systemctl start nginx
	systemctl start uwsgi
	systemctl start mongod
	chkconfig nginx on
	chkconfig uwsgi on
	chkconfig mongod on
}
_install_es(){
	mkdir -pv /usr/local/kencery/elasticsearch
	cd /usr/local/kencery/elasticsearch
	yum install -y git wget
	wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.6.3.tar.gz
	git clone git://github.com/mobz/elasticsearch-head.git
	git clone git://github.com/medcl/elasticsearch-analysis-ik.git
	tar axf elasticsearch-5.6.3.tar.gz 
	groupadd elasticsearch
	useradd  elasticsearch -g elasticsearch
	chown -R elasticsearch:elasticsearch /usr/local/kencery/elasticsearch/elasticsearch-5.6.3
	su elasticsearch -c "/usr/local/kencery/elasticsearch/elasticsearch-5.6.3/bin/elasticsearch -d"
}
_install_git(){
	# git --version && yum remove -y git 
	yum --version && yum install -y wget gcc curl-devel expat-devel gettext-devel openssl-devel zlib-devel perl-ExtUtils-MakeMaker
	apt --version && apt-get install -y wget gcc libcurl4-gnutls-dev libexpat1-dev gettext libz-dev libssl-dev perl-ExtUtils-MakeMaker
	wget https://www.kernel.org/pub/software/scm/git/git-2.7.2.tar.gz
	tar xzf git-2.7.2.tar.gz
	cd git-2.7.2
	make prefix=/usr/local/git all
	make prefix=/usr/local/git install
	echo 'export PATH=$PATH:/usr/local/git/bin' > /etc/profile.d/git.sh
	source /etc/profile.d/git.sh
	# curl: (35) SSL connect error - yum update -y nss
	# nvm - wget -qO- https://raw.githubusercontent.com/creationix/nvm/v0.33.8/install.sh | bash
}
_install_gitser(){	
	groupadd git;
	useradd git -g git -s /sbin/nologin 
	yum install -y git curl-devel expat-devel gettext-devel openssl-devel zlib-devel perl-devel
	
	cd /home/git/
	mkdir .ssh
	chmod 755 .ssh
	touch .ssh/authorized_keys
	chmod 644 .ssh/authorized_keys
	
	cd /home
	mkdir gitrepo
	chown git:git gitrepo/
	cd gitrepo
	git init --bare runoob.git
	git clone git@45.32.51.137:/home/gitrepo/runoob.git
}
_install_docker(){
	docker pull mariadb
	docker pull memcached
	docker run --name memcached -p 12001:11211 -d memcached memcached -m 64
	docker run --name mysql -d --rm -p 3306:3306 -v /data/docker/mariadb:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=TPG4ppk4rlncL3lO  mariadb
}
_install_golang(){
	wget https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz
	tar axf go1.10.1.linux-amd64.tar.gz -C /usr/local
	echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
	echo 'export GOPATH=$HOME/go:/data/web/golang' >> /etc/profile.d/go.sh
}
_main(){
	set -e
	_install_devel
	_install_nginx
	_install_webpy
}
