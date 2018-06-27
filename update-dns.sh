#!/bin/bash
#

if [ ! -z ${CF_KEY+x} ]; then
	./cloudflare-update-record.sh ${CF_MAIL} ${CF_KEY} ${DDNS_DOMAIN} ${DDNS_SUBDOMAIN}.${DDNS_DOMAIN}
fi

if [ ! -z ${OVH_AK+x} ]; then
	./ovh-update-record -domain=${DDNS_DOMAIN} -subdomain=${DDNS_SUBDOMAIN} -ak=${OVH_AK} -as=${OVH_AS} -ck=${OVH_CK}
fi

sleep ${DDNS_SLEEP:=8}h