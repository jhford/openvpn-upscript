package iptables

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	// DefaultIPTables is the default iptables command
	DefaultIPTables = "iptables"
)

// Config is configuration for ip tables modifications
type Config struct {
	TunDevice   string
	DestDevice  string
	FlushTables bool
	IPTables    string
}

// GenerateCommands returns a list of ordered argvs needed to configure
// the firewall
func (c Config) GenerateCommands() [][]string {
	commands := make([][]string, 0)

	if c.IPTables == "" {
		c.IPTables = DefaultIPTables
	}

	if c.FlushTables {
		commands = append(commands, []string{c.IPTables, "-F"})
		commands = append(commands, []string{c.IPTables, "-t", "nat", "-F"})
		commands = append(commands, []string{c.IPTables, "-X"})
	}

	commands = append(commands, []string{c.IPTables, "-t", "nat", "-A", "POSTROUTING", "-o", c.TunDevice, "-j", "MASQUERADE"})
	commands = append(commands, []string{c.IPTables, "-A", "FORWARD", "-i", c.TunDevice, "-o", c.DestDevice, "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT"})
	commands = append(commands, []string{c.IPTables, "-A", "FORWARD", "-i", c.DestDevice, "-o", c.TunDevice, "-j", "ACCEPT"})

	return commands
}

// Apply applies a configuration
func (c Config) Apply() error {
	commands := c.GenerateCommands()

	for _, command := range commands {
		cmd := exec.Command(command[0], command[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Libraries shouldn't write to stdio, but...
			fmt.Fprintf(os.Stderr, "ERROR: %v %s\n%s\n", command, err, output)
			return err
		}
	}

	return nil
}
