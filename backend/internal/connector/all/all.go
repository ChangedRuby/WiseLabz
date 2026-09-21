// Package all registers every connector implementation via blank imports.
// Import this package once (from the router or main) to activate all connectors.
// To add a connector: add one blank import line below and implement the
// connector package with an init() that calls connector.Register().
package all

// Register all connectors via side-effect imports.
import (
	_ "github.com/WiseLabz/wiselabz/internal/connector/adguardhome" // register AdGuard Home connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/cloudflare"  // register Cloudflare connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/custom"      // register custom connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/dnsresolver" // register DNS Resolver connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/docker"      // register docker connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/netbird"     // register Netbird connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/opnsense"    // register OPNsense connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/pfsense"     // register pfSense connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/pihole"      // register Pi-hole connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/portainer"   // register Portainer connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/proxmox"     // register Proxmox connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/traefik"     // register Traefik connector
	_ "github.com/WiseLabz/wiselabz/internal/connector/unifi"       // register UniFi connector
)
