package dbservice

import (
	"log"

	dbcommon "github.com/Ankr-network/dccn-common/db"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	// dropCollection for testing usage
	dropCollection()
}

// UserDB implements DBService
type DB struct {
	dbName         string
	collectionName string
	session        *mgo.Session
}

// New returns DBService.
func New(conf dbcommon.Config) (*DB, error) {
	session, err := dbcommon.CreateDBConnection(conf)
	if err != nil {
		return nil, err
	}

	return &DB{
		dbName:         conf.DB,
		collectionName: conf.Collection,
		session:        session,
	}, nil
}

func (p *DB) collection(session *mgo.Session) *mgo.Collection {
	return session.DB(p.dbName).C(p.collectionName)
}

// Get gets user item by id.
func (p *DB) Get(id int64) (*common_proto.DataCenter, error) {
	session := p.session.Clone()
	defer session.Close()

	var center common_proto.DataCenter
	err := p.collection(session).Find(bson.M{"id": id}).One(&center)
	return &center, err
}

// Get gets user item by name.
func (p *DB) GetByName(name string) (*common_proto.DataCenter, error) {
	session := p.session.Clone()
	defer session.Close()

	var center common_proto.DataCenter
	err := p.collection(session).Find(bson.M{"name": name}).One(&center)
	return &center, err
}

// Create creates a new data center item if it not exists
func (p *DB) Create(center *common_proto.DataCenter) error {
	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Insert(center)
}

// Update updates user item.
func (p *DB) Update(datacenter *common_proto.DataCenter) error {
	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Update(
		bson.M{"name": datacenter.Name},
		bson.M{"$set": bson.M{
			"Report":  datacenter.DcHeartbeatReport.Report,
			"Metrics": datacenter.DcHeartbeatReport.Metrics}})
}

func (p *DB) UpdateStatus(name string, status common_proto.DCStatus) error {
	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Update(
		bson.M{"name": name},
		bson.M{"$set": bson.M{"status": status}})
}

// Close closes the db connection.
func (p *DB) Close() {
	p.session.Close()
}

func (p *DB) dropCollection() {
	log.Println(p.session.DB(p.dbName).C(p.collectionName).DropCollection().Error())
}

func (p *DB) GetAll() (*[]*common_proto.DataCenter, error) {
	session := p.session.Clone()
	defer session.Close()

	var dcs []*common_proto.DataCenter

	if err := p.collection(session).Find(bson.M{}).All(&dcs); err != nil {
		return nil, err
	}
	return &dcs, nil
}
