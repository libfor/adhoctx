package adhoctx

import (
	"unsafe"
)

type PointerStack struct {
	OldVersion *PointerStack

	PointersByID   []*unsafe.Pointer
	LowestFreeSlot uint32
}

func (l *PointerStack) New() *PointerStack {
	newPs := &PointerStack{OldVersion: l}
	if l != nil {
		newPs.PointersByID = make([]*unsafe.Pointer, len(l.PointersByID), len(l.PointersByID))
		copy(newPs.PointersByID, l.PointersByID)
	}
	return newPs
}

func (l *PointerStack) Allocate(n unsafe.Pointer) uint32 {
	for i, ptr := range l.PointersByID[l.LowestFreeSlot:] {
		if ptr == nil {
			l.PointersByID[l.LowestFreeSlot+uint32(i)] = &n
			return l.LowestFreeSlot + uint32(i)
		}
	}
	newIDX := uint32(len(l.PointersByID))
	l.LowestFreeSlot = newIDX
	l.PointersByID = append(l.PointersByID, &n)
	return newIDX
}

func (l *PointerStack) Remove(n uint32) {
	if l.LowestFreeSlot < n {
		l.LowestFreeSlot = n
	}
	l.PointersByID[n] = nil
}

func (l PointerStack) Get(n uint32) unsafe.Pointer {
	if l.PointersByID == nil {
		return *l.OldVersion.PointersByID[n]
	}
	return *l.PointersByID[n]
}

type PointerStackView struct {
	atomicCommitter
}

func NewPointerStackView() *PointerStackView {
	return new(PointerStackView)
}

func (d *PointerStackView) Current() *PointerStack {
	return (*PointerStack)(d.getCurrent('p'))
}

func (d *PointerStackView) Commit(version *PointerStack) bool {
	oldVersion := version.OldVersion
	version.OldVersion = nil
	return d.tryCommit('p', unsafe.Pointer(oldVersion), unsafe.Pointer(version))
}

func (d *PointerStackView) AllocateID(p unsafe.Pointer) uint32 {
	for {
		t := d.Current().New()
		id := t.Allocate(p)
		if d.Commit(t) {
			return id
		}
	}
}

func (d *PointerStackView) RemoveID(id uint32) {
	for {
		t := d.Current().New()
		t.Remove(id)
		if d.Commit(t) {
			return
		}
	}
}

func (d *PointerStackView) GetPointer(id uint32) unsafe.Pointer {
	ptr := d.Current().Get(id)
	return ptr
}
