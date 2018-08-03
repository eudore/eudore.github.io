package ram

import (
	"strings"
	"strconv"
)

var Policys map[string]Policy
var methods map[string]int32

func init() {
	Policys = make(map[string]Policy)
	methods = make(map[string]int32)
	for i,m := range []string{"Get", "Post", "Put", "Delete", "Head", "Patch", "Options"} {
		methods[m] = int32(i)
	}
}

type Policy interface {
	ID() string
	Match(map[string]string) bool
	Effect() bool
}

type PolicyInfo struct {
	ID			string
	Type		string
	Config		string
	Condition	string
	Effect		bool
}
// type Role interface {
// 	ID() string
// 	Permit(Permission) bool
// }

func NewPolicy(pid string) (Policy, error) {
	var p PolicyInfo
	err := stmtQueryRamPolicy.QueryRow(pid).Scan(&p.Type, &p.Config, &p.Condition, &p.Effect)
	if err != nil {
		return  nil, err
	}
	p.ID = pid
	switch p.Type {
	case "http":
		return NewHttpPolicy(p),nil
	}
	return nil,nil
}


func GetPolicy(pid string) Policy {
	p, ok := Policys[pid]
	if !ok {
		p, _ = NewPolicy(pid)
		Policys[pid] = p
	}
	return p
}

func Match(ps []string,cons map[string]string) bool {
	for _,i := range ps {
		p := GetPolicy(i)
		if p != nil && p.Match(cons) {
			return p.Effect()
		}
	}
	return false
}




const (
	Get 	=	1 << iota
	Post
	Put
	Delete
	Head
	Patch
	Options
)


type HttpPolicy struct {
	id			string
	effect		bool
	Method		int32
	URL			[]string
	Condition	[]Condition
}
func NewHttpPolicy(p PolicyInfo) Policy {
	c := strings.Split(p.Config," ")
	n,_ := strconv.ParseInt(c[0], 16, 32)
	return &HttpPolicy{
		id:		p.ID,
		effect:	p.Effect,
		Method:	int32(n),
		URL:	c[1:],
	}
}
func (p *HttpPolicy) ID() string {
	return p.id
}
func (p *HttpPolicy) Effect() bool {
	return p.effect
}
func (p *HttpPolicy) Match(cons map[string]string) bool {
	me := cons["method"]
	if (p.Method & MethodID(me)) != p.Method {
		return false
	}
	// for _,a := range p.Condition {
	// 	a.Check("")
	// }
	url := cons["url"]
	for _,i := range p.URL {
		if MatchStar(i,url) {
			return true
		}
	}
	return false
}
func MethodID(me string) int32{
	id,ok := methods[me]
	if ok {
		return id
	}
	return -1
}