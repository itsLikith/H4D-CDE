// Copyright 2026 Likith Saragadam
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

package config

import (
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GRPCPort int `yaml:"grpc_port"`

	Separation struct {
		HorizontalNM    float64 `yaml:"horizontal_nm"`
		VerticalFt      float64 `yaml:"vertical_ft"`
		TemporalBufferS float64 `yaml:"temporal_buffer_s"`
	} `yaml:"separation"`

	Voxelization struct {
		H3Resolution int `yaml:"h3_resolution"`
		AltitudeBinFt int `yaml:"altitude_bin_ft"`
		TimeBinS      int `yaml:"time_bin_s"`
	} `yaml:"voxelization"`

	Risk struct {
		AdvisoryThreshold float64 `yaml:"advisory_threshold"`
	} `yaml:"risk"`

	Advisory struct {
		Weights struct {
			Delay               float64 `yaml:"delay"`
			PathDeviation       float64 `yaml:"path_deviation"`
			AltitudeChange      float64 `yaml:"altitude_change"`
			ConflictProbability float64 `yaml:"conflict_probability"`
		} `yaml:"weights"`
		DepartureDelayOptionsS []float64 `yaml:"departure_delay_options_s"`
		AltitudeOffsetFt       int       `yaml:"altitude_offset_ft"`
	} `yaml:"advisory"`

	AdaptiveDiscretization struct {
		Enabled          bool `yaml:"enabled"`
		FineResolution   int  `yaml:"fine_resolution"`
		FineTimeBinS     int  `yaml:"fine_time_bin_s"`
		DensityThreshold int  `yaml:"density_threshold"`
	} `yaml:"adaptive_discretization"`

	RedisAddr    string
	RedisTTL     time.Duration
	KafkaBrokers []string

	TrajectoryPredictorTarget string
	RiskScorerTarget          string
}

// Load reads config from separation_params.yaml and environment variables.
func Load() (*Config, error) {
	cfg := &Config{}

	// Sensible defaults
	cfg.GRPCPort = 50051
	cfg.Separation.HorizontalNM = 5.0
	cfg.Separation.VerticalFt = 1000.0
	cfg.Separation.TemporalBufferS = 60.0
	cfg.Voxelization.H3Resolution = 8
	cfg.Voxelization.AltitudeBinFt = 100
	cfg.Voxelization.TimeBinS = 10
	cfg.Risk.AdvisoryThreshold = 0.50
	cfg.Advisory.Weights.Delay = 0.30
	cfg.Advisory.Weights.PathDeviation = 0.30
	cfg.Advisory.Weights.AltitudeChange = 0.20
	cfg.Advisory.Weights.ConflictProbability = 0.20
	cfg.Advisory.DepartureDelayOptionsS = []float64{30, 60, 90, 120}
	cfg.Advisory.AltitudeOffsetFt = 500
	cfg.AdaptiveDiscretization.Enabled = false
	cfg.AdaptiveDiscretization.FineResolution = 9
	cfg.AdaptiveDiscretization.FineTimeBinS = 5
	cfg.AdaptiveDiscretization.DensityThreshold = 6

	yamlPath := os.Getenv("CONFIG_SEPARATION_PATH")
	if yamlPath == "" {
		yamlPath = "config/separation_params.yaml"
	}

	if data, err := os.ReadFile(yamlPath); err == nil {
		_ = yaml.Unmarshal(data, cfg)
	}

	// Environment variable overrides (zero hard-coding)
	if portStr := os.Getenv("VOXEL_ENGINE_GRPC_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			cfg.GRPCPort = p
		}
	} else if pStr := os.Getenv("PORT"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil {
			cfg.GRPCPort = p
		}
	}

	cfg.RedisAddr = os.Getenv("REDIS_ADDR")
	if cfg.RedisAddr == "" {
		cfg.RedisAddr = "localhost:6379"
	}
	cfg.RedisTTL = 120 * time.Second

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers != "" {
		cfg.KafkaBrokers = []string{brokers}
	} else {
		cfg.KafkaBrokers = []string{"localhost:9092"}
	}

	tpHost := os.Getenv("TRAJECTORY_PREDICTOR_HOST")
	if tpHost == "" {
		tpHost = "localhost"
	}
	tpPort := os.Getenv("TRAJECTORY_PREDICTOR_GRPC_PORT")
	if tpPort == "" {
		tpPort = "50052"
	}
	cfg.TrajectoryPredictorTarget = tpHost + ":" + tpPort

	rsHost := os.Getenv("RISK_SCORER_HOST")
	if rsHost == "" {
		rsHost = "localhost"
	}
	rsPort := os.Getenv("RISK_SCORER_GRPC_PORT")
	if rsPort == "" {
		rsPort = "50055"
	}
	cfg.RiskScorerTarget = rsHost + ":" + rsPort

	return cfg, nil
}
