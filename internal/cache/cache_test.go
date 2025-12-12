package cache

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCacheAddGet(t *testing.T) {
	tests := map[string]struct {
		key      string
		val      []byte
		expected []byte
	}{
		"valid": {
			key:      "https://example.com",
			val:      []byte("testdata"),
			expected: []byte("testdata"),
		},
		"urWithPath": {
			key:      "https://example.com/path",
			val:      []byte("moretestdata"),
			expected: []byte("moretestdata"),
		},
		"invalidWithoutKey": {
			key:      "",
			val:      []byte("value for no key"),
			expected: []byte(""),
		},
	}

	c := NewCache(5 * time.Second)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c.Add(tc.key, tc.val)
			b, _ := c.Get(tc.key)

			diff := cmp.Diff(tc.expected, b)
			if diff != "" {
				t.Fatalf("expected : %s , got : %s", string(tc.expected), string(b))
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
