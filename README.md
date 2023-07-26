# Activeledger Testnet Bootstrapper

The intention for this software is to allow the creation of a testnet setup that
is rapidly deployable.
Current processes either require manually sending transactions to the network or
using Activeledger IDE.
The IDE is helpful, but requires a manual approach and many clicks.

In an environment where you may want to reset the network quickly and frequently
both processes can be cumbersome.
This software will quickly create a testnet using configurable predefined settings
and onboard your smart contracts. This allows you to wipe out a testnetwork,
which may be polluted from development work, and start fresh with minimal effort.

Define once, update sometimes, redeploy often.

## Requirements

**NOTE:** Currently this software does not check if Activeledger is installed,
this is a future feature.

Activeledger is a Node application, you will need to install Node and
Activeledger to use this software.

See the documentation [here](https://github.com/activeledger/activeledger)

You also need to make sure you have the 2 required files and smartcontracts folder.
Defaults are included in this repo.

The 2 files are `config.json` and `contract-manifest.json`.
The folder is `smartcontracts`, which should contain the smart contracts listed
in the manifest file.

The name of the manifest and smartcontract folder/path is configurable in the
`config.json` file.

## Quickstart

1. Modify `config.json` and `contract-manifest.json` as needed.
2. Make sure your smart contracts are in the directory you specified in `config.json`,
   or path in the manifest
3. Run the deployer: `./deployer` (may require `chmod +x deployer` first)

## Config

The following is a sample config file, the software expects `config.json`
to be in the same local directory.

### Sample Config

```json
{
  "identity": "someiden",
  "namespace": "somenamespace",
  "contractDir": "smartcontracts",
  "contractManifest": "contract-manifest.json",
  "setupDataSaveFile": "setup-output.json",
  "testnetFolder": "sometestnet"
}
```

## Manifest

The manifest file contains a list of the smartcontracts you want to onboard to
the testnet.

You can also add ones that you don't yet want uploaded and use the `exclude: true`
pairing to ignore them.

### Sample Manifest

```json
{
  "contracts": [
    {
      "name": "contractname",
      "path": "smartcontracts/contract.ts",
      "version": "0.0.1",
      "exclude": false
    }
  ]
}
```

## Todo

A useful additonal feature would be to allow updating contracts or updating them.
The software would need to check for the existence of an output file, or perhaps
should maintain a hidden file to keep track of things.
Currently it will error on creating a testnet folder if one already exists.
For this feature it should check the hashes of the contracts to find updated ones,
update those, and upload any new ones.

## Changelog

### [Unreleased]

#### Features

- Ability to update contracts and upload new ones based on hidden file
- Check if Activeledger is installed (requires changes in Activeledger first)
- Ability to attempt to install Activeledger if it isn't installed
- Check if node/npm is installed

### [1.0.0] - 2023-06-30

#### Added

- Created the initial software
- Runs `activeledger --testnet`
- Onboards configured identity and namespace
- Uploads smartcontracts as definied in manifest file
