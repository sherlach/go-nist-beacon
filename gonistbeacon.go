//Package beacon implements an easy to use, but feature rich NIST Randomness Beacon API Wrapper in go
package beacon

import (
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
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

const beaconCertPEM = `
-----BEGIN CERTIFICATE-----
MIIHZTCCBk2gAwIBAgIESTWNPjANBgkqhkiG9w0BAQsFADBtMQswCQYDVQQGEwJV
UzEQMA4GA1UEChMHRW50cnVzdDEiMCAGA1UECxMZQ2VydGlmaWNhdGlvbiBBdXRo
b3JpdGllczEoMCYGA1UECxMfRW50cnVzdCBNYW5hZ2VkIFNlcnZpY2VzIFNTUCBD
QTAeFw0xNDA1MDcxMzQ4MzZaFw0xNzA1MDcxNDE4MzZaMIGtMQswCQYDVQQGEwJV
UzEYMBYGA1UEChMPVS5TLiBHb3Zlcm5tZW50MR8wHQYDVQQLExZEZXBhcnRtZW50
IG9mIENvbW1lcmNlMTcwNQYDVQQLEy5OYXRpb25hbCBJbnN0aXR1dGUgb2YgU3Rh
bmRhcmRzIGFuZCBUZWNobm9sb2d5MRAwDgYDVQQLEwdEZXZpY2VzMRgwFgYDVQQD
Ew9iZWFjb24ubmlzdC5nb3YwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQC/m2xcckaSYztt6/6YezaUmqIqY5CLvrfO2esEIJyFg+cv7S7exL3hGYeDCnQL
VtUIGViAnO9yCXDC2Kymen+CekU7WEtSB96xz/xGrY3mbwjS46QSOND9xSRMroF9
xbgqXxzJ7rL/0RMUkku3uurGb/cxUpzKt6ra7iUnzkk3BBk73kr2OXFyYYbtrN71
s0B9qKKJZuPQqmA5n80Xc3E2YbaoAW4/gesncFNL7Sdxw9NIA1L4feu/o8xp3FNP
pv2e25C0113x+yagvb1W0mw6ISwAKhJ+6G4t4hFejl7RujuiDfORgzIhHMR4CyWt
PZFVn2qxZuVooj1+mduLIXhDAgMBAAGjggPKMIIDxjAOBgNVHQ8BAf8EBAMCBsAw
FwYDVR0gBBAwDjAMBgpghkgBZQMCAQMHMIIBXgYIKwYBBQUHAQEEggFQMIIBTDCB
uAYIKwYBBQUHMAKGgatsZGFwOi8vc3NwZGlyLm1hbmFnZWQuZW50cnVzdC5jb20v
b3U9RW50cnVzdCUyME1hbmFnZWQlMjBTZXJ2aWNlcyUyMFNTUCUyMENBLG91PUNl
cnRpZmljYXRpb24lMjBBdXRob3JpdGllcyxvPUVudHJ1c3QsYz1VUz9jQUNlcnRp
ZmljYXRlO2JpbmFyeSxjcm9zc0NlcnRpZmljYXRlUGFpcjtiaW5hcnkwSwYIKwYB
BQUHMAKGP2h0dHA6Ly9zc3B3ZWIubWFuYWdlZC5lbnRydXN0LmNvbS9BSUEvQ2Vy
dHNJc3N1ZWRUb0VNU1NTUENBLnA3YzBCBggrBgEFBQcwAYY2aHR0cDovL29jc3Au
bWFuYWdlZC5lbnRydXN0LmNvbS9PQ1NQL0VNU1NTUENBUmVzcG9uZGVyMBsGA1Ud
CQQUMBIwEAYJKoZIhvZ9B0QdMQMCASIwggGHBgNVHR8EggF+MIIBejCB6qCB56CB
5IaBq2xkYXA6Ly9zc3BkaXIubWFuYWdlZC5lbnRydXN0LmNvbS9jbj1XaW5Db21i
aW5lZDEsb3U9RW50cnVzdCUyME1hbmFnZWQlMjBTZXJ2aWNlcyUyMFNTUCUyMENB
LG91PUNlcnRpZmljYXRpb24lMjBBdXRob3JpdGllcyxvPUVudHJ1c3QsYz1VUz9j
ZXJ0aWZpY2F0ZVJldm9jYXRpb25MaXN0O2JpbmFyeYY0aHR0cDovL3NzcHdlYi5t
YW5hZ2VkLmVudHJ1c3QuY29tL0NSTHMvRU1TU1NQQ0ExLmNybDCBiqCBh6CBhKSB
gTB/MQswCQYDVQQGEwJVUzEQMA4GA1UEChMHRW50cnVzdDEiMCAGA1UECxMZQ2Vy
dGlmaWNhdGlvbiBBdXRob3JpdGllczEoMCYGA1UECxMfRW50cnVzdCBNYW5hZ2Vk
IFNlcnZpY2VzIFNTUCBDQTEQMA4GA1UEAxMHQ1JMNjY3MzArBgNVHRAEJDAigA8y
MDE0MDUwNzEzNDgzNlqBDzIwMTYwNjEyMTgxODM2WjAfBgNVHSMEGDAWgBTTzudb
iafNbJHGZzapWHIJ7OI58zAdBgNVHQ4EFgQUGIOcf6r7Z9wk+2/YuG5oTs7Qwk8w
CQYDVR0TBAIwADAZBgkqhkiG9n0HQQAEDDAKGwRWOC4xAwIEsDANBgkqhkiG9w0B
AQsFAAOCAQEASc+lZBbJWsHB2WnaBr8ZfBqpgS51Eh+wLchgIq7JHhVn+LagkR8C
XmvP57a0L/E+MRBqvH2RMqwthEcjXio2WIu/lyKZmg2go9driU6H3s89X8snblDF
1B+iL73vhkLVdHXgStMS8AHbm+3BW6yjHens1tVmKSowg1P/bGT3Z4nmamdY9oLm
9sCgFccthC1BQqtPv1XsmLshJ9vmBbYMsjKq4PmS0aLA59J01YMSq4U1kzcNS7wI
1/YfUrfeV+r+j7LKBgNQTZ80By2cfSalEqCe8oxqViAz6DsfPCBeE57diZNLiJmj
a2wWIBquIAXxvD8w2Bue7pZVeUHls5V5dA==
-----END CERTIFICATE-----
`

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

func beaconCertificate() *x509.Certificate {
	pemBlock, remainder := pem.Decode([]byte(beaconCertPEM))
	if len(remainder) > 0 {
		panic(fmt.Sprintf("have %d bytes left in PEM", len(remainder)))
	}
	certChain, err := x509.ParseCertificates(pemBlock.Bytes)
	if err != nil {
		panic(err.Error())
	}
	if len(certChain) != 1 {
		panic(fmt.Sprintf("have %d certificates in beacon PEM, not 1", len(certChain)))
	}
	return certChain[0]
}

func ValidateSignature(cert x509.Certificate, signed []byte, signature []byte) error {
	return cert.CheckSignature(x509.SHA512WithRSA, signed, signature)
}

func (d dirtyrecord) VerificationData() (signed, signature []byte, err error) {
	signature, err = hex.DecodeString(d.SignatureValue)
	if err != nil {
		return nil, nil, err
	}

	sigLimit := len(signature) - 1
	for i := 0; i <= sigLimit/2; i++ {
		signature[i], signature[sigLimit-i] = signature[sigLimit-i], signature[i]
	}

	b := new(bytes.Buffer)
	b.Grow(200)
	_, _ = b.WriteString(d.Version)
	binary.Write(b, binary.BigEndian, d.Frequency)
	binary.Write(b, binary.BigEndian, d.TimeStamp)
	seed, err := hex.DecodeString(d.SeedValue)
	if err != nil {
		return nil, nil, err
	}
	_, _ = b.Write(seed)
	prev, err := hex.DecodeString(d.PreviousOutputValue)
	if err != nil {
		return nil, nil, err
	}
	_, _ = b.Write(prev)
	binary.Write(b, binary.BigEndian, d.StatusCode)

	return b.Bytes(), signature, nil
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

	data, sig, err := drec.VerificationData()
	if err != nil {
		return Record{}, errors.New("Unable to extract verification data")
	}
	err = ValidateSignature(*beaconCertificate(), data, sig)
	if err != nil {
		return Record{}, errors.New("Unable to validate beacon signature")
	}
	if time.Now().Unix() - rec.TimeStamp.Unix() > 60 {
		return Record{}, errors.New("Beacon is stale")
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

// Rand saves the data pertinent to the random generator functions of the library
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

// NewUpdatedRand does the same as NewRand but ensures that the random numbers are generated always with the latest record
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

// SetSeed sets a new source for the randomness generator
func (r *Rand) SetSeed(n int64) {
	r.rand = rand.New(rand.NewSource(n))
	r.update = false
}

// Int randomly generates a new int from the given seed
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
