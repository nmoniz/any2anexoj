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

	rec, ok := rq.Pop()
	if !ok {
		t.Fatalf("Pop() should return (_,true) when the list is not empty")
	}

	if rec, ok := rec.(testRecord); ok {
		if rec.id != 1 {
			t.Fatalf("Pop() should return the first record pushed but returned %d", rec.id)
		}
	} else {
		t.Fatalf("Pop() should return the original record")
	}
}

func TestRecordQueueNilReceiver(t *testing.T) {
	var rq *RecordQueue

	if rq.Len() > 0 {
		t.Fatalf("nil receiver should have zero lenght")
	}

	_, ok := rq.Pop()
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
