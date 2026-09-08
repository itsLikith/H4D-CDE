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
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";

export function FaqSection() {
  return (
    <section
      id="faq"
      className="py-24 px-4 sm:px-6 lg:px-8 bg-[#080C14] border-b border-slate-800/80"
    >
      <div className="max-w-4xl mx-auto space-y-12">
        {/* Section Header */}
        <div className="text-center space-y-4">
          <Badge
            variant="outline"
            className="text-xs font-medium border-slate-700 text-slate-300 bg-slate-800/50 px-3 py-1"
          >
            Frequently Asked Questions
          </Badge>
          <h2 className="text-3xl sm:text-4xl font-bold tracking-tight text-white">
            Architecture & Integration FAQ
          </h2>
          <p className="text-base text-slate-400">
            Common questions regarding system connectivity, capacity, and deconfliction
            logic.
          </p>
        </div>

        {/* shadcn Accordion */}
        <Accordion className="w-full space-y-3">
          <AccordionItem
            value="item-1"
            className="border border-slate-800 bg-slate-900/50 rounded-xl px-5 transition-colors"
          >
            <AccordionTrigger className="text-left text-sm sm:text-base font-semibold text-white hover:text-blue-400 hover:no-underline py-4">
              How do fleet operators connect to H4D-CDE?
            </AccordionTrigger>
            <AccordionContent className="text-sm text-slate-300 leading-relaxed pb-4">
              Operators integrate via standard REST endpoints and bi-directional
              WebSockets at the API Gateway. Flight plans are submitted in standard ASTM
              F3548-21 format or via simple JSON waypoints. Clearances and deconfliction
              advisories are returned in sub-15ms directly to your dispatch system or
              Ground Control Station (GCS).
            </AccordionContent>
          </AccordionItem>

          <AccordionItem
            value="item-2"
            className="border border-slate-800 bg-slate-900/50 rounded-xl px-5 transition-colors"
          >
            <AccordionTrigger className="text-left text-sm sm:text-base font-semibold text-white hover:text-blue-400 hover:no-underline py-4">
              How does the 4D voxel discretization prevent bottlenecks?
            </AccordionTrigger>
            <AccordionContent className="text-sm text-slate-300 leading-relaxed pb-4">
              By mapping 3D coordinates and timestamps into discrete space-time cells
              (Uber H3 cells, 100 ft altitude bins, and 10-second time slices), H4D-CDE
              indexes flight paths in constant time. Instead of comparing every flight
              plan against every other flight plan pairwise, the engine only evaluates
              aircraft occupying identical or neighboring 4D voxels.
            </AccordionContent>
          </AccordionItem>

          <AccordionItem
            value="item-3"
            className="border border-slate-800 bg-slate-900/50 rounded-xl px-5 transition-colors"
          >
            <AccordionTrigger className="text-left text-sm sm:text-base font-semibold text-white hover:text-blue-400 hover:no-underline py-4">
              How are conflict resolution advisories computed?
            </AccordionTrigger>
            <AccordionContent className="text-sm text-slate-300 leading-relaxed pb-4">
              When a conflict is detected, the automated cascade evaluates three options
              in order of operational cost: first, a brief temporal ground hold (incurring
              zero fuel consumption); second, an altitude tier step change; and third, a
              lateral shortest-path reroute around the congested cell.
            </AccordionContent>
          </AccordionItem>

          <AccordionItem
            value="item-4"
            className="border border-slate-800 bg-slate-900/50 rounded-xl px-5 transition-colors"
          >
            <AccordionTrigger className="text-left text-sm sm:text-base font-semibold text-white hover:text-blue-400 hover:no-underline py-4">
              How does the SHA-256 audit ledger ensure non-repudiation?
            </AccordionTrigger>
            <AccordionContent className="text-sm text-slate-300 leading-relaxed pb-4">
              Every flight reservation, conflict detection event, and advisory is streamed
              through Kafka and hashed into an append-only SHA-256 cryptographic chain
              stored in TimescaleDB. Any attempt to modify or delete a historical record
              invalidates the hash chain, enabling instant detection of data tampering.
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </div>
    </section>
  );
}
