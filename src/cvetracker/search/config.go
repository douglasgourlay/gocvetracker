package search

import (
	"cvetracker/mongo"
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Enabled bool          `json:"enabled" yaml:"enabled" xml:"enabled"`
	Mongo   *mongo.Config `json:"mongo" yaml:"mongo" xml:"mongo"`
	Port    int           `json:"port" yaml:"port" xml:"port"`
}

// NewConfig ...
func NewConfig() *Config {
	c := &Config{}
	c.Mongo = mongo.NewConfig()
	return c
}

// ExampleConfig ...
func ExampleConfig() *Config {
	c := NewConfig()
	c.Port = 8080
	c.Mongo.URI = "uri"
	c.Mongo.Database = "database"
	c.Mongo.Collection = "collection"
	return c
}

// AssertValid ...
func (c *Config) AssertValid() {

	c.Mongo.AssertValid()

	if c.Port <= 0 {
		panic("Port must be greater then zero")
	}

}

// YamlString ...
func (c *Config) YamlString() string {
	pjson, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(pjson)
}

func (c *Config) String() string {
	pjson, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(pjson)
}
