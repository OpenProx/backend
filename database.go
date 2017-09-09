package backend

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
	"github.com/asdine/storm/q"
)

// InitDatabase inits the database
func (i *Instance) InitDatabase() error {
	dbCon, err := storm.Open("openprox.db", storm.Codec(msgpack.Codec))
	if err != nil {
		return err
	}
	i.Database = dbCon
	i.Database.Init(&User{})
	i.Database.Init(&Proxy{})
	return nil
}

// HasProxy checks if a proxy exists
func (i *Instance) HasProxy(identifier string) bool {
	var dbProx Proxy
	return i.Database.One("Identifier", identifier, &dbProx) != storm.ErrNotFound
}

// GetCheckableProxy returns a new CheckRequest for a proxy
func (i *Instance) GetCheckableProxy(uid int) (*CheckRequest, error) {
	ago := time.Now().Unix() - 60*5

	var found []Proxy
	i.Database.Select(q.Lt("LastCheck", ago), q.Lt("ChecksLength", 5)).OrderBy("ChecksLength").Find(&found)

	if len(found) == 0 {
		return nil, fmt.Errorf("No checkable proxy found")
	}

	selected := found[rand.Intn(len(found))]
	token, err := GenerateRequestToken(selected.ID, uid, selected.CheckID)
	if err != nil {
		return nil, err
	}

	return &CheckRequest{
		Token:    token,
		IP:       selected.IP,
		Port:     selected.Port,
		Protocol: selected.Protocol,
	}, nil
}
