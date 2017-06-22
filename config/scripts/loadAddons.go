package scripts

const loadAddons = `#! /bin/bash

# wait for the API server to start
sleep 60

while true
do 
  if (/opt/kubernetes/bin/kubectl apply -f /opt/kubernetes/addons --validate=false)
  then 
  	break
  fi

sleep 30
done
`
