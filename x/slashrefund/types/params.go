package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyAllowedTokens         = []byte("AllowedTokens")
	KeyMaxEntries            = []byte("MaxEntries")
	DefaultAllowedTokens     = []string{"stake"}
	DefaultMaxEntries uint32 = 7
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAllowedTokens, &p.AllowedTokens, validateAllowedTokens),
		paramtypes.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
	}
}

// NewParams creates a new Params instance
func NewParams(
	allowedTokens []string,
	maxEntries uint32,
) Params {
	return Params{
		AllowedTokens: allowedTokens,
		MaxEntries: maxEntries,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultAllowedTokens,
		DefaultMaxEntries,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateAllowedTokens(p.AllowedTokens); err != nil {
		return err
	}

	if err := validateMaxEntries(p.MaxEntries); err != nil {
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
	allowedTokens, ok := v.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if len(allowedTokens) == 0 {
		return fmt.Errorf("no allowed token has been set")
	}

	// TODO: remove this line when multiple tokens are allowed
	if len(allowedTokens) > 1 {
		return fmt.Errorf("only one allowed token is currently accepted")
	}

	// ensure each denom is only registered one time.
	registered := make(map[string]bool)
	for _, token := range allowedTokens {
		if _, exists := registered[token]; exists {
			return fmt.Errorf("duplicate allowed tokens found: '%s'", token)
		}
		if err := validateAllowedToken(token); err != nil {
			return err
		}

		registered[token] = true
	}

	return nil
}

func validateAllowedToken(t interface{}) error {
	token, ok := t.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", t)
	}

	if strings.TrimSpace(token) == "" {
		return errors.New("allowed denoms cannot be blank")
	}
	return sdk.ValidateDenom(token)
}

func validateMaxEntries(i interface{}) error {
	number, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if number == 0 {
		return fmt.Errorf("max entries must be positive: %d", number)
	}

	return nil
}