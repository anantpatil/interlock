package config

import (
	"testing"

	typesapi "github.com/ehazlett/interlock/api/types"
)

var (
	sampleConfig = `
ListenAddr = ":8080"
GRPCAddr = ":8081"
DockerURL = "unix:///var/run/docker.sock"
`
)

func TestParseConfig(t *testing.T) {
	cfg, err := ParseConfig(sampleConfig)
	if err != nil {
		t.Fatalf("error parsing config: %s", err)
	}

	if cfg.ListenAddr != ":8080" {
		t.Fatalf("expected listen addr :8080; received %s", cfg.ListenAddr)
	}

	if cfg.DockerURL != "unix:///var/run/docker.sock" {
		t.Fatalf("expected docker url unix:///var/run/docker.sock; received %s", cfg.DockerURL)
	}
}

func TestSetConfigDefaults(t *testing.T) {
	cfg := &typesapi.PluginConfig{
		Version: "1",
	}

	SetConfigDefaults(cfg)

	if cfg.User != "www-data" {
		t.Fatalf("expected default user of www-data; received %q", cfg.User)
	}

	if cfg.WorkerProcesses != 1 {
		t.Fatalf("expected default worker processes of 1; received %d", cfg.WorkerProcesses)
	}

	if cfg.RlimitNoFile != 65535 {
		t.Fatalf("expected default rlimit no file of 65535; received %d", cfg.RlimitNoFile)
	}

	if cfg.ConnectTimeout != 600 {
		t.Fatalf("expected default proxy connect timeout of 600; received %d", cfg.ConnectTimeout)
	}

	if cfg.SendTimeout != 600 {
		t.Fatalf("expected default proxy send timeout of 600; received %d", cfg.SendTimeout)
	}

	if cfg.ReadTimeout != 600 {
		t.Fatalf("expected default proxy read timeout of 600; received %d", cfg.ReadTimeout)
	}

	if cfg.SendTimeout != 600 {
		t.Fatalf("expected default send timeout of 600; received %d", cfg.SendTimeout)
	}

	if cfg.SslCiphers != "HIGH:!aNULL:!MD5" {
		t.Fatalf("expected default SSL ciphers of HIGH:!aNULL:!MD5; received %d", cfg.SslCiphers)
	}

	if cfg.SslProtocols != "SSLv3 TLSv1 TLSv1.1 TLSv1.2" {
		t.Fatalf("expected default SSL protocols of SSLv3 TLSv1 TLSv1.1 TLSv1.2; received %d", cfg.SslProtocols)
	}
}
