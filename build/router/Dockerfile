FROM python:3.11-slim-bullseye

RUN apt-get update && \
    apt-get upgrade -y

RUN groupadd snmp && \
    useradd -g snmp snmp && \
    mkdir /data && \
    chown -R snmp:snmp /data

RUN pip install pyasn1==0.4.8 && \
    pip install snmpsim-lextudio && \
    rm -rf /usr/local/snmpsim/data && \
    rm -rf /usr/local/snmpsim/variation && \
    rm -rf /usr/local/lib/python3.11/site-packages/snmpsim/data && \
    rm -rf /usr/local/lib/python3.11/site-packages/snmpsim/variation
   
EXPOSE 1161/udp

COPY entrypoint.sh /entrypoint.sh

USER snmp

CMD [ "snmpsim-command-responder", "--data-dir", "/data","--agent-udpv4-endpoint", "0.0.0.0:1161", "--v3-user", "simulator", "--v3-auth-key", "auctoritas", "--v3-auth-proto", "SHA256", "--v3-priv-key", "privatus", "--v3-priv-proto", "AES" ]
