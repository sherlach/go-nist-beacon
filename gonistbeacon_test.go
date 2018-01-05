package beacon

import (
	"fmt"
	"testing"
)

func TestLastRecord(t *testing.T) {
	resp, err := LastRecord()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
