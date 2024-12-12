package proto

import (
	"fmt"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"regexp"
	"strings"
)

var (
	emptyMessage     = Message{"Empty", nil, nil}
	arrayReg         = regexp.MustCompile(`\[\d+]`)
	fieldTypeNameMap = func() map[string]string {
		basic := map[string]string{
			"float64":     "double",
			"float32":     "float",
			"int":         "int32",
			"int8":        "int32",
			"int16":       "int32",
			"int32":       "int32",
			"int64":       "int64",
			"uint":        "uint32",
			"uint8":       "uint32",
			"uint16":      "uint32",
			"uint32":      "uint32",
			"uint64":      "uint64",
			"bool":        "bool",
			"string":      "string",
			"any":         "bytes", // todo: use google.protobuf.Any
			"interface{}": "bytes",
		}
		result := make(map[string]string)
		for goType, protoType := range basic {
			result[goType] = protoType
			for _, prefix := range []string{
				"*", "[]", "*[]", "[]*", "*[]*",
			} {
				result[prefix+goType] = protoType
			}
		}
		return result
	}()
	byteFieldNameMap = map[string]string{
		"byte":     "uint32",
		"*byte":    "uint32",
		"[]byte":   "bytes",
		"*[]byte":  "bytes",
		"[]*byte":  "bytes",
		"*[]*byte": "bytes",
	}
	mapKeyTypeNameMap = func() map[string]string {
		basic := map[string]string{
			"int":    "int32",
			"int8":   "int32",
			"int16":  "int32",
			"int32":  "int32",
			"int64":  "int64",
			"uint":   "uint32",
			"uint8":  "uint32",
			"uint16": "uint32",
			"uint32": "uint32",
			"uint64": "uint64",
			"string": "string",
			"byte":   "uint32",
		}
		result := make(map[string]string)
		for goType, protoType := range basic {
			result[goType] = protoType
			for _, prefix := range []string{
				"*",
			} {
				result[prefix+goType] = protoType
			}
		}
		return result
	}()
)

func Unmarshal(data any, multiple bool) (f *File, err error) {
	switch val := data.(type) {
	case *spec.ApiSpec:
		f = &File{
			Syntax:  Version3,
			Package: strings.ToLower(val.Service.Name),
			Options: []*Option{
				{
					Name:  "go_package",
					Value: "/protoc-gen-go",
				},
			},
			Services: []*Service{{Name: val.Service.Name}},
		}
		messageMap := make(map[string]*Message, len(val.Types))
		for _, typ := range val.Types {
			defineStruct, _ := typ.(spec.DefineStruct)
			var message Message
			message.Name = defineStruct.Name()
			message.Descs = defineStruct.Documents()
			for _, member := range defineStruct.Members {
				var field MessageField
				if err = field.Unmarshal(&member); err != nil {
					return nil, err
				}
				message.Fields = append(message.Fields, &field)
			}
			f.Messages = append(f.Messages, &message)
			messageMap[message.Name] = &message
		}
		for _, group := range val.Service.JoinPrefix().Groups {
			var srv *Service
			if groupName := group.GetAnnotation("group"); groupName != "" && multiple {
				srv = &Service{Name: val.Service.Name + "/" + groupName}
				f.Services = append(f.Services, srv)
			} else {
				srv = f.Services[0]
			}
			for _, route := range group.Routes {
				var rpc ServiceRpc
				rpc.Name = route.Handler
				rpc.Descs = []string{strings.Trim(route.JoinedDoc(), `"`)}
				if defineStruct, ok := route.RequestType.(spec.DefineStruct); ok {
					rpc.Request = messageMap[defineStruct.Name()]
				} else {
					if _, exist := messageMap[emptyMessage.Name]; !exist {
						messageMap[emptyMessage.Name] = &emptyMessage
						f.Messages = append([]*Message{&emptyMessage}, f.Messages...)
					}
					rpc.Request = &emptyMessage
				}

				if defineStruct, ok := route.ResponseType.(spec.DefineStruct); ok {
					rpc.Response = messageMap[defineStruct.Name()]
				} else {
					if _, exist := messageMap[emptyMessage.Name]; !exist {
						messageMap[emptyMessage.Name] = &emptyMessage
						f.Messages = append([]*Message{&emptyMessage}, f.Messages...)
					}
					rpc.Response = &emptyMessage
				}
				srv.Rpcs = append(srv.Rpcs, &rpc)
			}
		}
	default:
		return nil, fmt.Errorf("unsupported type %T, only supported *spec.ApiSpec", data)
	}
	return
}

