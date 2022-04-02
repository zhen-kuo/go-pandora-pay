package wizard

import (
	"pandora-pay/addresses"
	"pandora-pay/blockchain/transactions/transaction"
	"pandora-pay/blockchain/transactions/transaction/transaction_simple"
	"pandora-pay/blockchain/transactions/transaction/transaction_simple/transaction_simple_extra"
	"pandora-pay/blockchain/transactions/transaction/transaction_simple/transaction_simple_parts"
	"pandora-pay/blockchain/transactions/transaction/transaction_type"
	"pandora-pay/cryptography"
	"pandora-pay/helpers"
)

func CreateSimpleTx(transfer *WizardTxSimpleTransfer, validateTx bool, statusCallback func(string)) (tx2 *transaction.Transaction, err error) {

	dataFinal, err := transfer.Data.getData()
	if err != nil {
		return
	}

	spaceExtra := 0

	var txScript transaction_simple.ScriptType
	var extraFinal transaction_simple_extra.TransactionSimpleExtraInterface

	switch transfer.Extra.(type) {
	case nil:
	}

	txBase := &transaction_simple.TransactionSimple{
		TxScript:    txScript,
		DataVersion: transfer.Data.getDataVersion(),
		Data:        dataFinal,
		Nonce:       transfer.Nonce,
		Extra:       extraFinal,
		Vin:         make([]*transaction_simple_parts.TransactionSimpleInput, len(transfer.Vin)),
		Vout:        make([]*transaction_simple_parts.TransactionSimpleOutput, len(transfer.Vout)),
	}

	tx := &transaction.Transaction{
		Version:                  transaction_type.TX_SIMPLE,
		SpaceExtra:               uint64(spaceExtra),
		TransactionBaseInterface: txBase,
	}
	statusCallback("Transaction Created")

	extraBytes := len(transfer.Vin) * cryptography.SignatureSize
	fee := setFee(tx, extraBytes, transfer.Fee.Clone(), true)
	if err = helpers.SafeUint64Add(&transfer.Vin[0].Amount, fee); err != nil {
		return nil, err
	}

	statusCallback("Transaction Fee set")

	statusCallback("Transaction Signing...")
	for i, vin := range txBase.Vin {
		var privateKey *addresses.PrivateKey
		if privateKey, err = addresses.NewPrivateKey(transfer.Vin[i].Key); err != nil {
			return nil, err
		}
		if vin.Signature, err = privateKey.Sign(tx.SerializeForSigning()); err != nil {
			return nil, err
		}
	}
	statusCallback("Transaction Signed")

	if err = bloomAllTx(tx, statusCallback); err != nil {
		return
	}
	return tx, nil
}
