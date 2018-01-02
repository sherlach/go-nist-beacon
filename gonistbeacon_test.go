package beacon

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLastRecord(t *testing.T) {
	cli := &http.Client{}
	resp, err := LastRecord(cli)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
