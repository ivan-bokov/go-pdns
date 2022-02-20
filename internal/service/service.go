package service

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ivan-bokov/go-pdns/internal/stacktrace"
	"github.com/ivan-bokov/go-pdns/internal/storage"
	"go.uber.org/zap"
)

type Service struct {
	dnssec bool
	stg    storage.IStorage
	logger *zap.Logger
}

func New(stg storage.IStorage, dnssec bool) *Service {
	return &Service{
		dnssec: dnssec,
		stg:    stg,
		logger: zap.NewExample(),
	}
}

func (s *Service) SetNotified(domainID int, serial int) error {
	_, err := s.stg.Exec(
		"update-serial-query",
		"serial", serial,
		"domain_id", domainID,
	)
	if err != nil {
		return stacktrace.Wrap(err)
	}
	return nil
}

func (s *Service) setLastCheck(domainID int, lastcheck int64) error {
	_, err := s.stg.Exec(
		"update-lastcheck-query",
		"last_check", lastcheck,
		"domain_id", domainID,
	)
	if err != nil {
		return stacktrace.Wrap(err)
	}
	return nil
}

func (s *Service) SetFresh(domainID int) error {
	return s.setLastCheck(domainID, time.Now().UTC().Unix())
}

func (s *Service) Lookup(qtype string, qname string, zoneID int) ([]*DNSResourceRecord, error) {
	var err error
	listRR := make([]*DNSResourceRecord, 0)
	var rows storage.IResult
	if qtype != "ANY" {
		if zoneID < 0 {
			rows, err = s.stg.Query(
				"basic-query",
				"qtype", qtype,
				"qname", qname,
			)
		} else {
			rows, err = s.stg.Query(
				"id-query",
				"qtype", qtype,
				"qname", qname,
				"domain_id", zoneID,
			)
		}
	} else {
		if zoneID < 0 {
			rows, err = s.stg.Query(
				"any-query",
				"qname", qname,
			)
		} else {
			rows, err = s.stg.Query(
				"any-id-query",
				"qname", qname,
				"domain_id", zoneID,
			)
		}
	}
	for rows.Next() {
		rr := new(DNSResourceRecord)
		err = rows.Scan(&rr.Content, &rr.TTL, &rr.Prio, &rr.Qtype, &rr.DomainID, &rr.Disabled, &rr.Qname, &rr.Auth)
		if err != nil {
			//TODO Добавить логирование
			continue
		}
		listRR = append(listRR, rr)
	}
	return listRR, err
}

func (s *Service) List(zonename string, domainID int, includeDisabled bool) ([]*DNSResourceRecord, error) {
	listRR := make([]*DNSResourceRecord, 0)
	if domainID < 0 {
		rows, err := s.stg.Query(
			"get-domain-id",
			"domain", zonename,
		)
		if err != nil {
			return listRR, stacktrace.Wrap(err)
		}
		if rows.Next() {
			err = rows.Scan(domainID)
			if err != nil {
				return listRR, stacktrace.Wrap(err)
			}
		} else {
			return listRR, stacktrace.New(fmt.Sprintf("Domain not found: %s", zonename))
		}
	}
	rows, err := s.stg.Query(
		"list-query",
		"include_disabled", includeDisabled,
		"domain_id", domainID,
	)
	if err != nil {
		return listRR, stacktrace.Wrap(err)
	}
	for rows.Next() {
		rr := new(DNSResourceRecord)
		err = rows.Scan(&rr.Content, &rr.TTL, &rr.Prio, &rr.Qtype, &rr.DomainID, &rr.Disabled, &rr.Qname, &rr.Auth)
		if err != nil {
			//TODO Добавить логирование
			continue
		}
		listRR = append(listRR, rr)
	}
	return listRR, err
}

func (s *Service) GetBeforeAndAfterNamesAbsolute(id int, qname string) error {
	if !s.dnssec {
		return stacktrace.New("Only for DNSSEC")
	}
	return stacktrace.New("No implementation")
}

