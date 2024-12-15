package model

type BulksmsLogData struct {
	TID             string `json:"component,omitempty"`
	AppID           string `json:"appid,omitempty"`
	FeedID          string `json:"feedid,omitempty"`
	EntID           string `json:"entid,omitempty"`
	Keyword         string `json:"keyword,omitempty"`
	OrgTemplate     string `json:"orgtemplate,omitempty"`
	DotStarCnt      string `json:"dotstarcnt,omitempty"`
	SpaceFlag       string `json:"spaceflag,omitempty"`
	SpecialCharFlag string `json:"specialcharflag,omitempty"`
	CustomDomain    string `json:"customdomain,omitempty"`
	Token           string `json:"token,omitempty"`
	DMCheckStatus   string `json:"dmcheckstatus,omitempty"`
	BReqID          string `json:"breqid,omitempty"`
	BTID            string `json:"btid,omitempty"`
	TraiCategoryID  string `json:"traicategoryid,omitempty"`
	TraiMessageType string `json:"traimessagetype,omitempty"`
	TraiMessageMode string `json:"traimessagemode,omitempty"`
	BSMSInTime      string `json:"bsms_intime,omitempty"`
	TemplateID      string `json:"template_id,omitempty"`
	BMsgTag         string `json:"bmsgtag,omitempty"`
	Text            string `json:"text,omitempty"`
	To              string `json:"to,omitempty"`
	From            string `json:"from,omitempty"`
	DLTEntityID     string `json:"dltentityid,omitempty"`
	BSMSOutTime     string `json:"bsms_outtime,omitempty"`
}
