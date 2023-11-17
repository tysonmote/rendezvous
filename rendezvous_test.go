package rendezvous

import (
	"fmt"
	"reflect"
	"testing"
)

var sampleKeys = []string{
	"352DAB08-C1FD-4462-B573-7640B730B721",
	"382080D3-B847-4BB5-AEA8-644C3E56F4E1",
	"2B340C12-7958-4DBE-952C-67496E15D0C8",
	"BE05F82B-902E-4868-8CC9-EE50A6C64636",
	"C7ECC571-E924-4523-A313-951DFD5D8073",
}

type getTestcase struct {
	key          string
	expectedNode string
}

func TestHashGet(t *testing.T) {
	hash := New()

	gotNode := hash.Get("foo")
	if len(gotNode) != 0 {
		t.Errorf("got: %#v, expected: %#v", gotNode, "")
	}

	hash.Add("a", "b", "c", "d", "e")

	testcases := []getTestcase{
		{"", "d"},
		{"foo", "e"},
		{"bar", "c"},
	}

	for _, testcase := range testcases {
		gotNode := hash.Get(testcase.key)
		if gotNode != testcase.expectedNode {
			t.Errorf("got: %#v, expected: %#v", gotNode, testcase.expectedNode)
		}
	}
}

func BenchmarkHashGet_5nodes(b *testing.B) {
	hash := New("a", "b", "c", "d", "e")
	for i := 0; i < b.N; i++ {
		hash.Get(sampleKeys[i%len(sampleKeys)])
	}
}

func BenchmarkHashGet_10nodes(b *testing.B) {
	hash := New("a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	for i := 0; i < b.N; i++ {
		hash.Get(sampleKeys[i%len(sampleKeys)])
	}
}

type getNTestcase struct {
	count         int
	key           string
	expectedNodes []string
}

func Test_Hash_GetN(t *testing.T) {
	hash := New()

	gotNodes := hash.GetN(2, "foo")
	if len(gotNodes) != 0 {
		t.Errorf("got: %#v, expected: %#v", gotNodes, []string{})
	}

	hash.Add("a", "b", "c", "d", "e")

	testcases := []getNTestcase{
		{1, "foo", []string{"e"}},
		{2, "bar", []string{"c", "e"}},
		{3, "baz", []string{"d", "a", "b"}},
		{2, "biz", []string{"b", "a"}},
		{0, "boz", []string{}},
		{100, "floo", []string{"d", "a", "b", "c", "e"}},
	}

	for _, testcase := range testcases {
		gotNodes := hash.GetN(testcase.count, testcase.key)
		if !reflect.DeepEqual(gotNodes, testcase.expectedNodes) {
			t.Errorf("got: %#v, expected: %#v", gotNodes, testcase.expectedNodes)
		}
	}
}

func BenchmarkHashGetN3_5_nodes(b *testing.B) {
	hash := New("a", "b", "c", "d", "e")
	for i := 0; i < b.N; i++ {
		hash.GetN(3, sampleKeys[i%len(sampleKeys)])
	}
}

func BenchmarkHashGetN5_5_nodes(b *testing.B) {
	hash := New("a", "b", "c", "d", "e")
	for i := 0; i < b.N; i++ {
		hash.GetN(5, sampleKeys[i%len(sampleKeys)])
	}
}

func BenchmarkHashGetN3_10_nodes(b *testing.B) {
	hash := New("a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	for i := 0; i < b.N; i++ {
		hash.GetN(3, sampleKeys[i%len(sampleKeys)])
	}
}

func BenchmarkHashGetN5_10_nodes(b *testing.B) {
	hash := New("a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	for i := 0; i < b.N; i++ {
		hash.GetN(5, sampleKeys[i%len(sampleKeys)])
	}
}

func TestHashRemove(t *testing.T) {
	hash := New("a", "b", "c")

	var keyForB string
	for i := 0; i < 10000; i++ {
		randomKey := fmt.Sprintf("key-%d", i)
		if hash.Get(randomKey) == "b" {
			keyForB = randomKey
			break
		}
	}

	if keyForB == "" {
		t.Fatalf("Failed to find a key that maps to 'b'")
	}

	hash.Remove("b")

	// Check if the key now maps to a different node
	newNode := hash.Get(keyForB)
	if newNode == "b" {
		t.Errorf("Key %s still maps to removed node 'b'", keyForB)
	}
	if newNode == "" {
		t.Errorf("Key %s does not map to any node after removing 'b'", keyForB)
	}
}
