package main

import (
    "database/sql"
    "fmt"
    "log"
    "net"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./dhcp.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    addr := net.UDPAddr{
        Port: 67,
        IP:   net.ParseIP("0.0.0.0"),
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    buf := make([]byte, 1024)
    for {
        n, clientAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Println("Error receiving packet:", err)
            continue
        }

        // Handle DHCP DISCOVER or REQUEST
        macAddress := extractMACAddress(buf[:n])
        ipAddress, err := getIPAddressForMAC(db, macAddress)
        if err != nil {
            log.Println("MAC address not found:", macAddress)
            continue
        }

        response := buildDHCPResponse(buf[:n], ipAddress)
        _, err = conn.WriteToUDP(response, clientAddr)
        if err != nil {
            log.Println("Error sending response:", err)
        }
    }
}

func extractMACAddress(packet []byte) string {
    // Extract the MAC address from the packet, usually in a DHCP DISCOVER it's around bytes 28-33
    return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", packet[28], packet[29], packet[30], packet[31], packet[32], packet[33])
}

func getIPAddressForMAC(db *sql.DB, macAddress string) (string, error) {
    var ipAddress string
    err := db.QueryRow("SELECT ip_address FROM dhcp_leases WHERE mac_address = ?", macAddress).Scan(&ipAddress)
    if err != nil {
        return "", err
    }
    return ipAddress, nil
}

func buildDHCPResponse(request []byte, ipAddress string) []byte {
    // Build a DHCP OFFER or ACK response based on the request
    // Set the offered IP address in the response
    response := make([]byte, len(request))
    copy(response, request)

    // Populate necessary fields for DHCP response
    // E.g., set the IP address to offer in the correct position

    return response
}
