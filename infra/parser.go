package infra

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func elemtype(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Ptr && t.Kind() != reflect.Interface {
		return t
	}
	return elemtype(t.Elem())
}

func elemvalue(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Interface {
		return v
	}
	if v.IsNil() {
		return v
	}
	return elemvalue(v.Elem())
}

func EncodeConfig(f io.Writer, config ConfigEntity) error {
	w := bufio.NewWriter(f)
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	for i := range config {
		sectionType := elemtype(reflect.TypeOf(config[i]))
		sectionValue := elemvalue(reflect.ValueOf(config[i]))
		for j := 0; j < sectionType.NumField(); j++ {
			sectionName := sectionType.Field(j).Tag.Get("section")
			if sectionName != "" {
				fmt.Fprintf(w, "[%s]\n", sectionName)
				continue
			}
			optionKey := sectionType.Field(j).Tag.Get("option")
			optionValue := elemvalue(sectionValue.Field(j))
			fmt.Fprintf(w, "%s = %s\n", optionKey, optionValue)
		}
	}
	return w.Flush()
}

func DecodeConfig(f io.Reader, config ConfigEntity) ConfigEntity {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	sc := bufio.NewScanner(f)
	type kvpair struct {
		key   string
		value string
	}
	type sect struct {
		name  string
		pairs []kvpair
	}
	type cfgslice []sect
	var index int = -1
	var line string
	var k string
	var v string
	cfg := make(cfgslice, 0, 10)
	for sc.Scan() {
		line = strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			index++
			cfg = append(cfg, sect{
				name:  strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"),
				pairs: make([]kvpair, 0, 10),
			})
			continue
		}
		kv := strings.Split(line, "=")
		if len(kv) > 1 {
			k = strings.TrimSpace(kv[0])
			if len(kv) > 2 {
				v = strings.TrimSpace(strings.Join(kv[1:], ""))
			} else {
				v = strings.TrimSpace(kv[1])
			}
			if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
				cfg[index].pairs = append(cfg[index].pairs, kvpair{
					key:   k,
					value: strings.Trim(v, "\""),
				})
				continue
			}
			if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
				cfg[index].pairs = append(cfg[index].pairs, kvpair{
					key:   k,
					value: strings.Trim(v, "'"),
				})
				continue
			}
			cfg[index].pairs = append(cfg[index].pairs, kvpair{
				key:   k,
				value: v,
			})
			continue
		}
	}
	var newConfig ConfigEntity
	for _, s := range cfg {
		section := NewConfigSection(s.name)
		if section == nil {
			continue
		}
		sectionType := elemtype(reflect.TypeOf(section))
		sectionValue := elemvalue(reflect.ValueOf(section))
		secvals := make(map[string]reflect.Value)
		for i := 0; i < sectionType.NumField(); i++ {
			key := sectionType.Field(i).Tag.Get("option")
			if key == "" {
				continue
			}
			secvals[key] = sectionValue.Field(i)
		}
		for _, kvp := range s.pairs {
			fv := reflect.ValueOf(kvp.value)
			secvals[kvp.key].Set(fv)
		}
		newConfig = append(newConfig, section)
	}

	return newConfig
}
