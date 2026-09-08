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
import { Cpu, Sparkles, Zap } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function BenchmarkMetrics() {
  return (
    <Card className="h-full border-border/80 bg-card/60 backdrop-blur-md flex flex-col">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm font-semibold flex items-center gap-2">
            <Zap className="h-4 w-4 text-amber-400" />
            Empirical Benchmark (ICSPIS 2025 Reproduction)
          </CardTitle>
          <Badge
            variant="outline"
            className="text-[10px] font-mono border-emerald-500/40 text-emerald-400 bg-emerald-500/5"
          >
            99.19% Speedup
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-4 flex-1 overflow-y-auto text-xs p-3 font-mono">
        {/* Table I Reproduction: Computational Performance */}
        <div className="rounded-lg border border-border/80 bg-background/50 p-2.5 space-y-2">
          <div className="flex items-center justify-between text-[11px]">
            <span className="font-semibold text-foreground font-sans flex items-center gap-1.5">
              <Cpu className="h-3.5 w-3.5 text-sky-400" />
              Table I — Scalability: Naive O(n²) vs H4D-CDE O(n)
            </span>
            <span className="text-[10px] text-emerald-400 font-bold">
              13,964x Fewer Ops
            </span>
          </div>

          <div className="space-y-1.5 text-[10.5px]">
            <div className="flex items-center justify-between p-1.5 rounded bg-muted/40">
              <span className="text-muted-foreground">Small Fleet (N=2,793 points):</span>
              <span className="text-foreground font-semibold">
                0.0181s ➔ <strong className="text-emerald-400">0.0111s</strong> (-38.7%)
              </span>
            </div>

            <div className="flex items-center justify-between p-1.5 rounded bg-muted/40">
              <span className="text-muted-foreground">
                Large Fleet (N=27,930 points):
              </span>
              <span className="text-foreground font-semibold">
                1.6505s ➔ <strong className="text-emerald-400">0.0134s</strong> (-99.19%)
              </span>
            </div>

            <div className="flex items-center justify-between p-1.5 rounded bg-muted/40">
              <span className="text-muted-foreground">Theoretical Complexity:</span>
              <span className="text-sky-400 font-bold">
                O(n²) Quadratic ➔ O(n) Linear
              </span>
            </div>
          </div>
        </div>

        {/* Table III Reproduction: AI Augmentation Benchmarks */}
        <div className="rounded-lg border border-border/80 bg-background/50 p-2.5 space-y-2">
          <div className="flex items-center justify-between text-[11px]">
            <span className="font-semibold text-foreground font-sans flex items-center gap-1.5">
              <Sparkles className="h-3.5 w-3.5 text-amber-400" />
              Table III — AI Model Accuracy Benchmarks
            </span>
            <span className="text-[10px] text-emerald-400 font-bold">
              All Targets Met
            </span>
          </div>

          <div className="space-y-1.5 text-[10.5px]">
            <div className="flex items-center justify-between p-1.5 rounded bg-muted/40">
              <div>
                <span className="text-foreground font-semibold">
                  Trajectory Predictor (GBM):
                </span>
                <p className="text-[9.5px] text-muted-foreground font-sans">
                  Target: MAE ≤ 15.2 m
                </p>
              </div>
              <Badge
                variant="outline"
                className="text-[10px] border-emerald-500/40 text-emerald-400 font-bold"
              >
                10.05 m [PASS]
              </Badge>
            </div>

            <div className="flex items-center justify-between p-1.5 rounded bg-muted/40">
              <div>
                <span className="text-foreground font-semibold">
                  Risk Scorer (XGBoost):
                </span>
                <p className="text-[9.5px] text-muted-foreground font-sans">
                  Target: AUC-ROC ≥ 0.89
                </p>
              </div>
              <Badge
                variant="outline"
                className="text-[10px] border-emerald-500/40 text-emerald-400 font-bold"
              >
                0.900 [PASS]
              </Badge>
            </div>

            <div className="flex items-center justify-between p-1.5 rounded bg-muted/40">
              <div>
                <span className="text-foreground font-semibold">
                  Demand Forecaster (TCN):
                </span>
                <p className="text-[9.5px] text-muted-foreground font-sans">
                  Target: MAPE ≤ 8.70%
                </p>
              </div>
              <Badge
                variant="outline"
                className="text-[10px] border-emerald-500/40 text-emerald-400 font-bold"
              >
                4.94% [PASS]
              </Badge>
            </div>
          </div>
        </div>

        {/* Microservice Latency Metrics */}
        <div className="text-[10px] text-muted-foreground font-sans flex items-center justify-between pt-1 border-t border-border/40">
          <span>Target Throughput: 10,000 flight plans/sec</span>
          <span className="text-emerald-400 font-semibold font-mono">P99 &lt; 2.5ms</span>
        </div>
      </CardContent>
    </Card>
  );
}
