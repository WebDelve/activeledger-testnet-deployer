<img src="assets/webdelve-logo.png" alt="WebDelve"/><br/>
[WebDelve - Purpose Drive Software](https://webdelve.co)

*** 

<img src="https://github.com/activeledger/activeledger/blob/master/docs/assets/Asset-23.png" alt="Activeledger" width="250"/><br/>
This project is built to work with Activeledger<br/>
[Activeledger on GitHub]( https://github.com/activeledger/activeledger )<br/>
[Activeledger Website](https://activeledger.io)

***

# Activeledger Testnet Bootstrapper & Contract uploader
![GitHub](https://img.shields.io/github/license/WebDelve/activeledger-testnet-deployer)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/WebDelve/activeledger-testnet-deployer)
![GitHub release (with filter)](https://img.shields.io/github/v/release/WebDelve/activeledger-testnet-deployer)
![GitHub all releases](https://img.shields.io/github/downloads/WebDelve/activeledger-testnet-deployer/total)

The intention for this software is to allow the creation of a testnet setup that
is rapidly deployable.
Current processes either require manually sending transactions to the network or
using [ Activeledger IDE ](https://github.com/activeledger/ide).
The IDE is helpful, but requires a manual approach and many clicks.

In an environment where you may want to reset the network quickly and frequently
both processes can be cumbersome.
This software will quickly create a testnet using configurable predefined settings
and onboard your smart contracts. This allows you to wipe out a testnetwork,
which may be polluted from development work, and start fresh with minimal effort.

Define once, update sometimes, redeploy often.

## Quick links
[Requirements](#requirements)<br/>
[Quickstart](#quickstart)<br/>
[Installation](#installation)<br/>
[Config](#config)<br/>
[Manifest](#manifest)<br/>
[Setup Output](#setup-output)<br/>
[Todo](#todo)<br/>
[Changelog](#changelog)

## Requirements

**NOTE:** Currently this software does not check if Activeledger is installed,
this is a future feature.

Activeledger is a Node application, you will need to install Node and
Activeledger to use this software.

See the Activeledger documentation [here](https://github.com/activeledger/activeledger)

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

There are two flags available when running, if no flags are provided help will
be shown.

`./deployer -t` - Deploy a testnet<br/>
`./deployer -u` - Update contracts<br/>
`./deployer -v` - Enables logging to console<br/>
`./deployer -hl` - Enables headless mode, no logging, and no questions. This
will overwrite a testnet if used in conjunction with `-t`, it also will use
whatever version is set in the manifest<br/>

## Installation

You can find the latest release of this software on the releases page 
[ here ](https://github.com/WebDelve/activeledger-testnet-deployer/releases)

Currently the only build is available for linux, however compilation for other
operating systems is simple.

### Compiling from source

Requires GoLang to be installed

1. Download this repo
2. In the root directory run `go build -o deployer` (named whatever you want)
3. Make executable, Linux and Mac: `chmod +x deployer`
4. Run it as defined in the quickstart section

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
  "testnetFolder": "sometestnet",
  "logToFile": true,
  "logFolder": "logs"
}
```

There are two additional config elements that are set by flags: Verbose logging,
and Headless mode.

Verbose logging will enable the printing, and writing, of debug messages.<br/>
Headless mode will not output messages to the console, but will still write to files

## Manifest

The manifest file contains a list of the smartcontracts you want to onboard to
the testnet.

You can also add ones that you don't yet want uploaded and use the `exclude: true`
pairing to ignore them.

The ID and Hash values can be left blank, they are set by the deployer automatically,
the ID will be added when the contract/s are onboarded, and the deployer starts
by checking for blank hashes and updating the manifest. It also will update the
hashes as required. Hashes are used to check if any updates have been made
to the contracts, if it finds none it won't continue.

The onboarded flag is also set by the software and is used to check if a contract
has been uploaded to the ledger or not.
When adding a new contract, make sure to either not include this, or preferably
include it and make sure it is false.

### Sample Manifest

```json
{
  "contracts": [
    {
      "name": "contractname",
      "id": "contractstreamid",
      "path": "smartcontracts/contract.ts",
      "version": "0.0.1",
      "hash": "contractdatahash",
      "exclude": false,
      "onboarded": false
    }
  ]
}
```

## Setup Output

Upon initial deployment the deployer will create a `setup-output.json` file
(or the file name set in `config.json`). This file contains data that you will
need when running transactions against the contracts, or that the deployer needs
when updating contracts. The `"contractData"` array is purely for your reference
as the same data is also set in the manifest for the deployer to access.

The identity stored in this file is linked to the key, although the key can be
used with other identities (not recommended), the identity is unusable without
the correct key.
It is worth noting that you only really need the Private PEM, as the rest of the
key data can be recreated from that. The hashes can be used to verify that the
PEMs are correct.

The namespace will be set to the same one defined in `config.json`

The `contractData` array is mainly used to provide users with easy reference
to the name and ID.
The name will match the one set in the manifest, and the ID is the stream ID
Activeledger assigned when the contract was onboarded.
The hash is not relevant here and is not set, it is a carry over from using
the same struct internally and will likely be removed in future versions.

```json
{
    "identity": "onboardedstreamid",
    "namespace": "yournamespace",
    "keyData": {
        "publicPem": "publicpem",
        "publicHash": "publichash",
        "privatePem": "privatepem",
        "privateHash": "privatehash"
    },
    "contractData": [
        {
            "name":"contractname",
            "id": "contractstreamid",
            "hash": ""
        }
    ]
}
```

## Todo

- ~~When adding a new contract to the manifest, after others have been onboarded 
already, needs to onboard that, there should be a flag in manifest that references
this: `"onboarded": true`~~
Added in 2.0.0

- ~~A useful additonal feature would be to allow updating contracts or updating them.
The software would need to check for the existence of an output file, or perhaps
should maintain a hidden file to keep track of things.~~
Added in 2.0.0

- ~~For this feature it should check the hashes of the contracts to find updated ones,
update those, and upload any new ones.~~
Added in 2.0.0 - Stores a hash in the manifest

- ~~Currently it will error on creating a testnet folder if one already exists.~~
Added in 2.0.0 - Checks if folder exists, deletes and recreates it if it does
Confirms with user before doing so.

## Changelog

### [Unreleased]

#### Features

- Check if Activeledger is installed (requires changes in Activeledger first)
- Ability to attempt to install Activeledger if it isn't installed
- Check if node/npm is installed
- Make labelling optional
- Add automatic version incrementing capability, this will be enabled with a 
CLI flag (e.g `-iv`) or an option in the config (e.g `"incrementVersion": true`), 
and will follow a basic schema defined in the config 
e.g: `"versionIncrement": "-.-.i"` where `i` is the digit to increment

### [2.0.0] - 2023-07-27

#### Added

- Ability to update contracts and upload new ones ~~based on hidden file~~ hash
in manifest
- Check if given testnet folder exists, if it does ask if it should be removed
if yes delete and recreate it, if no terminate
- Refactored code into internal packages
- Added sample json files
- Create headless mode so this can be used to automate processes
- ~~Add verbose logging mode~~ Added custom logger which handles this, 
adds contextual colouring, and enables logging to file

### [1.0.0] - 2023-06-30

#### Added

- Created the initial software
- Runs `activeledger --testnet`
- Onboards configured identity and namespace
- Uploads smartcontracts as defined in manifest file


## License

[MIT](https://github.com/WebDelve/activeledger-testnet-deployer/blob/master/LICENCE)
