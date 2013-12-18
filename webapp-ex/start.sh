#!/bin/bash

function shut_down() {
	# deregister webapp (this will remove all records for "webapp" so won't work if running multiple copies)
	redis-cli -h $HIPACHE_PORT_6379_TCP_ADDR -p 6379 del frontend:$HOST

	# send a signal to the webapp so it can shut down appropriately too
	kill -SIGTERM $pid	

	exit
}

# on `docker stop` clean up a little
trap "shut_down" SIGTERM SIGHUP SIGINT


IP=$(ifconfig eth0 | grep "inet addr" | awk '{print $2}' | awk -F: '{print $2}')

# Start webapp
/opt/webapp &
pid=$!

# Register address in Redis for Hipache
redis-cli -h $HIPACHE_PORT_6379_TCP_ADDR -p 6379 rpush frontend:$HOST $HOSTNAME
redis-cli -h $HIPACHE_PORT_6379_TCP_ADDR -p 6379 rpush frontend:$HOST http://$IP:8080


wait