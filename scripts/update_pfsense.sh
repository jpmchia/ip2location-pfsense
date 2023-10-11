#!/bin/bash
# This script copies the files from the pfSense/www directory to the pfSense firewall.
# It is intended to be used for development purposes only.


ARGS=$@
VERBOSE=false
PFSENSE_ADDR="192.168.0.1"
PFSENSE_DEST="/usr/local"
FILE_LOG="update_pfsense.log"

MAGENTA='\e[1;35m'; YELLOW='\e[1;33m'; RED='\e[1;31m'; GREEN='\e[1;32m'; BLUE='\e[1;34m'; CYAN='\e[0;96m'; NORMAL='\e[0m'
CURRENT_DIR="$( pwd )"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_DIR="$( dirname "$SCRIPT_DIR" )"
PFSENSE_DIR="${REPO_DIR}/pfSense"

# Print usage information
function print_usage() 
{
    if [ -z "$ARGS" ]; then
        echo -e "\nUsage: $0 [PFSENSE_USER@]PFSENSE_ADDR\n"
        echo -e "   -h, --help:       Displays this help message"
        echo -e "   -r, --remove:     Removes the pfSense files 1from the pfSense firewall"
        echo -e "   -v, --verbose:    Displays verbose output\n"
        echo -e "   PFSENSE_USER:     The username to use when connecting to the pfSense firewall."
        echo -e "   PFSENSE_ADDR:     The hostname or IP of the pfSense firewall to update.\n"
        echo -e "NOTE: This script assumes that you have access to the pfSense firewall via SSH.\n"
        echo -e "Examples: "
        echo -e "    $0 192.168.0.1"
        echo -e "    $0 root@192.168.0.1"
        exit 1
    fi
}

# Check that the required arguments have been specified
function check_args()
{
    if [ -z "$PFSENSE_ADDR" ]; then
        echo -e "${RED}ERROR: pfSense address not specified${NORMAL}"
        print_usage
        exit 1
    fi
}



# Print a message if VERBOSE is true
function print_verbose() {
    if [ "$VERBOSE" = true ]; then
        echo -e "${*}"
    fi
}

# Execute a command, if VERBOSE is true then the command will be printed before it is executed 
# and the output will be displayed, otherwise the output will be suppressed
function execute_verbose() {
    echo -e "${YELLOW}Executing: ${NORMAL}${*}"
    eval "${*}"
}

# Remove files from pfSense firewall, this is done by reading the file log from the previous run
# if it exists, otherwise it will read the files from the pfSense/www directory and remove them
function remove_files() {
    echo -e "${YELLOW}Removing files from pfSense firewall...${NORMAL}"

    # Check for the file log from previous run
    if [ -f "${FILE_LOG}" ]; then
        echo -e "${YELLOW}Found file log from previous run: ${NORMAL}${FILE_LOG}"
        echo -e "${YELLOW}Removing files from pfSense firewall...${NORMAL}"
        while IFS= read -r file; do
            print_verbose "${YELLOW}Removing ${NORMAL}${PFSENSE_ADDR}:${PFSENSE_DEST}/${file}${NORMAL}"
            ssh "${PFSENSE_ADDR}" "rm -f ${PFSENSE_DEST}/${file}"
        done < "${FILE_LOG}"
    else
        find "${PFSENSE_DIR}/www" -type f -exec realpath --relative-to="${PFSENSE_DIR}" '{}' \; | while IFS= read -r file; do
            print_verbose "${YELLOW}Removing ${NORMAL}${PFSENSE_ADDR}:${PFSENSE_DEST}/${file}${NORMAL}"
            echo -e "${file}" >> "${FILE_LOG}"
            ssh "${PFSENSE_ADDR}" "rm -f ${PFSENSE_DEST}/${file}"
        done
    fi
}

