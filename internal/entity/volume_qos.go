package entity

type QosSpecs struct {
	Specs struct {
		ReadIopsSec            string        `json:"read_iops_sec"`
		ReadIopsSecMax         string        `json:"read_iops_sec_max"`
		WriteIopsSec           string        `json:"write_iops_sec"`
		WriteIopsSecMax        string        `json:"write_iops_sec_max"`
		ReadBytesSec           string        `json:"read_bytes_sec"`
        ReadBytesSecMax	       string        `json:"read_bytes_sec_max"`
		WriteBytesSec          string        `json:"write_bytes_sec"`
		WriteBytesSecMax       string        `json:"write_bytes_sec_max"`
	} `json:"specs"`
	Consumer string `json:"consumer"`
	Name     string `json:"name"`
	Id       string `json:"id"`
}

type QosSpecsMap struct {
	QosSpecs `json:"qos_specs"`
	Links    []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
}

type QosSpecss struct {
	Qss               []QosSpecs `json:"qos_specs"`
}
