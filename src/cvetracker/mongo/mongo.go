package mongo

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"cvetracker/cve"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// NewClient ...
func NewClient(m *Config) (*Client, error) {

	zap.L().Debug("Starting MongoDB Client")

	m.AssertValid()

	c := &Client{
		uri:        m.URI,
		database:   m.Database,
		collection: m.Collection,
	}

	return c, c.connect()
}

// Client ...
type Client struct {
	uri        string
	database   string
	collection string
	client     *mongo.Client
	ctx        context.Context
}

func (m *Client) connect() error {

	var err error

	if m.client != nil {
		err = m.client.Ping(context.TODO(), nil)
		if err == nil {
			// We are already connected
			return nil
		}
	}

	var cancel context.CancelFunc
	m.ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m.client, err = mongo.Connect(m.ctx, options.Client().ApplyURI(
		m.uri,
	))

	if err != nil {
		return err
	}

	return nil
}

// bson.D{}

// Search ...
func (m *Client) Search(filter *cve.DougCVE) ([]cve.DougCVE, error) {

	var result []cve.DougCVE

	filterJSON, err := json.Marshal(filter)
	if err != nil {
		return result, err
	}

	var bsonFilter bson.M
	err = json.Unmarshal([]byte(filterJSON), &bsonFilter)
	if err != nil {
		return nil, err
	}

	err = m.connect()
	if err != nil {
		return result, err
	}

	collection := m.client.Database(m.database).Collection(m.collection)
	findOptions := options.Find()
	cur, err := collection.Find(context.TODO(), bsonFilter, findOptions)

	if err != nil {
		return result, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem cve.DougCVE
		err := cur.Decode(&elem)
		if err != nil {
			return result, err
		}
		// NaN is not valid for GO JSON so we will change it to zero
		if math.IsNaN(elem.Score) {
			elem.Score = 0
		}
		result = append(result, elem)
	}

	return result, nil

}

// GetAll ...
func (m *Client) GetAll() ([]cve.DougCVE, error) {

	var result []cve.DougCVE

	err := m.connect()
	if err != nil {
		return result, err
	}

	collection := m.client.Database(m.database).Collection(m.collection)
	findOptions := options.Find()
	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)

	if err != nil {
		return result, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem cve.DougCVE
		err := cur.Decode(&elem)
		if err != nil {
			return result, err
		}
		// NaN is not valid for GO JSON so we will change it to zero
		if math.IsNaN(elem.Score) {
			elem.Score = 0
		}
		result = append(result, elem)
	}

	return result, nil

}

// Get ...
func (m *Client) Get(cveID string) (*cve.DougCVE, error) {

	var result cve.DougCVE

	filter := bson.D{primitive.E{Key: "cve", Value: cveID}}

	err := m.connect()
	if err != nil {
		return nil, err
	}

	collection := m.client.Database(m.database).Collection(m.collection)

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

// Delete ...
func (m *Client) Delete(cveID string) error {

	filter := bson.D{primitive.E{Key: "cve", Value: cveID}}

	err := m.connect()
	if err != nil {
		return err
	}

	collection := m.client.Database(m.database).Collection(m.collection)

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}

	return nil
}

// Insert ...
func (m *Client) Insert(c *cve.DougCVE) error {
	collection := m.client.Database(m.database).Collection(m.collection)
	_, err := collection.InsertOne(context.TODO(), c)
	return err
}

// Shutdown ...
func (m *Client) Shutdown() {
	zap.L().Debug("Shutting Down MongoDB Client")
	if m.client != nil {
		m.client.Disconnect(m.ctx)
	}
}
