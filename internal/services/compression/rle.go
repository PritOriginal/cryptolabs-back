package compression

import (
	"bytes"
	"fmt"
	"strconv"
)

type RLE interface {
	Compress(data []byte) ([]byte, error)
	CompressWithDetails(data []byte) (CompressionDetails, error)
	Decompress(compressedData []byte) ([]byte, error)
}

type RLEService struct {
}

type RLEDetails struct {
	CompressionRatio float32 `json:"compression_ratio"`
	Size             int     `json:"size"`
}

func NewRLEService() *RLEService {
	return &RLEService{}
}

func (s *RLEService) Compress(data []byte) ([]byte, error) {
	compressedData := make([]byte, 0)

	add := func(counter int, char byte) {
		compressedData = append(compressedData, []byte(strconv.Itoa(counter))...)
		compressedData = append(compressedData, char)
	}

	char := data[0]
	counter := 1
	for _, c := range data[1:] {
		if c == char {
			counter++
		} else {
			add(counter, char)
			char = c
			counter = 1
		}
	}
	add(counter, char)

	return compressedData, nil
}

func (s *RLEService) CompressWithDetails(data []byte) (CompressionDetails, error) {
	compressedData, err := s.Compress(data)
	if err != nil {
		return CompressionDetails{}, err
	}

	rleDetalis := CompressionDetails{
		Data: compressedData,
		Details: RLEDetails{
			CompressionRatio: 1 - float32(len(compressedData))/float32(len(data)),
			Size:             len(compressedData),
		},
	}
	return rleDetalis, nil
}

func (s *RLEService) Decompress(compressedData []byte) ([]byte, error) {
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

func (s *RLEService) DecompressWithDetails(compressedData []byte) (CompressionDetails, error) {
	data, err := s.Decompress(compressedData)
	if err != nil {
		return CompressionDetails{}, err
	}

	rleDetalis := CompressionDetails{
		Data: data,
		Details: RLEDetails{
			CompressionRatio: 1 - float32(len(compressedData))/float32(len(data)),
			Size:             len(compressedData),
		},
	}
	return rleDetalis, nil
}
