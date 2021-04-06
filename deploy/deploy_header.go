package deploy

import (
	"bytes"
	"encoding/hex"
	"fmt"
	cl "github/casper-go/clvalue"
	"strconv"
	"time"
)

type Header struct {
	Account      *cl.PublicKey `json:"-"`
	BodyHash     []byte        `json:"-"`
	ChainName    string        `json:"chain_name"`
	Dependencies [][]byte      `json:"-"`
	GasPrice     uint64        `json:"gas_price"`
	Timestamp    time.Time     `json:"-"`
	TTL          uint64        `json:"-"`

	JSONHeader
}

type JSONHeader struct {
	JSONAccount      string   `json:"account"`
	JSONBodyHash     string   `json:"body_hash"`
	JSONDependencies []string `json:"dependencies"`
	JSONTimestamp    string   `json:"timestamp"`
	JSONTTL          string   `json:"ttl"`
}

func NewHeader(publicKey *cl.PublicKey, bodyHash []byte, params *Params) *Header {
	d, _ := time.ParseDuration("-8h")
	return &Header{
		Account:      publicKey,
		BodyHash:     bodyHash,
		ChainName:    params.chainName,
		Dependencies: params.dependencies,
		GasPrice:     params.gasPrice,
		Timestamp:    params.timestamp,
		TTL:          params.ttl,

		JSONHeader: JSONHeader{
			JSONAccount:      hex.EncodeToString(publicKey.ToBytes()),
			JSONBodyHash:     hex.EncodeToString(bodyHash),
			JSONDependencies: []string{},
			//JSONTimestamp:    params.timestamp.Format("2006-01-02T15:04:05.000Z"),
			JSONTimestamp:    params.timestamp.Add(d).Format("2006-01-02T15:04:05.000Z"),
			//JSONTimestamp: "2021-04-06T08:49:54.323Z",
			JSONTTL: strconv.Itoa(int(params.ttl)) + "ms",
		},
	}
}

func (h *Header) ToBytes() []byte {
	t := h.Timestamp.UnixNano() / int64(time.Millisecond)
	//t := 1617698994323

	fmt.Printf("timestamp:%d\n", t)
	return bytes.Join([][]byte{
		h.Account.ToBytes(),
		cl.ToBytesU64(uint64(t)),
		cl.ToBytesU64(h.TTL),
		cl.ToBytesU64(h.GasPrice),
		cl.ToBytesBytesArray(h.BodyHash),
		cl.ToByteSlice(h.Dependencies),
		cl.ToBytesString(h.ChainName),
	}, []byte{})
}
