#!/bin/bash

set -e

INSTALL_DIR="/usr/local/bin/zabbix-map-builder"

if [[ -f $INSTALL_DIR ]]
then
    echo "[INFO] Removing '$INSTALL_DIR'"
    sudo rm $INSTALL_DIR

    echo ""
    echo "--------------------------------------------------------------------------------------------------"
    echo "'zabbix-map-builder' was succesfuly uninstalled."
    echo "Thanks you for using this cli and don't forget to check out the documentation at 'https://github.com/Spartan0nix/zabbix-map-builder-go#zabbix-map-builder-go'"
    echo "--------------------------------------------------------------------------------------------------"
else
    echo "[WARNING] No previous install detected"
    echo "[INFO] - '$INSTALL_DIR' -> Missing"
fi

exit 0