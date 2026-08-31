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
import { CheckCircle2, Gauge, ShieldAlert, Sparkles, TrendingUp } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ConflictRecord } from "@/types/airspace";

interface ConflictRiskDrawerProps {
  conflicts: ConflictRecord[];
  selectedConflictId: string | null;
  onSelectConflict: (id: string | null) => void;
}

export function ConflictRiskDrawer({
  conflicts,
  selectedConflictId,
  onSelectConflict,
}: ConflictRiskDrawerProps) {
  const activeConflict =
    conflicts.find((c) => c.id === selectedConflictId) || conflicts[0];

  return (
    <Card className="h-full border-border/80 bg-card/60 backdrop-blur-md flex flex-col overflow-hidden">
      <CardHeader className="pb-2.5">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm font-semibold flex items-center gap-2">
            <ShieldAlert className="h-4 w-4 text-rose-400" />
            AI Conflict Diagnostics (XGBoost & TCN)
          </CardTitle>
          <Badge
            variant={conflicts.length > 0 ? "destructive" : "outline"}
            className="text-[10px] font-mono"
          >
            {conflicts.length} Active {conflicts.length === 1 ? "Pair" : "Pairs"}
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-3.5 flex-1 overflow-y-auto text-xs p-3">
        {conflicts.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-8 text-center text-muted-foreground">
            <CheckCircle2 className="h-8 w-8 text-emerald-400/60 mb-2" />
            <p className="font-semibold text-xs text-foreground">Airspace 4D Nominal</p>
            <p className="text-[11px] max-w-[220px]">
              No same-voxel or 18-neighbor voxel separation violations detected.
            </p>
          </div>
        ) : (
          <>
            {/* Conflict Selection Pills */}
            <div className="flex items-center gap-1.5 overflow-x-auto pb-1">
              {conflicts.map((c) => {
                const isSelected = activeConflict?.id === c.id;
                return (
                  <button
                    key={c.id}
                    onClick={() => onSelectConflict(c.id)}
                    className={`flex items-center gap-1.5 px-2.5 py-1 rounded-md text-[11px] font-mono border transition-all ${
                      isSelected
                        ? "border-rose-500 bg-rose-500/20 text-rose-200 font-semibold"
                        : "border-border/60 bg-muted/30 text-muted-foreground hover:bg-muted"
                    }`}
                  >
                    <span className="h-1.5 w-1.5 rounded-full bg-rose-400" />
                    {c.flightIdA} ↔ {c.flightIdB}
                  </button>
                );
              })}
            </div>

            {activeConflict && (
              <div className="space-y-3">
                {/* Risk Score Gauge & Severity */}
                <div className="rounded-lg border border-border/80 bg-background/50 p-3 space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-[11px] font-mono text-muted-foreground">
                      XGBoost P(Conflict):
                    </span>
                    <Badge
                      variant={
                        activeConflict.severity === "CRITICAL" ? "destructive" : "outline"
                      }
                      className="text-[10px] font-mono font-bold"
                    >
                      {activeConflict.severity} (
                      {Math.round(activeConflict.riskScore * 100)}%)
                    </Badge>
                  </div>

                  {/* Visual Progress Bar */}
                  <div className="h-2 w-full rounded-full bg-muted overflow-hidden">
                    <div
                      className={`h-full rounded-full transition-all duration-300 ${
                        activeConflict.riskScore > 0.75
                          ? "bg-rose-500 shadow-[0_0_8px_rgba(244,63,94,0.6)]"
                          : activeConflict.riskScore > 0.45
                            ? "bg-amber-400"
                            : "bg-emerald-400"
                      }`}
                      style={{ width: `${activeConflict.riskScore * 100}%` }}
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-2 text-[10px] font-mono text-muted-foreground pt-1">
                    <div>
                      <span>Separation: </span>
                      <strong className="text-foreground">
                        {activeConflict.distanceMeters} m
                      </strong>
                    </div>
                    <div>
                      <span>Closure Rate: </span>
                      <strong className="text-foreground">
                        {activeConflict.closureRateMps} m/s
                      </strong>
                    </div>
                    <div>
                      <span>Voxel: </span>
                      <strong className="text-sky-400">
                        {activeConflict.conflictType === "SAME_VOXEL"
                          ? "Same Hex Cell"
                          : "18-Neighbor"}
                      </strong>
                    </div>
                    <div>
                      <span>Time to CPA: </span>
                      <strong className="text-foreground">
                        {activeConflict.timeToCPASeconds} s
                      </strong>
                    </div>
                  </div>
                </div>

                {/* 7-Dimensional ML Feature Vector */}
                <div className="rounded-lg border border-border/80 bg-background/50 p-2.5 space-y-2">
                  <div className="flex items-center justify-between text-[11px]">
                    <span className="font-semibold text-foreground flex items-center gap-1.5">
                      <Gauge className="h-3.5 w-3.5 text-sky-400" />
                      7-D ML Feature Vector (Eq. 10)
                    </span>
                    <span className="font-mono text-[10px] text-muted-foreground">
                      AUC = 0.900
                    </span>
                  </div>

                  <div className="grid grid-cols-2 gap-1.5 text-[10px] font-mono">
                    <div className="rounded bg-muted/40 p-1.5">
                      <span className="text-muted-foreground">Local Density:</span>{" "}
                      <strong className="text-foreground">
                        {activeConflict.features.localTrafficDensity} / hex
                      </strong>
                    </div>
                    <div className="rounded bg-muted/40 p-1.5">
                      <span className="text-muted-foreground">Heading Diff:</span>{" "}
                      <strong className="text-foreground">
                        {activeConflict.features.headingDiffDeg}°
                      </strong>
                    </div>
                    <div className="rounded bg-muted/40 p-1.5">
                      <span className="text-muted-foreground">Wind Shear:</span>{" "}
                      <strong className="text-foreground">
                        {activeConflict.features.windShearKt} kt/100ft
                      </strong>
                    </div>
                    <div className="rounded bg-muted/40 p-1.5">
                      <span className="text-muted-foreground">Visibility:</span>{" "}
                      <strong className="text-foreground">
                        {activeConflict.features.visibilityKm} km
                      </strong>
                    </div>
                  </div>
                </div>

                {/* Dilated TCN 15-Minute Forecast */}
                <div className="rounded-lg border border-sky-500/30 bg-sky-500/5 p-2.5 space-y-1.5">
                  <div className="flex items-center justify-between text-[11px]">
                    <span className="font-semibold text-sky-300 flex items-center gap-1">
                      <TrendingUp className="h-3.5 w-3.5 text-sky-400" />
                      TCN Demand Forecast (Eq. 9)
                    </span>
                    <Badge
                      variant="outline"
                      className="text-[9px] font-mono border-sky-400/40 text-sky-300"
                    >
                      MAPE: 4.94%
                    </Badge>
                  </div>
                  <p className="text-[10px] text-muted-foreground font-mono">
                    Predicted Sector Load:{" "}
                    <strong className="text-sky-300">
                      {activeConflict.features.sectorLoadForecast} flights/voxel
                    </strong>{" "}
                    in +15 min horizon.
                  </p>
                </div>

                {/* Automated Advisory Recommendation */}
                {activeConflict.advisory && (
                  <div className="rounded-lg border border-emerald-500/30 bg-emerald-500/5 p-2.5 space-y-1">
                    <div className="flex items-center justify-between">
                      <span className="font-semibold text-emerald-400 text-[11px] flex items-center gap-1">
                        <Sparkles className="h-3.5 w-3.5" /> Resolution Advisory
                      </span>
                      <span className="font-mono text-[9px] text-muted-foreground">
                        Target: {activeConflict.advisory.targetEntityId}
                      </span>
                    </div>
                    <p className="text-[10px] text-foreground font-sans leading-relaxed">
                      {activeConflict.advisory.rationale}
                    </p>
                  </div>
                )}
              </div>
            )}
          </>
        )}
      </CardContent>
    </Card>
  );
}
