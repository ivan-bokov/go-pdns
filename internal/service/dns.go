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
	ID        int    `json:"id,omitempty"`
	Flags     int    `json:"flags,omitempty"`
	Active    bool   `json:"active,omitempty"`
	Published bool   `json:"published,omitempty"`
	Content   string `json:"content,omitempty"`
}
type DomainInfo struct {
	ID        int      `json:"id,omitempty"`
	Zone      string   `json:"zone,omitempty"`
	Kind      string   `json:"kind,omitempty"`
	Serial    int64    `json:"serial,omitempty"`
	Master    []string `json:"masters,omitempty"`
	LastCheck int64    `json:"last_check,omitempty"`
	Account   string   `json:"account,omitempty"`
}

type TSIGKey struct {
	Algorithm string `json:"algorithm,omitempty"`
	Content   string `json:"content,omitempty"`
}

type SOAData struct {
	Qname      string
	Nameserver string
	Hostmaster string
	TTL        int
	DomainID   int
	Serial     int64
	Refresh    int64
	Retry      int64
	Expire     int64
	Minimum    int64
}
