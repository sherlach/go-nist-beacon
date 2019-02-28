package beacon

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
	"time"
)

func TestLastRecord(t *testing.T) {
	resp, err := LastRecord()
	if err != nil {
		panic(err)
	}
	spew.Dump(resp)
}

//If if returns an error the test was a sucess, since there is no current record as new as now. If it fails, check if your system clock is running correctly
func TestCurrentRecord(t *testing.T) {
	resp, err := NextRecord(time.Now())
	if err == nil {
		t.Fail()
	}
	spew.Dump(resp, err, time.Now().Unix())

	resp, err = NextRecord(time.Now().AddDate(0, 0, -1))
	if err != nil {
		panic(err)
	}
	spew.Dump(resp)
}

func TestPreviousRecord(t *testing.T) {
	resp, err := PreviousRecord(time.Now())
	if err != nil {
		panic(err)
	}
	spew.Dump(resp)
}

//If if returns an error the test was a sucess, since there is no next record to the latest one. If it fails, check if your system clock is running correctly
func TestNextRecord(t *testing.T) {
	resp, err := NextRecord(time.Now())
	if err == nil {
		t.Fail()
	}
	spew.Dump(resp, err, time.Now().Unix())
}
