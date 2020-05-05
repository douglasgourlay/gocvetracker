package updater

import (
	"cvetracker/mongo"
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Mongo      *mongo.Config `json:"mongo" yaml:"mongo" xml:"mongo"`
	SwitchOnly bool          `json:"switch_only" yaml:"switch_only" xml:"switch_only"`
	Init       bool          `json:"init" yaml:"init" xml:"init"`
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
	c.Mongo.URI = "uri"
	c.Mongo.Database = "database"
	c.Mongo.Collection = "collection"
	c.SwitchOnly = true
	c.Init = false
	return c
}

// AssertValid ...
func (c *Config) AssertValid() {
	c.Mongo.AssertValid()
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
