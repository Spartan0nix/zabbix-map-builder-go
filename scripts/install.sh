#!/bin/bash

set -e

TMP_INSTALL_DIR="/tmp/zabbix-map-builder"
TMP_INSTALL_ARCHIVE_PATH="$TMP_INSTALL_DIR/zabbix-map-builder.tar.gz"
INSTALL_DIR="/usr/local/bin"

if [[ -f "$INSTALL_DIR/zabbix-map-builder" ]]
then
    echo "[WARNING] Previous install detected at '$INSTALL_DIR/zabbix-map-builder'"
    echo "Exiting..."
    exit 1
fi


# Retrieve the list of release objects from github API
releases=$(curl --silent https://api.github.com/repos/Spartan0nix/zabbix-map-builder-go/releases)

# ----------------------------------------------------
# TAG
# ----------------------------------------------------
echo "[INFO] Listing available version"
# Retrieve all tag_name from the available releases
tags=$(echo $releases | jq -r '.[].tag_name')
# Replace space by newline to render a better ouput in the shell
echo $tags | tr '[[:space:]]' '\n'

selectedTag=""
re=[0-9]\.[0-9]\.[0-9]
# Loop while user as not provided a proper release format (vX.Y.Z)
while [[ ! "$selectedTag" =~ $re ]]
do
    echo "Select one version for the list above :"
    read selectedTag
done

# ----------------------------------------------------
# RELEASE
# ----------------------------------------------------
# Extract the selected release
selectedRelease=$(echo $releases | jq --arg v $selectedTag '.[] | select(.tag_name == $v)')
if [[ $selectedRelease == "" ]]
then
    echo "[ERROR] No release '$selectedTag' found"
    exit 1
fi

# ----------------------------------------------------
# ASSET
# ----------------------------------------------------
assets=$(echo $selectedRelease | jq -r '.assets[].name')

echo ""
echo "[INFO] Listing available assets for the selected release"
echo $assets | tr '[[:space:]]' '\n'

selectedAsset=""
re=^.+tar\.gz$

while [[ ! "$selectedAsset" =~ $re ]]
do
    echo "Select one asset for the list above :"
    read selectedAsset
done

# ----------------------------------------------------
# ASSET URL
# ----------------------------------------------------
# Retrive the asset URL
assetUrl=$(echo $selectedRelease | jq -r --arg v $selectedAsset '.assets[] | select(.name == $v) | .browser_download_url')
if [[ $assetUrl == "" ]]
then
    echo "[ERROR] No assets '$selectedAsset' found for release"
    exit 1
fi

# ----------------------------------------------------
# INSTALL
# ----------------------------------------------------
echo ""
echo "[INFO] Creating '$TMP_INSTALL_DIR'"
if [[ ! -d $TMP_INSTALL_DIR ]]
then
    mkdir $TMP_INSTALL_DIR
fi

echo "[INFO] Downloading asset to '$TMP_INSTALL_ARCHIVE_PATH'"
wget $assetUrl -O $TMP_INSTALL_ARCHIVE_PATH -q --show-progress

echo "[INFO] Extracting archive"
tar -C $TMP_INSTALL_DIR -xzf $TMP_INSTALL_ARCHIVE_PATH

echo "[INFO] Moving binary to '$INSTALL_DIR'"
sudo mv "$TMP_INSTALL_DIR/zabbix-map-builder" /usr/local/bin

echo "[INFO] Updating permissions"
sudo chown $(id -un):$(id -gn) /usr/local/bin/zabbix-map-builder

echo "[INFO] Removing '$TMP_INSTALL_DIR'"
rm -r $TMP_INSTALL_DIR

echo ""
echo "--------------------------------------------------------------------------------------------------"
echo "Release '$selectedTag' of the 'zabbix-map-builder' was succesfuly installed to '$INSTALL_DIR'."
echo "Thanks you for using this cli and don't forget to check out the documentation at 'https://github.com/Spartan0nix/zabbix-map-builder-go#zabbix-map-builder-go'"
echo "--------------------------------------------------------------------------------------------------"

exit 0