package address

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
)

const (
	// FamilyName - used to hash
	FamilyName = "LEOTEC"
)

var (
	NS = HexDigest(FamilyName)[:6]
)

// Range range
type Range int

// Space space
type Space int

const (
	//OfferHistoryStart start
	OfferHistoryStart Range = 0
	// OfferHistoryEnd end
	OfferHistoryEnd Range = 1
	// AssetStart asset start
	AssetStart Range = 1
	// AssetEnd asset end
	AssetEnd Range = 50
	// HoldingStart start
	HoldingStart Range = 50
	// HoldingEnd end
	HoldingEnd Range = 125
	// AccountStart start
	AccountStart Range = 125
	// AccountEnd end
	AccountEnd Range = 200
	// OfferStart start
	OfferStart Range = 200
	// OfferEnd end
	OfferEnd Range = 256
)

const (
	// SpaceAsset  asset
	SpaceAsset Space = iota
	// SpaceHolding Space
	SpaceHolding
	//SpaceAccount Space
	SpaceAccount
	// SpaceOffer Space
	SpaceOffer
	// SpaceOfferHistory Space
	SpaceOfferHistory
	// AssetSpaceOtherFamily other family
	AssetSpaceOtherFamily Space = 100
)

func hash512(identifier string) string {
	shaBytes := sha512.Sum512_256([]byte(identifier))
	return hex.EncodeToString(shaBytes[:])
}

func compress(address string, start Range, stop Range) string {
	digest := address[:4]
	addressHex, err := strconv.ParseInt(digest, 16, 64)
	if err != nil {
		log.Fatal("Invliad address" + err.Error())
	}
	return fmt.Sprintf("%.2x", addressHex%int64(stop-start)+int64(start))
}

func MakeOfferAccountAddress(offerId, account string) string {
	offerHash := hash512(offerId)
	accountHash := hash512(account)
	return NS + "00" + offerHash[:60] + compress(accountHash, 1, 256)
}

func MakeOfferHistoryAddress(offerId string) string {
	offerHash := hash512(offerId)
	return NS + "00" + offerHash[:60] + "00"
}

func MakeAssetAddress(assetId string) string {
	assetHash := hash512(assetId)
	return NS + compress(assetHash, AssetStart, AssetEnd) + assetHash[:62]
}

func MakeHoldingAddress(holdingId string) string {
	holdingHash := hash512(holdingId)
	return NS + compress(holdingHash, HoldingStart, HoldingEnd) + holdingHash[:62]
}

func MakeAccountAddress(accountId string) string {
	hash := hash512(accountId)
	compressedKey := compress(hash, AccountStart, AccountEnd)
	return NS + compressedKey + hash[:62]
}

func MakeOfferAddress(offerId string) string {
	offerHash := hash512(offerId)
	return NS + compress(offerHash, OfferStart, OfferEnd) + offerHash[:62]
}

func contains(num, start, end Range) bool {
	return start < num && num < end
}

// HexDigest calculates hex digest of a string
func HexDigest(text string) string {
	d := sha512.Sum512_256([]byte(text))
	return hex.EncodeToString(d[:])
}

func AddressOf(address string) Space {
	if address[:len(NS)] != NS {
		return AssetSpaceOtherFamily
	}
	infix, err := strconv.ParseInt(address[6:8], 16, 64)
	if err != nil {
		log.Fatal("invalid address")
	}
	if contains(Range(infix), OfferHistoryStart, OfferHistoryEnd) {
		return SpaceOfferHistory
	} else if contains(Range(infix), OfferStart, OfferEnd) {
		return SpaceOffer
	} else if contains(Range(infix), AssetStart, AssetEnd) {
		return SpaceAsset
	} else if contains(Range(infix), AccountStart, AccountEnd) {
		return SpaceAccount
	} else if contains(Range(infix), HoldingStart, HoldingEnd) {
		return SpaceHolding
	} else {
		return AssetSpaceOtherFamily
	}
}
