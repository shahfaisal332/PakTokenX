"use client";
import { FormEvent, useEffect, useState } from "react";
import SiteHeader from "./components/SiteHeader";
import { CATEGORIES } from "./categories";

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

export default function SponsorPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [projectName, setProjectName] = useState("");
  const [totalSupply, setTotalSupply] = useState("");
  const [tokenPrice, setTokenPrice] = useState("");
  const [category, setCategory] = useState("");
  const [projectId, setProjectId] = useState("");
  const [revenueAmount, setRevenueAmount] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  async function loadProjects() {
    try {
      const response = await fetch(`${BACKEND_URL}/projects`);
      const data = await response.json();
      setProjects(data.project || []);
    } catch (error) {
      console.error(error);
      setMessage("Could not connect to the backend.");
    }
  }
  useEffect(() => { loadProjects(); }, []);

  async function createProject(event: FormEvent) {
    event.preventDefault();
    if (!category) {
      setMessage("Please select a field for this asset before continuing.");
      return;
    }
    setLoading(true);
    setMessage("");
    try {
      const response = await fetch(`${BACKEND_URL}/projects`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: projectName,
          totalSupply: Number(totalSupply),
          tokenPrice: Number(tokenPrice),
          category: category,
        }),
      });
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || "Project creation failed");
      setMessage(`Project created successfully. Transaction: ${data.txHash}`);
      setProjectName("");
      setTotalSupply("");
      setTokenPrice("");
      setCategory("");
      await loadProjects();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Project creation failed");
    } finally {
      setLoading(false);
    }
  }

  async function distributeRevenue(event: FormEvent) {
    event.preventDefault();
    setLoading(true);
    setMessage("");
    try {
      const response = await fetch(`${BACKEND_URL}/distribute`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ projectId: Number(projectId), amount: Number(revenueAmount) }),
      });
      const data = await response.json();
      if (!response.ok) throw new Error(data.error || "Revenue distribution failed");
      const revenueEvent = data.events?.find((e: { type: string }) => e.type === "revenue_distributed");
      setMessage(
        revenueEvent
          ? `Revenue distributed successfully. Investor payout: ${revenueEvent.attributes?.find((a: { key: string }) => a.key === "amount")?.value || "0"}`
          : `Revenue distributed successfully. Transaction: ${data.txHash}`
      );
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Revenue distribution failed");
    } finally {
      setLoading(false);
    }
  }

  function categoryLabel(value: string) {
    return CATEGORIES.find((c) => c.value === value)?.label || "Uncategorized";
  }

  return (
    <main className="min-h-screen bg-[#F4F6F8] text-[#1E293B]">
      <SiteHeader active="sponsor" />
      <div className="mx-auto max-w-6xl px-6 py-10">
        <div className="mb-10">
          <p style={{ fontFamily: "var(--font-mono, monospace)" }} className="text-xs uppercase tracking-widest text-[#008751]">
            Sponsor Portal
          </p>
          <h2 style={{ fontFamily: "var(--font-sora)" }} className="mt-2 text-3xl font-bold text-[#0B192C]">
            Tokenize Your Asset
          </h2>
          <p className="mt-2 text-[#5B6B7F]">Create a tokenized project and distribute revenue to investors.</p>
        </div>

        {message && (
          <div className="mb-8 rounded-xl border border-[#E4E8ED] bg-white p-4 text-sm shadow-sm">{message}</div>
        )}

        <div className="grid gap-8 md:grid-cols-2">
          <section className="rounded-2xl bg-white p-6 shadow-sm">
            <h3 style={{ fontFamily: "var(--font-sora)" }} className="mb-2 text-xl font-bold text-[#0B192C]">
              Create Project
            </h3>
            <p className="mb-6 text-sm text-[#5B6B7F]">Select the field this asset belongs to, then fill in the details.</p>
            <form onSubmit={createProject} className="space-y-5">
              <div>
                <label className="mb-2 block text-sm font-semibold text-[#0B192C]">Select Field</label>
                <div className="grid grid-cols-2 gap-3">
                  {CATEGORIES.map((c) => (
                    <button
                      type="button"
                      key={c.value}
                      onClick={() => setCategory(c.value)}
                      className={`flex flex-col items-center gap-2 rounded-xl border-2 p-4 text-sm font-semibold transition ${
                        category === c.value
                          ? "border-[#00A859] bg-[#E4F7EE] text-[#008751]"
                          : "border-[#E4E8ED] bg-[#F4F6F8] text-[#5B6B7F] hover:border-[#00A859]"
                      }`}
                    >
                      <span className="text-2xl">{c.icon}</span>
                      {c.label}
                    </button>
                  ))}
                </div>
              </div>
              <div>
                <label className="mb-2 block text-sm font-semibold text-[#0B192C]">Project Name</label>
                <input
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  placeholder="Gold Mining Project"
                  required
                  className="w-full rounded-lg border border-[#E4E8ED] px-4 py-3 outline-none focus:border-[#00A859]"
                />
              </div>
              <div>
                <label className="mb-2 block text-sm font-semibold text-[#0B192C]">Total Token Supply</label>
                <input
                  type="number" min="1"
                  value={totalSupply}
                  onChange={(e) => setTotalSupply(e.target.value)}
                  placeholder="1000"
                  required
                  className="w-full rounded-lg border border-[#E4E8ED] px-4 py-3 outline-none focus:border-[#00A859]"
                />
              </div>
              <div>
                <label className="mb-2 block text-sm font-semibold text-[#0B192C]">Token Price</label>
                <input
                  type="number" min="1"
                  value={tokenPrice}
                  onChange={(e) => setTokenPrice(e.target.value)}
                  placeholder="100"
                  required
                  className="w-full rounded-lg border border-[#E4E8ED] px-4 py-3 outline-none focus:border-[#00A859]"
                />
              </div>
              <button
                type="submit"
                disabled={loading}
                className="w-full rounded-lg bg-[#00A859] px-4 py-3 font-semibold text-white transition hover:bg-[#008751] disabled:opacity-50"
              >
                {loading ? "Processing..." : "Create Project"}
              </button>
            </form>
          </section>

          <section className="rounded-2xl bg-white p-6 shadow-sm">
            <h3 style={{ fontFamily: "var(--font-sora)" }} className="mb-6 text-xl font-bold text-[#0B192C]">
              Distribute Revenue
            </h3>
            <form onSubmit={distributeRevenue} className="space-y-5">
              <div>
                <label className="mb-2 block text-sm font-semibold text-[#0B192C]">Project ID</label>
                <input
                  type="number" min="0"
                  value={projectId}
                  onChange={(e) => setProjectId(e.target.value)}
                  placeholder="0"
                  required
                  className="w-full rounded-lg border border-[#E4E8ED] px-4 py-3 outline-none focus:border-[#0B192C]"
                />
              </div>
              <div>
                <label className="mb-2 block text-sm font-semibold text-[#0B192C]">Revenue Amount</label>
                <input
                  type="number" min="1"
                  value={revenueAmount}
                  onChange={(e) => setRevenueAmount(e.target.value)}
                  placeholder="500"
                  required
                  className="w-full rounded-lg border border-[#E4E8ED] px-4 py-3 outline-none focus:border-[#0B192C]"
                />
              </div>
              <button
                type="submit"
                disabled={loading}
                className="w-full rounded-lg bg-[#0B192C] px-4 py-3 font-semibold text-white transition hover:bg-[#121F33] disabled:opacity-50"
              >
                {loading ? "Processing..." : "Distribute Revenue"}
              </button>
            </form>
          </section>
        </div>

        <section className="mt-8 rounded-2xl bg-white p-6 shadow-sm">
          <h3 style={{ fontFamily: "var(--font-sora)" }} className="mb-6 text-xl font-bold text-[#0B192C]">
            Projects
          </h3>
          {projects.length === 0 ? (
            <p className="text-[#5B6B7F]">No projects found.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm">
                <thead>
                  <tr className="border-b border-[#E4E8ED] text-[#5B6B7F]">
                    <th className="px-3 py-3">ID</th>
                    <th className="px-3 py-3">Project</th>
                    <th className="px-3 py-3">Field</th>
                    <th className="px-3 py-3">Supply</th>
                    <th className="px-3 py-3">Token Price</th>
                  </tr>
                </thead>
                <tbody>
                  {projects.map((project) => (
                    <tr key={project.id} className="border-b border-[#F4F6F8]">
                      <td className="px-3 py-3">{project.id}</td>
                      <td className="px-3 py-3 font-semibold text-[#0B192C]">{project.name}</td>
                      <td className="px-3 py-3">
                        <span className="rounded-full bg-[#E4F7EE] px-3 py-1 text-xs font-semibold text-[#008751]">
                          {categoryLabel(project.category)}
                        </span>
                      </td>
                      <td className="px-3 py-3">{project.total_supply}</td>
                      <td className="px-3 py-3">{project.token_price}</td>
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
