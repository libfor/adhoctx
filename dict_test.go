package adhoctx

import (
	"testing"
)

func TestDict(t *testing.T) {
	dict := NewDictView()
	curr := dict.Reader()
	t.Log("curr version", curr)
	rw := dict.ReadWriter(curr)
	t.Log("tx version", rw)
	dict.SetKey(rw, "foo", "bar")
	t.Log("expect bar, got", dict.GetKey(rw, "foo"))
	committed := dict.Commit(curr, rw)
	t.Log("expect true, got", committed)

	newState := dict.Reader()
	t.Log("new version", newState)
	t.Log("expected bar, got", dict.GetKey(newState, "foo"))

	t.Log("stats", dict.String())
}

func BenchmarkNormalMaps(b *testing.B) {
	m := make(map[string]string)
	for i := 0; i < b.N; i++ {
		m["hello"] = m["goodbye"]
		m["goodbye"] = "forever"
	}
	b.Log(m)
}

func BenchmarkViewMaps(b *testing.B) {
	m := NewDictView()
	curr := m.Reader()
	rw := m.ReadWriter(curr)
	for i := 0; i < b.N; i++ {
		m.SetKey(rw, "hello", m.GetKey(rw, "goodbye"))
		m.SetKey(rw, "goodbye", "forever")
		m.Commit(curr, rw)
		curr = rw
		rw = m.ReadWriter(curr)
	}
	b.Log(m.String())
}
