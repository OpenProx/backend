package backend

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/msgpack"
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

func (i *Instance) AddChecks() {

}
