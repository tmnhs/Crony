#! /bin/bash
# go-server - this script starts and stops the go-server daemon
#

CURDIR=$(dirname "$0")
cd $CURDIR

PROJECTNAME=$2
if [ "$PROJECTNAME" != "admin" ] && [ "$PROJECTNAME" != "node" ]; then
  echo $"Usage: $0 {start|stop|restart} {admin|node} {testing|production}"
  exit 2
fi
PROJECTBASE="."
PROJECTBIN="$PROJECTBASE"/bin
PROJECTLOGS="$PROJECTNAME"/logs

cd $PROJECTBIN

prog="crony"$PROJECTNAME

ENVIRONMENT=$3

if [ "$ENVIRONMENT" == "" ]; then
  ENVIRONMENT="testing"
fi

start() {
  # shellcheck disable=SC2028
  echo -n $"Starting $PROJECTNAME: "
  if [ -x "$PROJECTBIN/$PROJECTNAME " ]; then
    echo -e $"no execute program"
    exit 5
  fi
  nohup ./$prog -e=$ENVIRONMENT >$PROJECTLOGS/run.log 2>&1 &
  echo -e $"ok"
}

stop() {
  echo -e $"Stopping $prog: "
  pid=$(ps -ef | grep $prog | grep -v grep | awk '{print $2}')
  if [ "$pid" ]; then
    echo -n $"kill process pid:$pid "
    kill $pid
    ret=0
    for((i=1;i<=15;i++));
    do
      sleep 1
      pid=$(ps -ef | grep $prog | grep -v grep | awk '{print $2}')
      if [ "$pid" ]; then
        ret=0
       else
         ret=1
         break
      fi
    done

      if [ "$ret" ]; then
        echo -e $"ok"
      else
        echo -e $"no"
      fi
  else
    echo -e $"no program process to stop"
  fi
}


restart() {
  stop
  sleep 2
  start
}

case "$1" in
start)
  $1
  ;;
stop)
  $1
  ;;
restart)
  $1
  ;;
*)
  echo $"Usage: $0 {start|stop|restart} {admin|node} {testing|production}"
  exit 2
  ;;
esac
