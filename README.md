# AICP app for SFF'23 demo

This repo is simplified demo of what the AICP back-end app can be. This project only focus on the customer side. Not on the merchant redeem side.
The role of this back-end server is to connect your application to an EVM blockchain node (Ethereum, Polygon, ...), and create an order on behalf of customers for which you manage their wallet.

# Features

- [x] Backend service to derive customer wallets from a main wallet.
  - Code in [internal/wallet/customer_waller.go](internal/wallet/customer_waller.go)
- [x] Support Alchemy, Infura, ... providers
- [x] Provide the [smart contract source code](contract/RetailOrderEscrow.sol)
- [x] APIs to interact with the RetailOrderEscrow smart contract
  - Microservice public API defined in [cmd/main.go router()](cmd/main.go)
  - APIs code to interact with the smart contract in [internal/server/api.go](internal/server/api.go)

**APIs for Customer actions**

- [x] `/api/createOrder`: Create an order on behalf of a customer (from their wallet).
  - Note: the customer wallet needs some Matic to execute this onchain transaction.
- [x] `/api/getOrderStatus`: Return the status of an order.
- [x] `/api/getCustomerAvailableBalance`: Provide the number of the ERC-1155 fungible token the customer has on their wallet.

**APIs for Merchant actions**

- [x] `/api/cancelOrder`: Cancel a order.
  - Note: the customer wallet needs some Matic to execute this onchain transaction.
- [ ] Redeem an order. This merchant feature is not implemented in this demo.

**Utilities APIs**

- [x] `/api/getCustomerPublicAddress`: Return the Public wallet address from a customer ID.
- [x] `/api/getTransactionStatus`: Utility API to get the onchain status of a transaction.
- [x] `/api/mintTokens`: Mint and Add the ERC-1155 fungible token to the customer wallet.
  - Note: the customer wallet needs some Matic to execute this onchain transaction.

# Table of content

