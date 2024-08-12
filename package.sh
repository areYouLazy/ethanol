#!/bin/bash

##
# Helper script to generate a .deb package
# To be completed
#
# Execution Flow
# - Generate folder structure
# - Generate a CONTROL file
# - Generate a systemd .service file
# - Generate preinst script
# - Generate postinst script
# - Generate prerm script
# - Generate postrm script
# - Generate .deb package
##

if [ "$#" -lt 1 ]; then
  echo "No arguments provided"
  exit 1
elif [ "$#" -gt 1 ]; then
  echo "Too much arguments"
  exit 1
fi

VERSION="$1"
echo "Generating ethanol package for version: "$VERSION""

# get current dir
ETHANOLDIR=$(pwd)

# setup variables
BUILDDIR="$ETHANOLDIR"/build/ethanol
RELEASEDIR="$ETHANOLDIR"/build/release

DEBIANDIR=/DEBIAN
BINDIR=/opt/ethanol/bin
SYSTEMDFOLDER=/etc/systemd/system
MANPAGESDIR=/usr/share/man
CONFIGDIR=/etc/ethanol/
PLUGINSDIR="$CONFIGDIR"/search_providers
LIBSDIR=/usr/local/lib
WORKINGDIR=/opt/ethanol

# cleanup working directory
rm -rf "$BUILDDIR"

# create directory structure
mkdir "$BUILDDIR"
mkdir -p "$RELEASEDIR"
mkdir -p "$BUILDDIR""$DEBIANDIR"
mkdir -p "$BUILDDIR""$BINDIR"
mkdir -p "$BUILDDIR""$SYSTEMDFOLDER"
mkdir -p "$BUILDDIR""$MANPAGESDIR"
mkdir -p "$BUILDDIR""$CONFIGDIR"
mkdir -p "$BUILDDIR""$LIBSDIR"
mkdir -p "$BUILDDIR""$PLUGINSDIR"
mkdir -p "$BUILDDIR""$WORKINGDIR"

# copy files
cp ./ethanol "$BUILDDIR""$PROJECTBINDIR"
cp ./config.template.yml "$BUILDDIR""$CONFIGDIR"
cp ./config.template.json "$BUILDDIR""$CONFIGDIR"
cp -r ./search_providers/* "$BUILDDIR""$PLUGINSDIR"

# generate control file
echo 'Package: Ethanol
Version: 0.1
Section: MetaSearch Engines
Priority: optional
Architecture: all
Homepage: https://github.com/areYouLazy/ethanol
Maintainer: areYouLazy
Description: Ethanol is a search aggregator,
    it is capable of querying multiple applications and return aggregate results.
    It is like a MetaSearch Engine but for applications like Jira, MediaWiki,
    OTRS, Check_MK and other enterprise applications
' > "$BUILDDIR""$DEBIANDIR"/control

# generate systemd service file
echo '[Unit]
Description=Ethanol
After=network.target

[Service]
Type=simple
User=ethanol
WorkingDirectory=/opt/ethanol
ExecStart=/opt/ethanol/bin/ethanol
Restart=on-failure

[Install]
WantedBy=multi-user.target
' > "$BUILDDIR""$SYSTEMDFOLDER"/ethanol.service

# generate preinst script
echo '#!/bin/bash
 
if  [ "$1" = install ] 
then
  echo "Setting up ethanol group and user"

  # creating ethanol group if he isn t already there
  if ! getent group ethanol >/dev/null; then
    # Adding system group: ethanol.
    addgroup --system ethanol >/dev/null
  fi
 
  # creating ethanol user if he isn t already there
  if ! getent passwd ethanol >/dev/null; then
    # Adding system user: ethanol.
    adduser \
      --system \
      --disabled-login \
      --ingroup ethanol \
      --no-create-home \
      --home /nonexistent \
      --gecos "Ethanol Server" \
      --shell /bin/false \
      ethanol  >/dev/null
  fi
fi
' > "$BUILDDIR""$DEBIANDIR"/preinst

# generate postinst script
echo '#!/bin/bash
 
if [ "$1" = configure ]; then
  echo "Take ownership of directory..."
  chown -R ethanol:ethanol /opt/ethanol

  echo "Symlink binary to path..."
  ln -s /opt/ethanol/bin/ethanol /usr/local/bin

  echo "Reload systemd..."
  systemctl daemon-reload

  echo "Enable ethanol.service..."
  systemctl enable ethanol.service

  echo "Start ethanol.service..."
  systemctl start ethanol.service
fi
' > "$BUILDDIR""$DEBIANDIR"/postinst


# generate prerm script
echo '#!/bin/bash
 
if [ -e "/etc/systemd/system/ethanol.service" ]; then
  echo "Stopping ethanol.service..."
  systemctl stop ethanol.service

  echo "Disabling ethanol.service..."
  systemctl disable ethanol.service
fi
 
if pgrep ethanol >/dev/null
then
  echo "Stopping ethanol..."
  killall ethanol
fi
' > "$BUILDDIR""$DEBIANDIR"/prerm

# generate postrm script
echo '#!/bin/bash
 
if [ $1 = remove ]; then
  echo "Removing ethanol user..."
  userdel ethanol -r -f
 
  echo "Deleting /opt/ethanol..."
  rm -rf /opt/ethanol
   
  echo "Removing symlinks..."
  unlink /usr/local/bin/ethanol
   
  echo "Removing cfg files and service..."
  rm /etc/systemd/system/ethanol.service
  rm -rf /etc/ethanol
fi
' > "$BUILDDIR""$DEBIANDIR"/postrm

chmod 0755 "$BUILDDIR""$DEBIANDIR"/preinst
chmod 0755 "$BUILDDIR""$DEBIANDIR"/postinst
chmod 0755 "$BUILDDIR""$DEBIANDIR"/prerm
chmod 0755 "$BUILDDIR""$DEBIANDIR"/postrm

# Build
dpkg -b "$BUILDDIR" "$RELEASEDIR"/ethanol_0.0.1_amd64.deb
echo ""
dpkg -I "$RELEASEDIR"/ethanol_0.0.1_amd64.deb
