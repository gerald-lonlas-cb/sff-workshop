require('dotenv').config();
const HDWalletProvider = require("@truffle/hdwallet-provider");
const { NODE_URL, MNEMONIC } = process.env;

module.exports = {
  networks: {
    goerli: {
      provider: function () {
        return new HDWalletProvider(MNEMONIC, NODE_URL)
      },
      network_id: "5",
      port: 8545,
      gas: 4465030
    }
  },
  contracts_directory: "./contract/",
  contracts_build_directory: "./contract/abi",
  compilers: {
    solc: {
      version: "^0.8.0",
       optimizer: {
         enabled: false,
         runs: 200
       }
    }
  }
};
