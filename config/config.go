package config

import (
	"io/ioutil"

	"github.com/polygonledger/edn"
)

type Configuration struct {
	DelegateName   string
	PeerAddresses  []string
	NodePort       int
	WebPort        int
	DelgateEnabled bool
	CreateGenesis  bool
	//TODO
	Verbose bool
}

// func (c *Configuration) UnmarshalEDN(bs []byte) error {

// 	//input.Log = &c.Log
// 	err := edn.Unmarshal(bs, &input)
// 	c.Env = string(input.Env)
// 	return err
// }

func ReadConf(fname string) (*Configuration, error) {
	bs, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	var c Configuration
	err = edn.Unmarshal(bs, &c)
	return &c, err
}
