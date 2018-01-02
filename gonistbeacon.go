package beacon

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

type Record struct {
	Version             string
	Frequency           int
	TimeStamp           time.Time
	SeedValue           big.Int
	PreviousOutputValue big.Int
	SignatureValue      big.Int
	OutputValue         big.Int
}

type dirtyrecord struct {
	version             string
	frequency           int
	timeStamp           int64
	seedValue           string
	previousOutputValue string
	signatureValue      string
	outputValue         string
}

func setString(s string, base int) big.Int {
	i := new(big.Int)
	i.SetString(s, base)
	return (*i)
}

func LastRecord(cli *http.Client) (Record, error) {
	r, err := cli.Get("https://beacon.nist.gov/rest/record/last")
	if err != nil {
		return Record{}, err
	}

	buf, err := ioutil.ReadAll(r.Body)
	fmt.Println(buf)

	var drec dirtyrecord
	d := xml.NewDecoder(r.Body)
	err = d.Decode(&drec)
	if err != nil {
		//return Record{}, err
		panic(err)
	}

	fmt.Println(drec)

	rec := Record{
		Version:             drec.version,
		Frequency:           drec.frequency,
		TimeStamp:           time.Unix(drec.timeStamp, 0),
		SeedValue:           setString(drec.seedValue, 16),
		PreviousOutputValue: setString(drec.previousOutputValue, 16),
		SignatureValue:      setString(drec.signatureValue, 16),
		OutputValue:         setString(drec.outputValue, 16),
	}
	return rec, nil
}
