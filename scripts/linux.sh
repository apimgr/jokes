#!/usr/bin/env bash
#===============================================================================
# Jokes API - Linux Installation Script
#===============================================================================

set -e

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

# Variables
PROJECT_NAME="jokes"
BINARY_NAME="jokes-api"
INSTALL_DIR="/usr/local/bin"
DATA_DIR="/usr/share/apimgr/$PROJECT_NAME"
CONFIG_DIR="/etc/apimgr/$PROJECT_NAME"
LOG_DIR="/var/log/apimgr/$PROJECT_NAME"
SERVICE_FILE="/etc/systemd/system/$BINARY_NAME.service"

# Find unused UID/GID between 100-999
find_unused_uid() {
    for uid in {100..999}; do
        if ! id -u $uid >/dev/null 2>&1; then
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

    if ! id -u $PROJECT_NAME >/dev/null 2>&1; then
        useradd -r -u $UID_GID -g $UID_GID -d $DATA_DIR -s /bin/false $PROJECT_NAME 2>/dev/null || \
        useradd -r -d $DATA_DIR -s /bin/false $PROJECT_NAME
    fi
}

# Create directories
create_directories() {
    echo -e "${GREEN}→${NC} Creating directories..."

    mkdir -p $INSTALL_DIR
    mkdir -p $DATA_DIR
    mkdir -p $CONFIG_DIR
    mkdir -p $LOG_DIR

    chown -R $PROJECT_NAME:$PROJECT_NAME $DATA_DIR
    chown -R $PROJECT_NAME:$PROJECT_NAME $CONFIG_DIR
    chown -R $PROJECT_NAME:$PROJECT_NAME $LOG_DIR
}

# Install binary
install_binary() {
    echo -e "${GREEN}→${NC} Installing binary..."

    if [ -f "./binaries/$BINARY_NAME" ]; then
        cp "./binaries/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    elif [ -f "./$BINARY_NAME" ]; then
        cp "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    else
        echo "Error: Binary not found"
        exit 1
    fi

    chmod +x "$INSTALL_DIR/$BINARY_NAME"
}

# Install data files
install_data() {
    echo -e "${GREEN}→${NC} Installing data files..."

    if [ -d "./src/data" ]; then
        cp -r ./src/data/* $DATA_DIR/
    fi

    chown -R $PROJECT_NAME:$PROJECT_NAME $DATA_DIR
}

# Create systemd service
create_service() {
    echo -e "${GREEN}→${NC} Creating systemd service..."

    cat > $SERVICE_FILE <<EOF
[Unit]
Description=Jokes API Service
After=network.target
Documentation=https://jokes.apimgr.us

[Service]
Type=simple
User=$PROJECT_NAME
Group=$PROJECT_NAME
WorkingDirectory=$DATA_DIR
ExecStart=$INSTALL_DIR/$BINARY_NAME
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$BINARY_NAME

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$DATA_DIR $CONFIG_DIR $LOG_DIR

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable $BINARY_NAME
    systemctl start $BINARY_NAME
}

# Main installation
main() {
    echo "📦 Installing Jokes API for Linux..."

    create_user
    create_directories
    install_binary
    install_data
    create_service

    echo -e "${GREEN}✓${NC} Service installed and started"
    echo ""
    echo "Service commands:"
    echo "  Status:  sudo systemctl status $BINARY_NAME"
    echo "  Stop:    sudo systemctl stop $BINARY_NAME"
    echo "  Start:   sudo systemctl start $BINARY_NAME"
    echo "  Restart: sudo systemctl restart $BINARY_NAME"
    echo "  Logs:    sudo journalctl -u $BINARY_NAME -f"
}

main "$@"
