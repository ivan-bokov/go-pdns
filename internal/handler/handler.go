package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ivan-bokov/go-pdns/internal/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) noImplementation(g *gin.Context) {
	g.JSON(200, gin.H{"result": false})
}

func (h *Handler) logAllResponse() gin.HandlerFunc {
	return func(g *gin.Context) {
		uri := g.Request.URL.RequestURI()
		method := g.Request.Method
		body, err := io.ReadAll(g.Request.Body)
		if err != nil {
			log.Println(fmt.Sprintf("[ERROR]: %#v", err))
		}
		g.Request.Body.Close()
		log.Println(fmt.Sprintf("URI:%s Method:%s Headers: %#v Body:%s", uri, method, g.Request.Header, string(body)))
		g.Next()
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(h.logAllResponse())
	r.GET("lookup/:qname/:qtype", h.lookup)    //++++
	r.GET("list/:domain_id/:zonename", h.list) // ++++
	r.GET("getbeforeandafternamesabsolute/:domain_id/:qname", h.getbeforeandafternamesabsolute)
	r.GET("getalldomainmetadata/:name", h.getAllDomainMetadata)    // ++++
	r.GET("getdomainmetadata/:name/:kind", h.getDomainMetadata)    // ++++
	r.PATCH("setdomainmetadata/:name/:kind", h.setDomainMetadata)  //++++
	r.GET("getdomainkeys/:name", h.getDomainKeys)                  // ++++
	r.PUT("adddomainkey/:name", h.addDomainKey)                    //+++?
	r.DELETE("removedomainkey/:name/:id", h.removeDomainKey)       // ++++
	r.POST("activatedomainkey/:name/:id", h.activateDomainKey)     // ++++
	r.POST("deactivatedomainkey/:name/:id", h.deactivateDomainKey) // ++++
	r.POST("publishdomainkey/:name/:id", h.publishDomainKey)       // ++++
	r.POST("unpublishdomainkey/:name/:id", h.unpublishDomainKey)   // ++++
	r.GET("gettsigkey/:name", h.getTSIGKey)                        // ++++
	r.GET("getdomaininfo/:name", h.getDomainInfo)                  // ++++
	r.PATCH("setnotified/:id", h.setNotified)                      // ++++
	r.GET("isMaster/:name/:ip", h.noImplementation)
	r.POST("supermasterbackend/:ip/:domain", h.noImplementation)
	r.POST("createslavedomain/:ip/:domain", h.createSlaveDomain) //++++
	r.PATCH("replacerrset/:domain_id/:qname/:qtype", h.noImplementation)
	r.PATCH("feedrecord/:trxid", h.feedRecord) //++--
	r.PATCH("feedents/:domain_id", h.noImplementation)
	r.PATCH("feedEnts3/:domain_id/:domain", h.noImplementation)
	r.POST("starttransaction/:domain_id/:domain", h.noImplementation)
	r.POST("committransaction/:trxid", h.noImplementation)
	r.POST("aborttransaction/:trxid", h.noImplementation)
	r.POST("calculatesoaserial/:domain", h.noImplementation)
	r.POST("directBackendCmd", h.noImplementation)
	r.GET("getAllDomains", h.getAllDomains) // ++++
	r.GET("searchRecords", h.noImplementation)
	r.GET("getUpdatedMasters", h.noImplementation)
	r.GET("getUnfreshSlaveInfos", h.getUnfreshSlaveInfos) // ++++
	r.PATCH("setFresh/:id", h.setFresh)                   // ++++

	return r
}
func (h *Handler) getTSIGKey(g *gin.Context) {
	name := g.Param("name")
	key, err := h.svc.GetTSIGKey(name)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": key})
}

func (h *Handler) unpublishDomainKey(g *gin.Context) {
	name := g.Param("name")
	var id int
	var err error
	if g.Param("id") != "" {
		id, err = strconv.Atoi(g.Param("id"))
		if err != nil {
			g.JSON(200, gin.H{"result": false})
			return
		}
	}
	err = h.svc.UnpublishDomainKey(name, id)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}
func (h *Handler) publishDomainKey(g *gin.Context) {
	name := g.Param("name")
	var id int
	var err error
	if g.Param("id") != "" {
		id, err = strconv.Atoi(g.Param("id"))
		if err != nil {
			g.JSON(200, gin.H{"result": false})
			return
		}
	}
	err = h.svc.PublishDomainKey(name, id)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}
func (h *Handler) deactivateDomainKey(g *gin.Context) {
	name := g.Param("name")
	var id int
	var err error
	if g.Param("id") != "" {
		id, err = strconv.Atoi(g.Param("id"))
		if err != nil {
			g.JSON(200, gin.H{"result": false})
			return
		}
	}
	err = h.svc.DeactivateDomainKey(name, id)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}

func (h *Handler) activateDomainKey(g *gin.Context) {
	name := g.Param("name")
	var id int
	var err error
	if g.Param("id") != "" {
		id, err = strconv.Atoi(g.Param("id"))
		if err != nil {
			g.JSON(200, gin.H{"result": false})
			return
		}
	}
	err = h.svc.ActivateDomainKey(name, id)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}

func (h *Handler) removeDomainKey(g *gin.Context) {
	name := g.Param("name")
	var id int
	var err error
	if g.Param("id") != "" {
		id, err = strconv.Atoi(g.Param("id"))
		if err != nil {
			g.JSON(200, gin.H{"result": false})
			return
		}
	}
	err = h.svc.RemoveDomainKey(name, id)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}

func (h *Handler) getDomainKeys(g *gin.Context) {
	name := g.Param("name")
	keys, err := h.svc.GetDomainKeys(name)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": keys})
}

func (h *Handler) getDomainMetadata(g *gin.Context) {
	name := g.Param("name")
	kind := g.Param("kind")
	meta, err := h.svc.GetDomainMetadata(name, kind)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": meta})

}

func (h *Handler) getUnfreshSlaveInfos(g *gin.Context) {
	di, err := h.svc.GetUnfreshSlaveInfos()
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": di})

}

func (h *Handler) getAllDomains(g *gin.Context) {
	disabled := false
	var err error
	if g.Query("includeDisabled") != "" {
		disabled, err = strconv.ParseBool(g.Query("includeDisabled"))
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"result": false})
			return
		}
	}
	di, err := h.svc.GetAllDomains(disabled)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": di})
}

