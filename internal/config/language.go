package config

type Limits struct {
	WallTimeS    int `yaml:"wall_time_s" json:"wall_time_s"`
	MemoryKB     int `yaml:"memory_kb" json:"memory_kb"`
	MaxProcesses int `yaml:"max_processes" json:"max_processes"`
}

type BuildConfig struct {
	Cmd           string   `yaml:"cmd"`
	Args          []string `yaml:"args"`
	Limits        Limits   `yaml:"limits"`
	FlagAllowlist []string `yaml:"flag_allowlist"`
}

type RunConfig struct {
	Cmd           string   `yaml:"cmd"`
	Args          []string `yaml:"args"`
	Limits        Limits   `yaml:"limits"`
	FlagAllowlist []string `yaml:"flag_allowlist"`
}

type Language struct {
	ID                     string       `yaml:"id"`
	Name                   string       `yaml:"name"`
	SourceFilename         string       `yaml:"source_filename"`
	SourceFilenameStrategy string       `yaml:"source_filename_strategy"`
	Artifact               string       `yaml:"artifact"`
	ArtifactFilenameStrategy string       `yaml:"artifact_filename_strategy"`
	VersionProbe           string       `yaml:"version_probe"`
	Build                  *BuildConfig `yaml:"build"`
	Run                    RunConfig    `yaml:"run"`
}
