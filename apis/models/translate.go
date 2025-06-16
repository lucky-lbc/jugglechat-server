package models

type TransItem struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

type TransReq struct {
	Items      []*TransItem `json:"items"`
	SourceLang string       `json:"source_lang"`
	TargetLang string       `json:"target_lang"`
}
