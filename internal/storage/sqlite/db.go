package sqlite

import (
	"database/sql"

	"github.com/ivan-bokov/go-pdns/internal/stacktrace"
	"github.com/ivan-bokov/go-pdns/internal/storage"
	sqlex "github.com/ivan-bokov/go-pdns/internal/storage/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	db       *sql.DB
	bindType int
	declare  map[string]string
}

func New(dataSource string) *Sqlite {
	db, err := sql.Open("sqlite3", dataSource) //":memory:"
	if err != nil {
		panic(stacktrace.Wrap(err))
	}
	return &Sqlite{
		db:       db,
		bindType: sqlex.BindType("sqlite3"),
		declare:  declareSQL(),
	}
}

func declareSQL() map[string]string {
	dec := make(map[string]string)
	record_query := "SELECT content,ttl,prio,type,domain_id,disabled,name,auth FROM records WHERE"

	dec["basic-query"] = record_query + " disabled=0 and type=:qtype and name=:qname"
	dec["id-query"] = record_query + " disabled=0 and type=:qtype and name=:qname and domain_id=:domain_id"
	dec["any-query"] = record_query + " disabled=0 and name=:qname"
	dec["any-id-query"] = record_query + " disabled=0 and name=:qname and domain_id=:domain_id"
	dec["list-query"] = "SELECT content,ttl,prio,type,domain_id,disabled,name,auth,ordername FROM records WHERE (disabled=0 OR :include_disabled) and domain_id=:domain_id order by name, type"
	dec["list-subzone-query"] = record_query + " disabled=0 and (name=:zone OR name like :wildzone) and domain_id=:domain_id"

	dec["remove-empty-non-terminals-from-zone-query"] = "delete from records where domain_id=:domain_id and type is null"
	dec["delete-empty-non-terminal-query"] = "delete from records where domain_id=:domain_id and name=:qname and type is null"

	dec["info-zone-query"] = "select id,name,master,last_check,notified_serial,type,account from domains where name=:domain"

	dec["get-domain-id"] = "select id from domains where name=:domain"

	dec["info-all-slaves-query"] = "select id,name,master,last_check from domains where type='SLAVE'"
	dec["supermaster-query"] = "select account from supermasters where ip=:ip and nameserver=:nameserver"
	dec["supermaster-name-to-ips"] = "select ip,account from supermasters where nameserver=:nameserver and account=:account"
	dec["supermaster-add"] = "insert into supermasters (ip, nameserver, account) values (:ip,:nameserver,:account)"
	dec["autoprimary-remove"] = "delete from supermasters where ip = :ip and nameserver = :nameserver"
	dec["list-autoprimaries"] = "select ip,nameserver,account from supermasters"

	dec["insert-zone-query"] = "insert into domains (type,name,master,account,last_check,notified_serial) values(:type, :domain, :masters, :account, null, null)"

	dec["insert-record-query"] = "insert into records (content,ttl,prio,type,domain_id,disabled,name,ordername,auth) values (:content,:ttl,:priority,:qtype,:domain_id,:disabled,:qname,:ordername,:auth)"
	dec["insert-empty-non-terminal-order-query"] = "insert into records (type,domain_id,disabled,name,ordername,auth,ttl,prio,content) values (null,:domain_id,0,:qname,:ordername,:auth,null,null,null)"

	dec["get-order-first-query"] = "select ordername from records where disabled=0 and domain_id=:domain_id and ordername is not null order by 1 asc limit 1"
	dec["get-order-before-query"] = "select ordername, name from records where disabled=0 and ordername <= :ordername and domain_id=:domain_id and ordername is not null order by 1 desc limit 1"
	dec["get-order-after-query"] = "select min(ordername) from records where disabled=0 and ordername > :ordername and domain_id=:domain_id and ordername is not null"
	dec["get-order-last-query"] = "select ordername, name from records where disabled=0 and ordername != '' and domain_id=:domain_id and ordername is not null order by 1 desc limit 1"

	dec["update-ordername-and-auth-query"] = "update records set ordername=:ordername,auth=:auth where domain_id=:domain_id and name=:qname and disabled=0"
	dec["update-ordername-and-auth-type-query"] = "update records set ordername=:ordername,auth=:auth where domain_id=:domain_id and name=:qname and type=:qtype and disabled=0"
	dec["nullify-ordername-and-update-auth-query"] = "update records set ordername=NULL,auth=:auth where domain_id=:domain_id and name=:qname and disabled=0"
	dec["nullify-ordername-and-update-auth-type-query"] = "update records set ordername=NULL,auth=:auth where domain_id=:domain_id and name=:qname and type=:qtype and disabled=0"

	dec["update-master-query"] = "update domains set master=:master where name=:domain"
	dec["update-kind-query"] = "update domains set type=:kind where name=:domain"
	dec["update-account-query"] = "update domains set account=:account where name=:domain"
	dec["update-serial-query"] = "update domains set notified_serial=:serial where id=:domain_id"
	dec["update-lastcheck-query"] = "update domains set last_check=:last_check where id=:domain_id"
	dec["info-all-master-query"] = "select domains.id, domains.name, domains.notified_serial, records.content from records join domains on records.domain_id=domains.id and records.name=domains.name where records.type='SOA' and records.disabled=0 and domains.type='MASTER'"
	dec["delete-domain-query"] = "delete from domains where name=:domain"
	dec["delete-zone-query"] = "delete from records where domain_id=:domain_id"
	dec["delete-rrset-query"] = "delete from records where domain_id=:domain_id and name=:qname and type=:qtype"
	dec["delete-names-query"] = "delete from records where domain_id=:domain_id and name=:qname"

	dec["add-domain-key-query"] = "insert into cryptokeys (domain_id, flags, active, published, content) select id, :flags, :active, :published, :content from domains where name=:domain"
	dec["get-last-inserted-key-id-query"] = "select last_insert_rowid()"
	dec["list-domain-keys-query"] = "select cryptokeys.id, flags, active, published, content from domains, cryptokeys where cryptokeys.domain_id=domains.id and name=:domain"
	dec["get-all-domain-metadata-query"] = "select kind,content from domains, domainmetadata where domainmetadata.domain_id=domains.id and name=:domain"
	dec["get-domain-metadata-query"] = "select content from domains, domainmetadata where domainmetadata.domain_id=domains.id and name=:domain and domainmetadata.kind=:kind"
	dec["clear-domain-metadata-query"] = "delete from domainmetadata where domain_id=(select id from domains where name=:domain) and domainmetadata.kind=:kind"
	dec["clear-domain-all-metadata-query"] = "delete from domainmetadata where domain_id=(select id from domains where name=:domain)"
	dec["set-domain-metadata-query"] = "insert into domainmetadata (domain_id, kind, content) select id, :kind, :content from domains where name=:domain"
	dec["activate-domain-key-query"] = "update cryptokeys set active=1 where domain_id=(select id from domains where name=:domain) and  cryptokeys.id=:key_id"
	dec["deactivate-domain-key-query"] = "update cryptokeys set active=0 where domain_id=(select id from domains where name=:domain) and  cryptokeys.id=:key_id"
	dec["publish-domain-key-query"] = "update cryptokeys set published=1 where domain_id=(select id from domains where name=:domain) and  cryptokeys.id=:key_id"
	dec["unpublish-domain-key-query"] = "update cryptokeys set published=0 where domain_id=(select id from domains where name=:domain) and  cryptokeys.id=:key_id"
	dec["remove-domain-key-query"] = "delete from cryptokeys where domain_id=(select id from domains where name=:domain) and cryptokeys.id=:key_id"
	dec["clear-domain-all-keys-query"] = "delete from cryptokeys where domain_id=(select id from domains where name=:domain)"
	dec["get-tsig-key-query"] = "select algorithm, secret from tsigkeys where name=:key_name"
	dec["set-tsig-key-query"] = "replace into tsigkeys (name,algorithm,secret) values(:key_name,:algorithm,:content)"
	dec["delete-tsig-key-query"] = "delete from tsigkeys where name=:key_name"
	dec["get-tsig-keys-query"] = "select name,algorithm, secret from tsigkeys"

	dec["get-all-domains-query"] = "select domains.id, domains.name, records.content, domains.type, domains.master, domains.notified_serial, domains.last_check, domains.account from domains LEFT JOIN records ON records.domain_id=domains.id AND records.type='SOA' AND records.name=domains.name WHERE records.disabled=0 OR :include_disabled"

	dec["list-comments-query"] = "SELECT domain_id,name,type,modified_at,account,comment FROM comments WHERE domain_id=:domain_id"
	dec["insert-comment-query"] = "INSERT INTO comments (domain_id, name, type, modified_at, account, comment) VALUES (:domain_id, :qname, :qtype, :modified_at, :account, :content)"
	dec["delete-comment-rrset-query"] = "DELETE FROM comments WHERE domain_id=:domain_id AND name=:qname AND type=:qtype"
	dec["delete-comments-query"] = "DELETE FROM comments WHERE domain_id=:domain_id"
	dec["search-records-query"] = record_query + " name LIKE :value ESCAPE '\\' OR content LIKE :value2 ESCAPE '\\' LIMIT :limit"
	dec["search-comments-query"] = "SELECT domain_id,name,type,modified_at,account,comment FROM comments WHERE name LIKE :value ESCAPE '\\' OR comment LIKE :value2 ESCAPE '\\' LIMIT :limit"

	return dec
}

