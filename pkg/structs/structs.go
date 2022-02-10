package structs

type YAMLFields struct {
	Name string `yaml:"name"`
	Value string `yaml:"value"`
}

type YAMLBruteforce struct {
	From string `yaml:"from"`
	File string `yaml:"file"`
	List []string `yaml:"list"`
	Field string `yaml:"field"`
	Threads int `yaml:"threads"`
	NoVerbose bool `yaml:"no_verbose"`
	Output	string	`yaml:"output"`
}

type YAMLSite struct {
	Host string `yaml:"host"`
	Method string `yaml:"method"`
}

type YAMLCrawl struct {
	Name   string `yaml:"name"`
	Search string `yaml:"search"`
}

type YAMLConfig struct {
	F []YAMLFields `yaml:"fields"`
	B YAMLBruteforce `yaml:"bruteforce"`
	OF YAMLOn_fail `yaml:"on_fail"`
	OP YAMLOn_pass `yaml:"on_pass"`
	H []YAMLHeaders `yaml:"headers"`
	S YAMLSite `yaml:"site"`
	C YAMLCrawl `yaml:"crawl"`
}

type YAMLOn_fail struct {
	Message string `yaml:"message"`
	Length int `yaml:"content_length"`
}

type YAMLOn_pass struct {
	Message string `yaml:"message"`
	Length int `yaml:"content_length"`
}

type YAMLHeaders struct {
	Name string `yaml:"name"`
	Value string `yaml:"value"`
}