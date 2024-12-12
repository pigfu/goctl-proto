package proto

import (
	"bytes"
	"strings"
	"text/template"
)

func (f *File) Marshal() ([]byte, error) {
	buf := bytes.NewBufferString("")
	if err := template.Must(template.New("").Funcs(funcMap).Parse(fileTemplate)).
		Execute(buf, f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Refine todo: remove unused message in message fields
func (f *File) Refine(includeRpcs, excludeRpcs []string) *File {
	if f == nil || f.Service == nil || len(f.Service.Rpcs) == 0 {
		return f
	}
	if len(includeRpcs) > 0 {
		includes := make(map[string]bool, len(includeRpcs))
		for _, rpc := range includeRpcs {
			includes[rpc] = true
		}
		newRpcs := make([]*ServiceRpc, 0, len(includeRpcs))
		for _, rpc := range f.Service.Rpcs {
			if !includes[rpc.Name] {
				continue
			}
			newRpcs = append(newRpcs, rpc)
		}
		f.Service.Rpcs = newRpcs
	} else {
		excludes := make(map[string]bool, len(excludeRpcs))
		for _, rpc := range excludeRpcs {
			excludes[rpc] = true
		}
		// find rpc which desc contains @goctl-proto
		preRpcs := make([]*ServiceRpc, 0, len(f.Service.Rpcs))
		for _, rpc := range f.Service.Rpcs {
			if strings.Contains(strings.Join(rpc.Descs, " "), "@goctl-proto") {
				for i := range rpc.Descs {
					rpc.Descs[i] = strings.TrimSpace(strings.ReplaceAll(rpc.Descs[i], "@goctl-proto", ""))
				}
				preRpcs = append(preRpcs, rpc)
			}
		}
		if len(preRpcs) > 0 {
			f.Service.Rpcs = preRpcs
		}
		for i := 0; i < len(f.Service.Rpcs); i++ {
			if !excludes[f.Service.Rpcs[i].Name] {
				continue
			}
			f.Service.Rpcs = append(f.Service.Rpcs[:i:cap(f.Service.Rpcs)-1], f.Service.Rpcs[i+1:]...)
			i--
		}
	}

	messageUsage := make(map[string]bool)
	for _, rpc := range f.Service.Rpcs {
		messageUsage[rpc.Request.Name], messageUsage[rpc.Response.Name] = true, true
		for _, field := range append(rpc.Request.Fields, rpc.Response.Fields...) {
			for _, name := range field.BaseTypeNames {
				messageUsage[name] = true
			}
		}
	}

	for i := 0; i < len(f.Messages); i++ {
		if messageUsage[f.Messages[i].Name] {
			continue
		}
		f.Messages = append(f.Messages[:i:cap(f.Messages)-1], f.Messages[i+1:]...)
		i--
	}
	return f
}
