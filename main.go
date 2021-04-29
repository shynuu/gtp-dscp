package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/shynuu/gtp-dscp/u32"
	"github.com/urfave/cli/v2"
)

func FindInterface(ipV4 string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error reading interfaces: %+v\n", err.Error())
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Printf("localAddresses: %+v\n", err.Error())
			return "", err
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				if v.IP.To4().String() == ipV4 {
					return i.Name, nil
				}
			}

		}
	}

	return "", errors.New("ERROR interface not found")
}

func run(ipv4 string, offset int) {

	var dscp []uint8 = []uint8{0x0a, 0x0c, 0x2e}

	for _, d := range dscp {

		protocols := []u32.Protocol{
			&u32.IPV4Header{
				Protocol: u32.PROTO_UDP,
				Set: &u32.IPV4Fields{
					Protocol: true,
				},
			},
			&u32.UDPHeader{
				SourcePort:      2152,
				DestinationPort: 2152,
				Set: &u32.UDPFields{
					SourcePort:      true,
					DestinationPort: true,
				},
			},
			&u32.GTPv1Header{
				HeaderOffset: offset,
				Version:      1,
				ProtocolType: 1,
				Set: &u32.GTPv1Fields{
					Version:      true,
					ProtocolType: true,
				},
			},
			&u32.IPV4Header{
				Version: 4,
				DSCP:    d,
				Set: &u32.IPV4Fields{
					DSCP:    true,
					Version: true,
				},
			},
		}

		iface, _ := FindInterface(ipv4)

		var m = u32.NewU32(&protocols, d)
		m.RunIface(iface)

	}

}

func main() {

	var ipv4 string = ""
	var offset int

	app := &cli.App{
		Name:  "gtp-qos",
		Usage: "Copy DSCP field of inner packet to outer packet",
		Authors: []*cli.Author{
			{Name: "Youssouf Drif"},
		},
		Copyright: "Copyright (c) 2021 Youssouf Drif",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ipv4",
				Usage:       "Select the interface",
				Destination: &ipv4,
				Required:    true,
				DefaultText: "unset",
			},
			&cli.IntFlag{
				Name:        "offset",
				Usage:       "Set the offset for GTP Header",
				Destination: &offset,
				Required:    true,
				DefaultText: "unset",
			},
		},
		Action: func(c *cli.Context) error {
			run(ipv4, offset)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
