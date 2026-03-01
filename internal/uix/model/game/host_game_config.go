package model

type FieldType int

const (
	FieldFlag FieldType = iota
	FieldString
	FieldInt
	FieldFloat
	FieldEnum
)

type HostConfigField struct {
	Name        string
	Description string
	Argument    string
	Required    bool
	Enabled     bool
	Type        FieldType
	Value       string
	EnumValues  map[string]string
	MinInt      *int64
	MaxInt      *int64
	MinFloat    *float64
	MaxFloat    *float64
}
