package internal

import (
	"testing"
)

func TestRecordQueue_Push(t *testing.T) {
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

type testRecord struct {
	Record

	id int
}
