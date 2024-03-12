package types

import (
	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// staking message types
const (
	TypeMsgUndelegate                = "begin_unbonding"
	TypeMsgCancelUnbondingDelegation = "cancel_unbond"
	TypeMsgEditValidator             = "edit_validator"
	TypeMsgCreateValidator           = "create_validator"
	TypeMsgDelegate                  = "delegate"
	TypeMsgBeginRedelegate           = "begin_redelegate"
)

// const (
// 	rpcURL      = "http://127.0.0.1:8545"
// 	contractABI = `[
//     {
//       "inputs": [
//         {
//           "internalType": "address",
//           "name": "owner",
//           "type": "address"
//         }
//       ],
//       "name": "balanceOf",
//       "outputs": [
//         {
//           "internalType": "uint256",
//           "name": "",
//           "type": "uint256"
//         }
//       ],
//       "stateMutability": "view",
//       "type": "function"
//     },
//   ]`
// )

var (
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateValidator)(nil)
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ sdk.Msg                            = &MsgEditValidator{}
	_ sdk.Msg                            = &MsgDelegate{}
	_ sdk.Msg                            = &MsgUndelegate{}
	_ sdk.Msg                            = &MsgBeginRedelegate{}
	_ sdk.Msg                            = &MsgCancelUnbondingDelegation{}
)

// NewMsgCreateValidator creates a new MsgCreateValidator instance.
// Delegator address and validator address are the same.
func NewMsgCreateValidator(
	valAddr sdk.ValAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	selfDelegation sdk.Coin, description Description, commission CommissionRates, minSelfDelegation math.Int,
) (*MsgCreateValidator, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateValidator{
		Description:       description,
		DelegatorAddress:  sdk.AccAddress(valAddr).String(),
		ValidatorAddress:  valAddr.String(),
		Pubkey:            pkAny,
		Value:             selfDelegation,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateValidator) Type() string { return TypeMsgCreateValidator }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delegator, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	addrs := []sdk.AccAddress{delegator}
	valAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)

	valAccAddr := sdk.AccAddress(valAddr)
	if !delegator.Equals(valAccAddr) {
		addrs = append(addrs, valAccAddr)
	}

	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// func Convert2EthAddress(address string, bech32Prefix string) (string, error) {
// 	data, err := sdk.GetFromBech32(address, bech32Prefix)
// 	if err != nil {
// 		return "", err
// 	}

// 	ethAddress := "0x" + hex.EncodeToString(data[len(data)-20:])
// 	return ethAddress, nil
// }

// func getBalances(address string) (int64, int64, error) {
// 	client, err := ethclient.Dial(rpcURL)

// 	if err != nil {
// 		return 0, 0, sdkerrors.ErrInvalidRequest.Wrapf("RPC is not valid")
// 	}

// 	ethAddress, err := Convert2EthAddress(address, "nxq")
// 	if err != nil {
// 		return 0, 0, sdkerrors.ErrInvalidAddress.Wrapf("invalid address: %s", err)
// 	}

// 	walletAddress := common.HexToAddress((ethAddress))
// 	tokenBalance, err := client.BalanceAt(context.Background(), walletAddress, nil)
// 	if err != nil {

// 	}

// 	contractAddress := common.HexToAddress((nxqconfig.ContractAddress))
// 	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
// 	if err != nil {
// 		return tokenBalance.Int64(), 0, sdkerrors.ErrInvalidRequest.Wrapf("invalid abi: %s", err)
// 	}
// 	data, err := parsedABI.Pack("balanceOf", walletAddress)
// 	if err != nil {
// 		return tokenBalance.Int64(), 0, sdkerrors.ErrInvalidRequest.Wrapf("Failed to pack data for balanceOf: %s", err)
// 	}
// 	callMsg := ethereum.CallMsg{
// 		To:   &contractAddress,
// 		Data: data,
// 	}
// 	result, err := client.CallContract(context.Background(), callMsg, nil)
// 	if err != nil {
// 		return tokenBalance.Int64(), 0, sdkerrors.ErrInvalidRequest.Wrapf("Failed to call contract: %s", err)
// 	}

// 	if len(result) == 0 {
// 		return tokenBalance.Int64(), 0, sdkerrors.ErrInvalidAddress.Wrapf("Invalid address: %s", err)
// 	}

