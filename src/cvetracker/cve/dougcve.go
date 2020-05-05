package cve

import (
	"encoding/json"
	"strconv"

	"gopkg.in/yaml.v2"
)

// DougCVE Simplified Doug CVE Format
type DougCVE struct {
	Cve      string  `json:"cve,omitempty" xml:"cve"`
	Product  string  `json:"product,omitempty"`
	Score    float64 `json:"score,omitempty"`
	Arista   string  `json:"arista,omitempty"`
	Aruba    string  `json:"aruba,omitempty"`
	Hp       string  `json:"hp,omitempty"`
	Ios      string  `json:"ios,omitempty"`
	IosXe    string  `json:"ios_xe,omitempty"`
	IosXr    string  `json:"ios_xr,omitempty"`
	Junos    string  `json:"junos,omitempty"`
	NxOs     string  `json:"nx_os,omitempty"`
	Platform string  `json:"platform,omitempty"`
}

// SetFromMap Set attributes from a map.
func (c *DougCVE) SetFromMap(m map[string]string) {

	if m["cve"] != "" {
		c.Cve = m["cve"]
	}

	if m["product"] != "" {
		c.Product = m["product"]
	}

	if m["arista"] != "" {
		c.Arista = m["arista"]
	}

	if m["aruba"] != "" {
		c.Aruba = m["aruba"]
	}

	if m["hp"] != "" {
		c.Hp = m["hp"]
	}

	if m["ios"] != "" {
		c.Ios = m["ios"]
	}

	if m["iosxe"] != "" {
		c.IosXe = m["iosxe"]
	}

	if m["iosxr"] != "" {
		c.IosXr = m["iosxr"]
	}

	if m["junos"] != "" {
		c.Junos = m["junos"]
	}

	if m["nxos"] != "" {
		c.NxOs = m["nxos"]
	}

	if m["platform"] != "" {
		c.Platform = m["platform"]
	}

	if m["score"] != "" {
		c.Score, _ = strconv.ParseFloat(m["score"], 32)
	}

}

// YamlString ...
func (c *DougCVE) YamlString() string {
	pjson, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(pjson)
}

func (c *DougCVE) String() string {
	pjson, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(pjson)
}
