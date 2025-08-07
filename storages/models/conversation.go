package models

type ConverConf struct {
	ID         int64
	ConverId   string
	ConverType int32
	SubChannel string
	ItemKey    string
	ItemValue  string
	ItemType   int32
	AppKey     string
}

type IConverConfStorage interface {
	Upsert(item ConverConf) error
	BatchUpsert(items []ConverConf) error
	QryConverConfs(appkey, converId, subChannel string, converType int32) (map[string]*ConverConf, error)
}
