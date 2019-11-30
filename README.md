# OpenVPN Helper

Run as an OpenVPN --up script, this program will create a resolv.conf file with
any given DNS servers.  It will also set up iptables to route traffic from the
tunnel created by OpenVPN to the specified network device (e.g. tun0-wlan0)
