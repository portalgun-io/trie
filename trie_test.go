package trie

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func addFromFile(t *Trie, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewScanner(file)

	for reader.Scan() {
		t.Add(reader.Text())
	}

	if reader.Err() != nil {
		log.Fatal(err)
	}
}

func TestTrieAdd(t *testing.T) {
	trie := CreateTrie()

	i := trie.Add("foo")

	if i != 3 {
		t.Errorf("Expected 3, got: %d", i)
	}
}

func TestRemove(t *testing.T) {
	trie := CreateTrie()
	initial := []string{"football", "foostar", "foosball"}

	for _, key := range initial {
		trie.Add(key)
	}

	trie.Remove("foosball")
	keys := trie.Keys()

	if len(keys) != 2 {
		t.Errorf("Expected 2 keys got %d", len(keys))
	}

	for _, k := range keys {
		if k != "football" && k != "foostar" {
			t.Errorf("key was: %s", k)
		}
	}

	keys = trie.FuzzySearch("foo")

	if len(keys) != 2 {
		t.Errorf("Expected 2 keys got %d", len(keys))
	}

	for _, k := range keys {
		if k != "football" && k != "foostar" {
			t.Errorf("Expected football got: %#v", k)
		}
	}
}

func TestTrieKeys(t *testing.T) {
	trie := CreateTrie()
	expected := []string{"bar", "foo"}

	for _, key := range expected {
		trie.Add(key)
	}

	kl := len(trie.Keys())
	if kl != 2 {
		t.Errorf("Expected 2 keys, got %d, keys were: %v", kl, trie.Keys())
	}

	for i, key := range trie.Keys() {
		if key != expected[i] {
			t.Errorf("Expected %#v, got %#v", expected[i], key)
		}
	}
}

func TestPrefixSearch(t *testing.T) {
	trie := CreateTrie()
	expected := []string{"foosball", "football", "foreboding", "forementioned", "foretold", "foreverandeverandeverandever", "forbidden"}
	defer func() {
		r := recover()
		if r != nil {
			t.Error(r)
		}
	}()

	trie.Add("bar")
	for _, key := range expected {
		trie.Add(key)
	}

	tests := []struct {
		pre      string
		expected []string
		length   int
	}{
		{"fo", expected, len(expected)},
		{"foosbal", []string{"foosball"}, 1},
		{"abc", []string{}, 0},
	}

	for _, test := range tests {
		actual := trie.PrefixSearch(test.pre)
		if len(actual) != test.length {
			t.Errorf("Expected len(actual) to == %d for pre %s", test.length, test.pre)
		}

		for i, key := range actual {
			if key != test.expected[i] {
				t.Errorf("Expected %v got: %v", test.expected[i], key)
			}
		}
	}

	trie.PrefixSearch("fsfsdfasdf")
}

func TestFuzzySearch(t *testing.T) {
	trie := CreateTrie()
	setup := []string{
		"foosball",
		"football",
		"bmerica",
		"ked",
		"kedlock",
		"frosty",
		"bfrza",
		"foo/bar/baz.go",
	}
	tests := []struct {
		partial string
		length  int
	}{
		{"fsb", 1},
		{"footbal", 1},
		{"football", 1},
		{"fs", 2},
		{"oos", 1},
		{"kl", 1},
		{"ft", 2},
		{"fy", 1},
		{"fz", 2},
		{"a", 5},
	}

	for _, key := range setup {
		trie.Add(key)
	}

	for _, test := range tests {
		actual := trie.FuzzySearch(test.partial)

		if len(actual) != test.length {
			t.Errorf("Expected len(actual) to == %d, was %d for %s", test.length, len(actual), test.partial)
		}
	}
}

func BenchmarkTieKeys(b *testing.B) {
	trie := CreateTrie()
	keys := []string{"bar", "foo", "baz", "bur", "zum", "burzum", "bark", "barcelona", "football", "foosball", "footlocker"}

	for _, key := range keys {
		trie.Add(key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Keys()
	}
}

func BenchmarkPrefixSearch(b *testing.B) {
	trie := CreateTrie()
	addFromFile(trie, "/usr/share/dict/words")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.PrefixSearch("fo")
	}
}

func BenchmarkFuzzySearch(b *testing.B) {
	trie := CreateTrie()
	addFromFile(trie, "/usr/share/dict/words")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.FuzzySearch("fs")
	}
}