func (db *Sqlite) Close() {
	_ = db.db.Close()
}

func (db *Sqlite) Query(stmt string, args ...interface{}) (storage.IResult, error) {
	if _, ok := db.declare[stmt]; !ok {
		return nil, stacktrace.New("Нет информации о запросе: " + stmt)
	}
	qs, names, err := sqlex.CompileNamedQuery(db.declare[stmt], db.bindType)
	if err != nil {
		return nil, stacktrace.Wrap(err)
	}
	arg, err := sqlex.ArgToMap(args...)
	if err != nil {
		return nil, stacktrace.Wrap(err)
	}

	parametrs := make([]interface{}, 0, len(names))
	for _, name := range names {
		if value, ok := arg[name]; !ok {
			parametrs = append(parametrs, nil)
		} else {
			parametrs = append(parametrs, value)
		}
	}

	rows, err := db.db.Query(qs, parametrs...)
	return rows, stacktrace.New("No implement")
}

func (db *Sqlite) Exec(stmt string, args ...interface{}) (int, error) {
	if _, ok := db.declare[stmt]; !ok {
		return 0, stacktrace.New("Нет информации о запросе: " + stmt)
	}
	qs, names, err := sqlex.CompileNamedQuery(db.declare[stmt], db.bindType)
	if err != nil {
		return 0, stacktrace.Wrap(err)
	}
	arg, err := sqlex.ArgToMap(args...)
	if err != nil {
		return 0, stacktrace.Wrap(err)
	}

	parametrs := make([]interface{}, 0, len(names))
	for _, name := range names {
		if value, ok := arg[name]; !ok {
			parametrs = append(parametrs, nil)
		} else {
			parametrs = append(parametrs, value)
		}
	}

	rows, err := db.db.Exec(qs, parametrs...)
	if err != nil {
		return 0, stacktrace.Wrap(err)
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return 0, stacktrace.Wrap(err)
	}
	return int(rowsAffected), nil
}

