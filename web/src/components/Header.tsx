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
import { AlertTriangle, Clock, Layers, Pause, Play, RotateCcw } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";

interface HeaderProps {
  isPlaying: boolean;
  onTogglePlay: () => void;
  onReset: () => void;
  speed: number;
  onSpeedChange: (speed: number) => void;
  activeAircraftCount: number;
  conflictCount: number;
  isGatewayOnline: boolean;
}

export function Header({
  isPlaying,
  onTogglePlay,
  onReset,
  speed,
  onSpeedChange,
  activeAircraftCount,
  conflictCount,
  isGatewayOnline,
}: HeaderProps) {
  const [utcTime, setUtcTime] = useState<string>("");

  useEffect(() => {
    const updateTime = () => {
      const now = new Date();
      setUtcTime(
        now.toISOString().slice(11, 19) + " UTC (" + now.toISOString().slice(0, 10) + ")"
      );
    };
    updateTime();
    const interval = setInterval(updateTime, 1000);
    return () => clearInterval(interval);
  }, []);

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/80 bg-background/95 backdrop-blur-md px-4 py-2.5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        {/* Brand & System Identifier */}
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 font-mono font-bold text-sm shadow-sm">
            H4D
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-sm font-semibold tracking-tight text-foreground">
                H4D-CDE
              </h1>
              <span className="text-xs text-muted-foreground hidden sm:inline">|</span>
              <span className="text-xs font-medium text-muted-foreground hidden sm:inline">
                Hexagonal 4D Conflict Detection Engine
              </span>
              <Badge
                variant="outline"
                className="text-[10px] font-mono px-1.5 py-0 h-4 border-emerald-500/40 text-emerald-400 bg-emerald-500/5"
              >
                ICSPIS 2025
              </Badge>
            </div>
            <div className="flex items-center gap-3 text-[11px] text-muted-foreground font-mono">
              <span className="flex items-center gap-1">
                <Clock className="h-3 w-3 text-muted-foreground/80" />
                {utcTime || "00:00:00 UTC"}
              </span>
              <span className="text-muted-foreground/40">•</span>
              <span>ASTM F3548-21 Compliant</span>
            </div>
          </div>
        </div>

        {/* Real-time Airspace Telemetry Status Pills */}
        <div className="flex items-center gap-2">
          {/* Active Entities */}
          <div className="flex items-center gap-1.5 rounded-md border border-border/60 bg-muted/30 px-2.5 py-1 text-xs">
            <Layers className="h-3.5 w-3.5 text-sky-400" />
            <span className="text-muted-foreground font-medium">UAVs:</span>
            <span className="font-mono font-semibold text-foreground">
              {activeAircraftCount}
            </span>
          </div>

          {/* Active Conflicts */}
          <div
            className={`flex items-center gap-1.5 rounded-md border px-2.5 py-1 text-xs transition-colors ${
              conflictCount > 0
                ? "border-rose-500/40 bg-rose-500/10 text-rose-300 animate-pulse"
                : "border-border/60 bg-muted/30 text-emerald-400"
            }`}
          >
            <AlertTriangle
              className={`h-3.5 w-3.5 ${
                conflictCount > 0 ? "text-rose-400" : "text-emerald-400"
              }`}
            />
            <span className="font-medium">Conflicts:</span>
            <span className="font-mono font-bold">{conflictCount}</span>
          </div>

          {/* Microservice Gateway Pulse */}
          <Tooltip>
            <TooltipTrigger>
              <div className="flex items-center gap-1.5 rounded-md border border-border/60 bg-muted/30 px-2.5 py-1 text-xs cursor-pointer">
                <span
                  className={`h-2 w-2 rounded-full ${
                    isGatewayOnline
                      ? "bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]"
                      : "bg-amber-400"
                  }`}
                />
                <span className="font-mono text-[11px] text-muted-foreground">
                  GW :8080
                </span>
              </div>
            </TooltipTrigger>
            <TooltipContent side="bottom">
              <div className="text-xs space-y-1">
                <p className="font-semibold">
                  API Gateway:{" "}
                  {isGatewayOnline ? "Connected (Healthy)" : "Standalone Demo"}
                </p>
                <p className="text-muted-foreground text-[11px]">
                  Fiber REST + WebSocket Hub + gRPC Backends
                </p>
              </div>
            </TooltipContent>
          </Tooltip>
        </div>

        {/* Mission Playback & Speed Controls */}
        <div className="flex items-center gap-1.5 bg-muted/40 p-1 rounded-lg border border-border/60">
          <Button
            variant={isPlaying ? "destructive" : "default"}
            size="sm"
            onClick={onTogglePlay}
            className="h-7 px-3 text-xs gap-1.5 font-medium"
          >
            {isPlaying ? (
              <>
                <Pause className="h-3.5 w-3.5" /> Pause
              </>
            ) : (
              <>
                <Play className="h-3.5 w-3.5" /> Run Scenario
              </>
            )}
          </Button>

          <Button
            variant="outline"
            size="sm"
            onClick={onReset}
            className="h-7 w-7 p-0"
            title="Reset to initial UAE 3-UAV scenario"
          >
            <RotateCcw className="h-3.5 w-3.5 text-muted-foreground" />
          </Button>

          <div className="flex items-center pl-1 pr-1 gap-1 text-[11px] font-mono text-muted-foreground border-l border-border/50 ml-1">
            {[1, 2, 5].map((s) => (
              <button
                key={s}
                onClick={() => onSpeedChange(s)}
                className={`px-1.5 py-0.5 rounded text-[10px] font-semibold transition-colors ${
                  speed === s
                    ? "bg-foreground/15 text-foreground"
                    : "hover:bg-muted text-muted-foreground"
                }`}
              >
                {s}x
              </button>
            ))}
          </div>
        </div>
      </div>
    </header>
  );
}
