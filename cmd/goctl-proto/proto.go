package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/pigfu/goctl-proto/internal/proto"
	"github.com/urfave/cli/v3"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"os"
	"path/filepath"
	"strings"
)

func checkAndGenDir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(dir, 0777)
}

func protoGen(_ context.Context, command *cli.Command) (err error) {
	output := command.String("output")
	defer func() {
		fmt.Print("Generate proto file")
		if output != "" {
			fmt.Printf(" %s", output)
		}
		if err != nil {
			fmt.Printf(" [FAILED]\n")
		} else {
			fmt.Printf(" [OK]\n")
		}
	}()
	var goctlPlugin plugin.Plugin
	if goctlPlugin.ApiFilePath = command.String("input"); goctlPlugin.ApiFilePath != "" {
		if goctlPlugin.Api, err = parser.Parse(goctlPlugin.ApiFilePath, ""); err != nil {
			return err
		}
	} else if plug, err := plugin.NewPlugin(); err == nil {
		goctlPlugin = *plug
	} else {
		return errors.New("api file not found, must set one of goctl -api or --input")
	}
	apiFile := filepath.Base(goctlPlugin.ApiFilePath)
	if goctlPlugin.Dir != "" {
		output = goctlPlugin.Dir
	}
	if err = checkAndGenDir(output); err != nil {
		return err
	}
	output = filepath.Join(output, strings.TrimSuffix(apiFile, filepath.Ext(apiFile))+".proto")
	pf, err := proto.Unmarshal(goctlPlugin, command.Bool("multiple"))
	if err != nil {
		return err
	}
	pd, err := pf.Refine(command.StringSlice("include-handler"), command.StringSlice("exclude-handler")).Marshal()
	if err != nil {
		return err
	}
	if err = os.WriteFile(output, pd, 0666); err != nil {
		return err
	}
	return
}
