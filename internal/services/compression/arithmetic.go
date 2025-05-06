package compression

import (
	"bytes"
	"encoding/binary"
	"maps"
	"math"
	"math/big"
	"slices"
	"sort"
	"unicode/utf8"
)

type Arithmetic interface {
	Compress([]byte) ([]byte, error)
	CompressWithDetails(data []byte) (CompressionDetails, error)
	Decompress(compressedData []byte) ([]byte, error)
}

type ArithmeticData struct {
	data           []byte
	frequencyTable FrequencyTable
}

type ArithmeticDetails struct {
	FrequencyTable   []FrequencyTableItem `json:"frequency_table"`
	CompressionRatio float32              `json:"compression_ratio"`
	Size             int                  `json:"size"`
}

type FrequencyTableItem struct {
	Val       string `json:"value"`
	Frequency int    `json:"frequency"`
}

type Interval struct {
	val   rune
	left  *big.Float
	right *big.Float
}

type ArithmeticService struct {
}

func NewArithmeticService() *ArithmeticService {
	return &ArithmeticService{}
}

type FrequencyTable map[rune]uint16

func (a *ArithmeticService) Compress(data []byte) ([]byte, error) {
	arithmeticData, err := a.compressData(data)
	if err != nil {
		return nil, err
	}
	return arithmeticData.data, nil
}

func (a *ArithmeticService) CompressWithDetails(data []byte) (CompressionDetails, error) {
	arithmeticData, err := a.compressData(data)
	if err != nil {
		return CompressionDetails{}, err
	}

	frequencyList := a.frequencyTableToList(arithmeticData.frequencyTable)

	arithmeticDetails := CompressionDetails{
		Data: arithmeticData.data,
		Details: ArithmeticDetails{
			FrequencyTable:   frequencyList,
			CompressionRatio: 1 - float32(len(arithmeticData.data))/float32(len(data)),
			Size:             len(arithmeticData.data),
		},
	}

	return arithmeticDetails, nil
}

func (a *ArithmeticService) compressData(data []byte) (ArithmeticData, error) {
	dataStr := string(data)
	dataLength := uint32(utf8.RuneCountInString(dataStr))

	var precision uint = calcPrecision(dataLength)

	frequencyTable := a.frequencyTable(dataStr)
	probabilityIntervals := a.probabilityIntervals(precision, dataLength, frequencyTable)
	n := a.compress(precision, dataStr, probabilityIntervals)

	dataPayload, err := n.GobEncode()
	if err != nil {
		return ArithmeticData{}, err
	}

	compressedData, err := a.allCompressedData(dataLength, frequencyTable, dataPayload)
	if err != nil {
		return ArithmeticData{}, err
	}

	arithmeticData := ArithmeticData{
		data:           compressedData,
		frequencyTable: frequencyTable,
	}
	return arithmeticData, nil
}

func calcPrecision(length uint32) uint {
	return uint(math.Round(float64(length)*math.Log2(10)) * 1.42)
}

func (a *ArithmeticService) frequencyTable(dataStr string) FrequencyTable {
	frequencyTable := make(FrequencyTable)
	for _, ch := range dataStr {
		frequencyTable[ch] += 1
	}

	return frequencyTable
}

