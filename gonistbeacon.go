package beacon

import(
  "net/http"
  "encoding/xml"
  "math/big"
  "time"
)

type Record struct{
  Version string
  Frequency int
  TimeStamp time.Time
  SeedValue big.Int
  PreviousOutputValue big.Int
  SignatureValue big.Int
  OutputValue big.Int
}

type dirtyrecord struct{
  version  string
  frequency int
  timeStamp int
  seedValue string
  previousOutputValue string
  signatureValue string
  outputValue string
}

func(c *http.Client) LastRecord() Record, err {
  r, err := c.Get("https://beacon.nist.gov/rest/record/last")
  if err != nil {
    return nil, err
  }

  var drec dirtyrecord
  drec, err := xml.Unmarshal(r.Body)
  if err != nil {
    return nil, err
  }

  var i big.Int
  rec := Record{
    Version : drec.version,
    Frequency : drec.frequency,
    TimeStamp : time.Unix(drec.timestamp, 0),
    SeedValue : i.String(drec.seedValue),
    PreviousOutputValue : i.String(drec.previousOutputValue),
    SignatureValue : i.String(drec.signatureValue),
    OutputValue : i.String(drec.outputValue),
  }
  return rec, nil
}
