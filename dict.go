package adhoctx

import (
	"unsafe"
)

type dictType map[string]string

type DictView struct {
	dictsByID     PointerStackView
	currentDictID atomicCommitter
}

func NewDictView() *DictView {
	dv := new(DictView)
	dict := make(dictType)
	dPtr := unsafe.Pointer(&dict)

	currVer := dv.dictsByID.AllocateID(dPtr)
	dv.currentDictID.unsafe_race_condition_commit('d', unsafe.Pointer(&currVer))
	return dv
}

func (d *DictView) Reader() uint32 {
	return *(*uint32)(d.currentDictID.getCurrent('d'))
}

func (d *DictView) ReadWriter(u uint32) uint32 {
	currDict := *(*dictType)(d.dictsByID.GetPointer(u))
	newDict := make(dictType, len(currDict))
	for k, v := range currDict {
		newDict[k] = v
	}
	return d.dictsByID.AllocateID(unsafe.Pointer(&newDict))
}

func (d *DictView) Commit(old, new uint32) bool {
	if !d.currentDictID.tryCommit('d', d.currentDictID.getCurrent('d'), unsafe.Pointer(&new)) {
		return false
	}

	d.dictsByID.RemoveID(old)
	return true
}

func (d *DictView) SetKey(u uint32, key, val string) {
	dict := *(*dictType)(d.dictsByID.GetPointer(u))
	dict[key] = val
}

func (d *DictView) GetKey(u uint32, key string) string {
	dict := *(*dictType)(d.dictsByID.GetPointer(u))
	return dict[key]
}

func (d *DictView) String() string {
	return "dictview[currentDictID:" + d.currentDictID.String() + ",dictsByID:" + d.dictsByID.String() + "]"
}
