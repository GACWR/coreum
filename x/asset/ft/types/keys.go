package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"

	"github.com/CoreumFoundation/coreum/v2/pkg/store"
)

const (
	// ModuleName defines the module name.
	ModuleName = "assetft"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

// Store key prefixes.
var (
	// TokenKeyPrefix defines the key prefix for the fungible token.
	TokenKeyPrefix = []byte{0x01}
	// SymbolKeyPrefix defines the key prefix for the fungible token by Symbol.
	SymbolKeyPrefix = []byte{0x02}
	// FrozenBalancesKeyPrefix defines the key prefix to track frozen balances.
	FrozenBalancesKeyPrefix = []byte{0x03}
	// GlobalFreezeKeyPrefix defines the key prefix to track global freezing of a Fungible Token.
	GlobalFreezeKeyPrefix = []byte{0x04}
	// WhitelistedBalancesKeyPrefix defines the key prefix to track whitelisted balances.
	WhitelistedBalancesKeyPrefix = []byte{0x05}
	// PendingTokenUpgradeKeyPrefix defines the key prefix for the pending version upgrades.
	PendingTokenUpgradeKeyPrefix = []byte{0x06}
	// TokenUpgradeStatusesKeyPrefix defines the key prefix for the fungible token upgrade statuses.
	TokenUpgradeStatusesKeyPrefix = []byte{0x07}
)

// CreateTokenKey creates the key for the fungible token.
func CreateTokenKey(issuer sdk.AccAddress, subunit string) []byte {
	return store.JoinKeys(CreateIssuerTokensPrefix(issuer), []byte(strings.ToLower(subunit)))
}

// CreatePendingTokenUpgradeKey creates the key for the fungible token version upgrade.
func CreatePendingTokenUpgradeKey(denom string) []byte {
	return store.JoinKeys(PendingTokenUpgradeKeyPrefix, []byte(denom))
}

// CreateIssuerTokensPrefix creates the key for the fungible token issued by account.
func CreateIssuerTokensPrefix(issuer sdk.AccAddress) []byte {
	return store.JoinKeys(TokenKeyPrefix, address.MustLengthPrefix(issuer))
}

// CreateSymbolKey creates the key for a ft symbol.
func CreateSymbolKey(addr []byte, symbol string) []byte {
	return store.JoinKeys(store.JoinKeys(SymbolKeyPrefix, addr), []byte(symbol))
}

// CreateFrozenBalancesKey creates the key for an account's frozen balances.
func CreateFrozenBalancesKey(addr []byte) []byte {
	return store.JoinKeys(FrozenBalancesKeyPrefix, address.MustLengthPrefix(addr))
}

// CreateGlobalFreezeKey creates the key for fungible token global freeze key.
func CreateGlobalFreezeKey(denom string) []byte {
	return store.JoinKeys(GlobalFreezeKeyPrefix, []byte(denom))
}

// CreateWhitelistedBalancesKey creates the key for an account's whitelisted balances.
func CreateWhitelistedBalancesKey(addr []byte) []byte {
	return store.JoinKeys(WhitelistedBalancesKeyPrefix, address.MustLengthPrefix(addr))
}

// CreateTokenUpgradeStatusesKey creates the key for the fungible token upgrade statuses.
func CreateTokenUpgradeStatusesKey(denom string) []byte {
	return store.JoinKeys(TokenUpgradeStatusesKeyPrefix, []byte(denom))
}

// AddressFromBalancesStore returns an account address from a balances prefix
// store. The key must not contain the prefix BalancesPrefix as the prefix store
// iterator discards the actual prefix.
//
// If invalid key is passed, AddressFromBalancesStore returns ErrInvalidKey.
func AddressFromBalancesStore(key []byte) (sdk.AccAddress, error) {
	if len(key) == 0 {
		return nil, ErrInvalidKey
	}
	addrLen := key[0]
	bound := int(addrLen)
	if len(key)-1 < bound {
		return nil, ErrInvalidKey
	}
	return key[1 : bound+1], nil
}
