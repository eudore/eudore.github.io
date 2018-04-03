#!/bin/bash


#生成 CA 证书
certtool --generate-privkey --outfile ca-key.pem
cat >ca.tmpl <<EOF
cn = "VPN CA"
organization = "Big Corp"
serial = 1
expiration_days = 3650
ca
signing_key
cert_signing_key
crl_signing_key
EOF
certtool --generate-self-signed --load-privkey ca-key.pem \
--template ca.tmpl --outfile ca-cert.pem
#生成本地服务器证书
certtool --generate-privkey --outfile server-key.pem
cat >server.tmpl <<EOF
cn = "weer.top"
organization = "MyCompany"
serial = 2
expiration_days = 3650
encryption_key
signing_key
tls_www_server
EOF
certtool --generate-certificate --load-privkey server-key.pem \
--load-ca-certificate ca-cert.pem --load-ca-privkey ca-key.pem \
--template server.tmpl --outfile server-cert.pem



#生成客户端证书
vim gen-client-cert.sh
#!/bin/bash

USER=$1
CA_DIR=$2
SERIAL=`date +%s`

certtool --generate-privkey --outfile $USER-key.pem

cat << _EOF_ >user.tmpl
cn = "$USER"
unit = "users"
serial = "$SERIAL"
expiration_days = 9999
signing_key
tls_www_client
_EOF_

certtool --generate-certificate --load-privkey $USER-key.pem --load-ca-certificate $CA_DIR/ca-cert.pem --load-ca-privkey $CA_DIR/ca-key.pem --template user.tmpl --outfile $USER-cert.pem

openssl pkcs12 -export -inkey $USER-key.pem -in $USER-cert.pem -name "$USER VPN Client Cert" -certfile $CA_DIR/ca-cert.pem -out $USER.p12



mkdir user
cd user
# user 指的是用户名，.. 指的是 ca 证书所在的目录
../gen-client-cert.sh user ..
# 按提示设置证书使用密码，或直接回车不设密码
