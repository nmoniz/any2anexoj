package internal

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestFillerQueue(t *testing.T) {
	var recCount int
	newRecord := func() Record {
		recCount++
		return testRecord{
			id: recCount,
		}
	}

	var rq FillerQueue

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

	rq.Push(NewFiller(newRecord()))
	if rq.Len() != 1 {
		t.Fatalf("pushing 1st record should result in lenght of 1")
	}

	rq.Push(NewFiller(newRecord()))
	if rq.Len() != 2 {
		t.Fatalf("pushing 2nd record should result in lenght of 2")
	}

	peekFiller, ok := rq.Peek()
	if !ok {
		t.Fatalf("Peek() should return (_,true) when the list is not empty")
	}

	if rec, ok := peekFiller.Record.(testRecord); ok {
		if rec.id != 1 {
			t.Fatalf("Peek() should return the 1st record pushed but returned %d", rec.id)
		}
	} else {
		t.Fatalf("Peek() should return the original record type")
	}

	if rq.Len() != 2 {
		t.Fatalf("Peek() should not affect the list length")
	}

	popFiller, ok := rq.Pop()
	if !ok {
		t.Fatalf("Pop() should return (_,true) when the list is not empty")
	}

	if rec, ok := popFiller.Record.(testRecord); ok {
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

func TestFillerQueueNilReceiver(t *testing.T) {
	var rq *FillerQueue

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

		expMsg := "Push to nil FillerQueue"
		if msg, ok := r.(string); !ok || msg != expMsg {
			t.Fatalf(`want panic message %q but got "%v"`, expMsg, r)
		}
	}()
	rq.Push(NewFiller(nil))
}

type testRecord struct {
	Record

	id       int
	quantity decimal.Decimal
}

func (tr testRecord) Quantity() decimal.Decimal {
	return tr.quantity
}

func TestFiller_Fill(t *testing.T) {
	tests := []struct {
		name     string
		r        Record
		quantity decimal.Decimal
		want     decimal.Decimal
		wantBool bool
	}{
		{
			name:     "fills 0 of zero quantity",
			r:        &testRecord{quantity: decimal.NewFromFloat(0.0)},
			quantity: decimal.Decimal{},
			want:     decimal.Decimal{},
			wantBool: true,
		},
		{
			name:     "fills 0 of positive quantity",
			r:        &testRecord{quantity: decimal.NewFromFloat(100.0)},
			quantity: decimal.Decimal{},
			want:     decimal.Decimal{},
			wantBool: false,
		},
		{
			name:     "fills 10 out of 100 and no previous fills",
			r:        &testRecord{quantity: decimal.NewFromFloat(100.0)},
			quantity: decimal.NewFromFloat(10),
			want:     decimal.NewFromFloat(10),
			wantBool: false,
		},
		{
			name:     "fills 10 out of 10 and no previous fills",
			r:        &testRecord{quantity: decimal.NewFromFloat(10.0)},
			quantity: decimal.NewFromFloat(10),
			want:     decimal.NewFromFloat(10),
			wantBool: true,
		},
		{
			name:     "filling 100 fills 10 out of 10 and no previous fills",
			r:        &testRecord{quantity: decimal.NewFromFloat(10.0)},
			quantity: decimal.NewFromFloat(100),
			want:     decimal.NewFromFloat(10),
			wantBool: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFiller(tt.r)
			got, gotBool := f.Fill(tt.quantity)
			if !tt.want.Equal(got) {
				t.Errorf("want 1st return value to be %v but got %v", tt.want, got)
			}
			if tt.wantBool != gotBool {
				t.Errorf("want 2nd return value to be %v but got %v", tt.wantBool, gotBool)
			}
		})
	}
}