func (s *Service) SetDomainMetadata(name string, kind string, meta []string) error {
	if !s.dnssec {
		return stacktrace.New("Only for DNSSEC")
	}
	_, err := s.stg.Exec("clear-domain-metadata-query",
		"domain", name,
		"kind", kind,
	)
	if err != nil {
		return stacktrace.Wrap(err)
	}
	errors := make([]error, 0)
	if len(meta) != 0 {
		for _, m := range meta {
			_, err = s.stg.Exec("set-domain-metadata-query",
				"kind", kind,
				"content", m,
				"domain", name,
			)
			if err != nil {
				errors = append(errors, stacktrace.Newf("%v", err))
			}
		}
	}
	if len(errors) != 0 {
		return stacktrace.New(fmt.Sprintf("Unable to set metadata kind %s for domain %s", kind, name))
	}
	return nil
}

func (s *Service) AddDomainKey(name string, key *KeyData) error {
	if !s.dnssec {
		return stacktrace.New("Only for DNSSEC")
	}
	_, err := s.stg.Exec("add-domain-key-query",
		"domain", name,
		"flags", key.Flags,
		"active", key.Active,
		"published", key.Published,
		"content", key.Content,
	)
	if err != nil {
		return stacktrace.Wrap(err)
	}
	return nil
}

func (s *Service) FeedRecord(rr *DNSResourceRecord, ordername string) error {
	var oName interface{}
	prio := 0
	auth := true
	content := rr.Content
	if rr.Qtype == "MX" || rr.Qtype == "SRV" {
		pos := FindFirstNotOf(content, "0123456789")
		if pos != -1 {
			//TODO Сделать очистку до первых цифр
		}
	}
	if s.dnssec {
		auth = rr.Auth
	}
	if ordername == "" {
		oName = nil
	} else {
		oName = []byte(strings.ToLower(ordername))
	}
	_, err := s.stg.Exec("insert-record-query",
		"content", content,
		"ttl", rr.TTL,
		"priority", prio,
		"qtype", rr.Qtype,
		"domain_id", rr.DomainID,
		"disabled", rr.Disabled,
		"qname", rr.Qname,
		"auth", auth,
		"ordername", oName,
	)

	return err
}

func (s *Service) CreateSlaveDomain(ip string, domain string) error {
	_, err := s.stg.Exec("insert-zone-query",
		"domain", domain,
		"account", "",
		"masters", fmt.Sprintf("%s:53", ip),
		"type", "SLAVE",
	)
	return err
}

func (s *Service) GetAllDomainMetadata(name string) (map[string][]string, error) {
	meta := make(map[string][]string)
	rows, err := s.stg.Query(
		"get-all-domain-metadata-query",
		"domain", name,
	)
	if err != nil {
		return meta, stacktrace.Wrap(err)
	}
	for rows.Next() {
		var m1, m2 string
		err = rows.Scan(m1, m2)
		if err != nil {
			return make(map[string][]string), stacktrace.Wrap(err)
		}
		if _, ok := meta[m1]; !ok {
			meta[m1] = make([]string, 0, 10)
		}
		meta[m1] = append(meta[m1], m2)
	}
	return meta, nil
}

func (s *Service) GetDomainInfo(name string) (*DomainInfo, error) {
	rows, err := s.stg.Query(
		"info-zone-query",
		"domain", name,
	)
	if err != nil {
		return new(DomainInfo), stacktrace.Wrap(err)
	}
	di := new(DomainInfo)
	if rows.Next() {
		master := ""
		err = rows.Scan(&di.ID, &di.Zone, &master, &di.LastCheck, &di.Serial, &di.Kind, &di.Account)
		if err != nil {
			log.Println("[ERROR] " + err.Error())
			return new(DomainInfo), stacktrace.Wrap(err)
		}
		if master != "" {
			di.Master = StringTok(master, " ,\t")
		}
	}
	return di, nil
}

func (s *Service) GetAllDomains(includeDisabled bool) ([]*DomainInfo, error) {
	rows, err := s.stg.Query(
		"get-all-domains-query",
		"include_disabled", includeDisabled,
	)
	if err != nil {
		return nil, stacktrace.Wrap(err)
	}
	dis := make([]*DomainInfo, 0, 10)
	for rows.Next() {
		di := new(DomainInfo)
		master := ""
		err = rows.Scan(&di.ID, &di.Zone, &master, &di.LastCheck, &di.Serial, &di.Kind, &di.Account)
		if err != nil {
			log.Println("[ERROR] " + err.Error())
			return nil, stacktrace.Wrap(err)
		}
		if master != "" {
			di.Master = StringTok(master, " ,\t")
		}
		dis = append(dis, di)
	}
	return dis, nil
}
