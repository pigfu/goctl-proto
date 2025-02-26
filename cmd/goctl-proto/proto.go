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
	if apiFile := filepath.Base(goctlPlugin.ApiFilePath); goctlPlugin.Dir != "" {
		output = filepath.Join(goctlPlugin.Dir, strings.TrimSuffix(apiFile, filepath.Ext(apiFile))+".proto")
	} else {
		fi, err := os.Stat(output)
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return errors.New("output is not a directory")
		}
		output = filepath.Join(output, strings.TrimSuffix(apiFile, filepath.Ext(apiFile))+".proto")
	}
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
