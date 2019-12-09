package main

import (
	"errors"
	"fmt"
	"github.com/jhford/openvpn-up/pkg"
	"github.com/jhford/openvpn-up/pkg/config"
	"github.com/jhford/openvpn-up/pkg/iptables"
	"github.com/jhford/openvpn-up/pkg/resolv"
	"github.com/urfave/cli"
	"log"
	"net"
	"os"
	"strings"
)

func bail(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "ERROR: %s\n", msg)
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func main() {
	app := cli.NewApp()

	app.Version = pkg.Version + "-" + pkg.Commit
	app.Name = pkg.Name
	app.Description = "an OpenVPN --up helper program to build a resolv.conf file and route traffic to ethernet device"
	app.Usage = "use as an openvpn --up script"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "tundev",
			EnvVar: "TUNDEV",
			Usage:  "tunnel device source",
		},
		cli.StringFlag{
			Name:   "destdev",
			EnvVar: "DESTDEV",
			Value:  "br0",
			Usage:  "destination device to route traffic to",
		},
		cli.StringFlag{
			Name:   "resolvconf",
			EnvVar: "RESOLVCONF",
			Value:  "/tmp/openvpn-up-resolv.conf",
			Usage:  "file path to the output resolv.conf file",
		},
		cli.BoolFlag{
			Name:   "flush",
			EnvVar: "FLUSH",
			Usage:  "flush iptables before setting up routing",
		},
		cli.StringFlag{
			Name:   "iptables",
			EnvVar: "IPTABLES",
			Value:  "iptables",
			Usage:  "custom command for iptables.  PATH resolution is not supported",
		},
		cli.StringSliceFlag{
			Name:   "append-nameserver",
			EnvVar: "APPEND_NAMESERVER",
			Usage:  "manually appended name servers to add to resolv.conf",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		var tundev string
		var destdev string

		if ctx.IsSet("tundev") && ctx.NArg() < 1 {
			tundev = ctx.String("tundev")
		} else if ctx.NArg() > 0 {
			tundev = ctx.Args()[0]
		} else {
			return errors.New("no tunnel device specified")
		}

		destdev = ctx.String("destdev")

		log.Printf("tunnel device: '%s' destination device: '%s' args: '%s'",
			tundev, destdev, strings.Join(ctx.Args(), "', '"))

		cfg := config.Config{
			DestDevice: ctx.String("destdev"),
		}

		err := cfg.ParseEnv()
		if err != nil {
			bail(err.Error())
		}

		if ctx.IsSet("nameservers") {
			for _, nameserver := range ctx.StringSlice("nameservers") {
				nsip := net.ParseIP(nameserver)
				if nsip == nil {
					return fmt.Errorf("invalid nameserver ip specified on command line: %s", nameserver)
				}
				cfg.NameServers = append(cfg.NameServers, pkg.NameServer(nsip))
			}
		}

		resolvConf := resolv.Config{
			NameServers: cfg.NameServers,
		}

		f, err := os.Create(ctx.String("resolvconf"))
		if err != nil {
			return err
		}
		_, err = f.Write(resolvConf.GenerateFile())
		if err != nil {
			return err
		}
		log.Printf("wrote resolv.conf to %s", ctx.String("resolvconf"))

		rules := iptables.Config{
			TunDevice:   cfg.TunDevice,
			DestDevice:  cfg.DestDevice,
			FlushTables: ctx.Bool("flush"),
			IPTables:    ctx.String("iptables"),
		}

		if err = rules.Apply(); err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		bail(err.Error())
	}
}
