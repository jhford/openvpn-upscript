package iptables_test

import (
	"testing"

	"github.com/jhford/openvpn-up/pkg/iptables"
	"github.com/stretchr/testify/assert"
)

func TestIPTablesGenerate(t *testing.T) {
	tests := []struct {
		input    iptables.Config
		expected [][]string
	}{
		{
			input: iptables.Config{
				FlushTables: true,
				TunDevice:   "tun0",
				DestDevice:  "br0",
				IPTables:    "testipt",
			},
			expected: [][]string{
				{"testipt", "-F"},
				{"testipt", "-t", "nat", "-F"},
				{"testipt", "-X"},
				{"testipt", "-t", "nat", "-A", "POSTROUTING", "-o", "tun0", "-j", "MASQUERADE"},
				{"testipt", "-A", "FORWARD", "-i", "tun0", "-o", "br0", "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT"},
				{"testipt", "-A", "FORWARD", "-i", "br0", "-o", "tun0", "-j", "ACCEPT"},
			},
		},
		{
			input: iptables.Config{
				TunDevice:  "tun0",
				DestDevice: "br0",
				IPTables:   "testipt",
			},
			expected: [][]string{
				{"testipt", "-t", "nat", "-A", "POSTROUTING", "-o", "tun0", "-j", "MASQUERADE"},
				{"testipt", "-A", "FORWARD", "-i", "tun0", "-o", "br0", "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT"},
				{"testipt", "-A", "FORWARD", "-i", "br0", "-o", "tun0", "-j", "ACCEPT"},
			},
		},
		{
			input: iptables.Config{
				TunDevice:  "tun0",
				DestDevice: "br0",
			},
			expected: [][]string{
				{"iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "tun0", "-j", "MASQUERADE"},
				{"iptables", "-A", "FORWARD", "-i", "tun0", "-o", "br0", "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT"},
				{"iptables", "-A", "FORWARD", "-i", "br0", "-o", "tun0", "-j", "ACCEPT"},
			},
		},
	}

	for _, test := range tests {
		test := test

		actual := test.input.GenerateCommands()

		assert.Equal(t, test.expected, actual)
	}
}
