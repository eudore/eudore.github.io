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


func (c *Config) help() {
	getcomment(c,"  --")
	fmt.Println("  --help\t\tShow help")
}

func (c *Config) SetData(arg string) error {
	return SetData(c,arg)
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
		if sv.Type().Kind() == reflect.Interface {
			getcomment(sv.Interface(),fmt.Sprintf("%s%s.",prefix,strings.ToLower(pt.Field(i).Name)))
		}
		if sv.Kind() == reflect.Struct {
			getcomment(sv.Addr().Interface(),fmt.Sprintf("%s%s.",prefix,strings.ToLower(pt.Field(i).Name)))
		}
	}
}

func SetData(p interface{} , arg string) error {
	kv := append(strings.SplitN(arg,"=",2),"")
	k,v := kv[0],kv[1]
	fs := strings.Split(k,".")
	len := len(fs) - 1
	var sv reflect.Value
	for i,_ := range fs {
		pv := reflect.ValueOf(p)
		for pv.Kind() == reflect.Ptr || pv.Kind() == reflect.Interface {
			pv= pv.Elem()
		}
		fmt.Println("\n----- ",pv.Kind(),pv.Type().Kind(),fs[i])
		switch pv.Kind() {
		case reflect.Struct:
			f,ok := pv.Type().FieldByName(strings.Title(fs[i]))
			if !ok{
				// error
				fmt.Println(errArg)
				return errArg
			}
			sv = pv.Field(f.Index[0])
			if i == len {
				return setvalue(sv,v)
			}
		case reflect.Map:
			fmt.Println("map1:",i)
			fmt.Println("ass= ",sv.Kind())
			sv1 := reflect.ValueOf(pv.Interface().(map[string]interface{})) 
			fmt.Println(sv1.Type(),sv1.Type().Key(),sv1.Kind())
			if i==len {
				sv1.SetMapIndex(reflect.ValueOf(fs[i]), reflect.ValueOf(v))
				return nil
			}else {
				fmt.Println("-",sv1,sv1.MapIndex(reflect.ValueOf(fs[i])))
				if sv1.MapIndex(reflect.ValueOf(fs[i])).Kind() == reflect.Invalid {
					b := make(map[string]interface{})
					sv1.SetMapIndex(reflect.ValueOf(fs[i]), reflect.ValueOf(&b))	
					sv = reflect.ValueOf(b)
					//sv = sv1.MapIndex(reflect.ValueOf(fs[i]))
					//sv = reflect.ValueOf(sv)
					// var aa interface{}
					// aa = b
			// 		sv = reflect.New(reflect.TypeOf(sv1))
			// 		sv.Elem().Set(reflect.ValueOf(make(map[string]interface{})))
			 fmt.Println("---",sv)
			// // fmt.Println(sv,sv.Elem(),sv.Addr())
					// sv1.SetMapIndex(reflect.ValueOf(fs[i]), reflect.MakeMap(reflect.TypeOf(b)) )
					
				}else {
					sv = reflect.ValueOf(sv.Interface().(map[string]interface{})[fs[i]])
				}
				//sv = sv1.MapIndex(reflect.ValueOf(fs[i]))
				//fmt.Println(reflect.TypeOf(b))
				fmt.Println(sv.Kind(),sv.Type())
			}
		}
		// is pointer
		if sv.Type().Kind() == reflect.Ptr {
			// is null
			if sv.Elem().Kind() == reflect.Invalid {
				fmt.Println(sv.IsNil(),"ssssssssssssssssssssssss")
				sv.Set(reflect.New(sv.Type().Elem()))	
			}
			fmt.Println("ssssssss")
			sv=sv.Elem()
		}
		if sv.Type().Kind() == reflect.Interface && sv.IsNil() {
			fmt.Println("new map")
			sv.Set(reflect.ValueOf(make(map[string]interface{})))
		}
		// next
		fmt.Println("Addr",sv.CanAddr())
		if sv.CanAddr() {
			p = sv.Addr().Interface()	
		}else {
			p= sv.Interface()
		}
	}
	fmt.Println("--end")
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
	case reflect.Interface:
		fmt.Println("Interface")
	case reflect.Map:
		fmt.Println("Interface Map")
		return json.Unmarshal([]byte(s),v.Interface())
	default:
		fmt.Println("default")
		return errValue
	}
	return nil
}







