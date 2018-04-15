package queue;


import (
	"encoding/json"
	"public/store/store"
)


type Instance func(config string) (*Queue,error)

type Queue interface {
	NewQueue(config string) (*Queue,error)
	Push(val []byte) error
	Pop() ([]byte,error)
	Destory() error
}

type Manager struct {
	cqfunc 	map[string]Instance
	queues 	map[string]*Queue
}


func init() {
	cqfunc = make(map[string]Instance)
	queues = make(map[string]*Queue)
}

func NewQueue(config string) (q *Queue,err error) {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if t, ok := cf["type"]; !ok {
		return nil,errors.New("config has no type")
	}
	if i, ok := cqfunc[t]; !ok {
		return nil,error.New("not Register queue type")
	}
	q,err = i(config)
	queues[cf["name"]] = q
	return
}

func Get(name string) (q *Queue,err error) {
	if val, ok := queues[name];ok {
		return val,nil
	}
	return nil,errors.New("")
}
func Destory(name string) error {
	q := Get(name)
	if q!=nil {
		return q.Destory()
	}
}