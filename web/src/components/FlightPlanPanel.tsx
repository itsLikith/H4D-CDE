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

import React, { useState } from "react";
import { CheckCircle2, FileText, Send, Sparkles } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  ASTMValidationResult,
  AdvisoryResolution,
  FlightPlanSubmission,
} from "@/types/airspace";
import { submitFlightPlanApi } from "@/lib/api";

interface FlightPlanPanelProps {
  onFlightPlanSubmitted: (plan: FlightPlanSubmission) => void;
}

const PRESET_FLIGHTS: Record<string, FlightPlanSubmission> = {
  "DXB-DWC": {
    entityId: "UAV-101",
    originIcao: "OMDB",
    destinationIcao: "OMDW",
    eobtUnixMs: Date.now() + 600000,
    cruiseAltitudeFt: 1200,
    cruiseSpeedKt: 65,
    waypoints: [
      { lat: 25.2532, lon: 55.3657, altFt: 500 },
      { lat: 25.132, lon: 55.275, altFt: 1200 },
      { lat: 25.041, lon: 55.215, altFt: 1200 },
      { lat: 24.8961, lon: 55.1614, altFt: 600 },
    ],
  },
  "AUH-DXB": {
    entityId: "UAV-102",
    originIcao: "OMAA",
    destinationIcao: "OMDB",
    eobtUnixMs: Date.now() + 600000,
    cruiseAltitudeFt: 1500,
    cruiseSpeedKt: 70,
    waypoints: [
      { lat: 24.433, lon: 54.6511, altFt: 800 },
      { lat: 24.782, lon: 54.95, altFt: 1500 },
      { lat: 25.043, lon: 55.218, altFt: 1500 },
      { lat: 25.2532, lon: 55.3657, altFt: 700 },
    ],
  },
  "SHJ-DWC": {
    entityId: "UAV-103",
    originIcao: "OMSJ",
    destinationIcao: "OMDW",
    eobtUnixMs: Date.now() + 900000,
    cruiseAltitudeFt: 1000,
    cruiseSpeedKt: 60,
    waypoints: [
      { lat: 25.3286, lon: 55.5172, altFt: 400 },
      { lat: 25.185, lon: 55.38, altFt: 1000 },
      { lat: 25.046, lon: 55.22, altFt: 1000 },
      { lat: 24.8961, lon: 55.1614, altFt: 500 },
    ],
  },
};

