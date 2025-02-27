package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/emiago/diago"
	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
	"gopkg.in/ini.v1"
)

type Connection struct {
	hostname string
	port     int
	username string
	password string
}

func ParseClientAccount(fn string) (*Connection, *diago.RegisterOptions, error) {
	raw, err := ini.Load(fn)
	if err != nil {
		return nil, nil, err
	}

	port, err := raw.Section("client").Key("port").Int()
	if err != nil {
		return nil, nil, err
	}
	con := &Connection{
		hostname: raw.Section("client").Key("hostname").String(),
		port:     port,
		username: raw.Section("client").Key("username").String(),
		password: raw.Section("client").Key("password").String(),
	}

	regOpts := &diago.RegisterOptions{
		Username: con.username,
		Password: con.password,
	}

	return con, regOpts, nil
}

func ParseForwardIP(fn string) (*string, error) {
	raw, err := ini.Load(fn)
	if err != nil {
		return nil, err
	}

	result := raw.Section("forwarding").Key("hostname").String()
	if result == "" {
		result = "0.0.0.0"
	}

	return &result, nil
}

func Setup(ctx context.Context, con Connection, forwardingIP string, regOpts diago.RegisterOptions) error {
	recipient := sip.Uri{}
	err := sip.ParseUri(
		fmt.Sprintf(
			"sip:%s@%s:%d",
			// "sip:%s@%s",
			con.username,
			con.hostname,
			con.port,
		),
		&recipient,
	)
	if err != nil {
		return fmt.Errorf("Account Information Parse Error - %v", err)
	}

	ua, err := sipgo.NewUA(
		sipgo.WithUserAgent(con.username),
		sipgo.WithUserAgentHostname(fmt.Sprintf("%s", con.hostname)),
	)
	if err != nil {
		return fmt.Errorf("User Agent Parse Error - %v", err)
	}
	defer ua.Close()

	tu := diago.NewDiago(ua, diago.WithTransport(
		diago.Transport{
			Transport: "udp",
			BindHost:  forwardingIP,
			BindPort:  0,
		},
	))

	go func() {
		tu.Serve(ctx, func(inDialog *diago.DialogServerSession) {
			log.Printf("[New Dialog Req.] ID: %s", inDialog.ID)
			defer log.Printf("[Finish Dialog] ID: %s", inDialog.ID)
		})
	}()

	err = tu.Register(ctx, recipient, regOpts)
	if err != nil {
		return fmt.Errorf("Failed to register - %v", err)
	}

	return nil
}
