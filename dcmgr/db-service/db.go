package dbservice

import (
	"github.com/Ankr-network/dccn-common/protos/common"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type DBService interface {
	// Get gets a dc item by pb's id.
	Get(id int64) (*common_proto.DataCenter, error)
	// Get gets a dc item by pb's name.
	GetByName(name string) (*common_proto.DataCenter, error)
	// Create Creates a new dc item if not exits.
	Create(center *common_proto.DataCenter) error
	// GetAll gets all task related to user id.
	GetAll() (*[]*common_proto.DataCenter, error)
	// Update updates dc item
	Update(center *common_proto.DataCenter) error
	// UpdateStatus updates dc item
	UpdateStatus(name string, status common_proto.DCStatus) error
	// Close closes db connection
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
	p.Close()
}

// Get gets user item by id.
func (p *DB) Get(id int64) (*common_proto.DataCenter, error) {
	var center common_proto.DataCenter
	err := p.collection.Find(bson.M{"id": id}).One(&center)
	return &center, err
}

// Get gets user item by name.
func (p *DB) GetByName(name string) (*common_proto.DataCenter, error) {
	var center common_proto.DataCenter
	err := p.collection.Find(bson.M{"name": name}).One(&center)
	return &center, err
}

// Create creates a new data center item if it not exists
func (p *DB) Create(center *common_proto.DataCenter) error {
	return p.collection.Insert(center)
}

// Update updates user item.
func (p *DB) Update(datacenter *common_proto.DataCenter) error {
	return p.collection.Update(
		bson.M{"name": datacenter.Name},
		bson.M{"$set": bson.M{
			"Report":  datacenter.DcHeartbeatReport.Report,
			"Metrics": datacenter.DcHeartbeatReport.Metrics}})
}

func (p *DB) UpdateStatus(name string, status common_proto.DCStatus) error {
	return p.collection.Update(
		bson.M{"name": name},
		bson.M{"$set": bson.M{"status": status}})
}

func (p *DB) GetAll() (*[]*common_proto.DataCenter, error) {
	var dcs []*common_proto.DataCenter
	if err := p.collection.Find(bson.M{}).All(&dcs); err != nil {
		return nil, err
	}
	log.Printf("why there is nothing %+v", dcs)

	return &dcs, nil
}