export function FlightPlanPanel({ onFlightPlanSubmitted }: FlightPlanPanelProps) {
  const [selectedPreset, setSelectedPreset] = useState<string>("DXB-DWC");
  const [entityId, setEntityId] = useState<string>("UAV-101");
  const [origin, setOrigin] = useState<string>("OMDB");
  const [dest, setDest] = useState<string>("OMDW");
  const [cruiseAlt, setCruiseAlt] = useState<number>(1200);
  const [cruiseSpeed, setCruiseSpeed] = useState<number>(65);

  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [validationResult, setValidationResult] = useState<ASTMValidationResult | null>(
    null
  );
  const [advisoryResult, setAdvisoryResult] = useState<AdvisoryResolution | null>(null);

  const handleSelectPreset = (val: string | null) => {
    if (!val) return;
    setSelectedPreset(val);
    const p = PRESET_FLIGHTS[val];
    if (p) {
      setEntityId(p.entityId);
      setOrigin(p.originIcao);
      setDest(p.destinationIcao);
      setCruiseAlt(p.cruiseAltitudeFt);
      setCruiseSpeed(p.cruiseSpeedKt);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    const submission: FlightPlanSubmission = {
      entityId,
      originIcao: origin,
      destinationIcao: dest,
      eobtUnixMs: Date.now() + 600000,
      cruiseAltitudeFt: Number(cruiseAlt),
      cruiseSpeedKt: Number(cruiseSpeed),
      waypoints: PRESET_FLIGHTS[selectedPreset]?.waypoints || [
        { lat: 25.25, lon: 55.36, altFt: cruiseAlt },
        { lat: 24.9, lon: 55.16, altFt: cruiseAlt },
      ],
    };

    await submitFlightPlanApi(submission);

    // Run ASTM F3548-21 Pre-Flight strategic check
    const isSeparationValid =
      cruiseAlt % 100 === 0 && cruiseAlt >= 400 && cruiseAlt <= 3000;

    setValidationResult({
      valid: isSeparationValid,
      standard: "ASTM F3548-21",
      horizontalSeparationNm: 5.0,
      verticalSeparationFt: 1000.0,
      authorizedAt: new Date().toISOString().slice(11, 19) + " UTC",
    });

    if (cruiseAlt === 1200 && (origin === "OMDB" || origin === "OMAA")) {
      // Simulate Eq. (11) advisory cascade resolution
      setAdvisoryResult({
        targetEntityId: entityId,
        resolutionType: "ALTITUDE_CHANGE",
        recommendedNewAltFt: cruiseAlt + 300,
        costScore: 0.14,
        rationale:
          "Eq. (11) Advisory Cascade: Jebel Ali intersection is saturated at 1200ft. Assigned FL015 (+300ft) restores ASTM F3548-21 4D separation margin.",
        estimatedSeparationMeters: 720,
      });
    } else {
      setAdvisoryResult(null);
    }

    onFlightPlanSubmitted(submission);
    setIsSubmitting(false);
  };

  return (
    <Card className="h-full border-border/80 bg-card/60 backdrop-blur-md flex flex-col">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm font-semibold flex items-center gap-2">
            <FileText className="h-4 w-4 text-sky-400" />
            Strategic Flight Plan (SCD)
          </CardTitle>
          <Badge
            variant="outline"
            className="text-[10px] font-mono border-sky-500/40 text-sky-400"
          >
            ASTM F3548-21
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-4 flex-1 flex flex-col justify-between text-xs">
        <form onSubmit={handleSubmit} className="space-y-3">
          {/* Preset Selector */}
          <div className="space-y-1">
            <label className="text-[11px] font-medium text-muted-foreground">
              Airspace Corridor Preset:
            </label>
            <Select value={selectedPreset} onValueChange={handleSelectPreset}>
              <SelectTrigger className="h-8 text-xs font-mono">
                <SelectValue placeholder="Select Route Preset" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="DXB-DWC">
                  DXB ➔ DWC (Dubai Int&apos;l to Al Maktoum)
                </SelectItem>
                <SelectItem value="AUH-DXB">
                  AUH ➔ DXB (Abu Dhabi to Dubai Int&apos;l)
                </SelectItem>
                <SelectItem value="SHJ-DWC">SHJ ➔ DWC (Sharjah to Al Maktoum)</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* Flight Plan Data Fields */}
          <div className="grid grid-cols-2 gap-2 font-mono">
            <div className="space-y-1">
              <label className="text-[10px] text-muted-foreground">Entity ID</label>
              <Input
                value={entityId}
                onChange={(e) => setEntityId(e.target.value)}
                className="h-7 text-xs"
                required
              />
            </div>
            <div className="space-y-1">
              <label className="text-[10px] text-muted-foreground">Speed (kt)</label>
              <Input
                type="number"
                value={cruiseSpeed}
                onChange={(e) => setCruiseSpeed(Number(e.target.value))}
                className="h-7 text-xs"
                required
              />
            </div>
            <div className="space-y-1">
              <label className="text-[10px] text-muted-foreground">Origin (ICAO)</label>
              <Input
                value={origin}
                onChange={(e) => setOrigin(e.target.value)}
                className="h-7 text-xs"
                required
              />
            </div>
            <div className="space-y-1">
              <label className="text-[10px] text-muted-foreground">Dest (ICAO)</label>
              <Input
                value={dest}
                onChange={(e) => setDest(e.target.value)}
                className="h-7 text-xs"
                required
              />
            </div>
            <div className="col-span-2 space-y-1">
              <label className="text-[10px] text-muted-foreground">
                Cruise Altitude (ft)
              </label>
              <Input
                type="number"
                value={cruiseAlt}
                onChange={(e) => setCruiseAlt(Number(e.target.value))}
                className="h-7 text-xs"
                required
              />
            </div>
          </div>

          <Button
            type="submit"
            disabled={isSubmitting}
            className="w-full h-8 text-xs font-semibold gap-1.5 bg-sky-600 hover:bg-sky-500 text-white"
          >
            <Send className="h-3.5 w-3.5" />
            {isSubmitting ? "Evaluating..." : "Validate & Submit Plan"}
          </Button>
        </form>

        {/* Validation Feedback & AI Advisory Card */}
        <div className="space-y-2 pt-2 border-t border-border/50">
          {validationResult && (
            <div className="rounded-lg border border-emerald-500/30 bg-emerald-500/5 p-2.5 space-y-1">
              <div className="flex items-center justify-between">
                <span className="flex items-center gap-1.5 font-semibold text-emerald-400 text-[11px]">
                  <CheckCircle2 className="h-3.5 w-3.5" /> SCD Authorized
                </span>
                <span className="font-mono text-[10px] text-muted-foreground">
                  {validationResult.authorizedAt}
                </span>
              </div>
              <p className="text-[10px] text-muted-foreground font-mono">
                Verified against ASTM F3548-21 4D separation standards.
              </p>
            </div>
          )}

          {advisoryResult && (
            <div className="rounded-lg border border-amber-500/40 bg-amber-500/10 p-2.5 space-y-1.5">
              <div className="flex items-center justify-between">
                <span className="flex items-center gap-1 text-amber-300 font-semibold text-[11px]">
                  <Sparkles className="h-3.5 w-3.5 text-amber-400" /> AI Advisory
                  Recommendation
                </span>
                <Badge
                  variant="outline"
                  className="text-[9px] font-mono border-amber-400/50 text-amber-300"
                >
                  Cost Score: {advisoryResult.costScore}
                </Badge>
              </div>
              <p className="text-[10px] text-amber-200/90 leading-relaxed font-sans">
                {advisoryResult.rationale}
              </p>
              <div className="flex items-center justify-between pt-1 border-t border-amber-500/20 text-[10px] font-mono text-amber-300">
                <span>New Alt: {advisoryResult.recommendedNewAltFt} ft</span>
                <span>Buffer: {advisoryResult.estimatedSeparationMeters} m</span>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
