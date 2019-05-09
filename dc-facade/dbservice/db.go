package dbservice

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBService interface {
	// Get gets a dc item by pb's id.
	Get(id int64) (*common_proto.DataCenterStatus, error)
	// Get gets a dc item by pb's name.
	Close()
}

// UserDB implements DBService
type DB struct {
	collection *mgo.Collection
}

// New returns DBService.
func New() (*DB, error) {
	collection := micro2.GetCollection("datacenter")
	return &DB{
		collection: collection,
	}, nil
}

func (p *DB) Close() {
	//p.Close()
}

// Get gets user item by id.
func (p *DB) Get(id int64) (*common_proto.DataCenterStatus, error) {
	var center common_proto.DataCenterStatus
	err := p.collection.Find(bson.M{"dcid": id}).One(&center)
	return &center, err
}

