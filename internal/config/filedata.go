package config

type FileData struct {
	InputFile  string `env:"GRPC_SERVER_INPUT_FILE" env-required:"true"`
	OutputFile string `env:"GRPC_SERVER_OUTPUT_FILE" env-required:"true"`
}
