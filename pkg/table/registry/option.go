package registry

import (
	"sync"

	"framework/pkg/cap/msg/errors"
	cap "framework/pkg/table/proto"
)

type optionListStore struct {
	id string
	// ordered list
	list []*cap.OptionValue
	kv   map[int32]*cap.OptionValue
}

func (ols *optionListStore) lookup(id int32) (*cap.OptionValue, error) {
	if v, ok := ols.kv[id]; ok {
		return v, nil
	}
	return nil, errors.Wrap(ErrOptionNotFound).FillDebugArgs(id, ols.id)
}

// OptionReg option registries
type OptionReg struct {
	m sync.Map
}

// Lookup value for optionID in optionTypeID
func (o *OptionReg) Lookup(optionTypeID string, optionID int32) (*cap.OptionValue, error) {
	if ols, ok := o.m.Load(optionTypeID); ok {
		return ols.(*optionListStore).lookup(optionID)
	}
	return nil, errors.Wrap(ErrOptionTypeNotFound).FillDebugArgs(optionTypeID)
}

// GetOptions get option list by type id
func (o *OptionReg) GetOptions(optionTypeID string) ([]*cap.OptionValue, error) {
	if ols, ok := o.m.Load(optionTypeID); ok {
		return ols.(*optionListStore).list, nil
	}
	return nil, errors.Wrap(ErrOptionTypeNotFound).FillDebugArgs(optionTypeID)
}

// Register ...
func (o *OptionReg) Register(optionTypeID string, values []*cap.OptionValue) error {
	if _, ok := o.m.Load(optionTypeID); ok {
		return errors.Wrap(ErrDupplicateNodeID).FillDebugArgs(optionTypeID)
	}
	ols := &optionListStore{
		id:   optionTypeID,
		list: values,
		kv:   map[int32]*cap.OptionValue{},
	}
	for _, v := range values {
		ols.kv[v.Id] = v
	}
	o.m.Store(optionTypeID, ols)
	return nil
}

// Store ...
func (o *OptionReg) Store() *sync.Map {
	return &o.m
}

// NewOptionReg creates option registry
func NewOptionReg() *OptionReg {
	return &OptionReg{}
}
