package compression

import (
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"unicode/utf8"

	"github.com/PritOriginal/cryptolabs-back/pkg/bitsio"
)

type LZW interface {
	Compress(data []byte) ([]byte, error)
	CompressWithDetails(data []byte) (CompressionDetails, error)
	Decompress(compressedData []byte) ([]byte, error)
}

type LZWService struct {
}

type LZWData struct {
	data              []byte
	dictionary        map[string]int
	reverseDictionary map[int]string
}

type LZWDetails struct {
	Dictionary       []LZWDictionaryItem `json:"dictionary"`
	CompressionRatio float32             `json:"compression_ratio"`
	Size             int                 `json:"size"`
}

type LZWDictionaryItem struct {
	Val string `json:"value"`
	Num int    `json:"number"`
}

func NewLZWService() *LZWService {
	return &LZWService{}
}

func (l *LZWService) Compress(data []byte) ([]byte, error) {
	compressedData := l.compressData(data)
	return compressedData.data, nil
}

func (l *LZWService) CompressWithDetails(data []byte) (CompressionDetails, error) {
	lzwData := l.compressData(data)

	dictionaryList := l.dictionaryToList(lzwData.dictionary)

	lzwDetails := CompressionDetails{
		Data: lzwData.data,
		Details: LZWDetails{
			Dictionary:       dictionaryList,
			CompressionRatio: 1 - float32(len(lzwData.data))/float32(len(data)),
			Size:             len(lzwData.data),
		},
	}
	return lzwDetails, nil
}

func (l *LZWService) compressData(data []byte) LZWData {
	dataStr := string(data)
	dictionary := l.makeDictionary()
	compressedData := l.compress(dataStr, dictionary)

	return LZWData{
		data:       compressedData,
		dictionary: dictionary,
	}
}

func (l *LZWService) makeDictionary() map[string]int {
	dictionary := make(map[string]int)
	for i := range 256 {
		dictionary[string(rune(i))] = i
	}
	ruAlphabet := "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	for i, ch := range []rune(ruAlphabet) {
		dictionary[string(ch)] = 256 + i
	}
	return dictionary
}

func (l *LZWService) compress(dataStr string, dictionary map[string]int) []byte {
	bitWriter := bitsio.NewBitWriter()
	sizeBit := 9

	writeCode := func(s string) {
		codeStr := fmt.Sprintf("%b", dictionary[s])
		for range sizeBit - len(codeStr) {
			bitWriter.WriteBit(false)
		}
		for _, code_ch := range codeStr {
			if code_ch == '1' {
				bitWriter.WriteBit(true)
			} else {
				bitWriter.WriteBit(false)
			}
		}
	}

	s := ""
	for _, ch := range dataStr {
		newStr := s + string(ch)
		if _, exsist := dictionary[newStr]; exsist {
			s += string(ch)
		} else {
			writeCode(s)
			dictionary[newStr] = len(dictionary)
			if len(dictionary) > int(math.Pow(2, float64(sizeBit))) {
				sizeBit++
			}

			s = string(ch)
		}
	}
	writeCode(s)

	return bitWriter.Bytes()
}

func (l *LZWService) writeCode(bitWriter *bitsio.BitWriter, sizeBit, code int) {
	codeStr := fmt.Sprintf("%b", code)
	for range sizeBit - len(codeStr) {
		bitWriter.WriteBit(false)
	}
	for _, code_ch := range codeStr {
		if code_ch == '1' {
			bitWriter.WriteBit(true)
		} else {
			bitWriter.WriteBit(false)
		}
	}
}

func (l *LZWService) dictionaryToList(dictionary map[string]int) []LZWDictionaryItem {
	list := make([]LZWDictionaryItem, 0, len(dictionary))
	for val, num := range dictionary {
		item := LZWDictionaryItem{
			Val: val,
			Num: num,
		}
		list = append(list, item)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Num < list[j].Num
	})
	return list
}

func (l *LZWService) reverseDictionaryToList(dictionary map[int]string) []LZWDictionaryItem {
	list := make([]LZWDictionaryItem, 0, len(dictionary))
	for num := range len(dictionary) {
		item := LZWDictionaryItem{
			Val: dictionary[num],
			Num: num,
		}
		list = append(list, item)
	}
	return list
}

func (l *LZWService) Decompress(compressedData []byte) ([]byte, error) {
	lzwData, err := l.decompressData(compressedData)
	if err != nil {
		return nil, err
	}
	return lzwData.data, nil
}

func (l *LZWService) DecompressWithDetails(compressedData []byte) (CompressionDetails, error) {
	lzwData, err := l.decompressData(compressedData)
	if err != nil {
		return CompressionDetails{}, err
	}

	dictionaryList := l.reverseDictionaryToList(lzwData.reverseDictionary)

	lzwDetails := CompressionDetails{
		Data: lzwData.data,
		Details: LZWDetails{
			Dictionary:       dictionaryList,
			CompressionRatio: 1 - float32(len(compressedData))/float32(len(lzwData.data)),
			Size:             len(lzwData.data),
		},
	}
	return lzwDetails, nil
}

func (l *LZWService) decompressData(compressedData []byte) (LZWData, error) {
	dictionary := l.makeReverseDictionary()
	data, err := l.decompress(compressedData, dictionary)
	if err != nil {
		return LZWData{}, err
	}

	lzwData := LZWData{
		data:              data,
		reverseDictionary: dictionary,
	}
	return lzwData, nil
}

func (l *LZWService) makeReverseDictionary() map[int]string {
	dictionary := make(map[int]string)
	for i := range 256 {
		dictionary[i] = string(rune(i))
	}
	ruAlphabet := "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	for i, ch := range []rune(ruAlphabet) {
		dictionary[256+i] = string(ch)
	}
	return dictionary
}

func (l *LZWService) decompress(compressedData []byte, dictionary map[int]string) ([]byte, error) {
	bitReader := bitsio.NewBitReader(compressedData)
	sizeBit := 9

	prevcode, _ := l.readCode(bitReader, sizeBit)
	s := dictionary[prevcode]

	c, _ := utf8.DecodeRuneInString(s)
	// c := s[0]

	data := make([]byte, 0)
	data = append(data, []byte(s)...)
	// dataStr := s
	// sizeBit++
	for !bitReader.IsEmpty() {
		code, err := l.readCode(bitReader, sizeBit)
		if err != nil {
			break
		}

		if _, exist := dictionary[code]; !exist {
			s = dictionary[prevcode]
			s = s + string(c)
		} else {
			s = dictionary[code]
		}

		data = append(data, []byte(s)...)
		// dataStr += s
		c, _ = utf8.DecodeRuneInString(s)
		dictionary[len(dictionary)] = dictionary[prevcode] + string(c)
		if len(dictionary) >= int(math.Pow(2, float64(sizeBit))) {
			sizeBit++
		}

		prevcode = code
	}

	return data, nil
	// return []byte(dataStr), nil
}

func (l *LZWService) readCode(bitReader *bitsio.BitReader, sizeBit int) (int, error) {
	bitWriter := bitsio.NewBitWriter()

	for range (64 - sizeBit) / 8 {
		bitWriter.WriteByte(0x0)
	}

	numZeroBits := 8 - sizeBit%8
	if numZeroBits < 8 {
		for range numZeroBits {
			bitWriter.WriteBit(false)
		}
	}
	for range sizeBit {
		if bitReader.IsEmpty() {
			return 0, fmt.Errorf("err")
		}
		bitWriter.WriteBit(bitReader.ReadBit())
	}

	code := int(binary.BigEndian.Uint64(bitWriter.Bytes()))
	return code, nil
}
