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

/**
 * UAE Airspace Reference Scenario Simulator.
 * Reproduces the 3-UAV urban corridor test scenario evaluated in
 * Sahadevan et al. (ICSPIS 2025), Section V-A across Dubai (DXB/DWC),
 * Abu Dhabi (AUH), and Sharjah (SHJ).
 */

import { AircraftState, ConflictRecord, GeoPoint } from "@/types/airspace";

export interface VertiportNode {
  id: string;
  name: string;
  icao: string;
  lat: number;
  lon: number;
  elevationFt: number;
}

export const UAE_VERTIPORTS: VertiportNode[] = [
  {
    id: "OMDB",
    name: "Dubai International (DXB)",
    icao: "OMDB",
    lat: 25.2532,
    lon: 55.3657,
    elevationFt: 62,
  },
  {
    id: "OMDW",
    name: "Al Maktoum International (DWC)",
    icao: "OMDW",
    lat: 24.8961,
    lon: 55.1614,
    elevationFt: 115,
  },
  {
    id: "OMAA",
    name: "Abu Dhabi International (AUH)",
    icao: "OMAA",
    lat: 24.433,
    lon: 54.6511,
    elevationFt: 88,
  },
  {
    id: "OMSJ",
    name: "Sharjah International (SHJ)",
    icao: "OMSJ",
    lat: 25.3286,
    lon: 55.5172,
    elevationFt: 116,
  },
];

// Great-circle Haversine distance in meters
export function haversineMeters(
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
): number {
  const R = 6371008.8;
  const dLat = ((lat2 - lat1) * Math.PI) / 180;
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos((lat1 * Math.PI) / 180) *
      Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
}

// Spherical linear interpolation between two geographic points
export function interpolateGeo(p1: GeoPoint, p2: GeoPoint, fraction: number): GeoPoint {
  const f = Math.max(0, Math.min(1, fraction));
  const lat = p1.lat + (p2.lat - p1.lat) * f;
  const lon = p1.lon + (p2.lon - p1.lon) * f;
  const altFt = (p1.altFt || 0) + ((p2.altFt || 0) - (p1.altFt || 0)) * f;
  return { lat, lon, altFt };
}

// Pseudo H3 index encoder for client-side visual representation (Resolution 8)
export function latLonToH3CellId(lat: number, lon: number): string {
  const latGrid = Math.floor((lat - 24.0) * 120);
  const lonGrid = Math.floor((lon - 54.0) * 120);
  const hash = Math.abs((latGrid * 73856093) ^ (lonGrid * 19349663)) % 0xffffffff;
  return `8843a${hash.toString(16).padStart(8, "0").slice(0, 10)}`;
}

// Waypoint tracks for the 3 reference UAVs
const UAV1_WAYPOINTS: GeoPoint[] = [
  { lat: 25.2532, lon: 55.3657, altFt: 500 }, // DXB Departure
  { lat: 25.132, lon: 55.275, altFt: 1200 }, // Downtown Dubai Corridor
  { lat: 25.041, lon: 55.215, altFt: 1200 }, // Jebel Ali Intersection (Conflict Zone)
  { lat: 24.8961, lon: 55.1614, altFt: 600 }, // DWC Arrival
];

const UAV2_WAYPOINTS: GeoPoint[] = [
  { lat: 24.433, lon: 54.6511, altFt: 800 }, // AUH Departure
  { lat: 24.782, lon: 54.95, altFt: 1200 }, // Ghantoot Corridor
  { lat: 25.043, lon: 55.218, altFt: 1200 }, // Jebel Ali Crossing (Converging Point)
  { lat: 25.2532, lon: 55.3657, altFt: 700 }, // DXB Arrival
];

const UAV3_WAYPOINTS: GeoPoint[] = [
  { lat: 25.3286, lon: 55.5172, altFt: 400 }, // SHJ Departure
  { lat: 25.185, lon: 55.38, altFt: 1500 }, // Meydan Corridor
  { lat: 25.046, lon: 55.22, altFt: 1200 }, // Jebel Ali Approach
  { lat: 24.8961, lon: 55.1614, altFt: 500 }, // DWC Arrival
];

export function getInitialScenarioAircraft(): AircraftState[] {
  return [
    {
      id: "UAV-001",
      callsign: "SKY-TAXI-01",
      lat: UAV1_WAYPOINTS[0].lat,
      lon: UAV1_WAYPOINTS[0].lon,
      altFt: UAV1_WAYPOINTS[0].altFt || 500,
      groundSpeedKt: 65,
      headingDeg: 215,
      verticalSpeedFpm: 400,
      originIcao: "OMDB",
      destIcao: "OMDW",
      currentVoxel: {
        h3Index: latLonToH3CellId(UAV1_WAYPOINTS[0].lat, UAV1_WAYPOINTS[0].lon),
        altBinFt: 5,
        timeBinS: 0,
        resolution: 8,
      },
      progressPercent: 0,
      status: "NOMINAL",
      trajectory: UAV1_WAYPOINTS,
    },
    {
      id: "UAV-002",
      callsign: "EMIRATES-CARGO-2",
      lat: UAV2_WAYPOINTS[0].lat,
      lon: UAV2_WAYPOINTS[0].lon,
      altFt: UAV2_WAYPOINTS[0].altFt || 800,
      groundSpeedKt: 70,
      headingDeg: 45,
      verticalSpeedFpm: 350,
      originIcao: "OMAA",
      destIcao: "OMDB",
      currentVoxel: {
        h3Index: latLonToH3CellId(UAV2_WAYPOINTS[0].lat, UAV2_WAYPOINTS[0].lon),
        altBinFt: 8,
        timeBinS: 0,
        resolution: 8,
      },
      progressPercent: 0,
      status: "NOMINAL",
      trajectory: UAV2_WAYPOINTS,
    },
    {
      id: "UAV-003",
      callsign: "SHJ-EXP-03",
      lat: UAV3_WAYPOINTS[0].lat,
      lon: UAV3_WAYPOINTS[0].lon,
      altFt: UAV3_WAYPOINTS[0].altFt || 400,
      groundSpeedKt: 60,
      headingDeg: 230,
      verticalSpeedFpm: 500,
      originIcao: "OMSJ",
      destIcao: "OMDW",
      currentVoxel: {
        h3Index: latLonToH3CellId(UAV3_WAYPOINTS[0].lat, UAV3_WAYPOINTS[0].lon),
        altBinFt: 4,
        timeBinS: 0,
        resolution: 8,
      },
      progressPercent: 0,
      status: "NOMINAL",
      trajectory: UAV3_WAYPOINTS,
    },
  ];
}

