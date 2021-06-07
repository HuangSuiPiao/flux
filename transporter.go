package flux

type (
	// Transporter 表示某种特定协议的后端服务，例如Dubbo, gRPC, Http等协议的后端服务。
	// 默认实现了Dubbo(gRpc)和Http两种协议。
	Transporter interface {
		// DoInvoke 执行指定目标EndpointService的通讯，返回响应结果，并解析响应数据
		DoInvoke(*Context, Service) (*ServeResponse, *ServeError)
	}
	// TransportCodecFunc 解析 Transporter 返回的原始数据，生成响应对象
	TransportCodecFunc func(ctx *Context, result interface{}, att map[string]interface{}) (*ServeResponse, error)
)
