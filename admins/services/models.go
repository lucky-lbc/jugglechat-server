package services

type Accounts struct {
	Items   []*Account `json:"items"`
	HasMore bool       `json:"has_more"`
	Offset  string     `json:"offset"`
}
type Account struct {
	Account       string `json:"account"`
	State         int    `json:"state"`
	CreatedTime   int64  `json:"created_time"`
	UpdatedTime   int64  `json:"updated_time"`
	ParentAccount string `json:"parent_account"`
	RoleId        int    `json:"role_id"`
}
