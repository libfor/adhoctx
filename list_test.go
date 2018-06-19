package adhoctx

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestReadsAndWriters(t *testing.T) {
	v := NewList32View()

	{
		n := v.Current().New()
		n.Add(10)
		v.Commit(n)
	}

	{
		n := v.Current().New()
		n.Add(20)
		t.Log("found", n.All())
		m := v.Current().New()
		m.Clear()
		m.Add(30)
		t.Log("found", m.All())
		v.Commit(m)
	}

}

func TestEventualCommits(t *testing.T) {
	v := NewList32View()
	var wg sync.WaitGroup
	for i := uint32(0); i < 200; i++ {
		wg.Add(1)
		n := i
		t.Log("starting", i)
		go AddEventually(&wg, n, v)
	}
	t.Log("waiting")
	wg.Wait()
	t.Log(v.Current().All())
	t.Log(v.String())
}

func AddEventually(wg *sync.WaitGroup, n uint32, v *List32View) {
	if wg != nil {
		defer wg.Done()
	}
	List32Transaction{ToClear: true, ToAdd: n}.Exec(v)
}

func BenchmarkEventualCommits(b *testing.B) {
	v := NewList32View()
	n := uint32(0)
	para := func(pb *testing.PB) {
		for pb.Next() {
			n := atomic.AddUint32(&n, 1)
			AddEventually(nil, n, v)
		}
	}
	b.RunParallel(para)
	b.Log(v.String())
}

func BenchmarkReadsAndWriters(b *testing.B) {

	v := NewList32View()

	{
		n := v.Current().New()
		n.Add(10)
		v.Commit(n)
	}

	para := func(pb *testing.PB) {
		for pb.Next() {
			{
				n := v.Current().New()
				n.Add(20)
				if n.All()[1] != 20 {
					b.Error("wrong val")
				}
				if len(n.All()) != 2 {
					b.Error("wrong size")
				}
			}

			{
				n := v.Current().New()
				n.Add(30)
				if n.All()[1] != 30 {
					b.Error("wrong val")
				}
				if len(n.All()) != 2 {
					b.Error("wrong size")
				}
			}
		}
	}
	b.RunParallel(para)
}