# Copy files to pfSense firewall
function copy_files() {
    echo -e "${YELLOW}Copying files to pfSense firewall...${NORMAL}"

    # Check for the file log from previous run
    if [ -f "${FILE_LOG}" ]; then
        echo -e "${YELLOW}Found file log from previous run: ${NORMAL}${FILE_LOG}"
        mv "${FILE_LOG}" "${FILE_LOG}.$(date -r ${FILE_LOG} +%Y%m%d%H%M%S)"
    fi

    find "${PFSENSE_DIR}/www" -type f -exec realpath --relative-to="${PFSENSE_DIR}" '{}' \; | while IFS= read -r file; do
        print_verbose "${YELLOW}Copying ${NORMAL}${file}${YELLOW} => ${NORMAL}${PFSENSE_ADDR}:${PFSENSE_DEST}/${file}${NORMAL}"
        echo -e "${file}" >> "${FILE_LOG}"
        scp "${PFSENSE_DIR}/${file}" "${PFSENSE_ADDR}:${PFSENSE_DEST}/${file}"
    done
}


function create_directories () {
	unset -e

	echo -e "${YELLOW}Creating target directories on pfSense firewall ... ${NORMAL}"

	find "${PFSENSE_DIR}/www" -type d -exec realpath --relative-to="${PFSENSE_DIR}" '{}' \; | while IFS= read -r dirname; do
		let count++
		print_verbose "${YELLOW}Creating directory on pfSense [${count}] ${NORMAL}${dirname}"
		echo -e "${PFSENSE_DEST}/${dirname}"
		ssh ${PFSENSE_ADDR} "mkdir -p ${PFSENSE_DEST}/${dirname}" < /dev/null
	done
}


function copy_dev_files() {
    echo -e "${YELLOW}Copying files to pfSense firewall...${NORMAL}"

    execute_verbose scp "${PFSENSE_DIR}/www/widgets/widgets/*" "${PFSENSE_ADDR}:${PFSENSE_DEST}/www/widgets/widgets/"
    execute_verbose scp "${PFSENSE_DIR}/www/widgets/javascript/*" "${PFSENSE_ADDR}:${PFSENSE_DEST}/www/widgets/javascript/"
    execute_verbose scp "${PFSENSE_DIR}/www/css/*.css" "${PFSENSE_ADDR}:${PFSENSE_DEST}/www/css/"
    exit 0
}

# Main function
function main() {
    check_args
    print_verbose -e "${YELLOW}Current directory: ${NORMAL} $CURRENT_DIR"
    print_verbose -e "${YELLOW}Script directory: ${NORMAL} $SCRIPT_DIR"
    print_verbose -e "${YELLOW}Repo directory: ${NORMAL}$REPO_DIR"
    print_verbose -e "${YELLOW}pfSense directory: ${NORMAL}$PFSENSE_DIR"
    print_verbose -e "${YELLOW}pfSense user: ${NORMAL}$PFSENSE_USER"
    print_verbose -e "${YELLOW}pfSense address: ${NORMAL}$PFSENSE_ADDR"
    print_verbose -e "${YELLOW}pfSense destination: ${NORMAL}$PFSENSE_DEST"

    if [ "$DEVMODE" = true ]; then
        echo -e "${YELLOW}Running in development mode, will copy a subset of files to pfSense firewall${NORMAL}"
        copy_dev_files
    fi
    if [ "$REMOVE" = true ]; then
        remove_files
    else
	create_directories
        copy_files
    fi
}


# Set defaults
VERBOSE=false
REMOVE=false
DEVMODE=false
PFSENSE_ADDR=""
ARGS=$@

# Parse arguments and options, requires PFSENSE_ADDR to be specified as the last argument, and the others are optional arguments
while [ "$1" != "" ]; do
    case $1 in
        -h | --help )           print_usage
                                exit 0
                                ;;
        -r | --remove )         REMOVE=true
                                ;;
        -v | --verbose )        VERBOSE=true
                                ;;
        -d | --dev )            DEVMODE=true
                                ;;
        * )                     PFSENSE_ADDR=$1
                                ;;
    esac
    shift
done

print_verbose -e "${YELLOW}Verbose mode: ${NORMAL}${VERBOSE}"
print_verbose -e "${YELLOW}Remove mode: ${NORMAL}${REMOVE}"
print_verbose -e "${YELLOW}pfSense address: ${NORMAL}${PFSENSE_ADDR}"

main
