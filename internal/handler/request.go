package handler

type addDomainKeyForm struct {
	Flags     int    `form:"flags"`
	Active    bool   `form:"active"`
	Published bool   `form:"published"`
	Content   string `form:"content"`
}

type feedRecordForm struct {
}
