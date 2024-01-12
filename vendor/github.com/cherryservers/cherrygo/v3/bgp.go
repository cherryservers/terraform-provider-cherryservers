package cherrygo

type ProjectBGP struct {
	Enabled  bool `json:"enabled,omitempty"`
	LocalASN int  `json:"local_asn,omitempty"`
}

type ServerBGP struct {
	Enabled   bool       `json:"enabled"`
	Available bool       `json:"available,omitempty"`
	Status    string     `json:"status,omitempty"`
	Routers   int        `json:"routers,omitempty"`
	Connected int        `json:"connected,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Active    int        `json:"active,omitempty"`
	Routes    []BGPRoute `json:"routes,omitempty"`
	Updated   string     `json:"updated,omitempty"`
}

type BGPRoute struct {
	Subnet  string `json:"subnet,omitempty"`
	Active  bool   `json:"active,omitempty"`
	Router  string `json:"router,omitempty"`
	Age     string `json:"age,omitempty"`
	Updated string `json:"updated,omitempty"`
}

type RegionBGP struct {
	Hosts []string `json:"hosts,omitempty"`
	Asn   int      `json:"asn,omitempty"`
}
