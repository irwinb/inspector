package models

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
)

type OperationsList list.List

type Endpoint struct {
	Id         uint           `json:"id"`
	Target     *string        `json:"target"`
	Name       *string        `json:"name"`
	Operations OperationsList `json:"operations"`
	Protocol   *string        `json:"protocol"`
}

var validProtocols = []string{"http", "https"}

func (l *OperationsList) PushBack(op *Operation) {
	list := (*list.List)(l)
	list.PushBack(op)
	if list.Len() > 10 {
		list.Remove(list.Front())
	}
}

func (l *OperationsList) MarshalJSON() ([]byte, error) {
	list := (*list.List)(l)
	var buf bytes.Buffer
	buf.WriteByte('[')

	for e := list.Front(); e != nil; e = e.Next() {
		bytes, err := json.Marshal(e.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(bytes)
		if e.Next() != nil {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func (l *OperationsList) UnmarshalJSON(data []byte) error {
	opSlice := make([]*Operation, 0, 0)
	err := json.Unmarshal(data, &opSlice)
	if err != nil {
		return err
	}

	for _, op := range opSlice {
		if op == nil {
			continue
		}
		l.PushBack(op)
	}
	return nil
}

func (e *Endpoint) Validate() error {
	if e.Target == nil || len(*e.Target) == 0 {
		return errors.New("The endpoint's target cannot be empty.")
	}
	if len(*e.Target) > 100 {
		return errors.New("The endpoint's target is too long.  < 100 plz.")
	}

	if e.Name == nil || len(*e.Name) == 0 {
		return errors.New("The endpoint's name cannot be empty.")
	}
	if len(*e.Name) > 100 {
		return errors.New("The endpoint's name is too long.  < 100 plz.")
	}
	validProtocol := false
	for _, protocol := range validProtocols {
		if protocol == *e.Protocol {
			validProtocol = true
			break
		}
	}
	if !validProtocol {
		return errors.New("Invalid protocol [" + *e.Protocol + "].")
	}

	return nil
}
