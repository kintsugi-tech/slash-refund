package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyAllowedTokens = []byte("AllowedTokens")
	// TODO: Determine the default value
	DefaultAllowedTokens string = "stake"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	allowedTokens string,
) Params {
	return Params{
		AllowedTokens: allowedTokens,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAllowedTokens,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAllowedTokens, &p.AllowedTokens, validateAllowedTokens),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateAllowedTokens(p.AllowedTokens); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateAllowedTokens validates the AllowedTokens param
func validateAllowedTokens(v interface{}) error {
	allowedTokens, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = allowedTokens

	return nil
}
