// Copyright 2026 H4D-CDE Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

"use client";

import React from "react";
import Link from "next/link";
import { ArrowRight, ChevronRight, Terminal } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

export function HeroSection() {
  return (
    <section className="relative overflow-hidden bg-[#080C14] pt-16 pb-20 border-b border-slate-800/80">
      {/* Subtle Background Grid */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#1e293b10_1px,transparent_1px),linear-gradient(to_bottom,#1e293b10_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)] pointer-events-none" />

      {/* Subtle Blue Radial Glow */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[350px] bg-blue-600/10 blur-[120px] rounded-full pointer-events-none" />

      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col items-center text-center space-y-8 max-w-4xl mx-auto">
          {/* Status Badge */}
          <div className="inline-flex items-center gap-2 rounded-full border border-blue-500/20 bg-blue-500/10 px-4 py-1.5 text-xs font-medium text-blue-300">
            <span className="flex h-2 w-2 rounded-full bg-blue-400" />
            <span>Autonomous Airspace Management & UTM Engine</span>
          </div>

          {/* Main Headline */}
          <div className="space-y-4">
            <h1 className="text-4xl sm:text-6xl font-extrabold tracking-tight text-white leading-[1.12]">
              Autonomous 4D{" "}
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-sky-300 to-indigo-300">
                Airspace Orchestration
              </span>
            </h1>
            <p className="text-lg sm:text-xl text-slate-300 max-w-2xl mx-auto leading-relaxed font-normal">
              A fast, deterministic conflict detection and resolution engine for drone
              fleets, air taxis, and next-generation low-altitude operations.
            </p>
          </div>

          {/* Action CTAs */}
          <div className="flex flex-wrap items-center justify-center gap-4 pt-2">
            <Link href="/console">
              <Button
                size="lg"
                className="h-12 px-7 text-sm font-semibold gap-2 bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-600/25 border border-blue-400/30 transition-all cursor-pointer"
              >
                <Terminal className="h-4 w-4" />
                Launch Live Console
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>

            <a href="#how-it-works">
              <Button
                variant="outline"
                size="lg"
                className="h-12 px-7 text-sm font-medium gap-2 border-slate-700 bg-slate-900/80 hover:bg-slate-800 text-slate-200"
              >
                How It Works
              </Button>
            </a>
          </div>

          {/* Key Product Metrics */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 w-full pt-10 text-left">
            <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-5 space-y-1.5 shadow-sm">
              <div className="text-xs font-medium text-slate-400 uppercase tracking-wider">
                Throughput
              </div>
              <div className="text-3xl font-bold text-white tracking-tight">10,000+</div>
              <div className="text-xs text-slate-400">Flight plans processed / sec</div>
            </div>

            <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-5 space-y-1.5 shadow-sm">
              <div className="text-xs font-medium text-slate-400 uppercase tracking-wider">
                Detection Speed
              </div>
              <div className="text-3xl font-bold text-blue-400 tracking-tight">
                &lt; 15 ms
              </div>
              <div className="text-xs text-slate-400">Deterministic query latency</div>
            </div>

            <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-5 space-y-1.5 shadow-sm">
              <div className="text-xs font-medium text-slate-400 uppercase tracking-wider">
                Separation Assurance
              </div>
              <div className="text-3xl font-bold text-emerald-400 tracking-tight">
                100%
              </div>
              <div className="text-xs text-slate-400">
                Zero unmanaged proximity hazards
              </div>
            </div>

            <div className="rounded-xl border border-slate-800 bg-slate-900/50 p-5 space-y-1.5 shadow-sm">
              <div className="text-xs font-medium text-slate-400 uppercase tracking-wider">
                Audit Trail
              </div>
              <div className="text-3xl font-bold text-indigo-400 tracking-tight">
                SHA-256
              </div>
              <div className="text-xs text-slate-400">Immutable cryptographic ledger</div>
            </div>
          </div>
        </div>

        {/* Live Corridor Status Preview Card */}
        <div className="mt-16 rounded-2xl border border-slate-800 bg-slate-950/80 p-3 sm:p-4 shadow-2xl overflow-hidden">
          <div className="flex items-center justify-between border-b border-slate-800/80 pb-3 px-2">
            <div className="flex items-center gap-2">
              <div className="h-3 w-3 rounded-full bg-rose-500/80" />
              <div className="h-3 w-3 rounded-full bg-amber-500/80" />
              <div className="h-3 w-3 rounded-full bg-emerald-500/80" />
              <span className="ml-2 text-xs font-mono text-slate-400">
                Airspace Surveillance & Deconfliction Monitor
              </span>
            </div>
            <div className="flex items-center gap-3">
              <span className="flex items-center gap-1.5 text-xs text-emerald-400 font-medium">
                <span className="h-2 w-2 rounded-full bg-emerald-400 animate-pulse" />
                Live Monitoring Active
              </span>
              <Link href="/console">
                <Button
                  size="sm"
                  variant="ghost"
                  className="h-7 text-xs text-blue-400 hover:text-blue-300"
                >
                  Open Live Radar <ChevronRight className="h-3 w-3 ml-1" />
                </Button>
              </Link>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 pt-4">
            {/* Sector Card 1 */}
            <div className="rounded-xl border border-slate-800/80 bg-slate-900/60 p-4 space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium text-slate-400">
                  Corridor Density
                </span>
                <Badge
                  variant="outline"
                  className="text-[10px] border-blue-500/30 text-blue-400"
                >
                  Dubai ➔ Abu Dhabi
                </Badge>
              </div>
              <div className="space-y-2">
                <div className="flex justify-between text-xs">
                  <span className="text-slate-300 font-mono">OMDB ➔ OMDW</span>
                  <span className="text-emerald-400 font-mono">Nominal (14 UAVs)</span>
                </div>
                <div className="h-1.5 w-full bg-slate-800 rounded-full overflow-hidden">
                  <div className="h-full bg-emerald-500 rounded-full w-[38%]" />
                </div>
                <div className="flex justify-between text-xs">
                  <span className="text-slate-300 font-mono">OMAA ➔ OMDB</span>
                  <span className="text-blue-400 font-mono">Moderate (26 UAVs)</span>
                </div>
                <div className="h-1.5 w-full bg-slate-800 rounded-full overflow-hidden">
                  <div className="h-full bg-blue-500 rounded-full w-[64%]" />
                </div>
              </div>
            </div>

            {/* Sector Card 2 */}
            <div className="rounded-xl border border-slate-800/80 bg-slate-900/60 p-4 space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium text-slate-400">
                  Automated Deconfliction
                </span>
                <Badge
                  variant="outline"
                  className="text-[10px] border-emerald-500/30 text-emerald-400"
                >
                  Active
                </Badge>
              </div>
              <div className="space-y-1.5 text-xs text-slate-300">
                <div className="flex items-center justify-between">
                  <span className="text-slate-400">Resolution Algorithm:</span>
                  <span className="text-white font-mono">Multi-Objective Cascade</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-slate-400">Advisory Priority:</span>
                  <span className="text-emerald-400 font-mono">
                    Temporal Hold (0s delay)
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-slate-400">Cleared Flight Plans:</span>
                  <span className="text-white font-mono">1,420 today (100% Cleared)</span>
                </div>
              </div>
            </div>

            {/* Sector Card 3 */}
            <div className="rounded-xl border border-slate-800/80 bg-slate-900/60 p-4 space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium text-slate-400">Audit Ledger</span>
                <Badge
                  variant="outline"
                  className="text-[10px] border-indigo-500/30 text-indigo-400"
                >
                  Verified
                </Badge>
              </div>
              <div className="space-y-1.5 text-xs text-slate-300">
                <div className="flex items-center justify-between">
                  <span className="text-slate-400">Ledger Status:</span>
                  <span className="text-emerald-400 font-mono font-medium">
                    Chain Intact
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-slate-400">Recent Proof Hash:</span>
                  <span className="text-slate-400 font-mono text-[10px]">
                    7f8e...3a1c
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-slate-400">Compliance Standard:</span>
                  <span className="text-white font-mono">ASTM F3548-21</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
