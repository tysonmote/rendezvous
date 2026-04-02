// Package rendezvous implements rendezvous hashing (a.k.a. highest random
// weight hashing). See http://en.wikipedia.org/wiki/Rendezvous_hashing for
// more information.
//
// Hash is safe for concurrent use: multiple goroutines may call Get, GetN, Add,
// and Remove on the same instance.
package rendezvous

import (
	"bytes"
	"hash"
	"hash/crc32"
	"sort"
	"sync"
	"unsafe"
)

var crc32Table = crc32.MakeTable(crc32.Castagnoli)

var hasherPool = sync.Pool{
	New: func() any {
		return crc32.New(crc32Table)
	},
}

type Hash struct {
	mu    sync.RWMutex
	nodes nodeScores
}

type nodeScore struct {
	node  []byte
	score uint32
}

// New returns a new Hash ready for use with the given nodes.
func New(nodes ...string) *Hash {
	h := &Hash{}
	h.Add(nodes...)
	return h
}

// Add adds additional nodes to the Hash.
func (h *Hash) Add(nodes ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, node := range nodes {
		h.nodes = append(h.nodes, nodeScore{node: []byte(node)})
	}
}

// Get returns the node with the highest score for the given key. If this Hash
// has no nodes, an empty string is returned.
func (h *Hash) Get(key string) string {
	var maxScore uint32
	var maxNode []byte

	keyBytes := unsafeBytes(key)

	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, node := range h.nodes {
		score := h.hash(node.node, keyBytes)
		if score > maxScore || (score == maxScore && bytes.Compare(node.node, maxNode) < 0) {
			maxScore = score
			maxNode = node.node
		}
	}

	return string(maxNode)
}

// GetN returns no more than n nodes for the given key, ordered by descending
// score. It does not reorder the internal node list.
func (h *Hash) GetN(n int, key string) []string {
	if n < 0 {
		n = 0
	}
	keyBytes := unsafeBytes(key)
	h.mu.RLock()
	scored := make(nodeScores, len(h.nodes))
	copy(scored, h.nodes)
	h.mu.RUnlock()
	for i := range scored {
		scored[i].score = h.hash(scored[i].node, keyBytes)
	}
	sort.Sort(&scored)

	if n > len(scored) {
		n = len(scored)
	}

	nodes := make([]string, n)
	for i := range nodes {
		nodes[i] = string(scored[i].node)
	}
	return nodes
}

// Remove removes a node from the Hash, if it exists
func (h *Hash) Remove(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	nodeBytes := []byte(node)
	for i, ns := range h.nodes {
		if bytes.Equal(ns.node, nodeBytes) {
			h.nodes = append(h.nodes[:i], h.nodes[i+1:]...)
			return
		}
	}
}

type nodeScores []nodeScore

func (s *nodeScores) Len() int      { return len(*s) }
func (s *nodeScores) Swap(i, j int) { (*s)[i], (*s)[j] = (*s)[j], (*s)[i] }
func (s *nodeScores) Less(i, j int) bool {
	if (*s)[i].score == (*s)[j].score {
		return bytes.Compare((*s)[i].node, (*s)[j].node) < 0
	}
	return (*s)[j].score < (*s)[i].score // Descending
}

func (h *Hash) hash(node, key []byte) uint32 {
	digest := hasherPool.Get().(hash.Hash32)
	defer hasherPool.Put(digest)
	digest.Reset()
	digest.Write(key)
	digest.Write(node)
	return digest.Sum32()
}

// unsafeBytes maps a string to a byte slice without copying. Callers must not
// retain the returned slice after the callee returns: hash.Hash.Write must
// consume the bytes during the call and must not keep references to them.
func unsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
