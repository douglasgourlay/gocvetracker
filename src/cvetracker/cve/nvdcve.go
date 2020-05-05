package cve

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// NVDCVED National Vulnerability Database Format see
// https://csrc.nist.gov/schema/nvd/feed/1.1/nvd_cve_feed_json_1.1.schema
type NVDCVED struct {
	CVEDataType         string   `json:"CVE_data_type" yaml:"CVE_data_type" xml:"CVE_data_type"`
	CVEDataFormat       string   `json:"CVE_data_format"`
	CVEDataVersion      string   `json:"CVE_data_version"`
	CVEDataNumberOfCVEs string   `json:"CVE_data_numberOfCVEs"`
	CVEDataTimestamp    string   `json:"CVE_data_timestamp"`
	CVEItems            []NVDCVE `json:"CVE_Items"`
}

// NVDCVE National Vulnerability Database Format see
// https://csrc.nist.gov/schema/nvd/feed/1.1/nvd_cve_feed_json_1.1.schema
type NVDCVE struct {
	Cve struct {
		DataType    string `json:"data_type"`
		DataFormat  string `json:"data_format"`
		DataVersion string `json:"data_version"`
		CVEDataMeta struct {
			ID       string `json:"ID"`
			ASSIGNER string `json:"ASSIGNER"`
		} `json:"CVE_data_meta"`
		Problemtype struct {
			ProblemtypeData []struct {
				Description []struct {
					Lang  string `json:"lang"`
					Value string `json:"value"`
				} `json:"description"`
			} `json:"problemtype_data"`
		} `json:"problemtype"`
		References struct {
			ReferenceData []struct {
				URL       string   `json:"url"`
				Name      string   `json:"name"`
				Refsource string   `json:"refsource"`
				Tags      []string `json:"tags"`
			} `json:"reference_data"`
		} `json:"references"`
		Description struct {
			DescriptionData []struct {
				Lang  string `json:"lang"`
				Value string `json:"value"`
			} `json:"description_data"`
		} `json:"description"`
	} `json:"cve"`
	Configurations struct {
		CVEDataVersion string `json:"CVE_data_version"`
		Nodes          Nodes  `json:"nodes"`
	} `json:"configurations"`

	Impact struct {
		BaseMetricV3 struct {
			CvssV3 struct {
				Version               string  `json:"version"`
				VectorString          string  `json:"vectorString"`
				AttackVector          string  `json:"attackVector"`
				AttackComplexity      string  `json:"attackComplexity"`
				PrivilegesRequired    string  `json:"privilegesRequired"`
				UserInteraction       string  `json:"userInteraction"`
				Scope                 string  `json:"scope"`
				ConfidentialityImpact string  `json:"confidentialityImpact"`
				IntegrityImpact       string  `json:"integrityImpact"`
				AvailabilityImpact    string  `json:"availabilityImpact"`
				BaseScore             float64 `json:"baseScore"`
				BaseSeverity          string  `json:"baseSeverity"`
			} `json:"cvssV3"`
			ExploitabilityScore float64 `json:"exploitabilityScore"`
			ImpactScore         float64 `json:"impactScore"`
		} `json:"baseMetricV3"`
		BaseMetricV2 struct {
			CvssV2 struct {
				Version               string  `json:"version"`
				VectorString          string  `json:"vectorString"`
				AccessVector          string  `json:"accessVector"`
				AccessComplexity      string  `json:"accessComplexity"`
				Authentication        string  `json:"authentication"`
				ConfidentialityImpact string  `json:"confidentialityImpact"`
				IntegrityImpact       string  `json:"integrityImpact"`
				AvailabilityImpact    string  `json:"availabilityImpact"`
				BaseScore             float64 `json:"baseScore"`
			} `json:"cvssV2"`
			Severity                string  `json:"severity"`
			ExploitabilityScore     float64 `json:"exploitabilityScore"`
			ImpactScore             float64 `json:"impactScore"`
			AcInsufInfo             bool    `json:"acInsufInfo"`
			ObtainAllPrivilege      bool    `json:"obtainAllPrivilege"`
			ObtainUserPrivilege     bool    `json:"obtainUserPrivilege"`
			ObtainOtherPrivilege    bool    `json:"obtainOtherPrivilege"`
			UserInteractionRequired bool    `json:"userInteractionRequired"`
		} `json:"baseMetricV2"`
	} `json:"impact"`
	PublishedDate    string `json:"publishedDate"`
	LastModifiedDate string `json:"lastModifiedDate"`
}

// Nodes ...
type Nodes []struct {
	Operator string `json:"operator"`
	CpeMatch []struct {
		Vulnerable            bool   `json:"vulnerable"`
		Cpe23URI              string `json:"cpe23Uri"`
		VersionStartIncluding string `json:"versionStartIncluding"`
		VersionEndExcluding   string `json:"versionEndExcluding"`
	} `json:"cpe_match,omitempty"`
	Children []struct {
		Operator string `json:"operator"`
		CpeMatch []struct {
			Vulnerable bool   `json:"vulnerable"`
			Cpe23URI   string `json:"cpe23Uri"`
		} `json:"cpe_match"`
	} `json:"children,omitempty"`
}

func (c *Nodes) String() string {
	pjson, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(pjson)
}

// YamlString ...
func (c *NVDCVE) YamlString() string {
	pjson, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(pjson)
}

func (c *NVDCVE) String() string {
	pjson, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(pjson)
}
