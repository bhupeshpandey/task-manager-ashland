package msgqueue

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}
