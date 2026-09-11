"use client";

import Image from "next/image";

export default function SiteHeader({
  active,
}: {
  active: "sponsor" | "investor";
}) {
  return (
    <nav className="sticky top-0 z-50 border-b border-[#E4E8ED] bg-[#F4F6F8]/90 backdrop-blur-md">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <a href="/" className="flex items-center gap-3">
          <Image src="/logo.png" alt="PakTokenX" width={38} height={38} style={{ width: "auto", height: "38px" }} />

          <span
            style={{ fontFamily: "var(--font-sora)" }}
            className="text-xl font-bold text-[#0B192C]"
          >
            PakToken<span className="text-[#00A859]">X</span>
          </span>
        </a>

        <div className="flex items-center gap-3">
          <a
            href="/"
            className={`rounded-lg px-4 py-2 text-sm font-semibold transition ${
              active === "sponsor"
                ? "bg-[#0B192C] text-white"
                : "text-[#5B6B7F] hover:text-[#0B192C]"
            }`}
          >
            Sponsor Portal
          </a>

          <a
            href="/investor"
            className={`rounded-lg px-4 py-2 text-sm font-semibold transition ${
              active === "investor"
                ? "bg-[#00A859] text-white"
                : "text-[#5B6B7F] hover:text-[#0B192C]"
            }`}
          >
            Investor Portal
          </a>
        </div>
      </div>
    </nav>
  );
}
