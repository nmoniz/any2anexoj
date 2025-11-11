package internal

import (
	"container/list"
	"math/big"
	"time"
)

type Record interface {
	Symbol() string
	Side() Side
	Price() *big.Float
	Quantity() *big.Float
	Timestamp() time.Time
}

type RecordQueue struct {
	l *list.List
}

// Push inserts the Record at the back of the queue. If pushing a nil Record then it's a no-op.
func (rq *RecordQueue) Push(r Record) {
	if r == nil {
		return
	}

	if rq == nil {
		// This would cause a panic anyway so, we panic with a more meaningful message
		panic("Push to nil RecordQueue")
	}

	if rq.l == nil {
		rq.l = list.New()
	}

	rq.l.PushBack(r)
}

// Pop removes and returns the first element of the list as the first return value. If the list is
// empty returns falso on the 2nd return value, true otherwise.
func (rq *RecordQueue) Pop() (Record, bool) {
	if rq == nil || rq.l == nil {
		return nil, false
	}

	el := rq.l.Front()
	if el == nil {
		return nil, false
	}

	val := rq.l.Remove(el)

	return val.(Record), true
}

func (rq *RecordQueue) Len() int {
	if rq == nil || rq.l == nil {
		return 0
	}

	return rq.l.Len()
}
