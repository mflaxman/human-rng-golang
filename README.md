
# THIS REPOSITORY COMES WITH ZERO GUARANTEES! USE AT YOUR OWN RISK!

#### Always perform sensitive operations on an airgapped computer and securely wipe it after.

## Quickstart:

Basic:
```bash
g$ go run *.go -firstWords="zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo"
SECRET INFO:
Full mnemonic (with checksum word):  zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo abstract
Full mnemonic length (# words):  12
--------------------------------------------------------------------------------
PUBLIC INFO:
SLIP132 Extended Pubkey: Zpub74h4HCSddD3n8wX1JMQjfb12fQeLAeFAbtN3zVqyTSj387zTdnYR4GpuA2giJMEvp5nJ7L48uGmVawyfV3pkHY5d6rMVtXFCk6J3Aw81r5c
Root Fingerprint: f7d04090
Network: mainnet
Derivation Path: m/48'/0'/0'/2'
--------------------------------------------------------------------------------
```

Pass `-testnet` flag to run on testnet
```bash
$ go run *.go -firstWords="zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo" -testnet
SECRET INFO:
Full mnemonic (with checksum word):  zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo abstract
Full mnemonic length (# words):  12
--------------------------------------------------------------------------------
PUBLIC INFO:
SLIP132 Extended Pubkey: Vpub5myBTZx9knAMMmypC51gfX2RXnS8rJWKNMUTCwZNxwajq2tKcrj15SPpbJFdYmG5EgUVDA3Gt5UQgUDoCqc5XaYN3iZNZWhFjH9ScbVPnHh
Root Fingerprint: f7d04090
Network: testnet3
Derivation Path: m/48'/1'/0'/2'
--------------------------------------------------------------------------------
```

Confirm your output matches [SeedPicker](https://seedpicker.net/calculator/last-word.html) and [human-rng-electrum](https://github.com/mflaxman/human-rng-electrum) before trusting it with funds.

#### Build from Source

```bash
$ go build *.go
```

Then transfer the resulting `main` file to your airgapped machine and run the following:
```bash
$ ./main -firstWords="zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo"
SECRET INFO:
Full mnemonic (with checksum word):  zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo buddy
Full mnemonic length (# words):  24
--------------------------------------------------------------------------------
PUBLIC INFO:
SLIP132 Extended Pubkey: Zpub74sb5KB3Ak1RwabGr8SHQnMTkd2mC3boVDgPf1jBFNxcXh7Nx4KV3XakPDtWLN5RpszdM7qcBN4wm7xreh8Ys2xYUBqQ9GtkTN8h5kRVecc
Root Fingerprint: 669dce62
Network: mainnet
Derivation Path: m/48'/0'/0'/2'
--------------------------------------------------------------------------------

```
