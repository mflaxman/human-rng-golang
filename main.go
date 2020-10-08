package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

func main() {

	firstWordsPtr := flag.String("firstWords", "", "REQUIRED: First words of mnemonic (23 words by default)")
	testnetBoolPtr := flag.Bool("testnet", false, "Use testnet (default is mainnet)")
	verbosityBoolPtr := flag.Bool("verbose", false, "Verbose printout (default is quiet)")
	checksumIntPtr := flag.Int("checksum", 0, "EXPERTS ONLY: Which checksum word to append, using a 0-index.")

	// https://stackoverflow.com/questions/53082613/set-the-order-for-output-by-flag-printdefaults
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Println("HumanRNG - Don't Trust Your Random Number Generator\n")
		fmt.Println("Usage:")
		fmt.Println("  go run *.go -firstWords=\"add bag cat ...\"\n")
		order := []string{"firstWords", "testnet", "verbose", "checksum"}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			fmt.Printf("-%s\n", flag.Name)
			fmt.Printf("  %s\n", flag.Usage)
		}
	}

	flag.Parse()

	if *checksumIntPtr > 0 {
		fmt.Println("WARNING!")
		fmt.Println("You have selected a checksum # greater than 0.")
		fmt.Println("This should only be attempted by experts users.")
		fmt.Println("If you do not know what you're doing, DO NOT CONTINUE.")
		fmt.Println(strings.Repeat("-", 80))
	}

	if *verbosityBoolPtr == true {
		fmt.Println("Input Flags (DEBUG ONLY):")
		fmt.Println("checksum word number to pick:", *checksumIntPtr)
		fmt.Println("testnet:", *testnetBoolPtr)
		fmt.Println("mnemonicFirstWords:", *firstWordsPtr)
		fmt.Println(strings.Repeat("-", 80))
	}

	firstWords := strings.TrimSpace(*firstWordsPtr)

	if len(firstWords) == 0 {
		fmt.Println("ERROR: You didn't supply supply firstWords of your mnemonic\n")
		flag.Usage()
		os.Exit(1)
	}

	invalidMnemonicWords := GetInvalidMnemonicWords(firstWords)
	if len(invalidMnemonicWords) > 0 {
		errStr := fmt.Sprintf("%s", invalidMnemonicWords)
		fmt.Println("Invalid BIP39 Mnemonic Word(s): " + errStr)
		os.Exit(1)
	}

	validChecksums, err := FindAllChecksumWords(firstWords)
	if err != nil {
		fmt.Println("Could not find valid checksum words")
		os.Exit(1)
	}

	// Append checksum word (default to 0th word)
	if *checksumIntPtr > len(validChecksums)-1 {
		fmt.Printf("Not enough valid checksum words to append (0-indexed) result #%d \n", *checksumIntPtr)
		os.Exit(1)
	}
	checksumWordToUse := validChecksums[*checksumIntPtr]
	mnemonic := firstWords + " " + checksumWordToUse

	passphrase := ""
	seed := bip39.NewSeed(mnemonic, passphrase)

	// p2wsh version bytes
	var networkName string
	var network chaincfg.Params
	var derivationPath string
	// https://github.com/satoshilabs/slips/blob/master/slip-0132.md
	pubkeyBytesToUse := [4]byte{}
	privkeyBytesToUse := [4]byte{}

	if *testnetBoolPtr == true {
		network = chaincfg.TestNet3Params
		networkName = "testnet3"
		derivationPath = "m/48'/1'/0'/2'"
		privkeyBytesToUse = [4]byte{0x02, 0x57, 0x50, 0x48} // Vpriv
		pubkeyBytesToUse = [4]byte{0x02, 0x57, 0x54, 0x83}  // Vpub

	} else {
		network = chaincfg.MainNetParams
		networkName = "mainnet"
		derivationPath = "m/48'/0'/0'/2'"
		privkeyBytesToUse = [4]byte{0x02, 0xaa, 0x7a, 0x99} // Zpriv
		pubkeyBytesToUse = [4]byte{0x02, 0xaa, 0x7e, 0xd3}  // Zpub
	}

	derivationPathSpecter := strings.Replace(
		strings.ReplaceAll(derivationPath, "'", "h"),
		"m/", "", 1,
	)

	masterXpriv, err := hdkeychain.NewMaster(seed, &network)
	if err != nil {
		fmt.Println("Couldn't create seed", err)
		os.Exit(1)
	}

	// TODO: this should be handled by the btcd library (when supported) and not my custom code
	xfp, err := RootXPrivToFingerprint(masterXpriv)
	if err != nil {
		fmt.Println("Error calculating fingerpint", err)
		os.Exit(1)
	}

	childXpriv, err := DeriveChildKeyFromPath(masterXpriv, derivationPath)
	if err != nil {
		fmt.Println("Error deriving child private key", err)
		os.Exit(1)
	}

	childXpub, err := childXpriv.Neuter()
	if err != nil {
		fmt.Println("Error deriving child public key", err)
		os.Exit(1)
	}

	childZpriv, err := Slip132Encode(childXpriv, privkeyBytesToUse)
	if err != nil {
		fmt.Println("Error encoding SLIP132 version bytes on private key")
		os.Exit(1)
	}

	childZpub, err := Slip132Encode(childXpub, pubkeyBytesToUse)
	if err != nil {
		fmt.Println("Error encoding SLIP132 version bytes on public key")
		os.Exit(1)
	}

	// Output
	fmt.Println("SECRET INFO:")
	fmt.Println("Full mnemonic (with checksum word): ", mnemonic)
	fmt.Println("Full mnemonic length (# words): ", len(strings.Split(mnemonic, " ")))
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("PUBLIC INFO:")
	fmt.Println("SLIP132 Extended Pubkey:", childZpub)
	fmt.Println("Root Fingerprint:", xfp)
	fmt.Println("Network:", networkName)
	fmt.Println("Derivation Path:", derivationPath)
	fmt.Println("Specter-Desktop Input Format:")
	fmt.Printf("  [%s/%s]%s\n", xfp, derivationPathSpecter, childZpub)
	fmt.Println(strings.Repeat("-", 80))

	if *verbosityBoolPtr == true {
		fmt.Println("  Advanced Details:")
		fmt.Println("  childXpub:", childXpub)
		fmt.Println("  childXpriv:", childXpriv)
		fmt.Println("  childZpriv:", childZpriv)
		fmt.Println(" ", len(validChecksums), "valid checksums:", validChecksums)
	}

}
