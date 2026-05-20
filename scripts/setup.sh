#!/bin/bash
# open-source-vpn Setup Script
# This script helps set up open-source-vpn server or client

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "This script must be run as root"
        exit 1
    fi
}

check_dependencies() {
    print_info "Checking dependencies..."
    
    local missing_deps=()
    
    # Check for required commands
    for cmd in go git make ip iptables; do
        if ! command -v $cmd &> /dev/null; then
            missing_deps+=($cmd)
        fi
    done
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing dependencies: ${missing_deps[*]}"
        print_info "Install them with:"
        print_info "  Ubuntu/Debian: apt-get install -y ${missing_deps[*]}"
        print_info "  CentOS/RHEL: yum install -y ${missing_deps[*]}"
        exit 1
    fi
    
    print_info "All dependencies satisfied"
}

enable_ip_forwarding() {
    print_info "Enabling IP forwarding..."
    
    sysctl -w net.ipv4.ip_forward=1 > /dev/null
    sysctl -w net.ipv6.conf.all.forwarding=1 > /dev/null
    
    # Make permanent
    if ! grep -q "net.ipv4.ip_forward=1" /etc/sysctl.conf; then
        echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
    fi
    if ! grep -q "net.ipv6.conf.all.forwarding=1" /etc/sysctl.conf; then
        echo "net.ipv6.conf.all.forwarding=1" >> /etc/sysctl.conf
    fi
    
    print_info "IP forwarding enabled"
}

load_tun_module() {
    print_info "Loading TUN module..."
    
    if ! lsmod | grep -q "^tun"; then
        modprobe tun
    fi
    
    # Make permanent
    if [ ! -f /etc/modules-load.d/tun.conf ]; then
        echo "tun" > /etc/modules-load.d/tun.conf
    fi
    
    print_info "TUN module loaded"
}

configure_firewall() {
    local port=$1
    print_info "Configuring firewall for port $port/udp..."
    
    # Try ufw first
    if command -v ufw &> /dev/null; then
        ufw allow $port/udp
        print_info "UFW rule added"
    # Try firewalld
    elif command -v firewall-cmd &> /dev/null; then
        firewall-cmd --permanent --add-port=$port/udp
        firewall-cmd --reload
        print_info "Firewalld rule added"
    else
        print_warn "No firewall manager found, please manually allow port $port/udp"
    fi
}

setup_nat() {
    local vpn_network=$1
    local interface=$2
    
    print_info "Setting up NAT for $vpn_network on $interface..."
    
    # Add iptables rules
    iptables -t nat -A POSTROUTING -s $vpn_network -o $interface -j MASQUERADE
    iptables -A FORWARD -i tun0 -j ACCEPT
    iptables -A FORWARD -o tun0 -j ACCEPT
    
    # Try to save rules
    if command -v iptables-save &> /dev/null; then
        if [ -d /etc/iptables ]; then
            iptables-save > /etc/iptables/rules.v4
            print_info "IPtables rules saved"
        else
            print_warn "Could not save iptables rules, they will be lost on reboot"
        fi
    fi
}

build_binaries() {
    print_info "Building binaries..."
    
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Are you in the project directory?"
        exit 1
    fi
    
    make build
    
    if [ -f "bin/vpn-server" ] && [ -f "bin/vpn-client" ]; then
        print_info "Binaries built successfully"
    else
        print_error "Failed to build binaries"
        exit 1
    fi
}

generate_server_keys() {
    local output_file=$1
    print_info "Generating server keys..."
    
    ./bin/vpn-server keygen --output $output_file
    
    print_info "Server keys saved to $output_file"
}

generate_client_keys() {
    local output_file=$1
    print_info "Generating client keys..."
    
    ./bin/vpn-client keygen --output $output_file
    
    print_info "Client keys saved to $output_file"
}

create_systemd_service() {
    print_info "Creating systemd service..."
    
    cat > /etc/systemd/system/vpn-server.service <<EOF
[Unit]
Description=open-source-vpn Server
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/vpn-server start --config /etc/vpn/server.yaml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    print_info "Systemd service created"
}

setup_server() {
    print_info "Setting up VPN server..."
    
    check_root
    check_dependencies
    enable_ip_forwarding
    load_tun_module
    
    # Build binaries
    build_binaries
    
    # Generate keys
    generate_server_keys "server-keys.json"
    
    # Copy configuration
    mkdir -p /etc/vpn
    if [ ! -f /etc/vpn/server.yaml ]; then
        cp configs/server.example.yaml /etc/vpn/server.yaml
        print_warn "Configuration template copied to /etc/vpn/server.yaml"
        print_warn "Please edit it and add your private key"
    fi
    
    # Install binary
    cp bin/vpn-server /usr/local/bin/
    chmod +x /usr/local/bin/vpn-server
    
    # Configure firewall
    read -p "Enter VPN port (default: 51820): " port
    port=${port:-51820}
    configure_firewall $port
    
    # Setup NAT
    read -p "Setup NAT? (y/n): " setup_nat_choice
    if [ "$setup_nat_choice" = "y" ]; then
        read -p "Enter VPN network (default: 10.8.0.0/24): " vpn_network
        vpn_network=${vpn_network:-10.8.0.0/24}
        
        read -p "Enter external interface (default: eth0): " ext_interface
        ext_interface=${ext_interface:-eth0}
        
        setup_nat $vpn_network $ext_interface
    fi
    
    # Create systemd service
    read -p "Create systemd service? (y/n): " create_service
    if [ "$create_service" = "y" ]; then
        create_systemd_service
        print_info "Enable service with: systemctl enable vpn-server"
        print_info "Start service with: systemctl start vpn-server"
    fi
    
    print_info "Server setup complete!"
    print_info "Next steps:"
    print_info "  1. Edit /etc/vpn/server.yaml and add your private key"
    print_info "  2. Add client public keys to the configuration"
    print_info "  3. Start the server: systemctl start vpn-server"
}

setup_client() {
    print_info "Setting up VPN client..."
    
    check_root
    check_dependencies
    load_tun_module
    
    # Build binaries
    build_binaries
    
    # Generate keys
    generate_client_keys "client-keys.json"
    
    # Copy configuration
    mkdir -p /etc/vpn
    if [ ! -f /etc/vpn/client.yaml ]; then
        cp configs/client.example.yaml /etc/vpn/client.yaml
        print_warn "Configuration template copied to /etc/vpn/client.yaml"
        print_warn "Please edit it and add your keys and server address"
    fi
    
    # Install binary
    cp bin/vpn-client /usr/local/bin/
    chmod +x /usr/local/bin/vpn-client
    
    print_info "Client setup complete!"
    print_info "Next steps:"
    print_info "  1. Edit /etc/vpn/client.yaml and add your keys"
    print_info "  2. Add server address and public key"
    print_info "  3. Connect: vpn-client connect --config /etc/vpn/client.yaml"
}

# Main menu
main() {
    echo "================================"
    echo "  open-source-vpn Setup Script"
    echo "================================"
    echo ""
    echo "1) Setup Server"
    echo "2) Setup Client"
    echo "3) Exit"
    echo ""
    read -p "Choose an option: " choice
    
    case $choice in
        1)
            setup_server
            ;;
        2)
            setup_client
            ;;
        3)
            print_info "Exiting..."
            exit 0
            ;;
        *)
            print_error "Invalid option"
            exit 1
            ;;
    esac
}

main
