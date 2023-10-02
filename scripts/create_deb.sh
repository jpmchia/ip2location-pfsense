#!/bin/bash

PACKAGE_NAME="ip2location-pfsense"
DEBEMAIL="info@terra-net.uk"
DEBFULLNAME="TerraNet UK (www.terra-net.uk)"

export DEBEMAIL
export DEBFULLNAME

MAGENTA='\e[1;35m'; YELLOW='\e[1;33m'; RED='\e[1;31m'; GREEN='\e[1;32m'; BLUE='\e[1;34m'; CYAN='\e[0;96m'; NORMAL='\e[0m'
CURRENT_DIR="$( pwd )"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_DIR="$( dirname "$SCRIPT_DIR" )"
PACKAGE_DIR="${REPO_DIR}/package"
DATE_TIME="$( date "+%Y%m%d-%H%M%S" )"

echo -e "${YELLOW}Current directory: ${NORMAL} $CURRENT_DIR"
echo -e "${YELLOW}Script directory: ${NORMAL} $SCRIPT_DIR"
echo -e "${YELLOW}Repo directory: ${NORMAL}$REPO_DIR"

# Obtain the version number from the version.go file
cd $REPO_DIR
VERSION="$(grep "const Version" backend/version/version.go | sed -E 's/.*"([^"]+)".*/\1/')"

DEB_PACKAGE_NAME="${PACKAGE_NAME}-${VERSION}"
DEB_PACKAGE_DIR="${PACKAGE_DIR}/${DEB_PACKAGE_NAME}"

# Display the target package name, version and directory
echo -e "\n${YELLOW}Package name: ${NORMAL}${PACKAGE_NAME}"
echo -e "${YELLOW}Package version: ${NORMAL}${VERSION}"
echo -e "${YELLOW}Package directory: ${NORMAL}${DEB_PACKAGE_DIR}"

# Check if the package directory already exists, if it does then rename it
if [ -d "${DEB_PACKAGE_DIR}" ]; then
	previous_datetime="$( date -r ${DEB_PACKAGE_DIR} "+%Y%m%d-%H%M%S" )"
	echo -e "\n${YELLOW}Backing up previous package creation: ${YELLOW}=>${NORMAL} ${DEB_PACKAGE_DIR}.${previous_datetime}"
	mv "${DEB_PACKAGE_DIR}" "${DEB_PACKAGE_DIR}.${previous_datetime}"
fi

# Create the package directory
echo -e "\n${YELLOW}Creating package directory:${NORMAL} ${DEB_PACKAGE_DIR} "
mkdir -p "${DEB_PACKAGE_DIR}"
cd "${DEB_PACKAGE_DIR}"

# Install the prerequisite packages
echo -e "\n${YELLOW}Installing prerequisite packages...${NORMAL}"
sudo apt install build-essential binutils lintian debhelper dh-make devscripts -y

# Create the package
dh_make -s -y --createorig --indep --packagename "${PACKAGE_NAME}" --email "${DEBEMAIL}" --fullname "${DEBFULLNAME}" --native

