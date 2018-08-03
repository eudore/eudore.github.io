package ram

import (
	"strings"
	"strconv"
)


var (
	ConditionAddr		=	"Addr"
)

type DomainCondition struct {
	Domain	string
}

type AddrCondition struct {
	key		string
	Addr	int64
	Mark	int64
}

func NewAddrCondition(addr string) *AddrCondition{
	var c AddrCondition
	c.key = "addr"
	a := strings.Split(addr,"/")
	c.Addr = ip2int(a[0])
	if len(a)== 1 {
		a = append(a,"32")
	}
	mark ,_ := strconv.ParseInt(a[1], 10, 64)
	c.Mark = 0xffffffff << (32 - uint(mark) )
	return &c
}

func ip2int(ipString  string) int64 {     
	ipSegs := strings.Split(ipString, ".")
	var ipInt int64 = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.ParseInt(ipSeg, 10, 64)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

func (c *AddrCondition) Key() string {
	return c.key
}
func (c *AddrCondition) Check(cons map[string]string) bool{
	addint := ip2int(cons["addr"])
	return ((addint ^ c.Addr ) & c.Mark) == 0
}


type Condition interface {
	Key()	string
	Check(string)	bool
}

func NewCondition(con string) Condition {
	cons := strings.SplitN(con,":",2)
	switch cons[0] {
	case "addr":
		//return NewAddrCondition(cons[1])
	}
	return nil
}