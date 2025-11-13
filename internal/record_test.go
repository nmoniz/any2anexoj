package internal

import (
	"testing"
)

func TestRecordQueue(t *testing.T) {
	var recCount int
	newRecord := func() Record {
		recCount++
		return testRecord{
			id: recCount,
		}
	}

	var rq RecordQueue

	if rq.Len() != 0 {
		t.Fatalf("zero value should have zero lenght")
	}

	_, ok := rq.Pop()
	if ok {
		t.Fatalf("Pop() should return (_,false) on a zero value")
	}

	_, ok = rq.Peek()
	if ok {
		t.Fatalf("Peek() should return (_,false) on a zero value")
	}

	rq.Push(nil)
	if rq.Len() != 0 {
		t.Fatalf("pushing nil should be a no-op")
	}

	rq.Push(newRecord())
	if rq.Len() != 1 {
		t.Fatalf("pushing 1st record should result in lenght of 1")
	}

	rq.Push(newRecord())
	if rq.Len() != 2 {
		t.Fatalf("pushing 2nd record should result in lenght of 2")
	}

	peekRec, ok := rq.Peek()
	if !ok {
		t.Fatalf("Peek() should return (_,true) when the list is not empty")
	}

	if peekRec, ok := peekRec.(testRecord); ok {
		if peekRec.id != 1 {
			t.Fatalf("Peek() should return the 1st record pushed but returned %d", peekRec.id)
		}
	} else {
		t.Fatalf("Peek() should return the original record type")
	}

	if rq.Len() != 2 {
		t.Fatalf("Peek() should not affect the list length")
	}

	popRec, ok := rq.Pop()
	if !ok {
		t.Fatalf("Pop() should return (_,true) when the list is not empty")
	}

	if rec, ok := popRec.(testRecord); ok {
		if rec.id != 1 {
			t.Fatalf("Pop() should return the first record pushed but returned %d", rec.id)
		}
	} else {
		t.Fatalf("Pop() should return the original record")
	}

	if rq.Len() != 1 {
		t.Fatalf("Pop() should remove an element from the list")
	}
}

func TestRecordQueueNilReceiver(t *testing.T) {
	var rq *RecordQueue

	if rq.Len() > 0 {
		t.Fatalf("nil receiver should have zero lenght")
	}

	_, ok := rq.Peek()
	if ok {
		t.Fatalf("Peek() on a nil receiver should return (_,false)")
	}

	_, ok = rq.Pop()
	if ok {
		t.Fatalf("Pop() on a nil receiver should return (_,false)")
	}

	rq.Push(nil)
	if rq.Len() != 0 {
		t.Fatalf("Push(nil) on a nil receiver should be a no-op")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected a panic but got nothing")
		}

		expMsg := "Push to nil RecordQueue"
		if msg, ok := r.(string); !ok || msg != expMsg {
			t.Fatalf(`want panic message %q but got "%v"`, expMsg, r)
		}
	}()
	rq.Push(testRecord{})
}

type testRecord struct {
	Record

	id int
}
