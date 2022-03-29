package whatsapp

import (
	"context"
	"github.com/mattn/go-ieproxy"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Client struct {
	*whatsmeow.Client
	groups    []*types.GroupInfo
	autoReply AutoReply
}

var onlineClients []*Client

func NewClient(id, proxyStr string) (*Client, <-chan whatsmeow.QRChannelItem) {
	var device *store.Device
	if id == "" {
		device = container.NewDevice()
	} else {
		jid, err := types.ParseJID(id)
		if err != nil {
			panic(err)
		}
		device, err = container.GetDevice(jid)

		if err != nil {
			panic(err)
		}
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := &Client{Client: whatsmeow.NewClient(device, clientLog)}
	if proxyStr != "" {
		_ = client.SetProxyAddress(proxyStr)
	} else {
		client.SetProxy(ieproxy.GetProxyFunc())
	}
	client.EnableAutoReply()
	client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.ClientOutdated:
		case *events.LoggedOut:
			for i, c := range onlineClients {
				if c.Store.ID == client.Store.ID {
					onlineClients = append(onlineClients[:i], onlineClients[i+1:]...)
					return
				}
			}
		}
	})

	if client.Store.ID == nil {
		c, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		return client, c
	} else {
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		return client, nil
	}
}

func GetClients() []*Client {
	return onlineClients
}

func GetClient(id string) *Client {
	for _, client := range onlineClients {
		if client.Store.ID.String() == id {
			return client
		}
	}
	return nil
}

func (c *Client) Login() {
	onlineClients = append(onlineClients, c)
}

func (c *Client) Logout() {
	c.Disconnect()
	for i, client := range onlineClients {
		if client.Store.ID.String() == c.Store.ID.String() {
			onlineClients = append(onlineClients[:i], onlineClients[i+1:]...)
			return
		}
	}
}

func (c *Client) Phone() string {
	return c.Store.ID.User
}
