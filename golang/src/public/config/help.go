package config

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
	"reflect"
	"encoding/json"
)

var (
	errArg 			=	errors.New("undefined args")
	errType 		= 	errors.New("undefined type")
	errValue 		= 	errors.New("undefined value type")
	errUndefined	=	errors.New("setdata use undefined interface")
)

type SetConfig interface {
	SetData(arg string) error
}

// Output struct commment info
func Help(c interface{}) {
	getcomment(c,"  --")
	getcomment(configinfos[c],"  --")
}

func getcomment(p interface{},prefix string) {
	pt := reflect.TypeOf(p).Elem()
	pv := reflect.ValueOf(p).Elem()
	for i := 0; i < pt.NumField(); i++ {
		sv := pv.Field(i)
		if c :=  pt.Field(i).Tag.Get("comment");c != "" {
			fmt.Printf("%s%s\t%s\n",prefix,strings.ToLower(pt.Field(i).Name),c)
		}
		if c := pt.Field(i).Tag.Get("help");c == "-" {
			continue
		}
		if sv.Type().Kind() == reflect.Ptr {
			// is null
			if sv.Elem().Kind() == reflect.Invalid {
				sv.Set(reflect.New(sv.Type().Elem()))	
			}
			sv=sv.Elem()
		}
		if sv.Kind() == reflect.Struct {
			getcomment(sv.Addr().Interface(),fmt.Sprintf("%s%s.",prefix,strings.ToLower(pt.Field(i).Name)))
		}
	}
}

// Set Data
func SetData(p interface{},arg string) error{
	kv := append(strings.SplitN(arg,"=",2),"")
	d,err := Data(p,kv[0])
	setvalue(d,kv[1])
	return err
}

func GetData(p interface{},arg string) (interface{}, error){
	kv := append(strings.SplitN(arg,"=",2),"")
	d,err := Data(p,kv[0])
	return d.Interface(),err
}

func Data(p interface{}, arg string) (reflect.Value, error) {
	fs := strings.Split(arg,".")
	len := len(fs) - 1
	for i,_ := range fs {
		pv := reflect.ValueOf(p)
		for pv.Kind() == reflect.Ptr || pv.Kind() == reflect.Interface {
			pv= pv.Elem()
		}
		//fmt.Println("\n----- ",pv.Kind(),pv.Type().Kind(),fs[i])
		if pv.Kind() == reflect.Struct {
			f,ok := pv.Type().FieldByName(strings.Title(fs[i]))
			if !ok{
				// error
				return pv,errArg
			}
			pv = pv.Field(f.Index[0])
			if i == len {
				return pv,nil
			}
		}
		// is null
		switch pv.Kind() {
		case reflect.Ptr:
			if pv.IsNil() {
				fmt.Println("+new struct ptr")
				pv.Set(reflect.New(pv.Type().Elem()))	
			}
			pv=pv.Elem()
		case reflect.Map,reflect.Interface:
			return pv,errArg
		}
		p = pv.Addr().Interface()
	}
	return reflect.ValueOf(p),nil
}

func setvalue(v reflect.Value,s string) error {
	switch v.Kind() {
	case reflect.Bool,
	reflect.Float32, reflect.Float64,
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
	reflect.String:
		k ,_ := getvalue(v.Type(),s)
		v.Set(k)
	case reflect.Ptr:
		if v.Elem().Kind() == reflect.Invalid {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return setvalue(v.Elem(),s)
	case reflect.Array:
	case reflect.Slice:
		vs := strings.Split(s,",")
		v.Set(reflect.MakeSlice(v.Type(),len(vs),len(vs)))
		for i,n := range vs {
			setvalue(v.Index(i),n)
		}
	case reflect.Map:
		fmt.Println("Interface Map")
		return json.Unmarshal([]byte(s),v.Interface())
	case reflect.Interface:
		fmt.Println("Interface")
		v.Set(reflect.ValueOf(s))
	default:
		fmt.Println("default")
		return errValue
	}
	return nil
}

func getvalue(t reflect.Type,s string) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Bool:
		if s == "" {
			return reflect.ValueOf(true), nil
		}else {
			rb,_ := strconv.ParseBool(s)
			return reflect.ValueOf(rb), nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Zero(t), errType
		}
		return reflect.ValueOf(n).Convert(t),nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return reflect.Zero(t), errType
		}
		return reflect.ValueOf(n).Convert(t),nil
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return reflect.Zero(t), errType
		}
		return reflect.ValueOf(n).Convert(t),nil
	case reflect.String:
		return reflect.ValueOf(s),nil
	}
	return reflect.Zero(t), errType
}