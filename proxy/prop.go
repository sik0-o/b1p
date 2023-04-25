package proxy

import (
	"encoding/json"

	"github.com/sik0-o/b1p/bitwise"
)

const (
	PROP_DISABLED = 1 << uint(iota)
	PROP_BLACKLISTED
	PROP_WORTH
	PROP_RATE_LIMIT
)

type Prop bitwise.Flag

func (pp *Prop) Flag() *bitwise.Flag {
	f := bitwise.Flag(*pp)
	return &f
}

func (pp *Prop) Set(prop uint, on bool) {
	f := pp.Flag()
	f.Set(prop, on)
	*pp = Prop(*f)
}

func (pp *Prop) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"disabled":    bitwise.Flag(*pp).Has(PROP_DISABLED),
		"blacklisted": bitwise.Flag(*pp).Has(PROP_BLACKLISTED),
		"worth":       bitwise.Flag(*pp).Has(PROP_WORTH),
		"rate_limit":  bitwise.Flag(*pp).Has(PROP_RATE_LIMIT),
	})
}
