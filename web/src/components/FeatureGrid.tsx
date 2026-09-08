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
import { BrainCircuit, FileCheck2, Fingerprint, Layers } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";

export function FeatureGrid() {
  return (
    <section
      id="features"
      className="relative py-24 px-4 sm:px-6 lg:px-8 bg-[#0B0F19] border-b border-slate-800/80"
    >
      <div className="max-w-7xl mx-auto space-y-16">
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto space-y-4">
          <Badge
            variant="outline"
            className="text-xs font-medium border-emerald-500/30 text-emerald-400 bg-emerald-500/5 px-3 py-1"
          >
            Core Capabilities
          </Badge>
          <h2 className="text-3xl sm:text-4xl font-bold tracking-tight text-white">
            Engineered for Precision, Speed & Reliability
          </h2>
          <p className="text-base text-slate-400 leading-relaxed">
            Built from the ground up for sub-millisecond query performance, multi-fleet
            scaling, and strict ASTM F3548-21 compliance.
          </p>
        </div>

        {/* 4 Capabilities Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Capability 1 */}
          <Card className="border-slate-800 bg-slate-900/50 backdrop-blur-md p-6 sm:p-8 space-y-4 hover:border-blue-500/40 transition-all">
            <div className="flex items-center justify-between">
              <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-blue-500/10 border border-blue-500/20 text-blue-400">
                <Layers className="h-5 w-5" />
              </div>
              <Badge
                variant="outline"
                className="text-[10px] font-mono border-slate-700 text-slate-300"
              >
                Sub-Millisecond Engine
              </Badge>
            </div>

            <div className="space-y-2">
              <h3 className="text-xl font-bold text-white tracking-tight">
                High-Density 4D Spatial Discretization
              </h3>
              <p className="text-sm text-slate-300 leading-relaxed">
                Transforms complex 3D airspace volumes and continuous timestamps into
                discrete 4D space-time voxels. Eliminates quadratic computing bottlenecks
                and scales effortlessly to tens of thousands of concurrent flight plans.
              </p>
            </div>

            <div className="rounded-lg border border-slate-800 bg-slate-950 p-3.5 text-xs text-slate-300 space-y-1.5 font-mono">
              <div className="text-blue-400 font-semibold">Key Specifications:</div>
              <div>
                • Spatial Resolution:{" "}
                <span className="text-white">Uber H3 Level 8 (~0.74 km²)</span>
              </div>
              <div>
                • Vertical Floors:{" "}
                <span className="text-white">100 ft Altitude Bins</span>
              </div>
              <div>
                • Temporal Interval:{" "}
                <span className="text-white">10-Second Time Slices</span>
              </div>
            </div>
          </Card>

          {/* Capability 2 */}
          <Card className="border-slate-800 bg-slate-900/50 backdrop-blur-md p-6 sm:p-8 space-y-4 hover:border-emerald-500/40 transition-all">
            <div className="flex items-center justify-between">
              <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
                <FileCheck2 className="h-5 w-5" />
              </div>
              <Badge
                variant="outline"
                className="text-[10px] font-mono border-slate-700 text-slate-300"
              >
                ASTM F3548-21 Compliant
              </Badge>
            </div>

            <div className="space-y-2">
              <h3 className="text-xl font-bold text-white tracking-tight">
                Automated Strategic Deconfliction
              </h3>
              <p className="text-sm text-slate-300 leading-relaxed">
                Evaluates flight plans before departure to identify and resolve spatial
                overlap in advance. Generates instant, cost-optimized advisories that
                minimize fuel burn and operational delays.
              </p>
            </div>

            <div className="rounded-lg border border-slate-800 bg-slate-950 p-3.5 text-xs text-slate-300 space-y-1.5 font-mono">
              <div className="text-emerald-400 font-semibold">
                Advisory Priority Hierarchy:
              </div>
              <div>
                1. <span className="text-white">Temporal Hold (Δt)</span> ➔ Zero fuel
                penalty
              </div>
              <div>
                2. <span className="text-white">Flight Level Adjustment (Δh)</span> ➔
                Altitude separation
              </div>
              <div>
                3. <span className="text-white">Dynamic Lateral Reroute</span> ➔ Automated
                shortest path
              </div>
            </div>
          </Card>

          {/* Capability 3 */}
          <Card className="border-slate-800 bg-slate-900/50 backdrop-blur-md p-6 sm:p-8 space-y-4 hover:border-indigo-500/40 transition-all">
            <div className="flex items-center justify-between">
              <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-indigo-500/10 border border-indigo-500/20 text-indigo-400">
                <BrainCircuit className="h-5 w-5" />
              </div>
              <Badge
                variant="outline"
                className="text-[10px] font-mono border-slate-700 text-slate-300"
              >
                Predictive Intelligence
              </Badge>
            </div>

            <div className="space-y-2">
              <h3 className="text-xl font-bold text-white tracking-tight">
                Traffic & Risk Forecasting
              </h3>
              <p className="text-sm text-slate-300 leading-relaxed">
                Machine learning models forecast 15-minute sector density spikes, refine
                physical aircraft trajectory uncertainties, and calculate
                loss-of-separation probabilities in real time.
              </p>
            </div>

            <div className="grid grid-cols-3 gap-2.5 text-xs font-mono">
              <div className="rounded-lg border border-slate-800 bg-slate-950 p-2.5 text-center">
                <div className="text-blue-400 font-bold">Trajectory</div>
                <div className="text-[11px] text-slate-400">Path Refinement</div>
              </div>
              <div className="rounded-lg border border-slate-800 bg-slate-950 p-2.5 text-center">
                <div className="text-emerald-400 font-bold">Risk Scorer</div>
                <div className="text-[11px] text-slate-400">7-D Model</div>
              </div>
              <div className="rounded-lg border border-slate-800 bg-slate-950 p-2.5 text-center">
                <div className="text-indigo-400 font-bold">Flow Model</div>
                <div className="text-[11px] text-slate-400">15-Min Forecast</div>
              </div>
            </div>
          </Card>

          {/* Capability 4 */}
          <Card className="border-slate-800 bg-slate-900/50 backdrop-blur-md p-6 sm:p-8 space-y-4 hover:border-purple-500/40 transition-all">
            <div className="flex items-center justify-between">
              <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-purple-500/10 border border-purple-500/20 text-purple-400">
                <Fingerprint className="h-5 w-5" />
              </div>
              <Badge
                variant="outline"
                className="text-[10px] font-mono border-slate-700 text-slate-300"
              >
                Tamper-Proof Audit
              </Badge>
            </div>

            <div className="space-y-2">
              <h3 className="text-xl font-bold text-white tracking-tight">
                Cryptographic Event Ledger
              </h3>
              <p className="text-sm text-slate-300 leading-relaxed">
                Guarantees complete regulatory transparency. Every flight reservation,
                conflict alert, and deconfliction advisory is immutably cryptographically
                chained for post-flight verification.
              </p>
            </div>

            <div className="rounded-lg border border-slate-800 bg-slate-950 p-3.5 text-xs text-slate-300 space-y-1.5 font-mono">
              <div className="text-purple-400 font-semibold">Audit Guarantee:</div>
              <div>
                • Cryptographic Hash:{" "}
                <span className="text-white">SHA-256 Sequential Block Chain</span>
              </div>
              <div>
                • Non-Repudiation:{" "}
                <span className="text-emerald-400 font-bold">100% Tamper Detection</span>
              </div>
              <div>
                • Storage Architecture:{" "}
                <span className="text-white">Kafka Stream + TimescaleDB</span>
              </div>
            </div>
          </Card>
        </div>
      </div>
    </section>
  );
}
