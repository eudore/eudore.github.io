#!/bin/bash
_install_nginx(){
	id nginx || useradd nginx
	groups nginx || groupadd -g nginx nginx
	mkdir -pv /var/lib/nginx/tmp
	yum -y install gcc gcc-c++ git wget automake pcre pcre-devel zlib-devel openssl openssl-devel
	cd /tmp
	wget http://nginx.org/download/nginx-1.13.0.tar.gz
	tar axf /tmp/nginx-1.13.0.tar.gz
	git clone https://github.com/cuber/ngx_http_google_filter_module
	git clone https://github.com/yaoweibin/ngx_http_substitutions_filter_module
	cd /tmp/nginx-1.13.0
	./configure --prefix=/usr/share/nginx --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --error-log-path=/var/log/nginx/error.log --http-log-path=/var/log/nginx/access.log --http-client-body-temp-path=/var/lib/nginx/tmp/client_body --http-proxy-temp-path=/var/lib/nginx/tmp/proxy --http-fastcgi-temp-path=/var/lib/nginx/tmp/fastcgi --http-uwsgi-temp-path=/var/lib/nginx/tmp/uwsgi --http-scgi-temp-path=/var/lib/nginx/tmp/scgi --pid-path=/run/nginx.pid --lock-path=/run/lock/subsys/nginx --user=nginx --group=nginx --with-file-aio --with-ipv6 --with-http_ssl_module --with-http_realip_module --with-http_addition_module --with-http_auth_request_module --with-http_sub_module --with-http_dav_module --with-http_flv_module --with-http_mp4_module --with-http_gunzip_module --with-http_gzip_static_module --with-http_random_index_module --with-http_secure_link_module --with-http_degradation_module --with-http_stub_status_module --with-mail_ssl_module --with-pcre --with-pcre-jit --add-module=../ngx_http_google_filter_module --add-module=../ngx_http_substitutions_filter_module
	make && make install
}
_install_nginx


server {
	listen       80;
	server_name  localhost;
	proxy_headers_hash_max_size 51200;
	proxy_headers_hash_bucket_size 6400;
	resolver 8.8.8.8;
	location / {
		proxy_pass https://www.google.com;
		proxy_connect_timeout 120;
		proxy_read_timeout 600;
		proxy_send_timeout 600;
		send_timeout 600;
		proxy_redirect    off;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		google on;
		google_language "zh-CN";
	}
}
