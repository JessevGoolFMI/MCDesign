package CompressorStrategies

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
)

type CompressorV1 struct{}

func (c *CompressorV1) Execute(value interface{}) (interface{}, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		fmt.Println("Error encoding data:", err)
		return nil, errors.New("error encoding data")
	}

	compressedData := buf.Bytes()
	return compressedData, nil
}

// CompressorV2 represents the second version of the compressor.
type CompressorV2 struct{}

func (c *CompressorV2) Execute(value interface{}) (interface{}, error) {
	compressedData, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return nil, errors.New("error encoding data")
	}
	return compressedData, nil
}

type CompressorV3 struct{}

func (c *CompressorV3) Execute(value interface{}) (interface{}, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	defer gzipWriter.Close()

	data, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return nil, errors.New("error encoding data")
	}

	if _, err := gzipWriter.Write(data); err != nil {
		fmt.Println("Error compressing data:", err)
		return nil, errors.New("error compressing data")
	}

	return buf.Bytes(), nil
}

type CompressorV4 struct{}

func (c *CompressorV4) Execute(value interface{}) (interface{}, error) {
	// Convert the interface to a string
	str, ok := value.(string)
	if !ok {
		// If the value is not a string, attempt to convert it using fmt.Sprintf
		str = fmt.Sprintf("%v", value)
	}

	// Reverse the string
	reversed := reverseString(str)
	return []byte(reversed), nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