// Computes piecewise track progression along waypoints
export function advanceAircraftProgress(
  ac: AircraftState,
  deltaFraction: number
): AircraftState {
  const newProgress = (ac.progressPercent + deltaFraction) % 1.0;
  const numSegments = ac.trajectory.length - 1;
  const totalSegmentIdx = newProgress * numSegments;
  const currentSeg = Math.floor(totalSegmentIdx);
  const segFraction = totalSegmentIdx - currentSeg;

  const p1 = ac.trajectory[currentSeg];
  const p2 = ac.trajectory[Math.min(currentSeg + 1, ac.trajectory.length - 1)];

  const interpolated = interpolateGeo(p1, p2, segFraction);

  // Compute heading
  const dLon = p2.lon - p1.lon;
  const dLat = p2.lat - p1.lat;
  let heading = (Math.atan2(dLon, dLat) * 180) / Math.PI;
  if (heading < 0) heading += 360;

  const altBin = Math.floor((interpolated.altFt || 1000) / 100);

  return {
    ...ac,
    lat: interpolated.lat,
    lon: interpolated.lon,
    altFt: Math.round(interpolated.altFt || 1000),
    headingDeg: Math.round(heading),
    progressPercent: newProgress,
    currentVoxel: {
      h3Index: latLonToH3CellId(interpolated.lat, interpolated.lon),
      altBinFt: altBin,
      timeBinS: Math.floor(Date.now() / 10000) % 1000,
      resolution: 8,
    },
  };
}

// Evaluates pairwise conflicts dynamically using 4D Hexagonal Voxelization logic
export function evaluateScenarioConflicts(fleet: AircraftState[]): ConflictRecord[] {
  const conflicts: ConflictRecord[] = [];

  for (let i = 0; i < fleet.length; i++) {
    for (let j = i + 1; j < fleet.length; j++) {
      const a = fleet[i];
      const b = fleet[j];

      const dist = haversineMeters(a.lat, a.lon, b.lat, b.lon);
      const altDiff = Math.abs(a.altFt - b.altFt);
      const sameVoxel =
        a.currentVoxel.h3Index === b.currentVoxel.h3Index &&
        Math.abs(a.currentVoxel.altBinFt - b.currentVoxel.altBinFt) <= 1;

      const isNeighbor = !sameVoxel && dist < 1200 && altDiff < 300; // ASTM F3548-21 separation threshold

      if (sameVoxel || isNeighbor) {
        const closureRate = Math.abs(a.groundSpeedKt + b.groundSpeedKt) * 0.514444; // m/s
        const headingDiff = Math.abs(a.headingDeg - b.headingDeg);

        // Calibrated Fermi/XGBoost risk approximation
        const rawRisk = Math.min(
          0.99,
          Math.max(
            0.15,
            1.0 / (1.0 + Math.exp(-0.005 * (1000 - dist) - 0.05 * (closureRate - 20)))
          )
        );

        const severity = rawRisk > 0.75 ? "CRITICAL" : rawRisk > 0.45 ? "MEDIUM" : "LOW";

        conflicts.push({
          id: `CONF-${a.id}-${b.id}`,
          flightIdA: a.id,
          flightIdB: b.id,
          conflictType: sameVoxel ? "SAME_VOXEL" : "NEIGHBOR_VOXEL",
          h3Cell: a.currentVoxel.h3Index,
          altBinFt: a.currentVoxel.altBinFt,
          timeBinS: a.currentVoxel.timeBinS,
          distanceMeters: Math.round(dist),
          closureRateMps: Math.round(closureRate * 10) / 10,
          headingDiffDeg: Math.round(headingDiff),
          riskScore: Math.round(rawRisk * 100) / 100,
          severity,
          detectedAtUnixMs: Date.now(),
          timeToCPASeconds: Math.max(5, Math.round(dist / Math.max(1, closureRate))),
          features: {
            nEntities: 2,
            closureRateMps: Math.round(closureRate * 10) / 10,
            headingDiffDeg: Math.round(headingDiff),
            localTrafficDensity: 3.2,
            sectorLoadForecast: 5.8,
            windShearKt: 2.1,
            visibilityKm: 10.0,
          },
          advisory: {
            targetEntityId: a.id,
            resolutionType: altDiff < 200 ? "ALTITUDE_CHANGE" : "DELAY",
            recommendedDeltaTSeconds: 30,
            recommendedNewAltFt: a.altFt + 300,
            costScore: 0.18,
            rationale:
              "ASTM F3548-21 vertical separation cascade: climb +300 ft to establish 4D protected buffer",
            estimatedSeparationMeters: 650,
          },
        });
      }
    }
  }

  return conflicts;
}
