package models

type ConverConfItem struct {
	ItemKey   string `json:"item_key"`
	ItemValue string `json:"item_value"`
	ItemType  int32  `json:"item_type"`
}
