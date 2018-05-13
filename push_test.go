package stringid

import (
	"strings"
	"sync"
	"testing"
)

func TestPushGenerator(t *testing.T) {
	pg := NewPushGenerator(nil)

	a, b := pg.Generate(), pg.Generate()
	if n := len(a); n != 20 {
		t.Errorf("length of a should be 20, got: %d", n)
	}

	if n := len(b); n != 20 {
		t.Errorf("length of b should be 20, got: %d", n)
	}

	if a == b {
		t.Errorf("a (%q) and b (%q) are equal", a, b)
	}

	if !(strings.Compare(a, b) < 0) {
		t.Errorf("a (%q) is not less than than b (%q)", a, b)
	}
}

func TestPushGeneratorMany(t *testing.T) {
	pg := NewPushGenerator(nil)

	wg := new(sync.WaitGroup)
	for i := 0; i < 4; i++ {
		wg.Add(1)

		go func(t *testing.T, wg *sync.WaitGroup, pg *PushGenerator) {
			defer wg.Done()

			var id, prev string
			ids := make(map[string]bool)
			for i := 0; i < 1000000; i++ {
				id = pg.Generate()
				if n := len(id); n != 20 {
					t.Errorf("generated id length should be 20, got: %d", n)
				}

				if _, exists := ids[id]; exists {
					t.Errorf("generated duplicate id %q", id)
				}

				if !(strings.Compare(prev, id) < 0) {
					t.Errorf("previously generated id %q is not less than generated id %q", prev, id)
				}

				ids[id], prev = true, id
			}
		}(t, wg, pg)
	}
	wg.Wait()
}
