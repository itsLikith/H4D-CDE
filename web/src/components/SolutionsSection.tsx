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
import { ArrowRight, Building2, CheckCircle2, MapPin, Shield, Truck } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

export function SolutionsSection() {
  return (
    <section
      id="use-cases"
      className="py-24 px-4 sm:px-6 lg:px-8 bg-[#080C14] border-b border-slate-800/80"
    >
      <div className="max-w-7xl mx-auto space-y-16">
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto space-y-4">
          <Badge
            variant="outline"
            className="text-xs font-medium border-blue-500/30 text-blue-400 bg-blue-500/5 px-3 py-1"
          >
            Practical Applications
          </Badge>
          <h2 className="text-3xl sm:text-4xl font-bold tracking-tight text-white">
            Built for Real-World Low-Altitude Operations
          </h2>
          <p className="text-base text-slate-400 leading-relaxed">
            From commercial delivery fleets to regional vertiports and civil aviation
            authorities, H4D-CDE provides autonomous, scalable deconfliction out of the
            box.
          </p>
        </div>

        {/* Tabbed Use Cases Interface */}
        <Tabs defaultValue="logistics" className="w-full space-y-8">
          <div className="flex justify-center">
            <TabsList className="grid grid-cols-2 md:grid-cols-4 w-full max-w-3xl h-11 bg-slate-900/80 border border-slate-800 p-1">
              <TabsTrigger
                value="logistics"
                className="text-xs gap-2 data-[state=active]:bg-blue-600 data-[state=active]:text-white"
              >
                <Truck className="h-3.5 w-3.5" />
                Drone Delivery
              </TabsTrigger>
              <TabsTrigger
                value="ansp"
                className="text-xs gap-2 data-[state=active]:bg-blue-600 data-[state=active]:text-white"
              >
                <Building2 className="h-3.5 w-3.5" />
                Air Traffic & ANSPs
              </TabsTrigger>
              <TabsTrigger
                value="vertiports"
                className="text-xs gap-2 data-[state=active]:bg-blue-600 data-[state=active]:text-white"
              >
                <MapPin className="h-3.5 w-3.5" />
                Vertiports & eVTOL
              </TabsTrigger>
              <TabsTrigger
                value="defense"
                className="text-xs gap-2 data-[state=active]:bg-blue-600 data-[state=active]:text-white"
              >
                <Shield className="h-3.5 w-3.5" />
                Emergency Response
              </TabsTrigger>
            </TabsList>
          </div>

          {/* 1. Drone Delivery Content */}
          <TabsContent value="logistics" className="space-y-6">
            <Card className="border-slate-800 bg-slate-900/40 backdrop-blur-md p-6 sm:p-8">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
                <div className="space-y-4">
                  <Badge
                    variant="outline"
                    className="border-blue-500/30 text-blue-400 text-xs"
                  >
                    Autonomous Drone Delivery Fleets
                  </Badge>
                  <h3 className="text-2xl font-bold text-white tracking-tight">
                    Scale Last-Mile Drone Deliveries with Zero Flight Bottlenecks
                  </h3>
                  <p className="text-sm text-slate-300 leading-relaxed">
                    Automate high-density corridor scheduling, prevent near-misses before
                    takeoff, and dynamically reroute delivery drones around congested
                    urban zones in under 15 milliseconds.
                  </p>
                  <ul className="space-y-2.5 text-sm text-slate-300 pt-2">
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Instant automated pre-flight clearance against all active airspace
                        traffic
                      </span>
                    </li>
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Zero-fuel loss temporal delay recommendations for high-priority
                        dispatches
                      </span>
                    </li>
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        REST and WebSocket APIs for direct integration with GCS and
                        dispatch software
                      </span>
                    </li>
                  </ul>
                  <div className="pt-4">
                    <Link href="/console">
                      <Button className="bg-blue-600 hover:bg-blue-500 text-white gap-2">
                        Open Delivery Simulation <ArrowRight className="h-4 w-4" />
                      </Button>
                    </Link>
                  </div>
                </div>

                <div className="rounded-xl border border-slate-800 bg-slate-950 p-5 space-y-4 font-mono text-xs">
                  <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                    <span className="text-slate-400">Fleet Operations Monitor</span>
                    <span className="text-emerald-400">● 100% Nominal</span>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-slate-400">Active Fleet Size:</span>
                      <span className="text-white font-bold">1,250 UAVs</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Clearance Throughput:</span>
                      <span className="text-blue-400 font-bold">10,000 plans/sec</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Average Clearance Time:</span>
                      <span className="text-emerald-400 font-bold">12 ms</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Safety Standard:</span>
                      <span className="text-white">ASTM F3548-21 Pre-Flight SCD</span>
                    </div>
                  </div>
                </div>
              </div>
            </Card>
          </TabsContent>

          {/* 2. Air Traffic & ANSP Content */}
          <TabsContent value="ansp" className="space-y-6">
            <Card className="border-slate-800 bg-slate-900/40 backdrop-blur-md p-6 sm:p-8">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
                <div className="space-y-4">
                  <Badge
                    variant="outline"
                    className="border-emerald-500/30 text-emerald-400 text-xs"
                  >
                    Air Navigation Service Providers & Civil Aviation
                  </Badge>
                  <h3 className="text-2xl font-bold text-white tracking-tight">
                    Digital Airspace Management for Low-Altitude Traffic
                  </h3>
                  <p className="text-sm text-slate-300 leading-relaxed">
                    Safely integrate uncrewed aircraft into controlled low-altitude
                    airspace without increasing air traffic controller workload. Built to
                    global ASTM F3548-21 and ICAO Annex 11 specifications.
                  </p>
                  <ul className="space-y-2.5 text-sm text-slate-300 pt-2">
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Multi-tiered strategic deconfliction (same-altitude and adjacent
                        cell checks)
                      </span>
                    </li>
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Cryptographic SHA-256 audit ledger guarantees complete
                        non-repudiation
                      </span>
                    </li>
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        High availability deployment on standard cloud clusters or
                        on-premise infrastructure
                      </span>
                    </li>
                  </ul>
                  <div className="pt-4">
                    <Link href="/console">
                      <Button className="bg-emerald-600 hover:bg-emerald-500 text-white gap-2">
                        View Radar Scope <ArrowRight className="h-4 w-4" />
                      </Button>
                    </Link>
                  </div>
                </div>

                <div className="rounded-xl border border-slate-800 bg-slate-950 p-5 space-y-4 font-mono text-xs">
                  <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                    <span className="text-slate-400">Safety & Compliance</span>
                    <span className="text-emerald-400">● Verified</span>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-slate-400">Framework:</span>
                      <span className="text-white">ASTM F3548-21 / ICAO Annex 11</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Spatial Voxel Resolution:</span>
                      <span className="text-blue-400 font-bold">
                        Uber H3 Res 8 (0.74 km²)
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Vertical Separation:</span>
                      <span className="text-white font-bold">100 ft Discrete Floors</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Audit Proof:</span>
                      <span className="text-emerald-400 font-bold">
                        SHA-256 Event Chain
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </Card>
          </TabsContent>

          {/* 3. Vertiports Content */}
          <TabsContent value="vertiports" className="space-y-6">
            <Card className="border-slate-800 bg-slate-900/40 backdrop-blur-md p-6 sm:p-8">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
                <div className="space-y-4">
                  <Badge
                    variant="outline"
                    className="border-indigo-500/30 text-indigo-400 text-xs"
                  >
                    Vertiport Operators & Air Taxis (eVTOL)
                  </Badge>
                  <h3 className="text-2xl font-bold text-white tracking-tight">
                    Dynamic Space-Time Slot Reservations for Vertiport Hubs
                  </h3>
                  <p className="text-sm text-slate-300 leading-relaxed">
                    Coordinate arrival, departure, and approach funnels across competing
                    air taxi operators. Eliminate airborne holding patterns with
                    predictive 15-minute sector density forecasting.
                  </p>
                  <ul className="space-y-2.5 text-sm text-slate-300 pt-2">
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Guaranteed approach and departure corridor reservations without
                        overlap
                      </span>
                    </li>
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Predictive demand balancing smooths out peak morning and evening
                        traffic spikes
                      </span>
                    </li>
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Fair multi-operator slot allocation ensuring competitive equity
                      </span>
                    </li>
                  </ul>
                  <div className="pt-4">
                    <Link href="/console">
                      <Button className="bg-indigo-600 hover:bg-indigo-500 text-white gap-2">
                        View Vertiport Flow <ArrowRight className="h-4 w-4" />
                      </Button>
                    </Link>
                  </div>
                </div>

                <div className="rounded-xl border border-slate-800 bg-slate-950 p-5 space-y-4 font-mono text-xs">
                  <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                    <span className="text-slate-400">Vertiport Capacity Status</span>
                    <span className="text-emerald-400">● 92% Efficiency</span>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-slate-400">Target Hubs:</span>
                      <span className="text-white">Dubai (OMDB) & Al Maktoum (OMDW)</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Air Taxi Slots / Hour:</span>
                      <span className="text-indigo-400 font-bold">180 Operations</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Airborne Holding:</span>
                      <span className="text-emerald-400 font-bold">
                        Zero Holding Delay
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </Card>
          </TabsContent>

          {/* 4. Emergency Response Content */}
          <TabsContent value="defense" className="space-y-6">
            <Card className="border-slate-800 bg-slate-900/40 backdrop-blur-md p-6 sm:p-8">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
                <div className="space-y-4">
                  <Badge
                    variant="outline"
                    className="border-rose-500/30 text-rose-400 text-xs"
                  >
                    Emergency Services & Public Safety
                  </Badge>
                  <h3 className="text-2xl font-bold text-white tracking-tight">
                    Instant Geofencing & Priority Emergency Corridors
                  </h3>
                  <p className="text-sm text-slate-300 leading-relaxed">
                    Deploy emergency tactical no-fly zones in under a second.
                    Automatically notify and clear civilian drone traffic out of active
                    medical flights, firefighting, or emergency rescue zones.
                  </p>
                  <ul className="space-y-2.5 text-sm text-slate-300 pt-2">
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Sub-second dynamic geofence broadcasting to all active operators
                      </span>
                    </li>
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Automated evasive trajectory generation for rapid corridor
                        evacuation
                      </span>
                    </li>
                    <li className="flex items-center gap-2">
                      <CheckCircle2 className="h-4 w-4 text-emerald-400 shrink-0" />
                      <span>
                        Complete cryptographic chain of custody for incident investigation
                      </span>
                    </li>
                  </ul>
                  <div className="pt-4">
                    <Link href="/console">
                      <Button className="bg-rose-600 hover:bg-rose-500 text-white gap-2">
                        Simulate Geofence <ArrowRight className="h-4 w-4" />
                      </Button>
                    </Link>
                  </div>
                </div>

                <div className="rounded-xl border border-slate-800 bg-slate-950 p-5 space-y-4 font-mono text-xs">
                  <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                    <span className="text-slate-400">Emergency Geofence Protocol</span>
                    <span className="text-rose-400">● Ready</span>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span className="text-slate-400">Activation Speed:</span>
                      <span className="text-rose-400 font-bold">&lt; 200 ms</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Corridor Evacuation:</span>
                      <span className="text-emerald-400 font-bold">
                        Automated Reroute
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-slate-400">Authorization:</span>
                      <span className="text-white">Role-Based Access (RBAC)</span>
                    </div>
                  </div>
                </div>
              </div>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </section>
  );
}
