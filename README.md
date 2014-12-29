# yubikey-gate

External authenticator for [Apache mod-auth-external](http://code.google.com/p/mod-auth-external/)

# Usage

## install mod-auth-external
Follow the instruction: [http://code.google.com/p/mod-auth-external/wiki/Installation](http://code.google.com/p/mod-auth-external/wiki/Installation)

## add in your vhost file

``
DefineExternalAuth yubikey environment /home/name/gate/yubikey-gate

<Directory "/var/www/html/private">
	   AuthType Basic
	   AuthName "Private AREA"
	   AuthBasicProvider external
	   AuthExternal yubikey
	   Require valid-user
	   AuthExternalContext /home/name/gate/gate.conf
</Directory>
``

## add account the in conf file
You need to add a section by user, put the private key and init the counter to 0

``
[stumpy]
key=7d9001750a1d86423e06c1425027c742
counter=0
``