func (v *MessageField) Unmarshal(data any) error {
	switch val := data.(type) {
	case *spec.Member:
		v.Name = val.Name
		v.Descs = val.Docs
		if comment := strings.TrimSpace(strings.TrimPrefix(val.GetComment(), "//")); comment != "" {
			v.Descs = append(v.Descs, comment)
		}
		v.BaseTypeNames = parseBaseType(val.Type.Name())
		typeName := arrayReg.ReplaceAllString(val.Type.Name(), "[]")
		if name, exist := fieldTypeNameMap[typeName]; exist {
			v.TypeName = name
			v.Repeated = strings.Contains(typeName, "[]")
		} else if name, exist = byteFieldNameMap[typeName]; exist {
			v.TypeName = name
		} else if strings.Contains(typeName, "map[") {
			name, err := parseMapField(typeName)
			if err != nil {
				return fmt.Errorf("parse map field %s falied, %w", val.Name, err)
			}
			v.TypeName = name
		} else {
			v.TypeName = strings.ReplaceAll(strings.ReplaceAll(typeName, "*", ""), "[]", "")
			v.Repeated = strings.Contains(typeName, "[]")
		}
		// todo: parse member.Tag
	default:
		return fmt.Errorf("unsupported type %T, only supported *spec.Member", data)
	}
	return nil
}

func parseMapField(typeName string) (string, error) {
	before, after, found := strings.Cut(typeName, "map[")
	if !found {
		return "", fmt.Errorf("field type %s is not map", typeName)
	}
	if before != "" || !strings.Contains(after, "]") {
		return "", fmt.Errorf("unsupported field type %s", typeName)
	}
	mapKV := strings.SplitN(after, "]", 2)
	if mapKeyTypeNameMap[mapKV[0]] == "" {
		mapKV[0] = "string"
	}
	mapKV[1] = arrayReg.ReplaceAllString(mapKV[1], "[]")
	if byteFieldNameMap[mapKV[1]] != "" {
		mapKV[1] = byteFieldNameMap[mapKV[1]]
	} else if strings.Contains(mapKV[1], "[") {
		mapKV[1] = "bytes"
	} else if fieldTypeNameMap[mapKV[1]] != "" {
		mapKV[1] = fieldTypeNameMap[mapKV[1]]
	} else {
		mapKV[1] = strings.ReplaceAll(mapKV[1], "*", "")
	}

	return fmt.Sprintf("map<%s,%s>", mapKeyTypeNameMap[mapKV[0]], mapKV[1]), nil
}

func parseBaseType(typeStr string) (baseTypes []string) {
	typeStr = arrayReg.ReplaceAllString(typeStr, "[]")
	typeStr = strings.ReplaceAll(typeStr, "map[", " ")
	typeStr = strings.ReplaceAll(typeStr, "[]", " ")
	typeStr = strings.ReplaceAll(typeStr, "]", " ")
	typeStr = strings.ReplaceAll(typeStr, "*", " ")
	typMap := make(map[string]bool)
	for _, typ := range strings.Split(typeStr, " ") {
		if typ == "" || typMap[typ] {
			continue
		}
		baseTypes = append(baseTypes, typ)
		typMap[typ] = true
	}
	return
}
