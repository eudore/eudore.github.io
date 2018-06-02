package config

import (
	"fmt"
	"errors"
	"strings"
	"strconv"
	"reflect"
)

var (
	errArg 		=	errors.New("undefined args")
	errType 	= 	errors.New("undefined type")
)

func (c *Config) help() {
	getcomment(c,"  --")
	fmt.Println("  --help\t\tShow help")
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

func (c *Config) set(arg string) error {
	kv := append(strings.SplitN(arg,"=",2),"")
	k,v := kv[0],kv[1]
	var p interface{} = c
	pt := reflect.TypeOf(p).Elem()
	pv := reflect.ValueOf(p).Elem()
	fs := strings.Split(k,".")
	len := len(fs) - 1
	for i,_ := range fs {
		f,ok := pt.FieldByName(strings.Title(fs[i]))
		if !ok{
			// error
			return errArg
		}
		sv := pv.Field(f.Index[0])
		if i == len {
			return setvalue(sv,v)
		}
		// is pointer
		if sv.Type().Kind() == reflect.Ptr {
			// is null
			if sv.Elem().Kind() == reflect.Invalid {
				sv.Set(reflect.New(sv.Type().Elem()))	
			}
			sv=sv.Elem()
		}
		if sv.Type().Kind() != reflect.Struct {
			return errType
		}
		// next
		p = sv.Addr().Interface()
		pt = reflect.TypeOf(p).Elem()
		pv = reflect.ValueOf(p).Elem()
	}
	return nil
}
func setvalue(v reflect.Value,s string) error {
	switch v.Kind() {
	case reflect.Bool:
		if s == "" {
			v.SetBool(true)
		}else {
			rb,_ := strconv.ParseBool(s)
			v.SetBool(rb)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v.OverflowInt(n) {
			return errType
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || v.OverflowUint(n) {
			return errType
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			return errType
		}
		v.SetFloat(n)
	case reflect.String:
		v.SetString(s)
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
	default:
		fmt.Println("default")
	}
	return nil
}