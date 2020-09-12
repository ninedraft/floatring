// Package floatring provides a ring buffer of float64 values.
package floatring

// Buffer is the ring buffer of float64 values.
type Buffer struct {
	ring []float64
	full bool
	wc   int
	rc   int
}

// NewBuffer creates a new buffer with the given capacity.
func New(cp int) *Buffer {
	return &Buffer{
		ring: make([]float64, cp),
	}
}

// IsFull returns true if the buffer if filled with values.
func (buff *Buffer) IsFull() bool { return buff.full }

// IsEmpty returns true if the buffer contains no values.
func (buff *Buffer) IsEmpty() bool { return !buff.full && buff.wc == buff.rc }

// Cap returns the buffer capacity.
func (buff *Buffer) Cap() int { return len(buff.ring) }

// Len returns the number of values in the buffer.
func (buff *Buffer) Len() int {
	switch {
	case buff.IsEmpty():
		return 0
	case buff.IsFull():
		return len(buff.ring)
	case buff.wc >= buff.rc:
		return buff.wc - buff.rc
	default:
		return len(buff.ring) + buff.wc - buff.rc
	}
}

// WriteValue writes a single value to the buffer.
// If the buffer is overflowed, then the older value will be overwritten.
func (buff *Buffer) WriteValue(v float64) {
	if len(buff.ring) == 0 {
		return
	}
	var n = len(buff.ring)
	if buff.IsFull() {
		buff.rc = (buff.rc + 1) % n
	}
	buff.ring[buff.wc] = v
	buff.wc = (buff.wc + 1) % n
	buff.full = buff.wc == buff.rc
}

// Write appends provided floats to the buffer.
// If the buffer is overflowed, then the older value will be overwritten.
func (buff *Buffer) Write(vv []float64) {
	for _, v := range vv {
		buff.WriteValue(v)
	}
}

// ReadValue returns the oldest value in the buffer or false, if the buffer is empty.
func (buff *Buffer) ReadValue() (_ float64, ok bool) {
	if buff.IsEmpty() {
		return 0, false
	}
	var n = len(buff.ring)
	var x = buff.ring[buff.rc]
	buff.full = false
	buff.rc = (buff.rc + 1) % n
	return x, true
}

// Read reads values from the buffer to the given slice and return the number of read values.
func (buff *Buffer) Read(vv []float64) int {
	for i := range vv {
		var x, ok = buff.ReadValue()
		if !ok {
			return i
		}
		vv[i] = x
	}
	return len(vv)
}

// DumpTo returns the slice of values as [oldest...newest].
func (buff *Buffer) DumpTo(dst []float64) []float64 {
	dst = dst[:0]
	switch {
	case buff.IsEmpty():
		return dst
	case buff.wc > buff.rc:
		dst = append(dst, buff.ring[buff.rc:buff.wc]...)
	default:
		dst = append(dst, buff.ring[buff.rc:]...)
		dst = append(dst, buff.ring[:buff.wc]...)
	}
	return dst
}

// ForEach calls provided hook for each value in the buffer in the LIFO order.
func (buff *Buffer) ForEach(op func(x float64) bool) {
	if buff.IsEmpty() {
		return
	}
	if buff.wc > buff.rc {
		buff.forEach(op)
		return
	}
	buff.forEachChunked(op)
}

func (buff *Buffer) forEach(op func(x float64) bool) {
	for _, x := range buff.ring[buff.rc:buff.wc] {
		if !op(x) {
			return
		}
	}
}

func (buff *Buffer) forEachChunked(op func(x float64) bool) {
	for _, x := range buff.ring[buff.rc:] {
		if !op(x) {
			return
		}
	}
	for _, x := range buff.ring[:buff.wc] {
		if !op(x) {
			return
		}
	}
}

// Reset the buffer to a zero state. Doesn't change the capacity.
func (buff *Buffer) Reset() {
	buff.wc = 0
	buff.rc = 0
	buff.full = false
}
