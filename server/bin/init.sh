#! /bin/sh
### BEGIN INIT INFO
# Provides:          scriptname
# Required-Start:    $remote_fs $syslog
# Required-Stop:     $remote_fs $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Start daemon at boot time
# Description:       Enable service provided by daemon.
### END INIT INFO

# Author: Damien Lebrun <dinoboff@gmail.com>

PATH=/sbin:/usr/sbin:/bin:/usr/bin:/usr/local/bin
DESC="Run user script in docker container"
NAME=verifier-server
DAEMON=/usr/local/bin/$NAME
DAEMON_ARGS="-http 0.0.0.0:80"
PIDFILE=/var/run/$NAME.pid
SCRIPTNAME=/etc/init.d/"$NAME"-d
DAEMONLOGFILE=/var/log/$NAME.log
USER=verifier-server
GROUP=docker

# Exit if the package is not installed
[ -x "$DAEMON" ] || exit 0


[ -r /etc/default/$NAME ] && . /etc/default/$NAME
. /lib/init/vars.sh
. /lib/lsb/init-functions


do_start()
{
	touch $DAEMONLOGFILE
	chgrp $GROUP $DAEMONLOGFILE

	# Return
	#   0 if daemon has been started
	#   1 if daemon was already running
	#   2 if daemon could not be started
	start-stop-daemon --start --quiet --pidfile $PIDFILE --make-pidfile \
		--background --chuid $USER:$GROUP \
		--exec $DAEMON \
		--test > /dev/null \
            || return 1

	start-stop-daemon --start --quiet --pidfile $PIDFILE --make-pidfile \
	    --background --no-close --chuid $USER:$GROUP \
	    --exec $DAEMON -- $DAEMON_ARGS >> $DAEMONLOGFILE 2>&1 \
            || return 2
}


do_stop()
{
	# Return
	#   0 if daemon has been stopped
	#   1 if daemon was already stopped
	#   2 if daemon could not be stopped
	#   other if a failure occurred
	start-stop-daemon --stop --quiet --retry=TERM/30/KILL/5 --pidfile $PIDFILE --name $NAME
	RETVAL="$?"
	[ "$RETVAL" = 2 ] && return 2

	start-stop-daemon --stop --quiet --oknodo --retry=0/30/KILL/5 --exec $DAEMON
	[ "$?" = 2 ] && return 2
	rm -f $PIDFILE

	return "$RETVAL"
}


case "$1" in
  start)
	[ "$VERBOSE" != no ] && log_daemon_msg "Starting $DESC" "$NAME"
	do_start
	case "$?" in
		0|1) [ "$VERBOSE" != no ] && log_end_msg 0 ;;
		2) [ "$VERBOSE" != no ] && log_end_msg 1 ;;
	esac
	;;
  stop)
	[ "$VERBOSE" != no ] && log_daemon_msg "Stopping $DESC" "$NAME"
	do_stop
	case "$?" in
		0|1) [ "$VERBOSE" != no ] && log_end_msg 0 ;;
		2) [ "$VERBOSE" != no ] && log_end_msg 1 ;;
	esac
	;;
  status)
	status_of_proc "$DAEMON" "$NAME" && exit 0 || exit $?
	;;
  restart|force-reload)
	log_daemon_msg "Restarting $DESC" "$NAME"
	do_stop
	case "$?" in
	  0|1)
		do_start
		case "$?" in
			0) log_end_msg 0 ;;
			1) log_end_msg 1 ;; # Old process is still running
			*) log_end_msg 1 ;; # Failed to start
		esac
		;;
	  *)
		# Failed to stop
		log_end_msg 1
		;;
	esac
	;;
  *)
	echo "Usage: $SCRIPTNAME {start|stop|status|restart|force-reload}" >&2
	exit 3
	;;
esac

