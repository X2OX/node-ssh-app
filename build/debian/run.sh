#!/bin/zsh

/etc/init.d/ssh restart

run(){
  RESULT=$(curl -s "http://$SERVICE/api/v1/nsa/debian/$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/)/register")

  if [ "$RESULT" -eq "" ] ; then
    echo "again"
    sleep 1
    run
  else
    cloudflared service install "$RESULT"
    tail -f /dev/null
  fi
}

run
