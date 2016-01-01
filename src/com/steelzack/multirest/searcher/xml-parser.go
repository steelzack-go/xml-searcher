package searcher

import (
	"github.com/kokardy/saxlike"
	"encoding/xml"
	"github.com/gocql/gocql"
	"strings"
)

type NodeType int

const (
	Value1UUID NodeType = iota
	Value2UUID NodeType = iota
)

//VoidHandler is a implemented Handler that do nothing.
type PartialHandler struct {
	currentnode   NodeType

	value1 gocql.UUID
	value2 gocql.UUID

	value1Found bool
	value2Found bool

	saxlike.VoidHandler
}

func (h *PartialHandler) StartElement(element xml.StartElement) {
	if (strings.Contains(element.Name.Local, "value1") && len(element.Name.Local) == 5 && !(*h).value1Found) {
		(*h).currentnode = Value1UUID
	} else if (strings.Contains(element.Name.Local, "value2") && !(*h).value2Found) {
		(*h).currentnode = Value2UUID
	} else {
		(*h).currentnode = -1
	}
}

func (h *PartialHandler) EndElement(element xml.EndElement) {
	(*h).currentnode = -1
}

func (h* PartialHandler) CharData(char xml.CharData) {
	nodevalue := string(char)
	switch (*h).currentnode {
	case Value1UUID :
		(*h).value1, _ = gocql.ParseUUID(nodevalue)
		(*h).value1Found = true
	case Value2UUID :
		(*h).value2, _ = gocql.ParseUUID(nodevalue)
		(*h).value2Found = true
	}
}
