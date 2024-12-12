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
	messageUsage := make(map[*Message]int32, len(f.Messages))
	if len(includeRpcs) > 0 {
		includes := make(map[string]bool, len(includeRpcs))
		for _, rpc := range includeRpcs {
			includes[rpc] = true
		}
		newRpcs := make([]*ServiceRpc, 0, len(includeRpcs))
		for _, rpc := range f.Service.Rpcs {
			if !includes[rpc.Name] {
				messageUsage[rpc.Request] = max(messageUsage[rpc.Request], 0)
				messageUsage[rpc.Response] = max(messageUsage[rpc.Response], 0)
				continue
			}
			messageUsage[rpc.Request]++
			messageUsage[rpc.Response]++
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
			messageUsage[rpc.Request] = 0
			messageUsage[rpc.Response] = 0
		}
		if len(preRpcs) > 0 {
			f.Service.Rpcs = preRpcs
		}
		for _, rpc := range f.Service.Rpcs {
			messageUsage[rpc.Request]++
			messageUsage[rpc.Response]++
		}
		for i := 0; i < len(f.Service.Rpcs); i++ {
			rpc := f.Service.Rpcs[i]
			if !excludes[rpc.Name] {
				continue
			}
			messageUsage[rpc.Request]--
			messageUsage[rpc.Response]--
			f.Service.Rpcs = append(f.Service.Rpcs[:i:cap(f.Service.Rpcs)-1], f.Service.Rpcs[i+1:]...)
			i--
		}
	}

	for i := 0; i < len(f.Messages); i++ {
		if cnt, exist := messageUsage[f.Messages[i]]; !exist || cnt > 0 {
			continue
		}
		f.Messages = append(f.Messages[:i:cap(f.Messages)-1], f.Messages[i+1:]...)
		i--
	}
	return f
}
