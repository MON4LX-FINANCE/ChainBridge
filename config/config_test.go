// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/ChainSafe/ChainBridge/config/flags"

	"github.com/urfave/cli/v2"
)

func createTempConfigFile() (*os.File, *Config, error) {
	ethCfg := &RawChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       "1",
		Endpoint: "endpoint",
		From:     "0x0",
		Opts:     []byte(`{}`),
	}
	testConfig := &Config{
		Chains:       []*RawChainConfig{ethCfg},
		KeystorePath: "",
	}
	tmpFile, err := ioutil.TempFile(os.TempDir(), "*.json")
	if err != nil {
		return nil, nil, err
	}
	f, err := testConfig.ToJSON(tmpFile.Name())
	return f, testConfig, err
}

func TestLoadJSONConfig(t *testing.T) {
	file, cfg, err := createTempConfigFile()
	if err != nil {
		t.Fatal(err)
	}

	fset := flag.NewFlagSet("test", 1)
	fset.String(flags.ConfigFileFlag.Name, file.Name(), "")
	ctx := cli.NewContext(nil, fset, nil)

	res, err := GetConfig(ctx)
	if err != nil {
		t.Fatalf("failed to get config: %x", err)
	}

	if !reflect.DeepEqual(res, cfg) {
		t.Fatalf("did not match\ngot: %+v\nexpected: %+v", res.Chains[0], cfg.Chains[0])
	}
}

func TestValdiateConfig(t *testing.T) {
	valid := &RawChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       "1",
		Endpoint: "endpoint",
		From:     "0x0",
		Opts:     nil,
	}

	missingType := &RawChainConfig{
		Name:     "chain",
		Type:     "",
		Id:       "1",
		Endpoint: "endpoint",
		From:     "0x0",
		Opts:     nil,
	}

	missingEndpoint := &RawChainConfig{
		Name:     "chain",
		Type:     "ethereum",
		Id:       "1",
		Endpoint: "",
		From:     "0x0",
		Opts:     nil,
	}

	missingName := &RawChainConfig{
		Name:     "",
		Type:     "ethereum",
		Id:       "1",
		Endpoint: "endpoint",
		From:     "0x0",
		Opts:     nil,
	}

	cfg := Config{
		Chains:       []*RawChainConfig{valid},
		KeystorePath: "",
	}

	err := cfg.validate()
	if err != nil {
		t.Fatal(err)
	}

	cfg = Config{
		Chains:       []*RawChainConfig{missingType},
		KeystorePath: "",
	}

	err = cfg.validate()
	if err == nil {
		t.Fatal("must require type field")
	}

	cfg = Config{
		Chains:       []*RawChainConfig{missingEndpoint},
		KeystorePath: "",
	}

	err = cfg.validate()
	if err == nil {
		t.Fatal("must require endpoint field")
	}

	cfg = Config{
		Chains:       []*RawChainConfig{missingName},
		KeystorePath: "",
	}

	err = cfg.validate()
	if err == nil {
		t.Fatal("must require name field")
	}
}
