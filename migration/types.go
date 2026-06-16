package migration

import "time"

type Type string

const (
	TypeTable Type = "table"
	TypeSP    Type = "sp"
	TypeSeed  Type = "seed"
)

func (t Type) Order() int {
	switch t {
	case TypeTable:
		return 1
	case TypeSP:
		return 2
	case TypeSeed:
		return 3
	default:
		return 9
	}
}

func (t Type) Label() string {
	switch t {
	case TypeTable:
		return "tabla"
	case TypeSP:
		return "stored procedure"
	case TypeSeed:
		return "seed"
	default:
		return "desconocido"
	}
}

func (t Type) IsValid() bool {
	switch t {
	case TypeTable, TypeSP, TypeSeed:
		return true
	default:
		return false
	}
}

type Migration struct {
	Module   string
	Version  int64
	Name     string
	Type     Type
	Filename string
	Content  []byte
}

type RunResult struct {
	Migration Migration
	Applied   bool
	Skipped   bool
	Error     error
	Duration  time.Duration
}

type StatusResult struct {
	Module    string
	Version   int64
	Filename  string
	Type      string
	AppliedAt *time.Time
	Pending   bool
}

type ValidateResult struct {
	Module   string
	Filename string
	Valid    bool
	Error    string
}
