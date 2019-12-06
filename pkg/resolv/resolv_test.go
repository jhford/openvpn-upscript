package resolv_test

import (
	"testing"

	"github.com/jhford/openvpn-helper/pkg"
	"github.com/stretchr/testify/assert"

	"github.com/jhford/openvpn-helper/pkg/resolv"
)

func TestGenerateFile(t *testing.T) {
	tests := []struct {
		input    resolv.Config
		expected []byte
	}{
		{
			input: resolv.Config{
				NameServers: []pkg.NameServer{{127, 0, 0, 1}},
			},
			expected: []byte("# generated by openvpn-helper\nnameserver 127.0.0.1\n"),
		},
	}

	for _, test := range tests {
		test := test

		actual := test.input.GenerateFile()

		assert.Equal(t, test.expected, actual)
	}
}
