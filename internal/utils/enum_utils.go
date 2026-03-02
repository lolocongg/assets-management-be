package utils

import (
	"encoding/json"
)

type Enum interface {
	String() string
}

func MarshalEnum(e Enum) ([]byte, error) {
	return json.Marshal(e.String())
}
