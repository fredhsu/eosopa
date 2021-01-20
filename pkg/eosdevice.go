package eosparse

// EOSDevice is an EOS endpoint
type EOSDevice struct {
	Hostname      string                 `json:"id"` // using id to be consistent with OPA
	Management    map[string]interface{} `json:"management"`
	IPNameServers NameServers            `json:"ipNameServers"`
	Logging       Logging                `json:"logging"`
	SWVersion     SWVersion              `json:"swVersion"`
	HWVersion     HWVersion              `json:"hwVersion"`
}

func NewEOSDevice() EOSDevice {
	m := NewManagement()
	ns := NameServers{}
	logging := Logging{}
	swver := SWVersion{}
	hwver := ""
	ed := EOSDevice{"", m, ns, logging, swver, hwver}
	return ed
}
