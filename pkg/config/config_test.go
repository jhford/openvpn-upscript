package config_test

// this code is a little ugly to avoid importing anything

import (
	"net"
	"os"
	"testing"

	"github.com/jhford/openvpn-up/pkg"
	"github.com/jhford/openvpn-up/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestParseForeignOption(t *testing.T) {
	good := []string{
		"dhcp-option DNS 127.0.0.1",
		"dhcp-option DNS 127.0.0.1",
		"dhcp-option DNS 127.0.0.1",
		"  dhcp-option     DNS   127.0.0.1     ",
		"	dhcp-option	DNS	127.0.0.1	",
	}

	for _, v := range good {
		outcome, err := config.ParseForeignOption(v)
		if err != nil {
			t.Fatal(err)
		}

		if typedOutcome, ok := outcome.(pkg.NameServer); ok {
			assert.True(t, net.IP{127, 0, 0, 1}.Equal(net.IP(typedOutcome)))
		} else {
			t.Errorf("good input had invalid output: %s", v)
		}
	}
}

func TestConfig_ParseEnv(t *testing.T) {
	tests := []struct {
		envs     [][2]string
		args     []string
		expected config.Config
	}{
		{
			envs: [][2]string{{"foreign_option_1", "dhcp-option DNS 127.0.0.1"}},
			expected: config.Config{
				NameServers: []pkg.NameServer{{127, 0, 0, 1}},
			},
		},
	}

	for _, test := range tests {
		test := test

		for _, v := range test.envs {
			err := os.Setenv(v[0], v[1])
			if err != nil {
				t.Fatal(err)
			}
		}

		cfg := config.Config{}
		err := cfg.ParseEnv()
		assert.NoError(t, err)

		for _, v := range test.expected.NameServers {
			assert.True(t, net.IP(v).Equal(net.IP{127, 0, 0, 1}))
		}
		assert.Len(t, cfg.NameServers, len(test.expected.NameServers))
	}
}
