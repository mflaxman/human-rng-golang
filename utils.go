package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	bip39 "github.com/tyler-smith/go-bip39"
)

func FindAllChecksumWords(firstWords string) ([]string, error) {
	// Extra safety:
	firstWords = strings.TrimSpace(firstWords)

	var validChecksums []string

	wordlist := bip39.GetWordList()
	for _, s := range wordlist {
		// fmt.Println(i, "attempt", s)
		_, err := bip39.EntropyFromMnemonic(firstWords + " " + s)
		if err == nil {
			validChecksums = append(validChecksums, s)
		}
	}
	if len(validChecksums) > 0 {
		return validChecksums, nil
	}

	return validChecksums, bip39.ErrInvalidMnemonic
}

func DeriveChildKeyFromPath(master *hdkeychain.ExtendedKey, path string) (child *hdkeychain.ExtendedKey, e error) {
	// This was helpful:
	// https://github.com/btcsuite/btcutil/issues/112

	parts := strings.Split(path[2:], "/")
	curr := master
	for _, p := range parts {
		num, err := strconv.Atoi(p[:len(p)-1])
		if err != nil {
			curr.Zero()
			return curr, err
		}
		hardened := p[len(p)-1 : len(p)]
		if hardened == "H" || hardened == "h" || hardened == "'" {
			curr, err = curr.Child(uint32(num) + hdkeychain.HardenedKeyStart)
			if err != nil {
				curr.Zero()
				return curr, err
			}
		} else {
			curr, err = curr.Child(0)
			if err != nil {
				curr.Zero()
				return curr, err
			}

		}
	}
	return curr, nil
}

func i32tob(val uint32) []byte {
	b := make([]byte, 4)
	// https://golang.org/pkg/encoding/binary/
	binary.BigEndian.PutUint32(b, val)
	return b
}

func Slip132Encode(xkey *hdkeychain.ExtendedKey, OutputVersionBytes [4]byte) (encoded *hdkeychain.ExtendedKey, e error) {

	var serializedKey []byte
	if xkey.IsPrivate() {
		privkey, err := xkey.ECPrivKey()
		if err != nil {
			return nil, err
		}
		serializedKey = privkey.Serialize()

	} else {

		pubkey, err := xkey.ECPubKey()
		if err != nil {
			return nil, err
		}
		serializedKey = pubkey.SerializeCompressed()
	}

	return hdkeychain.NewExtendedKey(OutputVersionBytes[:], serializedKey, xkey.ChainCode(), i32tob(xkey.ParentFingerprint()),
		xkey.Depth(), xkey.ChildIndex(), xkey.IsPrivate()), nil
}

func GetInvalidMnemonicWords(m string) (nvalidWords []string) {
	parts := strings.Split(strings.TrimSpace(m), " ")

	wordList := bip39.GetWordList()
	wordMap := map[string]int{}
	for i, v := range wordList {
		wordMap[v] = i
	}

	invalidWords := []string{}

	// FIXME
	for _, word := range parts {
		_, ok := wordMap[word]
		if !ok {
			invalidWords = append(invalidWords, word)
		}
	}

	return invalidWords

}

func RootXPrivToFingerprint(masterXPriv *hdkeychain.ExtendedKey) (string, error) {
	if masterXPriv.Depth() != 0 {
		return "", fmt.Errorf("No reason to calculate root fingerprint of non-root key: %s", masterXPriv.Depth())
	}
	masterXpub, err := masterXPriv.ECPubKey()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", btcutil.Hash160(masterXpub.SerializeCompressed())[:4]), nil
}
