package models

type ConverConfItem struct {
	ItemKey   string `json:"item_key"`
	ItemValue string `json:"item_value"`
	ItemType  int32  `json:"item_type"`
}

type SetConverConfsReq struct {
	TargetId   string                 `json:"target_id"`
	SubChannel string                 `json:"sub_channel"`
	ConverType int                    `json:"conver_type"`
	Confs      map[string]interface{} `json:"confs"`
}
