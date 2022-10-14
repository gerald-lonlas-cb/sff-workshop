# sff-workshop
For SFF workshop

## Running locally
Create an env file with content

```
USERNAME=<Username of Coinbase Cloud account>
PASSWORD=<Password of Coinbase Cloud account>
PRIVATE_KEY=<Private key of the account currently holding the ERC1155 NFT>
CONTRACT_ADDRESS=<Contract address of the ERC1155 NFT>
```

Create .env file in the project root dir

Build the server using

```
make dpes
make build
```

Start the server using
```
./bin/main 
```
or
```
make run
```

Make a request to airdrop ERC1155 tokens
```
curl --url 'http://localhost:8081/gettoken?to=<the address to airdrop tokens to>&id=<id of the nft item>&quantity=<amount of the nft item>'
```
Example
```
curl --url 'http://localhost:8081/gettoken?to=0xF820cf368b4a798b676DE9DEA90f637A9CdEE572&id=2&quantity=3'
```
