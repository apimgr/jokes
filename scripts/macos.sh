#!/usr/bin/env bash
#===============================================================================
# Jokes API - macOS Installation Script
#===============================================================================

set -e

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

# Variables
PROJECT_NAME="jokes"
BINARY_NAME="jokes-api"
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/usr/local/share/apimgr/$PROJECT_NAME"
CONFIG_DIR="/usr/local/etc/apimgr/$PROJECT_NAME"
LOG_DIR="/usr/local/var/log/apimgr/$PROJECT_NAME"
PLIST_FILE="/Library/LaunchDaemons/com.apimgr.$PROJECT_NAME.plist"

# Find unused UID/GID between 100-999
find_unused_uid() {
    for uid in {100..999}; do
        if ! dscl . -read /Users/_$PROJECT_NAME_$uid >/dev/null 2>&1; then
            echo $uid
            return
        fi
    done
    echo 999
}

# Create system user
create_user() {
    echo -e "${GREEN}→${NC} Creating system user..."

    UID_GID=$(find_unused_uid)
    USERNAME="_${PROJECT_NAME}_${UID_GID}"

    if ! dscl . -read /Users/$USERNAME >/dev/null 2>&1; then
        sudo dscl . -create /Users/$USERNAME
        sudo dscl . -create /Users/$USERNAME UserShell /usr/bin/false
        sudo dscl . -create /Users/$USERNAME UniqueID $UID_GID
        sudo dscl . -create /Users/$USERNAME PrimaryGroupID $UID_GID
        sudo dscl . -create /Users/$USERNAME NFSHomeDirectory $DATA_DIR
    fi
}

# Create directories
create_directories() {
    echo -e "${GREEN}→${NC} Creating directories..."

    sudo mkdir -p $INSTALL_DIR
    sudo mkdir -p $DATA_DIR
    sudo mkdir -p $CONFIG_DIR
    sudo mkdir -p $LOG_DIR

    sudo chown -R $USERNAME:staff $DATA_DIR
    sudo chown -R $USERNAME:staff $CONFIG_DIR
    sudo chown -R $USERNAME:staff $LOG_DIR
}

# Install binary
install_binary() {
    echo -e "${GREEN}→${NC} Installing binary..."

    if [ -f "./binaries/$BINARY_NAME" ]; then
        sudo cp "./binaries/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    elif [ -f "./$BINARY_NAME" ]; then
        sudo cp "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    else
        echo "Error: Binary not found"
        exit 1
    fi

    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
}

# Install data files
install_data() {
    echo -e "${GREEN}→${NC} Installing data files..."

    if [ -d "./src/data" ]; then
        sudo cp -r ./src/data/* $DATA_DIR/
    fi

    sudo chown -R $USERNAME:staff $DATA_DIR
}

# Create LaunchDaemon
create_service() {
    echo -e "${GREEN}→${NC} Creating LaunchDaemon..."

    sudo tee $PLIST_FILE > /dev/null <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.apimgr.$PROJECT_NAME</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_DIR/$BINARY_NAME</string>
    </array>
    <key>UserName</key>
    <string>$USERNAME</string>
    <key>WorkingDirectory</key>
    <string>$DATA_DIR</string>
    <key>StandardOutPath</key>
    <string>$LOG_DIR/stdout.log</string>
    <key>StandardErrorPath</key>
    <string>$LOG_DIR/stderr.log</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

    sudo chown root:wheel $PLIST_FILE
    sudo chmod 644 $PLIST_FILE
    sudo launchctl load $PLIST_FILE
}

# Main installation
main() {
    echo "📦 Installing Jokes API for macOS..."

    create_user
    create_directories
    install_binary
    install_data
    create_service

    echo -e "${GREEN}✓${NC} Service installed and started"
    echo ""
    echo "Service commands:"
    echo "  Status:  sudo launchctl list | grep $PROJECT_NAME"
    echo "  Stop:    sudo launchctl unload $PLIST_FILE"
    echo "  Start:   sudo launchctl load $PLIST_FILE"
    echo "  Logs:    tail -f $LOG_DIR/stdout.log"
}

main "$@"
