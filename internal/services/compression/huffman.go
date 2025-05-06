package compression

import (
	"bytes"
	"container/heap"
	"fmt"

	"github.com/PritOriginal/cryptolabs-back/pkg/bitsio"
)

// type Huffman interface {
// 	Compress(data []byte) ([]byte, error)
// 	CompressWithDetails(data []byte) CompressionDetails
// 	Decompress(compressedData []byte) ([]byte, error)
// }

type HuffmanService struct {
}

func NewHuffmanService() *HuffmanService {
	return &HuffmanService{}
}

type HuffmanData struct {
	data           []byte
	frequencyTable map[rune]int
	rootNode       *Node
	huffmanCode    map[rune]string
}

type HuffmanDetails struct {
	Codes            []HuffmanCode `json:"codes"`
	CompressionRatio float32       `json:"compression_ratio"`
	Size             int           `json:"size"`
}

type HuffmanCode struct {
	Val       string `json:"value"`
	Frequency int    `json:"frequency"`
	Code      string `json:"code"`
}

type CompressionDetails struct {
	Details interface{} `json:"details"`
	Data    []byte      `json:"data"`
}

type Item struct {
	value    Node
	priority int

	index int // Индекс элемента в куче.
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // избежать утечки памяти
	item.index = -1 // для безопасности
	*pq = old[0 : n-1]
	return item
}

