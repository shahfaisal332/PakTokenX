require("dotenv").config();

const { DirectSecp256k1HdWallet } = require("@cosmjs/proto-signing");
const { SigningStargateClient, StargateClient } = require("@cosmjs/stargate");
const { registry } = require("./paktoken-registry");

const RPC_URL = process.env.RPC_URL;
const PREFIX = process.env.PREFIX;

async function getSigningClient(mnemonic) {
  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
    prefix: PREFIX,
  });

  const [account] = await wallet.getAccounts();

  const client = await SigningStargateClient.connectWithSigner(
    RPC_URL,
    wallet,
    { registry }
  );

  return { client, address: account.address };
}

async function createWallet() {
  const wallet = await DirectSecp256k1HdWallet.generate(24, {
    prefix: PREFIX,
  });

  const [account] = await wallet.getAccounts();
  const mnemonic = wallet.mnemonic;

  return { mnemonic, address: account.address };
}

async function getReadOnlyClient() {
  return await StargateClient.connect(RPC_URL);
}

async function fundNewWallet(newAddress) {
  const { client, address } = await getSigningClient(process.env.FAUCET_MNEMONIC);
  const amount = { denom: process.env.DENOM, amount: "5000" };
  const fee = { amount: [{ denom: process.env.DENOM, amount: "1000" }], gas: "200000" };
  await client.sendTokens(address, newAddress, [amount], fee);
}

module.exports = { getSigningClient, createWallet, getReadOnlyClient, fundNewWallet };
