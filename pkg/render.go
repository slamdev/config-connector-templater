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

	jsonStr, err := json.Marshal(templated)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct; %w", err)
	}

	rendered, err := renderString(string(jsonStr), params)
	if err != nil {
		return nil, fmt.Errorf("failed to render string; %w", err)
	}

	structType := reflect.TypeOf(templated)
	outPtr := reflect.New(structType).Interface()

	if err := json.Unmarshal([]byte(rendered), outPtr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal struct to map; %w", err)
	}

	out := reflect.ValueOf(outPtr).Elem().Interface()

	return out, nil
}

func renderString(str string, params map[string]interface{}) (string, error) {
	tpl := template.New("_").Funcs(sprig.FuncMap())
	parsed, err := tpl.Parse(str)
	if err != nil {
		return "", fmt.Errorf("failed to parse template; %w", err)
	}
	rendered := new(strings.Builder)
	if err := parsed.Execute(rendered, params); err != nil {
		return "", fmt.Errorf("failed to execute template; %w", err)
	}
	return rendered.String(), nil
}

func structToMap(in interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	jsonStr, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct; %w", err)
	}
	if err := json.Unmarshal(jsonStr, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal struct to map; %w", err)
	}
	return out, nil
}
