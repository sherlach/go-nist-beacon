//Package beacon implements an easy to use, but featurerich NIST Randomness Beacon API Wrapper in go
package beacon

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Record is the data the NIST api returns
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
	_, ok := i.SetString(s, base)
	if !ok {
		i.SetInt64(-1)
	}
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

// SetClient is useful if you want to use your own http client, it adds the possibility to use a proxy to fetch the data for example.
func SetClient(cli *http.Client) {
	defaultClient = cli
}

func getRecord(url string) (Record, error) {
	r, err := defaultClient.Get(url)
	if err != nil {
		err = errors.New("Couldn't get the record from the API: " + err.Error())
		return Record{}, err
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.New("Couldn't read the API's response: " + err.Error())
		return Record{}, err
	}

	var drec dirtyrecord
	err = xml.Unmarshal(buf, &drec)
	if err != nil {
		err = errors.New("Couldn't unmarshal the API's response: " + err.Error())
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

// LastRecord fetches the latest record from the beacon and returns the record
func LastRecord() (Record, error) {
	return getRecord("https://beacon.nist.gov/rest/record/last")
}

// CurrentRecord fetches the record closest to the given timestamp
func CurrentRecord(t time.Time) (Record, error) {
	return getRecord("https://beacon.nist.gov/rest/record/" + strconv.FormatInt(t.Unix(), 10))
}

// PreviousRecord fetches the record previous to the given timestamp
func PreviousRecord(t time.Time) (Record, error) {
	return getRecord("https://beacon.nist.gov/rest/record/previous/" + strconv.FormatInt(t.Unix(), 10))
}

// NextRecord fetches the record after the given timestamp
func NextRecord(t time.Time) (Record, error) {
	return getRecord("https://beacon.nist.gov/rest/record/next/" + strconv.FormatInt(t.Unix(), 10))
}

// StartChainRecord fetches the start chain record for the given timestamp
func StartChainRecord(t time.Time) (Record, error) {
	return getRecord("https://beacon.nist.gov/rest/record/start-chain/" + strconv.FormatInt(t.Unix(), 10))
}

type Rand struct {
	update     bool
	updateTime time.Time
	rand       *rand.Rand
}

// NewRand creates a new random number generator using the given record's SeedValue>>MaxInt64 as source
func NewRand(r Record) *Rand {
	ret := new(Rand)
	ret.rand = new(rand.Rand)
	i := new(big.Int)
	ret.SetSeed(i.Rsh(&r.SeedValue, 448).Int64())
	return ret
}

func NewUpdatedRand() (*Rand, error) {
	rec, err := LastRecord()
	if err != nil {
		return nil, err
	}
	r := NewRand(rec)
	r.update = true
	r.updateTime = rec.TimeStamp
	return r, nil
}

func (r *Rand) SetSeed(n int64) {
	r.rand = rand.New(rand.NewSource(n))
	r.update = false
}

func (r *Rand) Int() int {
	if r.update && time.Now().After(r.updateTime.Add(time.Duration(time.Minute*1))) {
		rec, err := LastRecord()
		if err != nil {
			panic(errors.New("Couldn't update to the last record: " + err.Error()))
		}
		r.updateTime = rec.TimeStamp
		i := new(big.Int)
		r.SetSeed(i.Rsh(&rec.SeedValue, 448).Int64())
	}
	return r.rand.Int()
}