- [AICP app for SFF'23 demo](#aicp-app-for-sff23-demo)
- [Features](#features)
- [Table of content](#table-of-content)
  - [1. Pre-requisite](#1-pre-requisite)
    - [1.1a MacOs Users](#11a-macos-users)
    - [1.1b Ubuntu/Debian Users](#11b-ubuntudebian-users)
  - [2. Setup the back-end service](#2-setup-the-back-end-service)
    - [2.1 Install the dependencies](#21-install-the-dependencies)
    - [2.2 Create the .env config file](#22-create-the-env-config-file)
  - [3. Build and run the server](#3-build-and-run-the-server)
  - [4. Demo: Test the server](#4-demo-test-the-server)
    - [4.1 Know your customer wallet and balance](#41-know-your-customer-wallet-and-balance)
    - [4.2 Create an order on behalf of your customer](#42-create-an-order-on-behalf-of-your-customer)
    - [4.3 Cancel the order](#43-cancel-the-order)
- [FAQ](#faq)
  - [Question 1: In this Proof of Concept, do customer owns their wallet?](#question-1-in-this-proof-of-concept-do-customer-owns-their-wallet)
  - [Question 2: How do you generate a wallet for a customers?](#question-2-how-do-you-generate-a-wallet-for-a-customers)
  - [Question 3: Can we transfer ownership of the wallet to the customers? Can the customer get the private key?](#question-3-can-we-transfer-ownership-of-the-wallet-to-the-customers-can-the-customer-get-the-private-key)
  - [Question 4: Can I see the content of a customer wallet into a self-custody wallet app like Coinbase Wallet?](#question-4-can-i-see-the-content-of-a-customer-wallet-into-a-self-custody-wallet-app-like-coinbase-wallet)
- [Appendix](#appendix)
  - [Appendix 1: Update the Go Smart contract representation](#appendix-1-update-the-go-smart-contract-representation)
    - [Using an existing ABI file or Generating a new one from the contract source](#using-an-existing-abi-file-or-generating-a-new-one-from-the-contract-source)
    - [Generate the ABI file](#generate-the-abi-file)
    - [Generate the Go binding file](#generate-the-go-binding-file)
  - [Appendix 2: Deploy the smart contract](#appendix-2-deploy-the-smart-contract)
  - [Appendix 3: Netlify vs Local](#appendix-3-netlify-vs-local)

## 1. Pre-requisite

To run this micro-service locally, you will need:

- An EVM Wallet (Private key and Mnemonic).
  - We recommend to use a temporary wallet. DO NOT use an important wallet.
  - Suggestion: Create a wallet with [Coinbase Wallet](https://www.coinbase.com/wallet) and use this mnemonic
- Golang ([Install page](https://go.dev/doc/install))
- Some faucet Matic on your wallet to execute the "write" transactions
  - [faucet.polygon.technology](https://faucet.polygon.technology/)
  - [coinbase.com/faucets](https://coinbase.com/faucets)

### 1.1a MacOs Users

```bash
# Install Brew
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
brew update && brew doctor
brew install go
```

### 1.1b Ubuntu/Debian Users

```bash
sudo apt install golang-go
```

## 2. Setup the back-end service

### 2.1 Install the dependencies

```bash
make deps
```

### 2.2 Create the .env config file

The `.env` file contains all sensitive information to connect to your blockchain node and manage the wallet.
Create the `.env` file at the root of this project

```bash
cp env.sample .env
```

Edit the `.env` file values

```bash
vi .env
```

```ini
NODE_URI="polygon-mumbai.g.alchemy.com/v2/XXXXXXXXXXX"
CONTRACT_ADDRESS=<Address of the ERC1155 Escrow contract>
MNEMONIC=<Mnemonic of the AICP wallet that will hold customers funds>
```

For the smart contract you can either use 0x533d64229F077159d2c49D4C646055233AD3c04f (On Polygon Mumbai Testnet), or deploy it yourself. The source code is in `contract/RetailOrderEscrow.sol`

## 3. Build and run the server

```bash
make build && make run
```

## 4. Demo: Test the server

Before to start unsure the service is running (see step 3. Build and run the server).  
To run this demo you can either use a Curl command on a terminal, or your favorite REST client.

### 4.1 Know your customer wallet and balance

**4.1.1: Get the wallet address of the customer ID "`CustomerID_ABCD123`"**  
Note: This customer ID can be anything, til it is a string.

```ssh
curl -X GET 'http://localhost:8081/api/getCustomerPublicAddress?customer_id=CustomerID_ABCD123'
```

**4.1.2: Transfer some Matic tokens to the customer wallet (ex: 0.01)**
You can use one of the following service to drop faucet to the test wallet. Otherwise do a regular transfer.

- [faucet.polygon.technology](https://faucet.polygon.technology/)
- [coinbase.com/faucets](https://coinbase.com/faucets)

**4.1.3: Verify the customer balance for the test PBM/Voucher token**  
This smart contract generates fungible ERC-1155 tokens. The token_id 1 will work. You will need to know what is the token id you are looking for.

```ssh
curl -X GET 'http://localhost:8081/api/getCustomerAvailableBalance?customer_id=CustomerID_ABCD123&token_id=1'
```

**Result:** You should see a balance of 0

**4.1.4: Mint and Add some ERC-1155 tokens to the customers wallet**  
This step should not exist in production. This will be replaced by the airdrop the customers will receive.  
4.1.4.1 Here we are using it to fund our test customer wallet with 10 tokens.

```ssh
curl -X GET 'http://localhost:8081/api/mintTokens?customer_id=CustomerID_ABCD123&token_id=1&quantity=10'
```

4.1.4.2 Verify the customer wallet:

```ssh
curl -X GET 'http://localhost:8081/api/getCustomerAvailableBalance?customer_id=CustomerID_ABCD123&token_id=1'
```

**Result:** You should see a balance of 10

### 4.2 Create an order on behalf of your customer

Note that the current smart contract does not transfer the amount of tokens to the smart contract. So the current smart contract "CreateOrder" function only add a new order entry in the smart contract.

**4.2.1 Create the order #1 on behalf of customer "`CustomerID_ABCD123`" for an amount of 5 tokens, and set the merchant address to `0x1Fedd98E643E58dcAc6b30cC089D158A719c6EB7`**

```ssh
curl -X GET 'http://localhost:8081/api/createOrder?customer_id=CustomerID_ABCD123&order_id=1&merchant_address=0x1Fedd98E643E58dcAc6b30cC089D158A719c6EB7&amount=5'
```

**4.2.2 Retrieve the order status**

```ssh
curl -X GET 'http://localhost:8081/api/getOrderStatus?order_id=1'
```

**Result:** The status returned should be 0

As per the smart contract the status enun is:
0: PENDING
1: REDEEMED
2: CANCELLED

### 4.3 Cancel the order

Note that the current smart contract allows anyone to cancel any order. No guardrail has been set yet.

**4.3.1 Customer "CustomerID_ABCD123" cancels the order #1**

```ssh
curl -X GET 'http://localhost:8081/api/cancelOrder?customer_id=CustomerID_ABCD123&order_id=1'
```

**4.3.2 Retrieve the order status**

```ssh
curl -X GET 'http://localhost:8081/api/getOrderStatus?order_id=1'
```

**Result:** The status returned should be 2

# FAQ

## Question 1: In this Proof of Concept, do customer owns their wallet?

In this proof of concept, the customer do not own their wallet. As a custodian, you manage all the customers wallets you have generated.

## Question 2: How do you generate a wallet for a customers?

This proof of concept is based on the ability that a Wallet Mnemonic (BIP39 wallet) can generates almost an infinite number of sub-wallets. Hence, we convert a customer ID (a string) into an Integer that is used as the customer Wallet index.

The code logic is on [internal/wallet/main_wallet.go](internal/wallet/main_wallet.go) and [internal/wallet/customer_wallet.go](internal/wallet/main_wallet.go).

In short each customer is a sub-wallet of the path m/44'/0'/0'/0/.  
Ex:

- Customer `1` will be `m/44'/0'/0'/0/1`
- Customer `2` will be `m/44'/0'/0'/0/2`
- Customer `100404` will be `m/44'/0'/0'/0/100404`

You can see this concept in action on this [Online Wallet generator page](https://iancoleman.io/bip39/). Try to generate a wallet then go to the "Derived Addresses" to see that each index generate a new couple of Public and Private keys.

## Question 3: Can we transfer ownership of the wallet to the customers? Can the customer get the private key?

You cannot completely transfer the ownership of the customer wallet. As it is derived from your main wallet. You will still have full control of it. However you can share to the customer the private key of their wallet (only their wallet).

## Question 4: Can I see the content of a customer wallet into a self-custody wallet app like Coinbase Wallet?

Yes you can. Given you will generate thousand of wallets, the easiest way it to get the private key of the specific wallet and import it in a self-custody wallet app.

You will find a line commented on `the internal/wallet/customer_wallet.go` to display the private key in the service logs. **It is strongly discouraged to use it.** This commented line is here to debug during dev and tests.

# Appendix

## Appendix 1: Update the Go Smart contract representation

For the backend to work, it needs to know what method the smart contract has. Hence in the folder `contract/` you will find a _.abi and _.go files that represent the onchain smart contract.

Here how to generate the ABI and the Go file.

### Using an existing ABI file or Generating a new one from the contract source

You can get the ABI file from the contract owner, or generate it if you have access to the smart contract source code.

### Generate the ABI file

To generate a Go file representing a Solidity smart contract, you need to use the abigen tool provided by Go Ethereum (go-ethereum or geth).
The abigen tool generates a Go binding from a Solidity contract, given its Application Binary Interface (ABI) and bytecode.

**Please follow the steps below:**

Install Solc

```bash
# Ubuntu/Debian users
sudo apt-get install solc

# MacOS
brew tap ethereum/ethereum
brew install solidity
```

Compile your Solidity contract (It must be flatten) to obtain the ABI and bytecode. You can use the Solidity compiler (solc) for this. Run the following command in the directory where your contract file (for example, RetailOrderEscrow.sol) is located:

```bash
solc --abi --bin RetailOrderEscrow_flattened.sol -o build
```

Note: That solc does not know how to resolve and download Solidity imports hence you Solidity contract must be flatten (means contains the smart contract code + all the dependencies code in a single file). To flatten a file, you can use https://remix.ethereum.org/. Paste the code of your contract, then right click on the contract file name and select "Flatten". A new file `XXX_flattened.sol` will be created.

This command will create two files in the build directory: `RetailOrderEscrow.abi` (the contract ABI) and `RetailOrderEscrow.bin` (the bytecode).

### Generate the Go binding file

Use abigen to generate the Go binding. Run the following command:

```bash
../tools/abigen --bin=./build/RetailOrderEscrow.bin --abi=./build/RetailOrderEscrow.abi --pkg=contract --out=contract.go
```

This command generates a `contract.go` file in the same directory. This Go file contains a binding to the contract that you can use to interact with it from a Go program.

## Appendix 2: Deploy the smart contract

You can choose to use the Remix-IDE to deploy the contract manually. Please follow the instruction [here](https://remix-ide.readthedocs.io/en/latest/create_deploy.html#deploy-the-contract)

You might want to use the same MNEMONIC that you specify in the .env file so that you can directly transfer from the same wallet.

## Appendix 3: Netlify vs Local

We can use import from Git function of Netlify for deployment. Netlify uses build.sh and netlify.toml files to build and publish the server.

When running on Netlify, we don't pass -port option to the run time argument (port will be defaulted to -1 in this case). The logic in main.go will transform the http server into a lambda to be run on Netlify. In config.go we won't call godotenv.Load(".env") as the environment variables are set from Netlify config instead of .env file.

When running locally, we need to pass -port option and a normal http server will be started on that port, allowing us to test locally without the need for AWS lambda simulator. In config.go we will call godotenv.Load(".env") to set environment variables using .env file.
