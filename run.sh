#!/bin/zsh

/etc/init.d/ssh restart

if [ "$CLOUDFLARED_TOKEN" -eq "" ]; then
else
  cloudflared service install "$CLOUDFLARED_TOKEN"
fi

tail -f /dev/null
