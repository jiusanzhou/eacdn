package provider

import (
	"encoding/json"
	"errors"
	"reflect"

	"go.zoe.im/eacdn/pkg/core"
)

var (
	// provider creator registry
	_registry = make(map[string]func(data json.RawMessage) (Provider, error))
)

// TODO: with more host instance with ssh

// Provider define a cdn provider
type Provider interface {
	// CDN interface, TODO: with options
	Init(...Option) error
	Create(*core.Site, ...Option) error
	Update(*core.Site, ...Option) error
	// TODO: add logs and metrics subscribe
}

// Register the Provider creator, auto build the gen function from config or instance struct
// TODO: move to x
func Register(p Provider, typs ...string) error {

	vtype := reflect.TypeOf(p)
	if vtype.Kind() == reflect.Ptr {
		vtype = vtype.Elem()
	}

	// build the generator function
	fn := func(data json.RawMessage) (Provider, error) {
		src := reflect.New(vtype).Interface().(Provider)
		err := json.Unmarshal(data, &src)
		return src, err
	}

	// typ := src.Type()
	for _, typ := range typs {
		if _, ok := _registry[typ]; ok {
			return errors.New("type exits: " + typ)
		}
		_registry[typ] = fn
	}

	return nil
}

// Object for config unmarshal
type Object struct {
	Type string `json:"type,omitempty" yaml:"type"`
	Name string `json:"name,omitempty" yaml:"name"`

	// ==================================
	// TODO: auto have this data
	// TODO: auto have this data
	_raw       json.RawMessage
	_rawfields map[string]json.RawMessage
	// ==================================

	Provider `json:"-" yaml:"-"`
}

// MarshalJSON ...
func (s Object) MarshalJSON() ([]byte, error) {
	// TODO: marshal from the real one
	return json.Marshal(s._rawfields)
}

// UnmarshalJSON ...
func (s *Object) UnmarshalJSON(data []byte) error {

	s._raw = data

	err := json.Unmarshal(data, &s._rawfields)
	if err != nil {
		return err
	}

	// TODO: auto have this fields unmarshalling
	// TODO: auto have this fields unmarshalling
	json.Unmarshal(s._rawfields["type"], &s.Type)
	json.Unmarshal(s._rawfields["name"], &s.Name)

	// auto create the wrapper
	if fn, ok := _registry[s.Type]; ok {
		s.Provider, err = fn(s._raw)
		return err
	}

	return errors.New("source type not supported: " + s.Type)
}
