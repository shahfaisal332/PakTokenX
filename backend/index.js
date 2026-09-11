require("dotenv").config();
const express = require("express");
const cors = require("cors");
const pool = require("./db");
const { getSigningClient, createWallet, fundNewWallet } = require("./chain");

const app = express();
app.use(cors());
app.use(express.json());

const PORT = process.env.PORT || 3000;

const TYPE_URL_CREATE_PROJECT = "/paktoken.tokenization.v1.MsgCreateProject";
const TYPE_URL_BUY_TOKENS = "/paktoken.tokenization.v1.MsgBuyTokens";
const TYPE_URL_DISTRIBUTE_REVENUE = "/paktoken.tokenization.v1.MsgDistributeRevenue";

const fee = {
  amount: [{ denom: process.env.DENOM, amount: "5000" }],
  gas: "200000",
};

app.get("/health", (req, res) => {
  res.json({ status: "ok", service: "paktoken-backend" });
});

app.post("/register", async (req, res) => {
  const { name, idNumber } = req.body;
  if (!name || !idNumber) {
    return res.status(400).json({ error: "name and idNumber are required" });
  }
  try {
    const { mnemonic, address } = await createWallet();
    await fundNewWallet(address);
    const result = await pool.query(
      `INSERT INTO investors (name, id_number, wallet_address, mnemonic)
       VALUES ($1, $2, $3, $4) RETURNING id, name, wallet_address`,
      [name, idNumber, address, mnemonic]
    );
    res.json(result.rows[0]);
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: "registration failed" });
  }
});

app.post("/projects", async (req, res) => {
  const { name, totalSupply, tokenPrice, category } = req.body;
  if (!category) {
    return res.status(400).json({ error: "category is required" });
  }
  try {
    const { client, address } = await getSigningClient(process.env.SPONSOR_MNEMONIC);
    const msg = {
      typeUrl: TYPE_URL_CREATE_PROJECT,
      value: {
        creator: address,
        name: name,
        totalSupply: totalSupply.toString(),
        owner: address,
        tokenPrice: tokenPrice.toString(),
        category: category,
      },
    };
    const result = await client.signAndBroadcast(address, [msg], fee);
    res.json({ txHash: result.transactionHash, code: result.code });
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: "failed to create project", details: err.message });
  }
});

app.post("/invest", async (req, res) => {
  const { investorId, projectId, amount } = req.body;
  try {
    const investorRes = await pool.query("SELECT * FROM investors WHERE id = $1", [investorId]);
    if (investorRes.rows.length === 0) {
      return res.status(404).json({ error: "investor not found" });
    }
    const investor = investorRes.rows[0];
    const { client, address } = await getSigningClient(investor.mnemonic);

    const msg = {
      typeUrl: TYPE_URL_BUY_TOKENS,
      value: {
        creator: address,
        projectId: projectId.toString(),
        amount: amount.toString(),
      },
    };
    const result = await client.signAndBroadcast(address, [msg], fee);
    res.json({ txHash: result.transactionHash, code: result.code });
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: "investment failed", details: err.message });
  }
});

app.post("/distribute", async (req, res) => {
  const { projectId, amount } = req.body;
  try {
    const { client, address } = await getSigningClient(process.env.SPONSOR_MNEMONIC);
    const msg = {
      typeUrl: TYPE_URL_DISTRIBUTE_REVENUE,
      value: {
        creator: address,
        projectId: projectId.toString(),
        amount: amount.toString(),
      },
    };
    const result = await client.signAndBroadcast(address, [msg], fee);
    res.json({ txHash: result.transactionHash, code: result.code, events: result.events });
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: "distribution failed", details: err.message });
  }
});

app.post("/release-supply", async (req, res) => {
  const { projectId, releasedSupply } = req.body;
  try {
    const { client, address } = await getSigningClient(process.env.SPONSOR_MNEMONIC);
    const msg = {
      typeUrl: "/paktoken.tokenization.v1.MsgSetReleasedSupply",
      value: {
        creator: address,
        projectId: projectId.toString(),
        releasedSupply: releasedSupply.toString(),
      },
    };
    const result = await client.signAndBroadcast(address, [msg], fee);
    res.json({ txHash: result.transactionHash, code: result.code });
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: "failed to update released supply", details: err.message });
  }
});


app.get("/projects", async (req, res) => {
  try {
    const response = await fetch("http://127.0.0.1:1317/paktoken/tokenization/v1/project");
    const data = await response.json();
    res.json(data);
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: "failed to fetch projects" });
  }
});

app.get("/holdings/:investorId", async (req, res) => {
  try {
    const investorRes = await pool.query("SELECT wallet_address FROM investors WHERE id = $1", [req.params.investorId]);
    if (investorRes.rows.length === 0) {
      return res.status(404).json({ error: "investor not found" });
    }
    const walletAddress = investorRes.rows[0].wallet_address;
    const response = await fetch("http://127.0.0.1:1317/paktoken/tokenization/v1/holding");
    const data = await response.json();
    const mine = (data.holding || []).filter((h) => h.investor === walletAddress);
    res.json(mine);
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: "failed to fetch holdings" });
  }
});

app.listen(PORT, () => {
  console.log(`PakToken backend running on http://localhost:${PORT}`);
});
