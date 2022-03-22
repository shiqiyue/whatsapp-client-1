package whatsapp

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"whatsapp-client/model"
)

var container *sqlstore.Container

func init() {
	name := "Windows"
	store.CompanionProps.Os = &name
	store.CompanionProps.PlatformType = waProto.CompanionProps_CHROME.Enum()

	dbLog := waLog.Stdout("Database", "DEBUG", true)

	db, err := model.DB.DB.DB()
	if err != nil {
		panic(err)
	}
	container = sqlstore.NewWithDB(db, "sqlite3", dbLog)
	err = container.Upgrade()
	if err != nil {
		panic(err)
	}
	go func() {
		devices, err := container.GetAllDevices()
		if err != nil || len(devices) == 0 {
			return
		}
		client, qrChan := NewClient(devices[0].ID.String())
		if qrChan == nil {
			client.Login()
		}
	}()
}
