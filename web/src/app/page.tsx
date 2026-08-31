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

import React, { useEffect, useState } from "react";
import { FileText, Fingerprint, ShieldAlert, Zap } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AirspaceRadar } from "@/components/AirspaceRadar";
import { AuditTrailViewer } from "@/components/AuditTrailViewer";
import { BenchmarkMetrics } from "@/components/BenchmarkMetrics";
import { ConflictRiskDrawer } from "@/components/ConflictRiskDrawer";
import { FlightPlanPanel } from "@/components/FlightPlanPanel";
import { Header } from "@/components/Header";
import {
  AircraftState,
  ConflictRecord,
  FlightPlanSubmission,
  VoxelKey,
} from "@/types/airspace";
import {
  advanceAircraftProgress,
  evaluateScenarioConflicts,
  getInitialScenarioAircraft,
  latLonToH3CellId,
} from "@/lib/scenario";
import { checkGatewayHealth } from "@/lib/api";

export default function MissionControlDashboard() {
  // Scenario state
  const [isPlaying, setIsPlaying] = useState<boolean>(true);
  const [speed, setSpeed] = useState<number>(1);
  const [aircraft, setAircraft] = useState<AircraftState[]>(() =>
    getInitialScenarioAircraft()
  );
  const [conflicts, setConflicts] = useState<ConflictRecord[]>([]);

  // Selection state
  const [selectedAircraftId, setSelectedAircraftId] = useState<string | null>(null);
  const [selectedConflictId, setSelectedConflictId] = useState<string | null>(null);
  const [selectedVoxel, setSelectedVoxel] = useState<VoxelKey | null>(null);
  const [activeTab, setActiveTab] = useState<string>("conflicts");

  // Gateway connection state
  const [isGatewayOnline, setIsGatewayOnline] = useState<boolean>(false);

  // Poll gateway health
  useEffect(() => {
    const checkHealth = async () => {
      const up = await checkGatewayHealth();
      setIsGatewayOnline(up);
    };
    checkHealth();
    const interval = setInterval(checkHealth, 5000);
    return () => clearInterval(interval);
  }, []);

  // Main simulation physics & conflict detection loop
  useEffect(() => {
    if (!isPlaying) return;

    const tickIntervalMs = 200;
    const progressStep = 0.003 * speed;

    const timer = setInterval(() => {
      setAircraft((prev) => {
        const nextFleet = prev.map((ac) => advanceAircraftProgress(ac, progressStep));
        const detected = evaluateScenarioConflicts(nextFleet);
        setConflicts(detected);
        return nextFleet;
      });
    }, tickIntervalMs);

    return () => clearInterval(timer);
  }, [isPlaying, speed]);

  const handleResetScenario = () => {
    setAircraft(getInitialScenarioAircraft());
    setSelectedAircraftId(null);
    setSelectedConflictId(null);
    setSelectedVoxel(null);
  };

  const handleFlightPlanSubmitted = (plan: FlightPlanSubmission) => {
    // Add custom aircraft to the simulation fleet
    const newAc: AircraftState = {
      id: plan.entityId,
      callsign: `FLT-${plan.entityId}`,
      lat: plan.waypoints[0].lat,
      lon: plan.waypoints[0].lon,
      altFt: plan.cruiseAltitudeFt,
      groundSpeedKt: plan.cruiseSpeedKt,
      headingDeg: 180,
      verticalSpeedFpm: 0,
      originIcao: plan.originIcao,
      destIcao: plan.destinationIcao,
      currentVoxel: {
        h3Index: latLonToH3CellId(plan.waypoints[0].lat, plan.waypoints[0].lon),
        altBinFt: Math.floor(plan.cruiseAltitudeFt / 100),
        timeBinS: Math.floor(Date.now() / 10000) % 1000,
        resolution: 8,
      },
      progressPercent: 0,
      status: "NOMINAL",
      trajectory: plan.waypoints,
    };

    setAircraft((prev) => [...prev.filter((a) => a.id !== plan.entityId), newAc]);
    setSelectedAircraftId(plan.entityId);
  };

  return (
    <div className="flex flex-col h-screen w-screen overflow-hidden bg-[#070A0F] text-foreground">
      {/* Top Mission Control Header */}
      <Header
        isPlaying={isPlaying}
        onTogglePlay={() => setIsPlaying(!isPlaying)}
        onReset={handleResetScenario}
        speed={speed}
        onSpeedChange={setSpeed}
        activeAircraftCount={aircraft.length}
        conflictCount={conflicts.length}
        isGatewayOnline={isGatewayOnline}
      />

      {/* Main Mission Control Body */}
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-3 p-3 min-h-0 overflow-hidden">
        {/* Left / Center 8 Columns: 4D Airspace Radar Scope Canvas */}
        <div className="lg:col-span-8 h-full flex flex-col min-h-[400px]">
          <AirspaceRadar
            aircraft={aircraft}
            conflicts={conflicts}
            selectedAircraftId={selectedAircraftId}
            onSelectAircraft={(id) => {
              setSelectedAircraftId(id);
              if (id) {
                const conf = conflicts.find(
                  (c) => c.flightIdA === id || c.flightIdB === id
                );
                if (conf) {
                  setSelectedConflictId(conf.id);
                  setActiveTab("conflicts");
                }
              }
            }}
            selectedVoxel={selectedVoxel}
            onSelectVoxel={setSelectedVoxel}
          />
        </div>

        {/* Right 4 Columns: Tabbed Telemetry & Strategic Control Panels */}
        <div className="lg:col-span-4 h-full flex flex-col min-h-0">
          <Tabs
            value={activeTab}
            onValueChange={setActiveTab}
            className="h-full flex flex-col"
          >
            {/* Tab Triggers */}
            <TabsList className="grid grid-cols-4 w-full h-9 bg-muted/30 border border-border/60 p-1 mb-2 font-mono text-[11px]">
              <TabsTrigger value="conflicts" className="gap-1 px-1 text-[10px]">
                <ShieldAlert className="h-3 w-3 text-rose-400" />
                Conflicts
              </TabsTrigger>
              <TabsTrigger value="plan" className="gap-1 px-1 text-[10px]">
                <FileText className="h-3 w-3 text-sky-400" />
                Plan (SCD)
              </TabsTrigger>
              <TabsTrigger value="audit" className="gap-1 px-1 text-[10px]">
                <Fingerprint className="h-3 w-3 text-emerald-400" />
                Audit
              </TabsTrigger>
              <TabsTrigger value="benchmark" className="gap-1 px-1 text-[10px]">
                <Zap className="h-3 w-3 text-amber-400" />
                Paper
              </TabsTrigger>
            </TabsList>

            {/* Tab Contents */}
            <div className="flex-1 min-h-0">
              <TabsContent
                value="conflicts"
                className="h-full m-0 data-[state=active]:flex data-[state=active]:flex-col"
              >
                <ConflictRiskDrawer
                  conflicts={conflicts}
                  selectedConflictId={selectedConflictId}
                  onSelectConflict={setSelectedConflictId}
                />
              </TabsContent>

              <TabsContent
                value="plan"
                className="h-full m-0 data-[state=active]:flex data-[state=active]:flex-col"
              >
                <FlightPlanPanel onFlightPlanSubmitted={handleFlightPlanSubmitted} />
              </TabsContent>

              <TabsContent
                value="audit"
                className="h-full m-0 data-[state=active]:flex data-[state=active]:flex-col"
              >
                <AuditTrailViewer />
              </TabsContent>

              <TabsContent
                value="benchmark"
                className="h-full m-0 data-[state=active]:flex data-[state=active]:flex-col"
              >
                <BenchmarkMetrics />
              </TabsContent>
            </div>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
