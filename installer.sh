#!/bin/bash

echo "installing jotun binary..."
cp ./bin/jotun /usr/local/bin/jotun
if [ $? -eq 0 ]; then
echo "continue.."
else 
echo "failed to install jotun's binary"
exit 1 
fi
echo "Installing jotun's man page"
gzip -c jotun.1 > jotun.1.gz
cp jotun.1.gz /usr/share/man/man1/
if [ $? -eq 0 ]; then
echo "Installation finished. Exiting.."
else 
echo "failed to install jotun's man page"
exit 1 
fi
exit 0
