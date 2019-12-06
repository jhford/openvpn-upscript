package config

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/jhford/openvpn-helper/pkg"
)

// ParseForeignOption parses foreign options.
// Example: foreign_option_1="dhcp-option DNS 209.222.18.222"
func ParseForeignOption(v string) (interface{}, error) {
	fields := strings.Fields(v)

	if len(fields) == 3 && fields[0] == "dhcp-option" && fields[1] == "DNS" {
		ip := net.ParseIP(fields[2])
		if ip == nil {
			return nil, fmt.Errorf("invalid ip for nameserver")
		}
		return pkg.NameServer(ip), nil
	}

	return nil, nil
}

// Config stores config
type Config struct {
	TunDevice   string
	DestDevice  string
	NameServers []pkg.NameServer
}

// Determine builds a new and complete Config.  Args must function and be structured
// as os.Args.  This would ideally not read the environment directly, but as there's
// no easy to lookup environment (e.g. map[string]string) copy in the stdlib, this will
// do.  In usage this makes no difference, but in testing it does limit to a single
// concurrent ParseEnv() test.  This isn't guarded with a mutex, so is racey
func (c *Config) ParseEnv() error {
	for i := 1; ; i++ {
		key := fmt.Sprintf("foreign_option_%d", i)
		if v, ok := os.LookupEnv(key); ok {
			rawOpt, err := ParseForeignOption(v)
			if err != nil {
				return err
			}

			switch opt := rawOpt.(type) {
			case pkg.NameServer:
				c.NameServers = append(c.NameServers, opt)
			}
		} else {
			break
		}
	}

	return nil
}
