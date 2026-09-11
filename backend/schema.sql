CREATE TABLE investors (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  id_number TEXT NOT NULL,
  wallet_address TEXT NOT NULL,
  mnemonic TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now()
);
