// @anchors
//   code: LDMOV
//   ref: LDQUS

//go:build !darwin && !linux && !windows

package load

// totalMemory cannot be read on this system: memory does not bound the queue.
func totalMemory() int64 { return 0 }
