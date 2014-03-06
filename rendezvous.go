package rendezvous

import (
	"hash"
	"hash/crc32"
	"sort"
)

var (
	crc32Table = crc32.MakeTable(crc32.Castagnoli)
)

type nodeScore struct {
	node  []byte
	score uint32
}

type byScore []nodeScore

func (s byScore) Len() int           { return len(s) }
func (s byScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byScore) Less(i, j int) bool { return s[i].score < s[j].score }

type Hash struct {
	nodes  []nodeScore
	hasher hash.Hash32
}

func New(nodes ...string) *Hash {
	hash := &Hash{}
	hash.hasher = crc32.New(crc32Table)
	hash.Add(nodes...)
	return hash
}

func (h *Hash) Add(nodes ...string) {
	for _, node := range nodes {
		h.nodes = append(h.nodes, nodeScore{[]byte(node), 0})
	}
}

func (h *Hash) Get(key string) string {
	keyBytes := []byte(key)

	var maxScore uint32
	var maxNode []byte
	var score uint32

	for _, node := range h.nodes {
		score = h.hash(node.node, keyBytes)
		if score > maxScore {
			maxScore = score
			maxNode = node.node
		}
	}

	return string(maxNode)
}

func (h *Hash) GetN(n int, key string) []string {
	if len(h.nodes) == 0 || n == 0 {
		return []string{}
	}

	if n > len(h.nodes) {
		n = len(h.nodes)
	}

	keyBytes := []byte(key)
	for i := 0; i < len(h.nodes); i++ {
		h.nodes[i].score = h.hash(h.nodes[i].node, keyBytes)
	}
	sort.Sort(sort.Reverse(byScore(h.nodes)))

	nodes := make([]string, n)
	for i := 0; i < n; i++ {
		nodes[i] = string(h.nodes[i].node)
	}
	return nodes
}

func (h *Hash) hash(node, key []byte) uint32 {
	h.hasher.Reset()
	h.hasher.Write(key)
	h.hasher.Write(node)
	return h.hasher.Sum32()
}
