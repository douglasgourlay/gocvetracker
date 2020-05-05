// Download CVE(s) from NVD, convert to Simple Doug Format and install
// into MongoDB if it does not already exist

package updater

import (
	"bytes"
	"compress/gzip"
	"cvetracker/cve"
	"cvetracker/mongo"
	"cvetracker/util"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const nvdURIPre = "https://nvd.nist.gov/feeds/json/cve/1.1/nvdcve-1.1-"
const nvdURIPost = ".json.gz"
const nvdURIRecent = "https://nvd.nist.gov/feeds/json/cve/1.1/nvdcve-1.1-recent.json.gz"
const nvdStartYear = 2002

func (updater *updater) run() {

	var err error

	zap.L().Debug("Starting")

	updater.mongo, err = mongo.NewClient(updater.config.Mongo)
	if err != nil {
		panic(err)
	}

	if updater.config.Init {
		zap.L().Debug("Processing Init")
		stopYear, _, _ := time.Now().Date()
		stopYear = stopYear + 1
		for i := nvdStartYear; i < stopYear; i++ {
			uri := nvdURIPre + strconv.Itoa(i) + nvdURIPost
			updater.process(uri)
		}
	}

	zap.L().Debug("Processing Update")
	updater.process(nvdURIRecent)

	updater.mongo.Shutdown()

	zap.L().Debug("Goodbye")

}

func (updater *updater) process(uri string) error {

	zap.L().Debug(fmt.Sprintf("Processing %s", uri))

	resp, err := http.Get(uri)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	gr, err := gzip.NewReader(bytes.NewBuffer(buf))
	defer gr.Close()

	data, err := ioutil.ReadAll(gr)
	if err != nil {
		return err
	}

	cve := &cve.NVDCVED{}

	err = json.Unmarshal(data, &cve)
	if err != nil {
		return err
	}

	if updater.config.SwitchOnly {
		zap.L().Debug("Processing switch platforms only")
	}

	for _, d := range cve.CVEItems {
		dcve := dougCVE(&d)

		if updater.config.SwitchOnly {
			if dcve.Platform == "switch" {
				zap.L().Debug(fmt.Sprintf("Processing %s", dcve.Cve))
				updater.insert(dcve)
			}
		} else {
			zap.L().Debug(fmt.Sprintf("Processing %s", dcve.Cve))
			updater.insert(dcve)
		}

	}

	return nil
}

func (updater *updater) insert(dcve *cve.DougCVE) error {

	existing, err := updater.mongo.Get(dcve.Cve)
	if err != nil {
		return err
	}

	if existing == nil {
		zap.L().Debug(fmt.Sprintf("Inserting Cve %s", dcve.Cve))
		err = updater.mongo.Insert(dcve)
		if err != nil {
			return err
		}
	} else {
		zap.L().Debug(fmt.Sprintf("Cve %s already exist", dcve.Cve))
	}

	// TODO Get and check if they are different

	return nil
}

type updater struct {
	mongo  *mongo.Client
	config *Config
}

func dougCVE(c *cve.NVDCVE) *cve.DougCVE {

	result := &cve.DougCVE{}

	result.Cve = c.Cve.CVEDataMeta.ID
	result.Score = c.Impact.BaseMetricV2.CvssV2.BaseScore
	result.Product = c.Configurations.Nodes.String()

	// This will parse out the vendor, os, etc
	for _, n := range c.Configurations.Nodes {

		for _, c := range n.CpeMatch {
			if c.Vulnerable {
				processCPE(result, c.Cpe23URI)
			}
		}

		for _, d := range n.Children {
			for _, c := range d.CpeMatch {
				if c.Vulnerable {
					processCPE(result, c.Cpe23URI)
				}
			}

		}
	}

	return result
}

func processCPE(cve *cve.DougCVE, s string) error {

	// cpe format
	// cpe:2.3:o:cisco:ios_xr:*:*:*:*:*:*:*:*

	sa, err := util.TokenizeString(s, ':')

	if err != nil {
		return err
	}

	if len(sa) < 5 {
		return errors.New("CPE length is to short")
	}

	cpeType := sa[2]
	vendor := sa[3]
	value := sa[4]

	switch vendor {

	case "juniper":
		cve.Junos = "true"
		cve.Platform = "switch"
		break

	case "arista":
		cve.Arista = "true"
		cve.Platform = "switch"
		break

	case "hp":
		cve.Hp = "true"
		cve.Platform = "switch"
		break

	case "aruba":
		cve.Aruba = "true"
		cve.Platform = "switch"
		break

	case "cisco":

		cve.Platform = "switch"

		if cpeType == "o" {
			switch value {

			case "ios":
				cve.Ios = "true"
				break

			case "ios_xr":
				cve.IosXr = "true"
				break

			case "ios_xve":
				cve.IosXe = "true"
				break

			case "nx-os":
				cve.NxOs = "true"
				break

			}

		}

	}

	return nil
}
