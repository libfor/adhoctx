package adhoctx

import (
	"sync"
	"testing"
	"unsafe"
)

func TestSimpleOps(t *testing.T) {
	v := NewPointerStackView()
	n := int(10)

	t.Log("adding a couple pointers, capturing the last one")
	v.AllocateID(unsafe.Pointer(&n))
	v.AllocateID(unsafe.Pointer(&n))
	v.AllocateID(unsafe.Pointer(&n))

	p := v.AllocateID(unsafe.Pointer(&n))
	t.Log("added, now getting", p)
	t.Log(v.String())
	nPtr := v.GetPointer(p)
	t.Log("got", *(*int)(nPtr), "now removing")
	t.Log(v.String())
	v.RemoveID(p)
	t.Log(v.String())
}

func TestEventualCommits(t *testing.T) {
	v := NewPointerStackView()
	var wg sync.WaitGroup
	for i := int32(0); i < 200; i++ {
		wg.Add(1)
		n := i
		t.Log("starting", i)
		go AddEventually(&wg, n, v)
	}
	t.Log("waiting")
	wg.Wait()
	t.Log(v.String())
}

func AddEventually(wg *sync.WaitGroup, n int32, v *PointerStackView) {
	if wg != nil {
		defer wg.Done()
	}
	p := v.AllocateID(unsafe.Pointer(&n))
	v.RemoveID(p)
}

func BenchmarkEventualCommits(b *testing.B) {
	v := NewPointerStackView()
	n := int32(0)
	para := func(pb *testing.PB) {
		for pb.Next() {
			AddEventually(nil, n, v)
		}
	}
	b.RunParallel(para)
	b.Log(v.String())
}