// 	results, err := parsedABI.Unpack("balanceOf", result)
// 	if err != nil {
// 		return tokenBalance.Int64(), 0, sdkerrors.ErrInvalidRequest.Wrapf("Failed to unpack result: %s", err)
// 	}

// 	nftBalance := results[0].(*big.Int)
// 	return tokenBalance.Int64(), nftBalance.Int64(), nil
// }

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures both non-empty and valid
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if !sdk.AccAddress(valAddr).Equals(delAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator address is invalid")
	}

	// tokenBalance, nftBalance, err := getBalances(msg.DelegatorAddress)
	// if err != nil {
	// 	return err
	// }

	// if nftBalance < nxqconfig.MinValidatorNFTBalance {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Insufficient NFT balance")
	// }

	// if tokenBalance < nxqconfig.MinValidatorToken {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Insufficient Token balance")
	// }

	if msg.Pubkey == nil {
		return ErrEmptyValidatorPubKey
	}

	if !msg.Value.IsValid() || !msg.Value.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid delegation amount")
	}

	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Commission == (CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}

	if err := msg.Commission.Validate(); err != nil {
		return err
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return ErrSelfDelegationBelowMinimum
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pubKey)
}

// NewMsgEditValidator creates a new MsgEditValidator instance
//
//nolint:interfacer
func NewMsgEditValidator(valAddr sdk.ValAddress, description Description, newRate *sdk.Dec, newMinSelfDelegation *math.Int) *MsgEditValidator {
	return &MsgEditValidator{
		Description:       description,
		CommissionRate:    newRate,
		ValidatorAddress:  valAddr.String(),
		MinSelfDelegation: newMinSelfDelegation,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEditValidator) Type() string { return TypeMsgEditValidator }

// GetSigners implements the sdk.Msg interface.
func (msg MsgEditValidator) GetSigners() []sdk.AccAddress {
	valAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgEditValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgEditValidator) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.MinSelfDelegation != nil && !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(sdk.OneDec()) || msg.CommissionRate.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "commission rate must be between 0 and 1 (inclusive)")
		}
	}

	return nil
}

// NewMsgDelegate creates a new MsgDelegate instance.
//
//nolint:interfacer
func NewMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgDelegate {
	return &MsgDelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDelegate) Type() string { return TypeMsgDelegate }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	return nil
}

// NewMsgBeginRedelegate creates a new MsgBeginRedelegate instance.
//
//nolint:interfacer
func NewMsgBeginRedelegate(
	delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, amount sdk.Coin,
) *MsgBeginRedelegate {
	return &MsgBeginRedelegate{
		DelegatorAddress:    delAddr.String(),
		ValidatorSrcAddress: valSrcAddr.String(),
		ValidatorDstAddress: valDstAddr.String(),
		Amount:              amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgBeginRedelegate) Type() string { return TypeMsgBeginRedelegate }

// GetSigners implements the sdk.Msg interface
func (msg MsgBeginRedelegate) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorSrcAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid source validator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorDstAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid destination validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	return nil
}

// NewMsgUndelegate creates a new MsgUndelegate instance.
//
//nolint:interfacer
func NewMsgUndelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgUndelegate {
	return &MsgUndelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUndelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUndelegate) Type() string { return TypeMsgUndelegate }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUndelegate) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUndelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	return nil
}

// NewMsgCancelUnbondingDelegation creates a new MsgCancelUnbondingDelegation instance.
//
//nolint:interfacer
func NewMsgCancelUnbondingDelegation(delAddr sdk.AccAddress, valAddr sdk.ValAddress, creationHeight int64, amount sdk.Coin) *MsgCancelUnbondingDelegation {
	return &MsgCancelUnbondingDelegation{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
		CreationHeight:   creationHeight,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgCancelUnbondingDelegation) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCancelUnbondingDelegation) Type() string { return TypeMsgCancelUnbondingDelegation }

// GetSigners implements the sdk.Msg interface.
func (msg MsgCancelUnbondingDelegation) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	return []sdk.AccAddress{delegator}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgCancelUnbondingDelegation) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCancelUnbondingDelegation) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid amount",
		)
	}

	if msg.CreationHeight <= 0 {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid height",
		)
	}

	return nil
}
