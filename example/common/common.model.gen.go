// Code generated by https://github.com/chenxyzl/gsgen; DO NOT EDIT.
// gen_tools version: 1.1.4
// generate time: 2024-06-19 16:46:34
package common

import (
	"encoding/json"
	"fmt"
)

func (s *Common) String() string {
	doc := struct {
	}{}
	return fmt.Sprintf("%v", &doc)
}
func (s *Common) MarshalJSON() ([]byte, error) {
	doc := struct {
	}{}
	return json.Marshal(doc)
}
func (s *Common) UnmarshalJSON(data []byte) error {
	doc := struct {
	}{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	return nil
}
func (s *Common) Clone() (*Common, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	ret := Common{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
