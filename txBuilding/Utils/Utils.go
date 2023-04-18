package Utils

import (
	"Salvionied/apollo/serialization"
	"Salvionied/apollo/serialization/TransactionOutput"
	"Salvionied/apollo/serialization/UTxO"
	"Salvionied/apollo/txBuilding/Backend/Base"
	"encoding/hex"
	"log"

	"github.com/Salvionied/cbor/v2"
)

func Contains[T UTxO.Container[any]](container []T, contained T) bool {
	for _, c := range container {
		if c.EqualTo(contained) {
			return true
		}
	}
	return false
}

func MinLovelacePostAlonzo(output TransactionOutput.TransactionOutput, context Base.ChainContext) int64 {
	constantOverhead := 160
	amt := output.GetValue()
	if amt.Coin == 0 {
		amt.Coin = 1_000_000
	}
	tmp_out := TransactionOutput.TransactionOutput{
		IsPostAlonzo: true,
		PostAlonzo: TransactionOutput.TransactionOutputAlonzo{
			Address:   output.GetAddress(),
			Amount:    output.GetValue(),
			Datum:     output.GetDatum(),
			ScriptRef: output.GetScriptRef(),
		},
	}
	encoded, err := cbor.Marshal(tmp_out)
	if err != nil {
		log.Fatal(err)
	}
	return int64((constantOverhead + len(encoded)) * context.GetProtocolParams().GetCoinsPerUtxoByte())
}

func ToCbor(x interface{}) string {
	bytes, err := cbor.Marshal(x)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(bytes)
}

func Fee(context Base.ChainContext, txSize int, steps int64, mem int64) int64 {
	pm := context.GetProtocolParams()
	fee := int64(txSize*2*pm.MinFeeCoefficient +
		pm.MinFeeConstant +
		int(float32(steps)*pm.PriceStep) +
		int(float32(mem)*pm.PriceMem) + 200_000)
	return fee
}

func Copy[T serialization.Clonable[T]](input []T) []T {
	res := make([]T, 0)
	for _, value := range input {
		res = append(res, value.Clone())
	}
	return res
}
