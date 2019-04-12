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
	GetByName(name string) (*common_proto.DataCenterStatus, error)
	// Create Creates a new dc item if not exits.
	Create(center *common_proto.DataCenterStatus) error
	// GetAll gets all task related to user id.
	GetAll() (*[]*common_proto.DataCenterStatus, error)
	// Update updates dc item
	Update(center *common_proto.DataCenterStatus) error
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
func (p *DB) Get(id int64) (*common_proto.DataCenterStatus, error) {
	var center common_proto.DataCenterStatus
	err := p.collection.Find(bson.M{"id": id}).One(&center)
	return &center, err
}

// Get gets user item by name.
func (p *DB) GetByName(name string) (*common_proto.DataCenterStatus, error) {
	var center common_proto.DataCenterStatus
	err := p.collection.Find(bson.M{"name": name}).One(&center)
	return &center, err
}

// Create creates a new data center item if it not exists
func (p *DB) Create(center *common_proto.DataCenterStatus) error {
	return p.collection.Insert(center)
}

// Update updates user item.
func (p *DB) Update(datacenter *common_proto.DataCenterStatus) error {
	return p.collection.Update(
		bson.M{"name": datacenter.Name},
		bson.M{"$set": bson.M{
			"Report":  datacenter.DcHeartbeatReport.Report,
			"Metrics": datacenter.DcHeartbeatReport.Metrics}})
		return nil
}

func (p *DB) UpdateStatus(name string, status common_proto.DCStatus) error {
	return p.collection.Update(
		bson.M{"name": name},
		bson.M{"$set": bson.M{"status": status}})
}

func (p *DB) GetAll() (*[]*common_proto.DataCenterStatus, error) {
	var dcs []*common_proto.DataCenterStatus
	if err := p.collection.Find(bson.M{}).All(&dcs); err != nil {
		return nil, err
	}

	return &dcs, nil
}
