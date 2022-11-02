#!/bin/bash

HOSTS=$(hssh l | cut -d " " -f 1 > /tmp/hosts && cat /tmp/hosts)
OUTPUT_FILE="/tmp/reboot-status.out"
rm ${OUTPUT_FILE} > /dev/null || true
touch ${OUTPUT_FILE}

for VM in $HOSTS; do
  ssh -o ConnectTimeout=5 ${VM} "ls -la /var/run/reboot-required" &> /dev/null
  EXIT_CODE=$?
  if [[ ${EXIT_CODE} -eq 0 ]]; then
    echo "${VM}: reboot required" >> ${OUTPUT_FILE}
  elif [[ ${EXIT_CODE} -eq 1 ]] || [[ ${EXIT_CODE} -eq 2 ]]; then
    echo "${VM}: reboot not required" >> ${OUTPUT_FILE}
  elif [[ ${EXIT_CODE} -eq 255 ]]; then
    echo "${VM}: ssh connection failed" >> ${OUTPUT_FILE}
  else
    echo ${VM}: ${EXIT_CODE}
  fi
done

cat /tmp/reboot-status.out
