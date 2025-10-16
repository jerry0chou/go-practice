package net

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type NetworkInfo struct {
	LocalAddr   string
	RemoteAddr  string
	Protocol    string
	IsConnected bool
}

func ResolveHostname(hostname string) ([]string, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve hostname %s: %w", hostname, err)
	}

	var ipStrings []string
	for _, ip := range ips {
		ipStrings = append(ipStrings, ip.String())
	}

	return ipStrings, nil
}

func ReverseDNS(ip string) ([]string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil, fmt.Errorf("failed to perform reverse DNS lookup for %s: %w", ip, err)
	}

	return names, nil
}

func CheckPortAvailability(host, port string) bool {
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

func GetLocalIPs() ([]string, error) {
	var ips []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}

	return ips, nil
}

func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

func IsValidIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To16() != nil && parsedIP.To4() == nil
}

func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	privateBlocks := []*net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
		{IP: net.IPv4(127, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
	}

	for _, block := range privateBlocks {
		if block.Contains(parsedIP) {
			return true
		}
	}

	return false
}

func PingHost(host string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", host, err)
	}
	defer conn.Close()
	return nil
}

func GetNetworkInterfaces() ([]map[string]interface{}, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	var result []map[string]interface{}
	for _, iface := range interfaces {
		info := map[string]interface{}{
			"name":          iface.Name,
			"index":         iface.Index,
			"mtu":           iface.MTU,
			"hardware_addr": iface.HardwareAddr.String(),
			"flags":         iface.Flags.String(),
			"up":            iface.Flags&net.FlagUp != 0,
			"loopback":      iface.Flags&net.FlagLoopback != 0,
			"multicast":     iface.Flags&net.FlagMulticast != 0,
		}

		addrs, err := iface.Addrs()
		if err == nil {
			var addresses []string
			for _, addr := range addrs {
				addresses = append(addresses, addr.String())
			}
			info["addresses"] = addresses
		}

		result = append(result, info)
	}

	return result, nil
}

func ParseCIDR(cidr string) (string, string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse CIDR %s: %w", cidr, err)
	}

	return ip.String(), ipnet.String(), nil
}

func IsIPInCIDR(ip, cidr string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	return ipnet.Contains(parsedIP)
}

func PrintNetworkInfo() {
	fmt.Println("🌐 Network Information")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Println("📍 Local IP Addresses:")
	localIPs, err := GetLocalIPs()
	if err != nil {
		fmt.Printf("❌ Error getting local IPs: %v\n", err)
	} else {
		for _, ip := range localIPs {
			ipType := "IPv4"
			if !IsValidIPv4(ip) {
				ipType = "IPv6"
			}
			privacy := "Public"
			if IsPrivateIP(ip) {
				privacy = "Private"
			}
			fmt.Printf("   %s (%s, %s)\n", ip, ipType, privacy)
		}
	}

	fmt.Println("\n🔌 Network Interfaces:")
	interfaces, err := GetNetworkInterfaces()
	if err != nil {
		fmt.Printf("❌ Error getting network interfaces: %v\n", err)
	} else {
		for _, iface := range interfaces {
			fmt.Printf("   %s (Index: %d, MTU: %d)\n",
				iface["name"], iface["index"], iface["mtu"])
			fmt.Printf("     Hardware: %s\n", iface["hardware_addr"])
			fmt.Printf("     Flags: %s\n", iface["flags"])

			if addresses, ok := iface["addresses"].([]string); ok {
				for _, addr := range addresses {
					fmt.Printf("     Address: %s\n", addr)
				}
			}
			fmt.Println()
		}
	}
}

func DemonstrateNetworkOperations() {
	fmt.Println("🔧 Network Operations Demo")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("🔍 Hostname Resolution:")
	fmt.Println(strings.Repeat("-", 30))

	hostnames := []string{
		"google.com",
		"github.com",
		"localhost",
		"invalid-hostname-12345.com",
	}

	for _, hostname := range hostnames {
		ips, err := ResolveHostname(hostname)
		if err != nil {
			fmt.Printf("❌ %s: %v\n", hostname, err)
		} else {
			fmt.Printf("✅ %s: %s\n", hostname, strings.Join(ips, ", "))
		}
	}

	fmt.Println("\n🔄 Reverse DNS Lookup:")
	fmt.Println(strings.Repeat("-", 30))

	ips := []string{
		"8.8.8.8",
		"1.1.1.1",
		"127.0.0.1",
		"192.168.1.1",
	}

	for _, ip := range ips {
		names, err := ReverseDNS(ip)
		if err != nil {
			fmt.Printf("❌ %s: %v\n", ip, err)
		} else {
			fmt.Printf("✅ %s: %s\n", ip, strings.Join(names, ", "))
		}
	}

	fmt.Println("\n✅ IP Address Validation:")
	fmt.Println(strings.Repeat("-", 30))

	testIPs := []string{
		"192.168.1.1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"256.256.256.256",
		"not-an-ip",
		"127.0.0.1",
	}

	for _, ip := range testIPs {
		valid := IsValidIP(ip)
		ipv4 := IsValidIPv4(ip)
		ipv6 := IsValidIPv6(ip)
		private := IsPrivateIP(ip)

		fmt.Printf("IP: %s\n", ip)
		fmt.Printf("  Valid: %t, IPv4: %t, IPv6: %t, Private: %t\n",
			valid, ipv4, ipv6, private)
		fmt.Println()
	}

	fmt.Println("🚪 Port Availability Check:")
	fmt.Println(strings.Repeat("-", 30))

	ports := []string{"80", "443", "22", "8080", "3000"}
	host := "localhost"

	for _, port := range ports {
		available := CheckPortAvailability(host, port)
		status := "❌ In Use"
		if available {
			status = "✅ Available"
		}
		fmt.Printf("%s: %s:%s\n", status, host, port)
	}

	fmt.Println("\n📊 CIDR Operations:")
	fmt.Println(strings.Repeat("-", 30))

	cidrs := []string{
		"192.168.1.0/24",
		"10.0.0.0/8",
		"172.16.0.0/12",
	}

	for _, cidr := range cidrs {
		ip, network, err := ParseCIDR(cidr)
		if err != nil {
			fmt.Printf("❌ Error parsing %s: %v\n", cidr, err)
		} else {
			fmt.Printf("✅ %s -> IP: %s, Network: %s\n", cidr, ip, network)
		}
	}

	fmt.Println("\n🎯 IP in CIDR Check:")
	fmt.Println(strings.Repeat("-", 30))

	testCases := []struct {
		ip   string
		cidr string
	}{
		{"192.168.1.100", "192.168.1.0/24"},
		{"10.0.0.5", "192.168.1.0/24"},
		{"172.16.1.1", "172.16.0.0/12"},
		{"8.8.8.8", "8.8.8.0/24"},
	}

	for _, tc := range testCases {
		inRange := IsIPInCIDR(tc.ip, tc.cidr)
		status := "❌ Not in range"
		if inRange {
			status = "✅ In range"
		}
		fmt.Printf("%s: %s in %s\n", status, tc.ip, tc.cidr)
	}

	fmt.Println("\n🏓 Simple Ping Test:")
	fmt.Println(strings.Repeat("-", 30))

	hosts := []string{
		"google.com:80",
		"github.com:443",
		"localhost:22",
		"invalid-host:80",
	}

	for _, host := range hosts {
		err := PingHost(host, 3*time.Second)
		status := "❌ Unreachable"
		if err == nil {
			status = "✅ Reachable"
		}
		fmt.Printf("%s: %s\n", status, host)
	}
}
