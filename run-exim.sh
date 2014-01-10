#!/bin/bash
if [ "$USER" != "root" ]
then
  echo Need to be root
  exit
fi

# check /etc/hosts
NS=$( grep mailbot /etc/hosts )
EXPECTNS="10.10.10.2 test.mailbot.net"
if [ "$NS" != "$EXPECTNS" ]; then
  echo "/etc/hosts does not contain $EXPECTNS"
  exit
fi

# all the files are copied to /tmp 
for MDIR in './' './mailbot' './Grant-Murray/mailbot' './github.com/Grant-Murray/mailbot' 'NOTFOUND'
do
  [ -f "$MDIR/usercreds.txt" ] && break
done

if [ "${MDIR}" = "NOTFOUND" ]; then
  echo "Unable to find the mailbot directory"
  exit 1
fi

cp /etc/ssl/GLM-Hosts/test.mailbot.net.key /tmp
cp /etc/ssl/GLM-Hosts/test.mailbot.net.pem /tmp
rm /tmp/exim.mailbot.*.log
rm -rf /tmp/mailbot.boxes
rm -r /var/spool/exim/*
chown grant:zm /var/spool/exim
install --owner grant --group zm -d /tmp/mailbot.boxes
cp $MDIR/usercreds.txt /tmp 
install --owner root --group root --mode 0644 -T $MDIR/exim.test.conf /etc/mail/exim.conf 

#clear
exim -bdf -d -v 

for MUSER in /tmp/mailbot.boxes/* ; do
  echo "===== $MUSER ====="
  cat $MUSER
  echo ; echo;
done
