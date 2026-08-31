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

import React, { useMemo, useState } from "react";
import { Layers, Radio } from "lucide-react";
import { Slider } from "@/components/ui/slider";
import { AircraftState, ConflictRecord, VoxelKey } from "@/types/airspace";
import { UAE_VERTIPORTS } from "@/lib/scenario";

interface AirspaceRadarProps {
  aircraft: AircraftState[];
  conflicts: ConflictRecord[];
  selectedAircraftId: string | null;
  onSelectAircraft: (id: string | null) => void;
  selectedVoxel: VoxelKey | null;
  onSelectVoxel: (voxel: VoxelKey | null) => void;
}

// Bounding box for UAE Northern Emirates & Abu Dhabi airspace
const BOUNDS = {
  minLat: 24.3,
  maxLat: 25.5,
  minLon: 54.4,
  maxLon: 55.7,
};

export function AirspaceRadar({
  aircraft,
  conflicts,
  selectedAircraftId,
  onSelectAircraft,
  selectedVoxel,
  onSelectVoxel,
}: AirspaceRadarProps) {
  const [altitudeFilterFt, setAltitudeFilterFt] = useState<number[]>([0, 3000]);
  const [showHexGrid, setShowHexGrid] = useState<boolean>(true);
  const [showFlightPaths, setShowFlightPaths] = useState<boolean>(true);

  // SVG coordinate transformation
  const project = (lat: number, lon: number) => {
    const x = ((lon - BOUNDS.minLon) / (BOUNDS.maxLon - BOUNDS.minLon)) * 800;
    const y = ((BOUNDS.maxLat - lat) / (BOUNDS.maxLat - BOUNDS.minLat)) * 600;
    return { x, y };
  };

  // Generate H3 Hexagonal Grid wireframe centers across the corridor
  const hexGridCells = useMemo(() => {
    const cells = [];
    const latStep = 0.08;
    const lonStep = 0.09;

    let rowIdx = 0;
    for (let lat = BOUNDS.minLat + 0.05; lat <= BOUNDS.maxLat - 0.05; lat += latStep) {
      const offset = (rowIdx % 2) * (lonStep / 2);
      for (let lon = BOUNDS.minLon + 0.05; lon <= BOUNDS.maxLon - 0.05; lon += lonStep) {
        const { x, y } = project(lat, lon + offset);
        const h3Id = `8843a${Math.abs(Math.floor(lat * 1000 + lon * 100))
          .toString(16)
          .padEnd(6, "0")}`;
        cells.push({ x, y, lat, lon: lon + offset, h3Id });
      }
      rowIdx++;
    }
    return cells;
  }, []);

  // Compute hexagon polygon vertices with radius R
  const getHexPoints = (cx: number, cy: number, r: number = 24) => {
    const points = [];
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 3) * i - Math.PI / 6;
      points.push(`${cx + r * Math.cos(angle)},${cy + r * Math.sin(angle)}`);
    }
    return points.join(" ");
  };

  // Filter aircraft within active altitude range
  const visibleAircraft = aircraft.filter(
    (ac) => ac.altFt >= altitudeFilterFt[0] && ac.altFt <= altitudeFilterFt[1]
  );

  return (
    <div className="relative flex flex-col h-full w-full rounded-xl border border-border/80 bg-[#070A0F] overflow-hidden select-none">
      {/* Top Map Toolbar */}
      <div className="absolute top-3 left-3 right-3 z-10 flex flex-wrap items-center justify-between gap-2 pointer-events-none">
        {/* Airspace Region Pill */}
        <div className="flex items-center gap-2 rounded-lg border border-border/60 bg-background/90 backdrop-blur-md px-3 py-1.5 shadow-md pointer-events-auto">
          <Radio className="h-4 w-4 text-emerald-400 animate-pulse" />
          <div className="flex items-center gap-1.5 text-xs font-mono">
            <span className="font-semibold text-foreground">UAE UAM-CORRIDOR</span>
            <span className="text-muted-foreground">•</span>
            <span className="text-muted-foreground text-[11px]">H3 Res 8 (0.74 km²)</span>
          </div>
        </div>

        {/* Radar Controls & Altitude Filter */}
        <div className="flex items-center gap-2 rounded-lg border border-border/60 bg-background/90 backdrop-blur-md px-3 py-1.5 shadow-md pointer-events-auto">
          <div className="flex items-center gap-1.5 text-xs">
            <Layers className="h-3.5 w-3.5 text-muted-foreground" />
            <span className="text-muted-foreground text-[11px] font-mono">
              Alt Slices:
            </span>
            <span className="font-mono text-xs font-semibold text-sky-400">
              {altitudeFilterFt[0]} - {altitudeFilterFt[1]} ft
            </span>
          </div>

          <div className="w-28 mx-1">
            <Slider
              value={altitudeFilterFt}
              min={0}
              max={3000}
              step={100}
              onValueChange={(val) => {
                if (Array.isArray(val)) {
                  setAltitudeFilterFt([...val]);
                } else if (typeof val === "number") {
                  setAltitudeFilterFt([val, 3000]);
                }
              }}
              className="py-1 cursor-pointer"
            />
          </div>

          <div className="flex items-center gap-1 pl-2 border-l border-border/50">
            <button
              onClick={() => setShowHexGrid(!showHexGrid)}
              className={`px-2 py-0.5 rounded text-[11px] font-mono transition-colors ${
                showHexGrid
                  ? "bg-sky-500/20 text-sky-300 border border-sky-500/40"
                  : "text-muted-foreground hover:bg-muted"
              }`}
            >
              HexGrid
            </button>
            <button
              onClick={() => setShowFlightPaths(!showFlightPaths)}
              className={`px-2 py-0.5 rounded text-[11px] font-mono transition-colors ${
                showFlightPaths
                  ? "bg-emerald-500/20 text-emerald-300 border border-emerald-500/40"
                  : "text-muted-foreground hover:bg-muted"
              }`}
            >
              Paths
            </button>
          </div>
        </div>
      </div>

      {/* Main SVG Tactical Airspace Canvas */}
      <div className="flex-1 w-full h-full min-h-[500px]">
        <svg
          viewBox="0 0 800 600"
          className="w-full h-full"
          style={{
            background: "radial-gradient(ellipse at center, #0B111A 0%, #05080C 100%)",
          }}
        >
          <defs>
            {/* Radar Scope Rings Grid */}
            <pattern id="radar-grid" width="80" height="80" patternUnits="userSpaceOnUse">
              <path
                d="M 80 0 L 0 0 0 80"
                fill="none"
                stroke="rgba(255,255,255,0.03)"
                strokeWidth="1"
              />
            </pattern>
            {/* Conflict Glow Filter */}
            <filter id="hazard-glow" x="-20%" y="-20%" width="140%" height="140%">
              <feGaussianBlur stdDeviation="3" result="blur" />
              <feComposite in="SourceGraphic" in2="blur" operator="over" />
            </filter>
          </defs>

          {/* Background Scope Grid */}
          <rect width="800" height="600" fill="url(#radar-grid)" />

          {/* Compass Rings */}
          <circle
            cx="400"
            cy="300"
            r="180"
            fill="none"
            stroke="rgba(56,189,248,0.04)"
            strokeDasharray="3 6"
          />
          <circle cx="400" cy="300" r="280" fill="none" stroke="rgba(56,189,248,0.03)" />

          {/* H3 Hexagonal Grid Overlay */}
          {showHexGrid && (
            <g className="hex-grid-layer">
              {hexGridCells.map((cell) => {
                const isSelected = selectedVoxel?.h3Index === cell.h3Id;
                const hasAircraft = aircraft.some(
                  (a) => a.currentVoxel.h3Index === cell.h3Id
                );

                return (
                  <polygon
                    key={cell.h3Id}
                    points={getHexPoints(cell.x, cell.y, 22)}
                    fill={
                      isSelected
                        ? "rgba(56,189,248,0.25)"
                        : hasAircraft
                          ? "rgba(16,185,129,0.08)"
                          : "none"
                    }
                    stroke={
                      isSelected
                        ? "#38bdf8"
                        : hasAircraft
                          ? "rgba(16,185,129,0.3)"
                          : "rgba(255,255,255,0.06)"
                    }
                    strokeWidth={isSelected ? "1.5" : "0.75"}
                    className="transition-colors cursor-pointer hover:fill-sky-500/20 hover:stroke-sky-400"
                    onClick={() =>
                      onSelectVoxel({
                        h3Index: cell.h3Id,
                        altBinFt: Math.floor(altitudeFilterFt[0] / 100),
                        timeBinS: 0,
                        resolution: 8,
                      })
                    }
                  />
                );
              })}
            </g>
          )}

          {/* Flight Path Lines */}
          {showFlightPaths && (
            <g className="flight-trajectories">
              {aircraft.map((ac) => {
                const pathPoints = ac.trajectory.map((w) => {
                  const p = project(w.lat, w.lon);
                  return `${p.x},${p.y}`;
                });
                const isSelected = selectedAircraftId === ac.id;

                return (
                  <g key={`path-${ac.id}`}>
                    <polyline
                      points={pathPoints.join(" ")}
                      fill="none"
                      stroke={isSelected ? "#38bdf8" : "rgba(255,255,255,0.15)"}
                      strokeWidth={isSelected ? "2" : "1"}
                      strokeDasharray={isSelected ? "none" : "4 4"}
                    />
                    {ac.trajectory.map((w, idx) => {
                      const p = project(w.lat, w.lon);
                      return (
                        <circle
                          key={`wp-${ac.id}-${idx}`}
                          cx={p.x}
                          cy={p.y}
                          r="2.5"
                          fill={isSelected ? "#38bdf8" : "rgba(255,255,255,0.3)"}
                        />
                      );
                    })}
                  </g>
                );
              })}
            </g>
          )}

          {/* Active Conflict Markers & Convergence Vectors */}
          {conflicts.map((conf) => {
            const acA = aircraft.find((a) => a.id === conf.flightIdA);
            const acB = aircraft.find((a) => a.id === conf.flightIdB);
            if (!acA || !acB) return null;

            const pA = project(acA.lat, acA.lon);
            const pB = project(acB.lat, acB.lon);
            const midX = (pA.x + pB.x) / 2;
            const midY = (pA.y + pB.y) / 2;

            return (
              <g key={conf.id} className="conflict-hazard">
                {/* Distance Connector Line */}
                <line
                  x1={pA.x}
                  y1={pA.y}
                  x2={pB.x}
                  y2={pB.y}
                  stroke="#f43f5e"
                  strokeWidth="1.5"
                  strokeDasharray="3 3"
                />

                {/* Animated Pulsing Hazard Rings at Conflict Epicenter */}
                <circle
                  cx={midX}
                  cy={midY}
                  r="26"
                  fill="rgba(244,63,94,0.15)"
                  stroke="#f43f5e"
                  strokeWidth="1"
                  className="animate-ping opacity-60"
                />
                <circle
                  cx={midX}
                  cy={midY}
                  r="14"
                  fill="rgba(244,63,94,0.3)"
                  stroke="#f43f5e"
                  strokeWidth="1.5"
                  filter="url(#hazard-glow)"
                />

                {/* Conflict Tag */}
                <g transform={`translate(${midX + 16}, ${midY - 14})`}>
                  <rect
                    width="110"
                    height="32"
                    rx="4"
                    fill="rgba(15,23,42,0.95)"
                    stroke="#f43f5e"
                    strokeWidth="1"
                  />
                  <text
                    x="6"
                    y="13"
                    fill="#fda4af"
                    fontSize="9"
                    fontWeight="bold"
                    fontFamily="monospace"
                  >
                    {conf.conflictType === "SAME_VOXEL"
                      ? "⚠ SAME-VOXEL"
                      : "⚠ 18-NEIGHBOR"}
                  </text>
                  <text x="6" y="25" fill="#fecdd3" fontSize="8.5" fontFamily="monospace">
                    Dist: {conf.distanceMeters}m | P: {Math.round(conf.riskScore * 100)}%
                  </text>
                </g>
              </g>
            );
          })}

          {/* Vertiports (Major Hubs) */}
          {UAE_VERTIPORTS.map((vp) => {
            const { x, y } = project(vp.lat, vp.lon);
            return (
              <g key={vp.id} transform={`translate(${x}, ${y})`}>
                <circle r="7" fill="#0f172a" stroke="#38bdf8" strokeWidth="1.5" />
                <circle r="2.5" fill="#38bdf8" />
                <rect
                  x="10"
                  y="-12"
                  width="70"
                  height="20"
                  rx="3"
                  fill="rgba(15,23,42,0.85)"
                  stroke="rgba(255,255,255,0.1)"
                  strokeWidth="0.5"
                />
                <text
                  x="14"
                  y="2"
                  fill="#e2e8f0"
                  fontSize="9"
                  fontWeight="bold"
                  fontFamily="monospace"
                >
                  {vp.icao} ({vp.name.split(" ")[0]})
                </text>
              </g>
            );
          })}

          {/* Active Aircraft Targets */}
          {visibleAircraft.map((ac) => {
            const { x, y } = project(ac.lat, ac.lon);
            const isSelected = selectedAircraftId === ac.id;
            const inConflict = conflicts.some(
              (c) => c.flightIdA === ac.id || c.flightIdB === ac.id
            );

            // Vector heading line length ~ speed
            const headingRad = ((ac.headingDeg - 90) * Math.PI) / 180;
            const vx = x + 24 * Math.cos(headingRad);
            const vy = y + 24 * Math.sin(headingRad);

            return (
              <g
                key={ac.id}
                className="cursor-pointer"
                onClick={() => onSelectAircraft(isSelected ? null : ac.id)}
              >
                {/* Velocity Vector */}
                <line
                  x1={x}
                  y1={y}
                  x2={vx}
                  y2={vy}
                  stroke={inConflict ? "#f43f5e" : "#34d399"}
                  strokeWidth="1.5"
                />

                {/* Target Symbol */}
                <circle
                  cx={x}
                  cy={y}
                  r={isSelected ? "7" : "5"}
                  fill={inConflict ? "#f43f5e" : isSelected ? "#38bdf8" : "#10b981"}
                  stroke="#ffffff"
                  strokeWidth={isSelected ? "2" : "1"}
                  className={inConflict ? "animate-pulse" : ""}
                />

                {/* Aircraft Data Block (Callsign, Alt, Speed) */}
                <g transform={`translate(${x + 10}, ${y + 8})`}>
                  <rect
                    width="78"
                    height="28"
                    rx="3"
                    fill="rgba(15,23,42,0.92)"
                    stroke={
                      isSelected
                        ? "#38bdf8"
                        : inConflict
                          ? "#f43f5e"
                          : "rgba(255,255,255,0.15)"
                    }
                    strokeWidth="0.75"
                  />
                  <text
                    x="5"
                    y="11"
                    fill="#f8fafc"
                    fontSize="8.5"
                    fontWeight="bold"
                    fontFamily="monospace"
                  >
                    {ac.id} ({ac.destIcao})
                  </text>
                  <text x="5" y="22" fill="#94a3b8" fontSize="7.5" fontFamily="monospace">
                    {ac.altFt}ft • {ac.groundSpeedKt}kt
                  </text>
                </g>
              </g>
            );
          })}
        </svg>
      </div>

      {/* Bottom Telemetry Bar */}
      <div className="absolute bottom-3 left-3 right-3 z-10 flex items-center justify-between text-[11px] font-mono text-muted-foreground bg-background/90 backdrop-blur-md px-3 py-1.5 rounded-lg border border-border/60">
        <div className="flex items-center gap-3">
          <span>Lat: 24.30°N - 25.50°N</span>
          <span>Lon: 54.40°E - 55.70°E</span>
          {selectedVoxel && (
            <span className="text-sky-400 font-semibold">
              Selected Voxel: {selectedVoxel.h3Index} [Alt: {selectedVoxel.altBinFt * 100}
              ft]
            </span>
          )}
        </div>

        <div className="flex items-center gap-2">
          <span className="flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-emerald-400" /> Nominal
          </span>
          <span className="flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-rose-500" /> Conflict Hazard
          </span>
        </div>
      </div>
    </div>
  );
}
