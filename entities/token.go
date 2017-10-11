package entities

type Token struct {
	Id      string   `json:"id"`
	User    string   `json:"user"`
	Expires int64    `json:"expires"`
	Roles   []string `json:"roles"`
}

func (t *Token) String() string {
	return t.Id
}
