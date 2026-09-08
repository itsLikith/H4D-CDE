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
import { FaqSection } from "@/components/FaqSection";
import { FeatureGrid } from "@/components/FeatureGrid";
import { Footer } from "@/components/Footer";
import { HeroSection } from "@/components/HeroSection";
import { HowItWorksSection } from "@/components/HowItWorksSection";
import { Navbar } from "@/components/Navbar";
import { SolutionsSection } from "@/components/SolutionsSection";

export default function LandingPage() {
  return (
    <main className="min-h-screen flex flex-col bg-[#080C14] text-foreground selection:bg-blue-500/30">
      {/* Product Navigation */}
      <Navbar />

      {/* Hero Section */}
      <HeroSection />

      {/* 4-Step Technical Workflow */}
      <HowItWorksSection />

      {/* Core Platform Features */}
      <FeatureGrid />

      {/* Practical Use Cases */}
      <SolutionsSection />

      {/* Product FAQ */}
      <FaqSection />

      {/* Product Footer */}
      <Footer />
    </main>
  );
}
