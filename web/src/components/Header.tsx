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
import Image from "next/image";
import Link from "next/link";
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
    <header className="sticky top-0 z-50 w-full border-b border-slate-800 bg-[#080C14]/95 backdrop-blur-md px-4 py-2.5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        {/* Brand & System Identifier */}
        <div className="flex items-center gap-3">
          <Link
            href="/"
            className="flex h-9 w-9 items-center justify-center rounded-lg border border-slate-700 bg-slate-900 p-1 shadow-sm hover:border-blue-500 transition-colors"
            title="Return to Overview"
          >
            <Image
              src="/logo.png"
              alt="H4D-CDE Logo"
              width={32}
              height={32}
              className="h-full w-full object-contain"
            />
          </Link>
          <div>
            <div className="flex items-center gap-2">
              <Link
                href="/"
                className="text-sm font-bold tracking-tight text-white hover:text-blue-400 transition-colors"
              >
                H4D-CDE
              </Link>
              <span className="text-xs text-slate-600 hidden sm:inline">|</span>
              <span className="text-xs font-medium text-slate-300 hidden sm:inline">
                Airspace Operations Console
              </span>
              <Badge
                variant="outline"
                className="text-[10px] font-mono px-1.5 py-0 h-4 border-blue-500/30 text-blue-400 bg-blue-500/5"
              >
                Live Radar
              </Badge>
            </div>
            <div className="flex items-center gap-3 text-[11px] text-slate-400 font-mono">
              <span className="flex items-center gap-1">
                <Clock className="h-3 w-3 text-slate-500" />
                {utcTime || "00:00:00 UTC"}
              </span>
              <span className="text-slate-700">•</span>
              <span>ASTM F3548-21 4D SCD</span>
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
