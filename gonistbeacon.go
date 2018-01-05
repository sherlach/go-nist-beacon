package beacon

import (
	"encoding/xml"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

// Record Chewed down version of the records the http server returns
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
	Record              xml.Name `xml:"record"`
	Version             string   `xml:"version"`
	Frequency           string   `xml:"frequency"`
	TimeStamp           string   `xml:"timeStamp"`
	SeedValue           string   `xml:"seedValue"`
	PreviousOutputValue string   `xml:"previousOutputValue"`
	SignatureValue      string   `xml:"signatureValue"`
	OutputValue         string   `xml:"outputValue"`
	StatusCode          string   `xml:"statusCode"`
}

func setString(s string, base int) big.Int {
	i := new(big.Int)
	i.SetString(s, base)
	return (*i)
}

func atoi(a string) int {
	b, err := strconv.Atoi(a)
	if err != nil {
		b = -1
	}
	return b
}

var defaultClient = &http.Client{}

// LastRecord Fetches the latest record from the beacon and returns the record
func LastRecord() (Record, error) {
	r, err := defaultClient.Get("https://beacon.nist.gov/rest/record/last")
	if err != nil {
		return Record{}, err
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Record{}, err
	}

	var drec dirtyrecord
	err = xml.Unmarshal(buf, &drec)
	if err != nil {
		return Record{}, err
	}

	rec := Record{
		Version:             drec.Version,
		Frequency:           atoi(drec.Frequency),
		TimeStamp:           time.Unix(int64(atoi(drec.TimeStamp)), 0),
		SeedValue:           setString(drec.SeedValue, 16),
		PreviousOutputValue: setString(drec.PreviousOutputValue, 16),
		SignatureValue:      setString(drec.SignatureValue, 16),
		OutputValue:         setString(drec.OutputValue, 16),
	}
	return rec, nil
}

// SetClient If you want to use your own client, to use a proxy to fetch the data for example.
func SetClient(cli *http.Client) {
	defaultClient = cli
}