func (h *Handler) lookup(g *gin.Context) {
	qtype := g.Param("qtype")
	qname := g.Param("qname")
	zoneID := -1
	var err error
	if g.Request.Header.Get("X-RemoteBackend-zone-id") != "" {
		zoneID, err = strconv.Atoi(g.Request.Header.Get("X-RemoteBackend-zone-id"))
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"result": false})
			return
		}
	}
	listRR, err := h.svc.Lookup(qtype, qname, zoneID)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": listRR})
}
func (h *Handler) getDomainInfo(g *gin.Context) {
	name := g.Param("name")
	di, err := h.svc.GetDomainInfo(name)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": di})
}

func (h *Handler) list(g *gin.Context) {
	zonename := g.Param("zonename")
	domainID := -1
	var err error
	if g.Request.Header.Get("X-RemoteBackend-domain-id") != "" {
		domainID, err = strconv.Atoi(g.Request.Header.Get("X-RemoteBackend-domain-id"))
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"result": false})
			return
		}
	}
	if g.Param("domain_id") != "" {
		domainID, err = strconv.Atoi(g.Param("domain_id"))
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"result": false})
			return
		}
	}
	listRR, err := h.svc.List(zonename, domainID, false)
	if err != nil {
		g.JSON(200, gin.H{"result": make([]string, 0)})
		return
	}
	g.JSON(200, gin.H{"result": listRR})
}
func (h *Handler) getAllDomainMetadata(g *gin.Context) {
	name := g.Param("name")
	var err error
	meta, err := h.svc.GetAllDomainMetadata(name)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": meta})
}

func (h *Handler) getbeforeandafternamesabsolute(g *gin.Context) {
	//TODO непонятно что делать с параметрами, разобраться когда все закончу либо осенит

}

func (h *Handler) setDomainMetadata(g *gin.Context) {
	name := g.Param("name")
	kind := g.Param("kind")
	type valueMetadata struct {
		Value []string `json:"value,omitempty" form:"value"`
	}
	values := new(valueMetadata)
	if err := g.Bind(values); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}
	err := h.svc.SetDomainMetadata(name, kind, values.Value)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}

func (h *Handler) addDomainKey(g *gin.Context) {
	name := g.Param("name")
	key := new(service.KeyData)
	var err error
	if flags, ok := g.GetPostForm("flags"); ok {
		key.Flags, err = strconv.Atoi(flags)
		if err != nil {
			g.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	if active, ok := g.GetPostForm("active"); ok {
		key.Active, err = strconv.ParseBool(active)
		if err != nil {
			g.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	if published, ok := g.GetPostForm("published"); ok {
		key.Published, err = strconv.ParseBool(published)
		if err != nil {
			g.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	if content, ok := g.GetPostForm("content"); ok {
		key.Content = content
	}

	err = h.svc.AddDomainKey(name, key)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}

func (h *Handler) feedRecord(g *gin.Context) {
	m := make(map[string]string)
	var ok bool
	if m, ok = g.GetPostFormMap("rr"); !ok {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}
	ttl, err := strconv.Atoi(m["ttl"])
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}
	var auth bool
	if auth, err = strconv.ParseBool(m["auth"]); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}
	err = h.svc.FeedRecord(&service.DNSResourceRecord{
		Qname:   m["qname"],
		Content: m["content"],
		TTL:     ttl,
		Qtype:   m["qtype"],
		Auth:    auth,
		Qclass:  m["qclass"],
	}, "")
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}

func (h *Handler) createSlaveDomain(g *gin.Context) {
	ip := g.Param("ip")
	domain := g.Param("domain")
	err := h.svc.CreateSlaveDomain(ip, domain)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}

func (h *Handler) setFresh(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}
	err = h.svc.SetFresh(id)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}

func (h *Handler) setNotified(g *gin.Context) {
	id, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"result": false})
		return
	}
	var serial int
	if s, ok := g.GetPostForm("serial"); ok {
		serial, err = strconv.Atoi(s)
		if err != nil {
			g.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	err = h.svc.SetNotified(id, serial)
	if err != nil {
		g.JSON(200, gin.H{"result": false})
		return
	}
	g.JSON(200, gin.H{"result": true})
}
