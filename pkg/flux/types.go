package flux

import (
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

const (
	// ScopePath 从动态Path参数中获取
	ScopePath = "PATH"
	// ScopePathMap 查询所有Path参数
	ScopePathMap = "PATH_MAP"
	// ScopeQuery 从Query参数中获取
	ScopeQuery      = "QUERY"
	ScopeQueryMulti = "QUERY_MUL"
	// ScopeQueryMap 获取全部Query参数
	ScopeQueryMap = "QUERY_MAP"
	// ScopeForm 只从Form表单参数参数列表中读取
	ScopeForm      = "FORM"
	ScopeFormMulti = "FORM_MUL"
	// ScopeFormMap 获取Form全部参数
	ScopeFormMap = "FORM_MAP"
	// ScopeParam 只从Query和Form表单参数参数列表中读取
	ScopeParam = "PARAM"
	// ScopeHeader 只从Header参数中读取
	ScopeHeader = "HEADER"
	// ScopeHeaderMap 获取Header全部参数
	ScopeHeaderMap = "HEADER_MAP"
	// ScopeAttr 获取Http Attributes的单个参数
	ScopeAttr = "ATTR"
	// ScopeAttrs 获取Http Attributes的Map结果
	ScopeAttrs = "ATTRS"
	// ScopeBody 获取Body数据
	ScopeBody = "BODY"
	// ScopeRequest 获取Request元数据
	ScopeRequest = "REQUEST"
	// ScopeAuto 自动查找数据源
	ScopeAuto = "AUTO"
)

// Support protocols
const (
	ProtoDubbo = "DUBBO"
	ProtoGRPC  = "GRPC"
	ProtoHttp  = "HTTP"
	ProtoEcho  = "ECHO"
	ProtoInApp = "INAPP"
)

const (
	SpecKindService  = "flux.go/ServiceSpec"
	SpecKindEndpoint = "flux.go/EndpointSpec"
)

// NamedValueSpec 定义KV键值对
type NamedValueSpec struct {
	Name  string      `json:"name" yaml:"name"`
	Value interface{} `json:"value" yaml:"value"`
}

func (a NamedValueSpec) GetString() string {
	rv := reflect.ValueOf(a.Value)
	if rv.Kind() == reflect.Slice {
		if rv.Len() == 0 {
			return ""
		}
		return cast.ToString(rv.Index(0).Interface())
	}
	return cast.ToString(a.Value)
}

func (a NamedValueSpec) GetStrings() []string {
	return cast.ToStringSlice(a.Value)
}

func (a NamedValueSpec) IsValid() bool {
	return a.Name != "" && a.Value != nil
}

// Annotations 注解，用于声明模型的固定有属性。注解不可被传递到后端服务
type Annotations map[string]interface{}

func (a Annotations) Exists(name string) bool {
	_, ok := a[name]
	return ok
}

func (a Annotations) Get(name string) NamedValueSpec {
	if v, ok := a[name]; ok {
		return NamedValueSpec{Name: name, Value: v}
	}
	return NamedValueSpec{}
}

func (a Annotations) GetEx(name string) (NamedValueSpec, bool) {
	if v, ok := a[name]; ok {
		return NamedValueSpec{Name: name, Value: v}, true
	}
	return NamedValueSpec{}, false
}

// EncodingType Golang内置参数类型
type EncodingType string

func (m EncodingType) Contains(s string) bool {
	return strings.Contains(string(m), s)
}

const (
	EncodingTypeGoObject      = EncodingType("go:object")
	EncodingTypeGoString      = EncodingType("go:string")
	EncodingTypeGoListString  = EncodingType("go:[]string")
	EncodingTypeGoListObject  = EncodingType("go:[]object")
	EncodingTypeGoMapString   = EncodingType("go:map[string]object")
	EncodingTypeMapStringList = EncodingType("go:map[string][]string")
)

// ValueObject 包含指示值的媒体类型和Value结构
type ValueObject struct {
	Valid    bool         // 是否有效
	Value    interface{}  // 原始值类型
	Encoding EncodingType // 数据媒体类型
}
