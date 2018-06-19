package adhoctx

import (
	"fmt"
	"log"
	"sync/atomic"
	"unsafe"
)

const DebugMemAtomic = false

type atomicCommitter struct {
	current  unsafe.Pointer
	Counters struct {
		Commits   uint64
		Rollbacks uint64
	}
}

func (d *atomicCommitter) String() string {
	return fmt.Sprintf(`list32[commits:%d,rollbacks:%d]`,
		atomic.LoadUint64(&d.Counters.Commits),
		atomic.LoadUint64(&d.Counters.Rollbacks),
	)
}

func (d *atomicCommitter) unsafe_race_condition_commit(id byte, new unsafe.Pointer) {
	if DebugMemAtomic {
		log.Printf("ac set %c %p %p", id, d.current, new)
	}

	atomic.StorePointer(&d.current, new)
}

func (d *atomicCommitter) getCurrent(id byte) unsafe.Pointer {
	if DebugMemAtomic {
		log.Printf("ac get %c %p", id, d.current)
	}

	return atomic.LoadPointer(&d.current)
}

func (d *atomicCommitter) tryCommit(id byte, old, new unsafe.Pointer) bool {
	if DebugMemAtomic {
		log.Printf("ac try %c %p %p %p", id, d.current, old, new)
	}

	swapped := atomic.CompareAndSwapPointer(&d.current, old, new)

	if swapped {
		atomic.AddUint64(&d.Counters.Commits, 1)
		return true
	}

	atomic.AddUint64(&d.Counters.Rollbacks, 1)
	return false
}
