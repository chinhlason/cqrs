package sync

type User struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	CreateAt  int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	Orders    []Order `json:"orders"`
}

type Order struct {
	Id        int64  `json:"id"`
	UserId    int64  `json:"user_id"`
	Product   string `json:"product"`
	CreateAt  int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Source struct {
	Version   string `json:"version"`
	Connector string `json:"connector"`
	name      string `json:"name"`
	Tsms      int64  `json:"ts_ms"`
	Snapshot  string `json:"snapshot"`
	Db        string `json:"db"`
	Schema    string `json:"schema"`
	Table     string `json:"table"`
	TxId      int64  `json:"txId"`
	LsSn      int64  `json:"lsn"`
}

type UserMessage struct {
	Before *User   `json:"before"`
	After  *User   `json:"after"`
	Source *Source `json:"source"`
	Op     string  `json:"op"`
	Tsms   int64   `json:"ts_ms"`
}

type OrderMessage struct {
	Before *Order  `json:"before"`
	After  *Order  `json:"after"`
	Source *Source `json:"source"`
	Op     string  `json:"op"`
	Tsms   int64   `json:"ts_ms"`
}

type ESRecord struct {
	Index  string      `json:"_index"`
	Id     string      `json:"_id"`
	Source interface{} `json:"_source"`
}
