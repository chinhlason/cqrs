package query

type Document struct {
	Id     string      `json: _id`
	Type   string      `json: _type`
	Index  string      `json: _index`
	Source interface{} `json: _source`
}