func (db *Sqlite) CreateTable() error {
	_, err := db.db.Exec(`
PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE domains (
  id                    INTEGER PRIMARY KEY,
  name                  VARCHAR(255) NOT NULL COLLATE NOCASE,
  master                VARCHAR(128) DEFAULT NULL,
  last_check            INTEGER DEFAULT NULL,
  type                  VARCHAR(6) NOT NULL,
  notified_serial       INTEGER DEFAULT NULL,
  account               VARCHAR(40) DEFAULT NULL
);

CREATE UNIQUE INDEX name_index ON domains(name);


CREATE TABLE records (
  id                    INTEGER PRIMARY KEY,
  domain_id             INTEGER DEFAULT NULL,
  name                  VARCHAR(255) DEFAULT NULL,
  type                  VARCHAR(10) DEFAULT NULL,
  content               VARCHAR(65535) DEFAULT NULL,
  ttl                   INTEGER DEFAULT NULL,
  prio                  INTEGER DEFAULT NULL,
  disabled              BOOLEAN DEFAULT 0,
  ordername             VARCHAR(255),
  auth                  BOOL DEFAULT 1,
  FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX records_lookup_idx ON records(name, type);
CREATE INDEX records_lookup_id_idx ON records(domain_id, name, type);
CREATE INDEX records_order_idx ON records(domain_id, ordername);


CREATE TABLE supermasters (
  ip                    VARCHAR(64) NOT NULL,
  nameserver            VARCHAR(255) NOT NULL COLLATE NOCASE,
  account               VARCHAR(40) NOT NULL
);

CREATE UNIQUE INDEX ip_nameserver_pk ON supermasters(ip, nameserver);


CREATE TABLE comments (
  id                    INTEGER PRIMARY KEY,
  domain_id             INTEGER NOT NULL,
  name                  VARCHAR(255) NOT NULL,
  type                  VARCHAR(10) NOT NULL,
  modified_at           INT NOT NULL,
  account               VARCHAR(40) DEFAULT NULL,
  comment               VARCHAR(65535) NOT NULL,
  FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX comments_idx ON comments(domain_id, name, type);
CREATE INDEX comments_order_idx ON comments (domain_id, modified_at);


CREATE TABLE domainmetadata (
 id                     INTEGER PRIMARY KEY,
 domain_id              INT NOT NULL,
 kind                   VARCHAR(32) COLLATE NOCASE,
 content                TEXT,
 FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX domainmetaidindex ON domainmetadata(domain_id);


CREATE TABLE cryptokeys (
 id                     INTEGER PRIMARY KEY,
 domain_id              INT NOT NULL,
 flags                  INT NOT NULL,
 active                 BOOL,
 published              BOOL DEFAULT 1,
 content                TEXT,
 FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX domainidindex ON cryptokeys(domain_id);


CREATE TABLE tsigkeys (
 id                     INTEGER PRIMARY KEY,
 name                   VARCHAR(255) COLLATE NOCASE,
 algorithm              VARCHAR(50) COLLATE NOCASE,
 secret                 VARCHAR(255)
);

CREATE UNIQUE INDEX namealgoindex ON tsigkeys(name, algorithm);
COMMIT;
`)
	return err
}
