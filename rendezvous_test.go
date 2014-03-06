package rendezvous

import (
	"reflect"
	"testing"
)

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
