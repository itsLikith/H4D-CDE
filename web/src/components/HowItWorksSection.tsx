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
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";

export function HowItWorksSection() {
  const steps = [
    {
      number: "01",
      title: "4D Space-Time Discretization",
      description:
        "Continuous 3D trajectories and timestamps (lat, lon, alt, time) are partitioned into discrete 4D voxels using Uber H3 Resolution 8 cells (~0.74 km²), 100 ft altitude floors, and 10-second temporal bins.",
      badge: "Spatial Partitioning",
      detail: "Eliminates quadratic pairwise checks",
    },
    {
      number: "02",
      title: "Sub-15ms Conflict Detection",
      description:
        "The engine indexes flight paths across both direct voxel occupancy (Stage A) and 18-neighbor adjacent proximity zones (Stage B), identifying potential loss-of-separation in milliseconds.",
      badge: "Real-Time Indexing",
      detail: "10,000+ flight plans / second",
    },
    {
      number: "03",
      title: "Automated Resolution Advisories",
      description:
        "When an intersection is detected, the multi-objective optimizer computes the lowest-cost resolution: prioritizing zero-fuel temporal holds first, altitude tier adjustments second, and lateral reroutes third.",
      badge: "Cost-Optimized",
      detail: "Minimizes delay & fuel burn",
    },
    {
      number: "04",
      title: "Immutable SHA-256 Audit Trail",
      description:
        "Every flight reservation, conflict detection event, and resolution advisory is cryptographically hashed and committed to an immutable append-only ledger for non-repudiation and regulatory compliance.",
      badge: "Non-Repudiation",
      detail: "100% Tamper-evident verification",
    },
  ];

  return (
    <section
      id="how-it-works"
      className="py-24 px-4 sm:px-6 lg:px-8 bg-[#0B0F19] border-b border-slate-800/80"
    >
      <div className="max-w-7xl mx-auto space-y-16">
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto space-y-4">
          <Badge
            variant="outline"
            className="text-xs font-medium border-blue-500/30 text-blue-400 bg-blue-500/5 px-3 py-1"
          >
            System Architecture
          </Badge>
          <h2 className="text-3xl sm:text-4xl font-bold tracking-tight text-white">
            How H4D-CDE Orchestrates Airspace
          </h2>
          <p className="text-base text-slate-400 leading-relaxed">
            From raw trajectory submission to automated clearance and cryptographic audit,
            the platform executes a four-stage deconfliction pipeline in sub-15ms.
          </p>
        </div>

        {/* 4 Steps Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {steps.map((step, idx) => (
            <Card
              key={idx}
              className="border-slate-800 bg-slate-900/50 backdrop-blur-md p-6 space-y-4 flex flex-col justify-between hover:border-blue-500/40 transition-all"
            >
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-2xl font-black font-mono text-blue-400">
                    {step.number}
                  </span>
                  <Badge
                    variant="outline"
                    className="text-[10px] font-mono border-slate-700 text-slate-300"
                  >
                    {step.badge}
                  </Badge>
                </div>

                <div className="space-y-2">
                  <h3 className="text-base font-bold text-white tracking-tight">
                    {step.title}
                  </h3>
                  <p className="text-xs text-slate-300 leading-relaxed">
                    {step.description}
                  </p>
                </div>
              </div>

              <div className="pt-3 border-t border-slate-800/80 flex items-center justify-between text-[11px] font-mono text-slate-400">
                <span>{step.detail}</span>
              </div>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
