/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"html/template"
	"reflect"
	"strings"
)

func Render(templated interface{}, data interface{}) (interface{}, error) {
	params, err := structToMap(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template params; %w", err)
	}
	if err := renderStruct(reflect.ValueOf(templated), params); err != nil {
		return nil, fmt.Errorf("failed to render spec; %w", err)
	}
	return templated, nil
}

func structToMap(in interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	inrec, _ := json.Marshal(in)
	if err := json.Unmarshal(inrec, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal struct to map; %w", err)
	}
	return out, nil
}

func renderStruct(v reflect.Value, data interface{}) error {
	s := v
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if err := renderValue(f, data); err != nil {
			return fmt.Errorf("failed to render struct; %w", err)
		}
	}
	return nil
}

func renderValue(v reflect.Value, data interface{}) error {
	if v.Kind() == reflect.String && v.CanSet() {
		tpl := template.New("_").Funcs(sprig.FuncMap())
		parsed, err := tpl.Parse(v.String())
		if err != nil {
			return fmt.Errorf("failed to parse template; %w", err)
		}
		rendered := new(strings.Builder)
		if err := parsed.Execute(rendered, data); err != nil {
			return fmt.Errorf("failed to execute template; %w", err)
		}
		v.SetString(rendered.String())
	} else if v.Kind() == reflect.Struct {
		x1 := reflect.ValueOf(v.Interface())
		if err := renderStruct(x1, data); err != nil {
			return fmt.Errorf("failed to render value; %w", err)
		}
	} else if v.Kind() == reflect.Ptr {
		el := v.Elem()
		if el.Kind() == reflect.Struct {
			if err := renderStruct(el, data); err != nil {
				return fmt.Errorf("failed to render value; %w", err)
			}
		} else {
			if err := renderValue(el, data); err != nil {
				return fmt.Errorf("failed to render value; %w", err)
			}
		}
	}
	return nil
}
