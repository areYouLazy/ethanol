#!/bin/bash

##
# Helper script to build core and plugins
# Ignores plugin folders starting with an underscore _
# .so files are placed in ./search_providers
##

# get current dir
ETHANOLDIR=$(pwd)

# setup some variables
SEARCHPROVIDERSDIR="$ETHANOLDIR"/search_providers
PLUGINSFOLDER="$ETHANOLDIR"/ethanol_plugins
GO=$(which go)

# upgrade modules
"$GO" get -u

# cleanup modules list
echo "[+] cleanup module list"
go mod tidy

# upgrade modules
echo "[+] upgrade go modules"
"$GO" get -u ./...

echo "[+] building ethanol core"
# build
"$GO" build

# enter ethanol plugins folder
cd "$PLUGINSFOLDER"

echo "[*] building plugins"
# iterate folders and build plugins
for f in *; do
    if [[ ! -L "$f" && -d "$f" ]]; then # ignore symlinks
        if [[ "$f" == _* ]]; then 
            echo "[-] ignoring folder $f"
        else
            echo "[+] building $f plugin"
            cd $f
            "$GO" build -buildmode=plugin -o "$SEARCHPROVIDERSDIR"/"$f".so ./*.go
            cd ..
        fi
    else
        echo "[-] $f is not a folder"
    fi
done

cd $ETHANOLDIR
