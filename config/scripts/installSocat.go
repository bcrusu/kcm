package scripts

// https://gist.github.com/cdemers/be415cb46327e56c5c47f9689a07a456
const installSocat = `#! /bin/bash

if [ -e /opt/bin/socat ]
then
  echo socat binary is already installed. Skipping install...
  exit 0
fi

# Make socat directories
mkdir -p /opt/bin/socat.d/bin /opt/bin/socat.d/lib

# Create socat wrapper
cat << EOF > /opt/bin/socat
#! /bin/bash
PATH=/usr/bin:/bin:/usr/sbin:/sbin:/opt/bin
LD_LIBRARY_PATH=/opt/bin/socat.d/lib:$LD_LIBRARY_PATH exec /opt/bin/socat.d/bin/socat "\$@"
EOF

chmod +x /opt/bin/socat

# Get socat and libraries from the CoreOS toolbox 
cat <<EOF | toolbox 
dnf install -y socat
cp -f /usr/bin/socat /media/root/opt/bin/socat.d/bin/socat
cp -f /usr/lib64/libssl.so.1.0.2h /media/root/opt/bin/socat.d/lib/libssl.so.10
cp -f /usr/lib64/libcrypto.so.1.0.2h /media/root/opt/bin/socat.d/lib/libcrypto.so.10
EOF
`