// update изменяет приоритет и значение Item в очереди.
func (pq *PriorityQueue) update(item *Item, value Node, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

type Node struct {
	value rune
	left  *Node
	right *Node
}

func (h *HuffmanService) Compress(data []byte) ([]byte, error) {
	huffmanData := h.compressData(data)
	return huffmanData.data, nil
}

func (h *HuffmanService) CompressWithDetails(data []byte) (CompressionDetails, error) {
	huffmanData := h.compressData(data)

	codes := h.makeHuffmanCodeList(huffmanData.frequencyTable, huffmanData.huffmanCode)

	huffmanDetails := CompressionDetails{
		Data: huffmanData.data,
		Details: HuffmanDetails{
			Codes:            codes,
			CompressionRatio: 1 - float32(len(huffmanData.data))/float32(len(data)),
			Size:             len(huffmanData.data),
		},
	}
	return huffmanDetails, nil
}

func (h *HuffmanService) compressData(data []byte) HuffmanData {
	dataStr := string(data)

	frequencyTable := h.frequencyTable(dataStr)
	rootNode := h.buildTree(frequencyTable)
	huffmanCode := h.makeHuffmanCode(rootNode)

	dataPayload, numSkipBits := h.compress(dataStr, huffmanCode)
	compressedData := h.allCompressedData(numSkipBits, rootNode, dataPayload)

	return HuffmanData{
		data:           compressedData,
		frequencyTable: frequencyTable,
		rootNode:       &rootNode,
		huffmanCode:    huffmanCode,
	}
}

func (h *HuffmanService) makeHuffmanCodeList(frequencyTable map[rune]int, huffmanCode map[rune]string) []HuffmanCode {
	codes := make([]HuffmanCode, 0, len(huffmanCode))
	for ch := range huffmanCode {
		huffmanCodeItem := HuffmanCode{
			Val:       string(ch),
			Frequency: frequencyTable[ch],
			Code:      huffmanCode[ch],
		}
		codes = append(codes, huffmanCodeItem)
	}
	return codes
}

func (h *HuffmanService) frequencyTable(dataStr string) map[rune]int {
	frequencyTable := make(map[rune]int)
	for _, ch := range dataStr {
		frequencyTable[ch] += 1
	}

	return frequencyTable
}

func (h *HuffmanService) buildTree(frequencyTable map[rune]int) Node {
	pq := make(PriorityQueue, len(frequencyTable))
	i := 0
	for ch, frequency := range frequencyTable {
		pq[i] = &Item{
			value:    Node{value: ch},
			priority: frequency,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	for pq.Len() != 1 {
		left := heap.Pop(&pq).(*Item)
		right := heap.Pop(&pq).(*Item)

		sum := left.priority + right.priority
		newItem := &Item{
			value: Node{
				left:  &left.value,
				right: &right.value,
			},
			priority: sum,
		}
		heap.Push(&pq, newItem)
		pq.update(newItem, newItem.value, sum)
	}

	return heap.Pop(&pq).(*Item).value
}

func (h *HuffmanService) makeHuffmanCode(rootNode Node) map[rune]string {
	huffmanCode := make(map[rune]string, 0)
	type StackItem struct {
		node Node
		way  string
	}
	stack := make([]StackItem, 0)
	stack = append(stack, StackItem{
		node: rootNode,
		way:  "",
	})
	for len(stack) > 0 {
		n := len(stack) - 1
		item := stack[n]
		currentNode := item.node
		stack = stack[:n]

		if currentNode.left == nil && currentNode.right == nil {
			huffmanCode[currentNode.value] = item.way
		}

		if currentNode.right != nil {
			stack = append(stack, StackItem{
				node: *currentNode.right,
				way:  item.way + "1",
			})
		}
		if currentNode.left != nil {
			stack = append(stack, StackItem{
				node: *currentNode.left,
				way:  item.way + "0",
			})
		}
	}

	return huffmanCode
}

func (h *HuffmanService) compress(dataStr string, huffmanCode map[rune]string) ([]byte, byte) {
	bitWriter := bitsio.NewBitWriter()
	for _, ch := range dataStr {
		code := huffmanCode[ch]
		for _, code_ch := range code {
			if code_ch == '1' {
				bitWriter.WriteBit(true)
			} else {
				bitWriter.WriteBit(false)
			}
		}
	}

	numSkipBits := bitWriter.BitsLeftToByte()
	if numSkipBits == 8 {
		numSkipBits = 0
	}

	return bitWriter.Bytes(), numSkipBits
}

func (h *HuffmanService) allCompressedData(numSkipBits byte, rootNode Node, dataPayload []byte) []byte {
	binaryTree := h.tree2binary(rootNode)

	compressedData := make([]byte, 0)
	compressedData = append(compressedData, numSkipBits)
	compressedData = append(compressedData, binaryTree...)
	compressedData = append(compressedData, dataPayload...)
	return compressedData
}

func (h *HuffmanService) tree2binary(rootNode Node) []byte {
	bitWriter := bitsio.NewBitWriter()

	isFirst := true
	stack := make([]Node, 0)
	stack = append(stack, rootNode)
	for len(stack) > 0 {
		n := len(stack) - 1
		currentNode := stack[n]
		stack = stack[:n]

		if !isFirst {
			if currentNode.left == nil && currentNode.right == nil {
				bitWriter.WriteBit(true)
				bitWriter.WtiteRune(currentNode.value)
			} else {
				bitWriter.WriteBit(false)
			}
		} else {
			isFirst = false
		}

		if currentNode.right != nil {
			stack = append(stack, *currentNode.right)
		}
		if currentNode.left != nil {
			stack = append(stack, *currentNode.left)
		}
	}

	return bitWriter.Bytes()
}

func (h *HuffmanService) Decompress(compressedData []byte) ([]byte, error) {
	huffmanData, err := h.decompressData(compressedData)
	if err != nil {
		return nil, err
	}

	return huffmanData.data, nil
}

func (h *HuffmanService) DecompressWithDetails(compressedData []byte) (CompressionDetails, error) {
	huffmanData, err := h.decompressData(compressedData)
	if err != nil {
		return CompressionDetails{}, err
	}

	codes := h.makeHuffmanCodeList(huffmanData.frequencyTable, huffmanData.huffmanCode)

	huffmanDetails := CompressionDetails{
		Data: huffmanData.data,
		Details: HuffmanDetails{
			Codes:            codes,
			CompressionRatio: 1 - float32(len(compressedData))/float32(len(huffmanData.data)),
			Size:             len(huffmanData.data),
		},
	}
	return huffmanDetails, nil
}

func (h *HuffmanService) decompressData(compressedData []byte) (HuffmanData, error) {
	bitReader := bitsio.NewBitReader(compressedData)

	numSkipBits := h.numSkipBits(bitReader)
	rootNode := h.restoreTree(bitReader)
	payloadData := compressedData[bitReader.NumReadByte():]

	huffmanCode := h.makeHuffmanCode(*rootNode)

	for ch, code := range huffmanCode {
		fmt.Printf("%s:%s\n", string(ch), code)
	}

	data, err := h.decompress(rootNode, payloadData, numSkipBits)
	if err != nil {
		return HuffmanData{}, err
	}

	huffmanData := HuffmanData{
		data:        data,
		rootNode:    rootNode,
		huffmanCode: huffmanCode,
	}
	return huffmanData, nil
}

func (h *HuffmanService) numSkipBits(bitReader *bitsio.BitReader) byte {
	// bitWriter := bitsio.NewBitWriter()
	// for range 4 {
	// 	bitWriter.WriteBit(false)
	// }
	// for range 4
	// 	bitWriter.WriteBit(bitReader.ReadBit())
	// }
	// return bitWriter.Bytes()[0]
	return bitReader.ReadByte()
}

func (h *HuffmanService) restoreTree(bitReader *bitsio.BitReader) *Node {
	rootNode := &Node{}

	stack := make([]*Node, 0)
	stack = append(stack, rootNode)

	for !h.checkIsFull(stack) {
		bit := bitReader.ReadBit()
		var newNode *Node
		if bit {
			r := bitReader.ReadRune()
			newNode = &Node{value: r}

		} else {
			newNode = &Node{}
		}

		for len(stack) > 0 {
			n := len(stack) - 1
			currentNode := stack[n]
			if currentNode.left == nil {
				currentNode.left = newNode
				break
			} else if currentNode.right == nil {
				currentNode.right = newNode
				break
			} else {
				stack = stack[:n]
			}
		}
		if !bit {
			stack = append(stack, newNode)
		}
	}
	bitReader.FinishByte()

	return rootNode
}

func (h *HuffmanService) checkIsFull(stack []*Node) bool {
	for i := len(stack) - 1; i >= 0; i-- {
		currentNode := stack[i]
		if currentNode.left == nil || currentNode.right == nil {
			return false
		}
	}
	return true
}

func (h *HuffmanService) decompress(rootNode *Node, compressedData []byte, numSkipBits byte) ([]byte, error) {
	var data bytes.Buffer
	bitReader := bitsio.NewBitReader(compressedData)
	numLastBitsInLastByte := 8 - numSkipBits
	node := rootNode
	for !bitReader.IsLastByte() || (bitReader.IsLastByte() && numLastBitsInLastByte > 0) {
		if bitReader.IsLastByte() {
			numLastBitsInLastByte--
		}

		currentBit := bitReader.ReadBit()
		if currentBit {
			if node.right != nil {
				node = node.right
			}
		} else {
			if node.left != nil {
				node = node.left
			}
		}

		if node.left == nil && node.right == nil {
			_, err := data.WriteRune(node.value)
			if err != nil {
				return nil, err
			}
			node = rootNode
		}
	}

	return data.Bytes(), nil
}
