package mongo

import "encoding/json"

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}

// Config ...
type Config struct {
	URI        string `json:"uri,omitempty" yaml:"uri,omitempty" xml:"uri"`
	Database   string `json:"database,omitempty" yaml:"database,omitempty" xml:"database"`
	Collection string `json:"collection,omitempty" yaml:"collection,omitempty" xml:"collection"`
}

// AssertValid ...
func (c *Config) AssertValid() {

	if c.URI == "" {
		panic("URI must be set")
	}

	if c.Database == "" {
		panic("Database must be set")
	}

	if c.Collection == "" {
		panic("Collection must be set")
	}

}

func (c *Config) String() string {
	pjson, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(pjson)
}
