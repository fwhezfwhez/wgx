package wgx

import "encoding/json"

type ReasonPanel struct {
	Scene string                 `json:"scene"`
	Args  map[string]interface{} `json:"args"`
}

func (rp ReasonPanel) JSONLine() string {
	rs, _ := json.Marshal(rp)
	return string(rs)
}
