package config

type Limits struct {
	WallTimeS    int `yaml:"wall_time_s"`
	MemoryKB     int `yaml:"memory_kb"`
	MaxProcesses int `yaml:"max_processes"`
}

type BuildConfig struct {
	Cmd           string   `yaml:"cmd"`
	Args          []string `yaml:"args"`
	Limits        Limits   `yaml:"limits"`
	FlagAllowlist []string `yaml:"flag_allowlist"`
}

type RunConfig struct {
	Cmd    string   `yaml:"cmd"`
	Args   []string `yaml:"args"`
	Limits Limits   `yaml:"limits"`
}

type Language struct {
	ID             string       `yaml:"id"`
	Name           string       `yaml:"name"`
	SourceFilename string       `yaml:"source_filename"`
	Artifact       string       `yaml:"artifact"`
	Build          *BuildConfig `yaml:"build"`
	Run            RunConfig    `yaml:"run"`
}
