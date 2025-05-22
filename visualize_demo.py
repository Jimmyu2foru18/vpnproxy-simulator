import time
import random
import os
import sys

COLORS = {
    'green': '\033[32m',
    'cyan': '\033[36m',
    'magenta': '\033[35m',
    'reset': '\033[0m',
    'yellow': '\033[33m',
    'bold': '\033[1m'
}

def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')

def print_header():
    print(f"{COLORS['bold']}VPNProxy Simulator - Visual Demo{COLORS['reset']}")
    print("=" * 40)
    print()

def print_connection(client, proxy, target):
    print(f"{COLORS['green']}Client [{client}]{COLORS['reset']} <---> {COLORS['cyan']}Proxy [{proxy}]{COLORS['reset']} <---> {COLORS['magenta']}Target [{target}]{COLORS['reset']}")
    print()

def visualize_data_flow(direction, data, size):
    if direction == "client->proxy":
        source_color = COLORS['green']
        dest_color = COLORS['cyan']
        arrow = "====>"
        source = "Client"
        dest = "Proxy"
    elif direction == "proxy->target":
        source_color = COLORS['cyan']
        dest_color = COLORS['magenta']
        arrow = "====>"
        source = "Proxy"
        dest = "Target"
    elif direction == "target->proxy":
        source_color = COLORS['magenta']
        dest_color = COLORS['cyan']
        arrow = "<===="
        source = "Target"
        dest = "Proxy"
    else:  
        source_color = COLORS['cyan']
        dest_color = COLORS['green']
        arrow = "<===="
        source = "Proxy"
        dest = "Client"

    print(f"{source_color}{source}{COLORS['reset']} {arrow} {dest_color}{dest}{COLORS['reset']} [{size} bytes]")

    hex_view = ' '.join([f"{b:02x}" for b in data[:16]])
    if len(data) > 16:
        hex_view += " ..."
    print(f"HEX: {hex_view}")
 
    text_view = ''
    for b in data[:32]:
        if 32 <= b <= 126:
            text_view += chr(b)
        else:
            text_view += '.'
    if len(data) > 32:
        text_view += "..."
    print(f"TXT: {text_view}\n")

def simulate_http_request():
    # Simulate HTTP request
    request = b"GET / HTTP/1.1\r\nHost: example.com\r\nUser-Agent: VPNProxy-Simulator\r\nConnection: close\r\n\r\n"
    visualize_data_flow("client->proxy", request, len(request))
    time.sleep(0.5)
    
    visualize_data_flow("proxy->target", request, len(request))
    time.sleep(1)
    
    # Simulate HTTP response (simplified)
    response_header = b"HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 1256\r\n\r\n"
    visualize_data_flow("target->proxy", response_header, len(response_header))
    time.sleep(0.5)
    
    response_body = b"<!DOCTYPE html><html><head><title>Example Domain</title></head><body><h1>Example Domain</h1><p>This domain is for use in illustrative examples in documents.</p></body></html>"
    visualize_data_flow("target->proxy", response_body, len(response_body))
    time.sleep(0.5)
    
    visualize_data_flow("proxy->client", response_header, len(response_header))
    time.sleep(0.3)
    
    visualize_data_flow("proxy->client", response_body, len(response_body))

def simulate_connection_metrics():
    print(f"\n{COLORS['yellow']}Connection Metrics:{COLORS['reset']}")
    print(f"  Total Bytes In:  {random.randint(1000, 5000)}")
    print(f"  Total Bytes Out: {random.randint(500, 2000)}")
    print(f"  Duration:        {random.randint(1, 5)}.{random.randint(100, 999)} seconds")
    print()

def main():
    clear_screen()
    print_header()
    
    print("Initializing VPN/Proxy Simulator...")
    time.sleep(1)
    
    print("Starting proxy server on port 8080...")
    time.sleep(1)
    
    print("Proxy server ready and listening for connections.")
    time.sleep(1)
    
    print("\nClient connecting to proxy...")
    time.sleep(1)
    
    client_ip = "192.168.1." + str(random.randint(100, 200))
    proxy_ip = "127.0.0.1:8080"
    target_ip = "93.184.216.34"
    
    print_connection(client_ip, proxy_ip, target_ip)
    time.sleep(1)
    
    print("Establishing secure TLS connection...")
    time.sleep(1)
    
    print("Connection established. Starting data transfer...\n")
    time.sleep(1)
    
    simulate_http_request()
    
    simulate_connection_metrics()
    
    print("Connection closed.")
    print("\nThank you for using VPNProxy Simulator!")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\nSimulation terminated by user.")
    except Exception as e:
        print(f"\n\nError: {e}")
    finally:
        print(COLORS['reset'])