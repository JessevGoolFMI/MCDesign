package DispenserStrategies

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

type DispenserV1 struct{}

func (d *DispenserV1) Execute(value interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return encoded, nil
}

type DispenserV2 struct{}

func (d *DispenserV2) Execute(value interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:]), nil
}

type DispenserV3 struct{}

func (d *DispenserV3) Execute(value interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)
	hash := sha512.Sum512([]byte(str))
	return hex.EncodeToString(hash[:]), nil
}

type DispenserV4 struct{}

func (d *DispenserV4) Execute(value interface{}) (interface{}, error) {
	str := fmt.Sprintf("%v", value)
	return strings.ToUpper(str), nil
}