func (a *ArithmeticService) frequencyTableToBinary(frequencyTable FrequencyTable) ([]byte, error) {
	buf := new(bytes.Buffer)
	for ch, frequency := range frequencyTable {
		if _, err := buf.WriteRune(ch); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, frequency); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (a *ArithmeticService) frequencyTableToList(frequencyTable FrequencyTable) []FrequencyTableItem {
	list := make([]FrequencyTableItem, 0, len(frequencyTable))
	for val, frequency := range frequencyTable {
		item := FrequencyTableItem{
			Val:       string(val),
			Frequency: int(frequency),
		}
		list = append(list, item)
	}
	return list
}

func (a *ArithmeticService) allCompressedData(dataLength uint32, frequencyTable FrequencyTable, dataPayload []byte) ([]byte, error) {
	frequencyTableSize := uint16(len(frequencyTable))
	binaryFrequencyTable, err := a.frequencyTableToBinary(frequencyTable)
	if err != nil {
		return nil, err
	}

	compressedData := new(bytes.Buffer)
	binary.Write(compressedData, binary.LittleEndian, dataLength)
	binary.Write(compressedData, binary.LittleEndian, frequencyTableSize)
	compressedData.Write(binaryFrequencyTable)
	compressedData.Write(dataPayload)

	return compressedData.Bytes(), nil
}

func (a *ArithmeticService) probabilityIntervals(precision uint, dataLength uint32, frequencyTable FrequencyTable) map[rune]Interval {
	probabilityIntervals := make(map[rune]Interval)

	keys := make([]rune, 0, len(probabilityIntervals))
	for k, _ := range frequencyTable {
		keys = append(keys, k)
	}

	sumFloat := new(big.Float).SetInt64(int64(dataLength))

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// var left *big.Rat = new(big.Rat).SetInt64(0)
	var left *big.Float = big.NewFloat(0).SetPrec(precision)
	for _, r := range keys {
		probability := new(big.Float).Quo(new(big.Float).SetInt64(int64(frequencyTable[r])), sumFloat)
		right := new(big.Float).Add(left, probability)

		probabilityIntervals[r] = Interval{
			val:   r,
			left:  left,
			right: right,
		}

		left = right

		// fmt.Printf("%v: %v - %v\n", string(r), probabilityIntervals[r].left, probabilityIntervals[r].right)
	}

	return probabilityIntervals
}

func (a *ArithmeticService) compress(precision uint, dataStr string, probabilityIntervals map[rune]Interval) *big.Float {
	var left, right *big.Float = big.NewFloat(0).SetPrec(precision), big.NewFloat(1).SetPrec(precision)

	for _, r := range dataStr {
		interval := new(big.Float).Sub(right, left)

		newRight := new(big.Float).Mul(interval, probabilityIntervals[r].right)
		newRight = new(big.Float).Add(left, newRight)

		newLeft := new(big.Float).Mul(interval, probabilityIntervals[r].left)
		newLeft = new(big.Float).Add(left, newLeft)

		right = newRight
		left = newLeft
	}

	// return (left + right) / 2
	return new(big.Float).Quo(new(big.Float).Add(left, right), new(big.Float).SetInt64(2))
}

func (a *ArithmeticService) Decompress(compressedData []byte) ([]byte, error) {
	arithmeticData, err := a.decompressData(compressedData)
	if err != nil {
		return nil, err
	}
	return arithmeticData.data, nil
}

func (a *ArithmeticService) DecompressWithDetails(compressedData []byte) (CompressionDetails, error) {
	arithmeticData, err := a.decompressData(compressedData)
	if err != nil {
		return CompressionDetails{}, err
	}

	frequencyList := a.frequencyTableToList(arithmeticData.frequencyTable)

	arithmeticDetails := CompressionDetails{
		Data: arithmeticData.data,
		Details: ArithmeticDetails{
			FrequencyTable:   frequencyList,
			CompressionRatio: 1 - float32(len(compressedData))/float32(len(arithmeticData.data)),
			Size:             len(arithmeticData.data),
		},
	}
	return arithmeticDetails, nil
}

func (a *ArithmeticService) decompressData(compressedData []byte) (ArithmeticData, error) {
	buf := bytes.NewBuffer(compressedData)
	var dataLength uint32
	err := binary.Read(buf, binary.LittleEndian, &dataLength)
	if err != nil {
		return ArithmeticData{}, err
	}

	var precision uint = calcPrecision(dataLength)

	var frequencyTableSize uint16
	err = binary.Read(buf, binary.LittleEndian, &frequencyTableSize)
	if err != nil {
		return ArithmeticData{}, err
	}

	frequencyTable, err := a.binaryToFrequencyTable(buf, frequencyTableSize)
	if err != nil {
		return ArithmeticData{}, err
	}
	probabilityIntervals := a.probabilityIntervals(precision, dataLength, frequencyTable)

	decode := new(big.Float)
	err = decode.GobDecode(buf.Bytes())
	if err != nil {
		return ArithmeticData{}, err
	}

	data := a.decompress(decode, int(dataLength), probabilityIntervals)

	arithmeticData := ArithmeticData{
		data:           data,
		frequencyTable: frequencyTable,
	}
	return arithmeticData, nil
}

func (a *ArithmeticService) binaryToFrequencyTable(buf *bytes.Buffer, size uint16) (FrequencyTable, error) {
	frequencyTable := make(FrequencyTable)
	for range size {
		r, _, err := buf.ReadRune()
		if err != nil {
			return nil, err
		}
		var frequency uint16
		if err := binary.Read(buf, binary.LittleEndian, &frequency); err != nil {
			return nil, err
		}

		frequencyTable[r] = frequency
	}
	return frequencyTable, nil
}

func (a *ArithmeticService) decompress(n *big.Float, dataLength int, probabilityIntervals map[rune]Interval) []byte {
	probabilityIntervalsSorted := slices.Collect(maps.Values(probabilityIntervals))
	sort.Slice(probabilityIntervalsSorted, func(i, j int) bool {
		return probabilityIntervalsSorted[i].left.Cmp(probabilityIntervalsSorted[j].left) < 0
	})

	buf := bytes.Buffer{}

	for range dataLength {
		interval := a.findInterval(probabilityIntervalsSorted, n)
		buf.WriteRune(interval.val)

		n = new(big.Float).Sub(n, interval.left)
		n = new(big.Float).Quo(n, new(big.Float).Sub(interval.right, interval.left))
	}
	return buf.Bytes()
}

func (a *ArithmeticService) findInterval(probabilityIntervals []Interval, n *big.Float) Interval {
	leftPointer := 0
	rightPointer := len(probabilityIntervals) - 1
	for leftPointer <= rightPointer {
		midPointer := int((leftPointer + rightPointer) / 2)
		midInterval := probabilityIntervals[midPointer]

		moreThanLeft := n.Cmp(midInterval.left) >= 0
		if moreThanLeft && n.Cmp(midInterval.right) < 0 {
			return midInterval
		} else if moreThanLeft {
			leftPointer = midPointer + 1
		} else {
			rightPointer = midPointer - 1
		}
	}
	return probabilityIntervals[leftPointer]
}
