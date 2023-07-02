#!/bin/bash

# --process-user=snmp \
# --process-group=snmp \

# -------------------------------------
#  Snmpsim configuration
# -------------------------------------
snmpsim-command-responder \
    --data-dir=/data \
    --agent-udpv4-endpoint=0.0.0.0:1161 \
    --v3-user=simulator \
    --v3-auth-key=auctoritas \
    --v3-auth-proto=SHA256 \
    --v3-priv-key=privatus \
    --v3-priv-proto=AES
    