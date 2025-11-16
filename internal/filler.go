package internal

import (
	"container/list"

	"github.com/shopspring/decimal"
)

type Filler struct {
	Record

	filled decimal.Decimal
}

func NewFiller(r Record) *Filler {
	return &Filler{
		Record: r,
	}
}

// Fill accrues some quantity. Returns how mutch was accrued in the 1st return value and whether
// it was filled or not on the 2nd return value.
func (f *Filler) Fill(quantity decimal.Decimal) (decimal.Decimal, bool) {
	unfilled := f.Record.Quantity().Sub(f.filled)
	delta := decimal.Min(unfilled, quantity)
	f.filled = f.filled.Add(delta)
	return delta, f.IsFilled()
}

// IsFilled returns true if the fill is equal to the record quantity.
func (f *Filler) IsFilled() bool {
	return f.filled.Equal(f.Quantity())
}

type FillerQueue struct {
	l *list.List
}

// Push inserts the Filler at the back of the queue.
func (fq *FillerQueue) Push(f *Filler) {
	if f == nil {
		return
	}

	if fq == nil {
		// This would cause a panic anyway so, we panic with a more meaningful message
		panic("Push to nil FillerQueue")
	}

	if fq.l == nil {
		fq.l = list.New()
	}

	fq.l.PushBack(f)
}

// Pop removes and returns the first Filler of the queue in the 1st return value. If the list is
// empty returns false on the 2nd return value, true otherwise.
func (fq *FillerQueue) Pop() (*Filler, bool) {
	el := fq.frontElement()
	if el == nil {
		return nil, false
	}

	val := fq.l.Remove(el)

	return val.(*Filler), true
}

// Peek returns the front Filler of the queue in the 1st return value. If the list is empty returns
// false on the 2nd return value, true otherwise.
func (fq *FillerQueue) Peek() (*Filler, bool) {
	el := fq.frontElement()
	if el == nil {
		return nil, false
	}

	return el.Value.(*Filler), true
}

func (fq *FillerQueue) frontElement() *list.Element {
	if fq == nil || fq.l == nil {
		return nil
	}

	return fq.l.Front()
}

// Len returns how many elements are currently on the queue
func (fq *FillerQueue) Len() int {
	if fq == nil || fq.l == nil {
		return 0
	}

	return fq.l.Len()
}
