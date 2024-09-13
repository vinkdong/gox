package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type NullInt64 struct {
	sql.NullInt64
}

func (n *NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

type NullInt32 struct {
	sql.NullInt32
}

func (n NullInt32) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int32)
	}
	return json.Marshal(nil)
}

type NullUInt64 struct {
	UInt64 uint64
	Valid  bool // Valid is true if Int64 is not NULL
}

func (n *NullUInt64) Scan(value interface{}) error {
	if value == nil {
		n.UInt64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	switch v := value.(type) {
	case uint64:
		n.UInt64 = v
	case int64:
		n.UInt64 = uint64(v)
	case *uint64:
		n.UInt64 = *v
	default:
		return fmt.Errorf("cannot scan type %T into NullUInt64: %v", value, value)
	}
	return nil
}

func (n NullUInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.UInt64)
	}
	return json.Marshal(nil)
}
