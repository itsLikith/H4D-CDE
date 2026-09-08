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
import { Fingerprint, Hash, RotateCcw } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AuditBlock } from "@/types/airspace";
import { generateSampleAuditTrail } from "@/lib/api";

export function AuditTrailViewer() {
  const [blocks, setBlocks] = useState<AuditBlock[]>(() => generateSampleAuditTrail());
  const [tamperedIndex, setTamperedIndex] = useState<number | null>(null);

  // Simulates malicious payload alteration to prove hash continuity
  const handleTamperBlock = (idx: number) => {
    setTamperedIndex(idx);
    setBlocks((prev) =>
      prev.map((b, i) => {
        if (i === idx) {
          return {
            ...b,
            eventType: "ALTERED_EVENT_UNAUTHORIZED",
            isValid: false,
            tampered: true,
          };
        }
        if (i > idx) {
          // Downstream blocks fail verification because prevHash no longer matches
          return {
            ...b,
            isValid: false,
          };
        }
        return b;
      })
    );
  };

  const handleResetAudit = () => {
    setTamperedIndex(null);
    setBlocks(generateSampleAuditTrail());
  };

  const isChainValid = blocks.every((b) => b.isValid);

  return (
    <Card className="h-full border-border/80 bg-card/60 backdrop-blur-md flex flex-col">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm font-semibold flex items-center gap-2">
            <Fingerprint className="h-4 w-4 text-emerald-400" />
            Cryptographic SHA-256 Audit Trail
          </CardTitle>
          <div className="flex items-center gap-2">
            <Badge
              variant={isChainValid ? "outline" : "destructive"}
              className={`text-[10px] font-mono ${
                isChainValid
                  ? "border-emerald-500/40 text-emerald-400 bg-emerald-500/5"
                  : ""
              }`}
            >
              {isChainValid ? "✓ Hash Chain Valid" : "⚠ Tamper Detected"}
            </Badge>
            {tamperedIndex !== null && (
              <Button
                variant="outline"
                size="sm"
                onClick={handleResetAudit}
                className="h-6 px-2 text-[10px] font-mono gap-1"
              >
                <RotateCcw className="h-3 w-3" /> Reset
              </Button>
            )}
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-3 flex-1 overflow-y-auto text-xs p-3 font-mono">
        <p className="text-[11px] font-sans text-muted-foreground leading-relaxed">
          Every spatial voxel reservation, conflict alert, and advisory decision is
          cryptographically hashed via SHA-256 (Kafka ➔ TimescaleDB durable ledger).
        </p>

        <div className="space-y-2">
          {blocks.map((block, idx) => (
            <div
              key={block.blockIndex}
              className={`rounded-lg border p-2.5 space-y-1.5 transition-colors ${
                block.tampered
                  ? "border-rose-500/80 bg-rose-500/15"
                  : !block.isValid
                    ? "border-amber-500/50 bg-amber-500/10 opacity-75"
                    : "border-border/80 bg-background/50"
              }`}
            >
              <div className="flex items-center justify-between text-[11px]">
                <span className="font-bold text-foreground flex items-center gap-1.5">
                  <Hash className="h-3.5 w-3.5 text-muted-foreground" />
                  Block #{block.blockIndex} • {block.entityId}
                </span>
                <span className="text-[10px] text-muted-foreground">
                  {block.timestamp.slice(11, 19)} UTC
                </span>
              </div>

              <div className="text-[10px] text-sky-400 font-semibold">
                Event: {block.eventType}
              </div>

              <div className="text-[9.5px] text-muted-foreground/80 space-y-0.5 pt-1 border-t border-border/40">
                <div className="truncate">
                  Prev:{" "}
                  <span className="text-foreground/70">
                    {block.prevHash.slice(0, 24)}...
                  </span>
                </div>
                <div className="truncate">
                  Hash:{" "}
                  <span
                    className={
                      block.tampered
                        ? "text-rose-300 font-bold"
                        : "text-emerald-400/90 font-medium"
                    }
                  >
                    {block.blockHash.slice(0, 24)}...
                  </span>
                </div>
              </div>

              {tamperedIndex === null && (
                <div className="pt-1 flex justify-end">
                  <button
                    onClick={() => handleTamperBlock(idx)}
                    className="text-[9px] text-rose-400/70 hover:text-rose-300 transition-colors underline"
                  >
                    [Simulate Tamper]
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
