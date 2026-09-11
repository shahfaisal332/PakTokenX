"use client";
import { FormEvent, useEffect, useState } from "react";
import SiteHeader from "../components/SiteHeader";
import { CATEGORIES } from "../categories";

const BACKEND_URL = "http://localhost:3000";

type Project = {
  id: string;
  name: string;
  total_supply: string;
  owner: string;
  token_price: string;
  creator: string;
  category: string;
};
type Investor = { id: number; name: string; wallet_address: string };
type Holding = { index: string; amount: string; project_id: string; investor: string; creator: string };

export default function InvestorPage() {
  const [investor, setInvestor] = useState<Investor | null>(null);
  const [name, setName] = useState("");
  const [idNumber, setIdNumber] = useState("");
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [holdings, setHoldings] = useState<Holding[]>([]);
  const [investmentAmounts, setInvestmentAmounts] = useState<{ [key: string]: string }>({});
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const savedInvestor = localStorage.getItem("paktoken_investor");
    if (savedInvestor) {
      try { setInvestor(JSON.parse(savedInvestor)); } catch { localStorage.removeItem("paktoken_investor"); }
    }
    loadProjects();
  }, []);
  useEffect(() => { if (investor) loadHoldings(investor.id); }, [investor]);

  async function loadProjects() {
    try {
      const response = await fetch(`${BACKEND_URL}/projects`);
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || "Failed to load projects");
      setProjects(data.project || []);
    } catch (error) {
      console.error(error);
      setMessage("Could not load projects.");
    }
  }
  async function loadHoldings(investorId: number) {
    try {
      const response = await fetch(`${BACKEND_URL}/holdings/${investorId}`);
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || "Failed to load holdings");
      setHoldings(data || []);
    } catch (error) {
      console.error(error);
      setHoldings([]);
    }
  }
  async function registerInvestor(event: FormEvent) {
    event.preventDefault();
    setLoading(true);
    setMessage("");
    try {
      const response = await fetch(`${BACKEND_URL}/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, idNumber }),
      });
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || "Registration failed");
      const newInvestor: Investor = { id: data.id, name: data.name, wallet_address: data.wallet_address };
      setInvestor(newInvestor);
      localStorage.setItem("paktoken_investor", JSON.stringify(newInvestor));
      setName("");
      setIdNumber("");
      setMessage(`Registration successful. Investor ID: ${data.id}`);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Registration failed");
    } finally {
      setLoading(false);
    }
  }
  async function invest(event: FormEvent, projectId: string) {
    event.preventDefault();
    if (!investor) { setMessage("Please register as an investor first."); return; }
    const amount = investmentAmounts[projectId];
    if (!amount || Number(amount) <= 0) { setMessage("Please enter a valid investment amount."); return; }
    setLoading(true);
    setMessage("");
    try {
      const response = await fetch(`${BACKEND_URL}/invest`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ investorId: investor.id, projectId: Number(projectId), amount: Number(amount) }),
      });
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || "Investment failed");
      setMessage(`Investment successful. Transaction: ${data.txHash}`);
      setInvestmentAmounts({ ...investmentAmounts, [projectId]: "" });
      await loadHoldings(investor.id);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Investment failed");
    } finally {
      setLoading(false);
    }
  }
  function logoutInvestor() {
    localStorage.removeItem("paktoken_investor");
    setInvestor(null);
    setHoldings([]);
    setMessage("Investor session cleared.");
  }
  function categoryLabel(value: string) {
    return CATEGORIES.find((c) => c.value === value)?.label || "Uncategorized";
  }

  const projectsInField = selectedCategory ? projects.filter((p) => p.category === selectedCategory) : [];

  return (
    <main className="min-h-screen bg-[#F4F6F8] text-[#1E293B]">
      <SiteHeader active="investor" />
      <div className="mx-auto max-w-6xl px-6 py-10">
        <div className="mb-10">
          <p style={{ fontFamily: "var(--font-mono, monospace)" }} className="text-xs uppercase tracking-widest text-[#008751]">
            Investor Portal
          </p>
          <h2 style={{ fontFamily: "var(--font-sora)" }} className="mt-2 text-3xl font-bold text-[#0B192C]">
            Invest In Real Assets
          </h2>
          <p className="mt-2 text-[#5B6B7F]">Register, choose a field, and invest in a listed project.</p>
        </div>

        {message && (
          <div className="mb-8 rounded-xl border border-[#E4E8ED] bg-white p-4 text-sm shadow-sm">{message}</div>
        )}

        {!investor ? (
          <section className="mb-8 max-w-xl rounded-2xl bg-white p-6 shadow-sm">
            <h3 style={{ fontFamily: "var(--font-sora)" }} className="mb-6 text-xl font-bold text-[#0B192C]">
              Register as Investor
            </h3>
            <form onSubmit={registerInvestor} className="space-y-4">
              <div>
                <label className="mb-2 block text-sm font-semibold text-[#0B192C]">Name</label>
                <input
                  type="text" value={name} onChange={(e) => setName(e.target.value)}
                  placeholder="Investor Name" required
                  className="w-full rounded-lg border border-[#E4E8ED] px-4 py-3 outline-none focus:border-[#00A859]"
                />
              </div>
              <div>
                <label className="mb-2 block text-sm font-semibold text-[#0B192C]">ID Number</label>
                <input
                  type="text" value={idNumber} onChange={(e) => setIdNumber(e.target.value)}
                  placeholder="Enter ID number" required
                  className="w-full rounded-lg border border-[#E4E8ED] px-4 py-3 outline-none focus:border-[#00A859]"
                />
              </div>
              <button
                type="submit" disabled={loading}
                className="w-full rounded-lg bg-[#00A859] px-4 py-3 font-semibold text-white transition hover:bg-[#008751] disabled:opacity-50"
              >
                {loading ? "Registering..." : "Register Investor"}
              </button>
            </form>
          </section>
        ) : (
          <section className="mb-8 rounded-2xl bg-white p-6 shadow-sm">
            <div className="flex flex-col justify-between gap-4 md:flex-row md:items-center">
              <div>
                <p className="text-sm text-[#5B6B7F]">Logged in as</p>
                <h3 style={{ fontFamily: "var(--font-sora)" }} className="text-xl font-bold text-[#0B192C]">
                  {investor.name}
                </h3>
                <p className="mt-1 text-sm text-[#5B6B7F]">Investor ID: {investor.id}</p>
                <p className="mt-1 break-all text-xs text-[#5B6B7F]">Wallet: {investor.wallet_address}</p>
              </div>
              <button
                onClick={logoutInvestor}
                className="rounded-lg border border-[#E4E8ED] px-5 py-2.5 font-semibold text-[#5B6B7F] hover:bg-[#F4F6F8]"
              >
                Clear Session
              </button>
            </div>
          </section>
        )}

        <section className="mb-8 rounded-2xl bg-white p-6 shadow-sm">
          <h3 style={{ fontFamily: "var(--font-sora)" }} className="mb-2 text-xl font-bold text-[#0B192C]">
            Step 1- Select a Field
          </h3>
          <p className="mb-6 text-sm text-[#5B6B7F]">Choose which type of asset you want to invest in.</p>
          <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
            {CATEGORIES.map((c) => (
              <button
                key={c.value}
                onClick={() => setSelectedCategory(c.value)}
                className={`flex flex-col items-center gap-2 rounded-xl border-2 p-5 text-sm font-semibold transition ${
                  selectedCategory === c.value
                    ? "border-[#00A859] bg-[#E4F7EE] text-[#008751]"
                    : "border-[#E4E8ED] bg-[#F4F6F8] text-[#5B6B7F] hover:border-[#00A859]"
                }`}
              >
                <span className="text-2xl">{c.icon}</span>
                {c.label}
              </button>
            ))}
          </div>
        </section>

        {selectedCategory && (
          <section className="mb-8 rounded-2xl bg-white p-6 shadow-sm">
            <div className="mb-6 flex items-center justify-between">
              <div>
                <h3 style={{ fontFamily: "var(--font-sora)" }} className="text-xl font-bold text-[#0B192C]">
                  Step 2- Projects in {categoryLabel(selectedCategory)}
                </h3>
                <p className="mt-1 text-sm text-[#5B6B7F]">Choose a project and enter how many tokens to buy.</p>
              </div>
              <button
                onClick={loadProjects}
                className="rounded-lg border border-[#E4E8ED] px-4 py-2 text-sm font-semibold text-[#5B6B7F] hover:bg-[#F4F6F8]"
              >
                Refresh
              </button>
            </div>
            {projectsInField.length === 0 ? (
              <p className="text-[#5B6B7F]">No projects registered in this field yet.</p>
            ) : (
              <div className="grid gap-6 md:grid-cols-2">
                {projectsInField.map((project) => (
                  <div key={project.id} className="rounded-xl border border-[#E4E8ED] p-5">
                    <div className="mb-4">
                      <p className="text-sm text-[#5B6B7F]">Project #{project.id}</p>
                      <h4 style={{ fontFamily: "var(--font-sora)" }} className="text-xl font-bold text-[#0B192C]">
                        {project.name}
                      </h4>
                    </div>
                    <div className="mb-5 grid grid-cols-2 gap-3 text-sm">
                      <div className="rounded-lg bg-[#F4F6F8] p-3">
                        <p className="text-[#5B6B7F]">Total Supply</p>
                        <p className="mt-1 font-bold text-[#0B192C]">{project.total_supply}</p>
                      </div>
                      <div className="rounded-lg bg-[#F4F6F8] p-3">
                        <p className="text-[#5B6B7F]">Token Price</p>
                        <p className="mt-1 font-bold text-[#0B192C]">{project.token_price}</p>
                      </div>
                    </div>
                    <form onSubmit={(e) => invest(e, project.id)} className="space-y-3">
                      <input
                        type="number" min="1"
                        value={investmentAmounts[project.id] || ""}
                        onChange={(e) => setInvestmentAmounts({ ...investmentAmounts, [project.id]: e.target.value })}
                        placeholder="Number of tokens"
                        disabled={!investor || loading}
                        required
                        className="w-full rounded-lg border border-[#E4E8ED] px-4 py-3 outline-none focus:border-[#00A859] disabled:bg-[#F4F6F8]"
                      />
                      <button
                        type="submit" disabled={!investor || loading}
                        className="w-full rounded-lg bg-[#00A859] px-4 py-3 font-semibold text-white transition hover:bg-[#008751] disabled:cursor-not-allowed disabled:opacity-50"
                      >
                        {!investor ? "Register First" : loading ? "Processing..." : "Invest"}
                      </button>
                    </form>
                  </div>
                ))}
              </div>
            )}
          </section>
        )}

        <section className="rounded-2xl bg-white p-6 shadow-sm">
          <div className="mb-6 flex items-center justify-between">
            <div>
              <h3 style={{ fontFamily: "var(--font-sora)" }} className="text-xl font-bold text-[#0B192C]">
                My Holdings
              </h3>
              <p className="mt-1 text-sm text-[#5B6B7F]">Your token holdings recorded on the PakTokenX blockchain.</p>
            </div>
            {investor && (
              <button
                onClick={() => loadHoldings(investor.id)}
                className="rounded-lg border border-[#E4E8ED] px-4 py-2 text-sm font-semibold text-[#5B6B7F] hover:bg-[#F4F6F8]"
              >
                Refresh
              </button>
            )}
          </div>
          {!investor ? (
            <p className="text-[#5B6B7F]">Register as an investor to view your holdings.</p>
          ) : holdings.length === 0 ? (
            <p className="text-[#5B6B7F]">You do not have any token holdings yet.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead>
                  <tr className="border-b border-[#E4E8ED] text-[#5B6B7F]">
                    <th className="px-3 py-3">Project ID</th>
                    <th className="px-3 py-3">Tokens Held</th>
                    <th className="px-3 py-3">Investor Wallet</th>
                  </tr>
                </thead>
                <tbody>
                  {holdings.map((holding) => (
                    <tr key={holding.index} className="border-b border-[#F4F6F8]">
                      <td className="px-3 py-3 font-semibold text-[#0B192C]">{holding.project_id}</td>
                      <td className="px-3 py-3 font-bold text-[#0B192C]">{holding.amount}</td>
                      <td className="max-w-xs break-all px-3 py-3 text-xs text-[#5B6B7F]">{holding.investor}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
