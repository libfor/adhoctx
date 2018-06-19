package adhoctx

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type List32 struct {
	OldVersion *List32
	Items      []uint32
}

func (l *List32) Pointer() unsafe.Pointer {
	return unsafe.Pointer(l)
}

func (l *List32) New() *List32 {
	return &List32{OldVersion: l}
}

func (l *List32) Clear() {
	l.Items = make([]uint32, 0)
}

func (l *List32) Add(n uint32) {
	if l.Items == nil {
		l.Items = append(l.Items, l.OldVersion.All()...)
	}
	l.Items = append(l.Items, n)
}

func (l *List32) All() []uint32 {
	if l == nil {
		return nil
	}
	return (*l).Items
}

type List32View struct {
	current  unsafe.Pointer
	Counters struct {
		Commits   uint64
		Rollbacks uint64
	}
}

func NewList32View() *List32View {
	return new(List32View)
}

func (d *List32View) String() string {
	return fmt.Sprintf(`list32[commits:%d,rollbacks:%d]`,
		atomic.LoadUint64(&d.Counters.Commits),
		atomic.LoadUint64(&d.Counters.Rollbacks),
	)
}

func (d *List32View) Current() *List32 {
	return (*List32)(atomic.LoadPointer(&d.current))
}

func (d *List32View) Commit(version *List32) bool {
	oldVersion := version.OldVersion
	version.OldVersion = nil
	swapped := atomic.CompareAndSwapPointer(&d.current, oldVersion.Pointer(), version.Pointer())

	if swapped {
		atomic.AddUint64(&d.Counters.Commits, 1)
		return true
	}

	version.OldVersion = oldVersion
	atomic.AddUint64(&d.Counters.Rollbacks, 1)
	return false
}

type List32Transaction struct {
	ToClear bool
	ToAdd   uint32
}

func (l List32Transaction) Exec(d *List32View) {
	for {
		t := d.Current().New()
		if l.ToClear {
			t.Clear()
		}
		if l.ToAdd != 0 {
			t.Add(l.ToAdd)
		}
		if d.Commit(t) {
			break
		}
	}
}
