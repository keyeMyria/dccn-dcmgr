package dbservice

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/golang/protobuf/ptypes/timestamp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type DBService interface {
	// Get gets a dc item by pb's id.
	Get(id string) (*DataCenterRecord, error)
	// Get gets a dc item by pb's name.
	GetByName(name string) (*DataCenterRecord, error)
	GetByUserID(uid string) (*DataCenterRecord, error)
	GetByID(id string) (*DataCenterRecord, error)
	// Create Creates a new dc item if not exits.
	Create(center *DataCenterRecord) error
	Reset(clusterID string, center *DataCenterRecord) error
	// GetAll gets all task related to user id.
	GetAll() (*[]*DataCenterRecord, error)
	GetAvaliableList() (*[]*DataCenterRecord, error)
	// Update updates dc item
	Update(center *DataCenterRecord) error
	// UpdateStatus updates dc item
	UpdateStatus(clusterID string, status common_proto.DCStatus) error

	UpdateClientCert(clusterID string, clientcert string) error
	// Close closes db connection
	Close()
}

type DataCenterRecord struct {
	DcId              string
	ClusterName       string
	GeoLocation       *common_proto.GeoLocation
	DcStatus          common_proto.DCStatus
	DcAttributes      *common_proto.DataCenterAttributes
	DcHeartbeatReport *common_proto.DCHeartbeatReport
	UserId            string
	Clientcert        string
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
func (p *DB) Get(id string) (*DataCenterRecord, error) {
	var center DataCenterRecord
	log.Printf("get datacetner %s \n", id)
	err := p.collection.Find(bson.M{"dcid": id}).One(&center)
	return &center, err
}

// Get gets user item by id.
func (p *DB) GetByID(id string) (*DataCenterRecord, error) {
	var record DataCenterRecord
	err := p.collection.Find(bson.M{"dcid": id}).One(&record)
	return &record, err
}

// Get gets user item by name.
func (p *DB) GetByName(name string) (*DataCenterRecord, error) {
	var center DataCenterRecord
	err := p.collection.Find(bson.M{"clustername": name}).One(&center)
	return &center, err
}

func (p *DB) GetByUserID(uid string) (*DataCenterRecord, error) {
	var center DataCenterRecord
	err := p.collection.Find(bson.M{"userid": uid}).One(&center)
	return &center, err
}

// Create creates a new data center item if it not exists
func (p *DB) Create(center *DataCenterRecord) error {
	return p.collection.Insert(center)
}

func (p *DB) Reset(clusterID string, datacenter *DataCenterRecord) error {
	p.setUpdateTime(clusterID)
	return p.collection.Update(
		bson.M{"dcid": clusterID},
		bson.M{"$set": bson.M{
			"clustername":       datacenter.ClusterName,
			"clientcert":        datacenter.Clientcert,
			"userid":            datacenter.UserId,
			"dcstatus":          datacenter.DcStatus,
			"dcheartbeatreport": datacenter.DcHeartbeatReport}})
}

func (p *DB) setUpdateTime(id string) {
	now := time.Now().Unix()
	dataeTime := &timestamp.Timestamp{Seconds: now}
	p.collection.Update(
		bson.M{"dcid": id},
		bson.M{"$set": bson.M{
			"dcattributes.lastmodifieddate": dataeTime,
		}})

}

func (p *DB) Update(record *DataCenterRecord) error {
	if record.GeoLocation != nil && len(record.GeoLocation.Lat) > 0 {
		return p.UpdateWithGEO(record)
	} else {
		return p.UpdateWithoutGEO(record)
	}
}

func (p *DB) UpdateWithGEO(record *DataCenterRecord) error {
	//log.Printf("UpdateWithGEO---------> %+v \n", record)
	p.setUpdateTime(record.DcId)
	return p.collection.Update(
		bson.M{"dcid": record.DcId},
		bson.M{"$set": bson.M{
			"dcstatus":          record.DcStatus,
			"geolocation":       record.GeoLocation,
			"dcheartbeatreport": record.DcHeartbeatReport,
		}})
}

func (p *DB) UpdateWithoutGEO(record *DataCenterRecord) error {
	//log.Printf("UpdateWithoutGEO---------> %+v \n", record)
	p.setUpdateTime(record.DcId)
	return p.collection.Update(
		bson.M{"dcid": record.DcId},
		bson.M{"$set": bson.M{
			"dcstatus":          record.DcStatus,
			"dcheartbeatreport": record.DcHeartbeatReport,
		}})
}

// Update updates user item.
func (p *DB) UpdateName(id string, name string) error {
	return p.collection.Update(
		bson.M{"dcid": id},
		bson.M{"$set": bson.M{
			"clustername": name}})
}

func (p *DB) UpdateStatus(clusterID string, status common_proto.DCStatus) error {
	log.Printf("update UpdateStatus %s %s \n", clusterID, status)
	return p.collection.Update(
		bson.M{"dcid": clusterID},
		bson.M{"$set": bson.M{"dcstatus": status}})
}

func (p *DB) UpdateClientCert(clusterID string, clientcert string) error {
	return p.collection.Update(
		bson.M{"dcid": clusterID},
		bson.M{"$set": bson.M{"clientcert": clientcert}})
}

func (p *DB) GetAll() (*[]*DataCenterRecord, error) {
	var dcs []*DataCenterRecord
	if err := p.collection.Find(bson.M{}).All(&dcs); err != nil {
		return nil, err
	}

	return &dcs, nil
}

func (p *DB) GetAvaliableList() (*[]*DataCenterRecord, error) {
	var dcs []*DataCenterRecord
	if err := p.collection.Find(bson.M{"dcstatus": common_proto.DCStatus_AVAILABLE}).All(&dcs); err != nil {
		return nil, err
	}

	return &dcs, nil
}
