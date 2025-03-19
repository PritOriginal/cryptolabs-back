package services

import (
	"bytes"
	"fmt"
	"strconv"
)

type Compression interface {
	Compress(data []byte) []byte
	Decompress(compressedData []byte) ([]byte, error)
}

type CompressionService struct {
}

func NewCompressionService() *CompressionService {
	return &CompressionService{}
}

func (s *CompressionService) Compress(data []byte) []byte {
	compressedData := make([]byte, 0)

	char := data[0]
	counter := 1
	for _, c := range data[1:] {
		if c == char {
			counter++
		} else {
			compressedData = append(compressedData, []byte(strconv.Itoa(counter))...)
			compressedData = append(compressedData, char)
			char = c
			counter = 1
		}
	}
	compressedData = append(compressedData, []byte(strconv.Itoa(counter))...)
	compressedData = append(compressedData, char)

	return compressedData
}

func (s *CompressionService) Decompress(compressedData []byte) ([]byte, error) {
	data := make([]byte, 0)

	const (
		GET_COUNTER = iota
		GET_BYTE
	)

	counter := 1
	stepDecompress := 0
	buffer := bytes.NewBufferString("")
	for _, c := range compressedData {
		if _, err := strconv.Atoi(string(c)); err == nil {
			if stepDecompress == GET_BYTE {
				stepDecompress = GET_COUNTER
			}
			buffer.WriteByte(c)
		} else if stepDecompress == GET_COUNTER {
			stepDecompress = GET_BYTE
			counter, err = strconv.Atoi(buffer.String())
			if err != nil {
				return data, err
			}
			buffer.Reset()

			for i := 0; i < counter; i++ {
				data = append(data, c)
			}
		} else {
			return data, fmt.Errorf("invalid data")
		}
	}
	if stepDecompress == GET_COUNTER {
		return data, fmt.Errorf("invalid data")
	}

	return data, nil
}
