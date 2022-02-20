package service

type DNSName struct {
	name string
}

type DNSPacket struct {
}

type DNSResourceRecord struct {
	Qname        string `json:"qname,omitempty"`
	OrderName    string `json:"order_name,omitempty"`
	WildcardName string `json:"wildcard_name,omitempty"`
	Content      string `json:"content,omitempty"`
	TTL          int    `json:"ttl,omitempty"`
	DomainID     int    `json:"domain_id,omitempty"`
	Qtype        string `json:"qtype,omitempty"`
	Auth         bool   `json:"auth,omitempty"`
	Disabled     bool   `json:"disabled,omitempty"`
	Qclass       string `json:"qclass,omitempty"`
	Prio         int    `json:"prio,omitempty"`
}

type ComboAddress struct {
	IP   string
	Port int
}

type KeyData struct {
	Flags     int
	Active    bool
	Published bool
	Content   string
}
