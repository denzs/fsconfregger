#!/bin/bash

if [ "$#" == "1" ]
then
cat > /opt/freeswitch/conf/sip_profiles/external/$1.xml << EOF
<include>
	<gateway name="conf_$1">
		<param name="realm" value="10.28.128.231"/>
		<param name="username" value="$1"/>
		<param name="password" value="very_secret"/>

        <param name="extension" value="auto_to_user"/>

		<!-- <param name="extension" value="auto_to_user"/> -->
		<param name="caller-id-in-from" value="true"/>

		<!-- If you are having a problem with the default registering as gw+gateway_name@ip you can set this to true to use extension@ip -->
		<param name="extension-in-contact" value="true"/>

		<!-- uncomment to disable registering at the gateway -->
		<param name="register" value="true"/>

		<!--send an options ping every x seconds, failure will unregister and/or mark it down-->
		<param name="ping" value="60"/>
	</gateway>
</include>
EOF
    exit 0
fi

if [ "$#" == "2" ]
then
    rm /opt/freeswitch/conf/sip_profiles/external/$1.xml
    exit 0
fi

echo "1 or 2 paramaters required"
echo "example: $0 1234"
echo "example: $0 1234 del"
exit 1
