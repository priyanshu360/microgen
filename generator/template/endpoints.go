package template

import (
	. "github.com/dave/jennifer/jen"
	"github.com/devimteam/microgen/parser"
	"github.com/devimteam/microgen/util"
)

const PackageAliasGoKit = "github.com/go-kit/kit/endpoint"

type EndpointsTemplate struct {
}

// Renders endpoints file.
//
//		// This file was automatically generated by "microgen" utility.
//		// Please, do not edit.
//		package stringsvc
//
//		import (
//		context "context"
//		endpoint "github.com/go-kit/kit/endpoint"
//		)
//
//		type Endpoints struct {
//			CountEndpoint endpoint.Endpoint
//		}
//
//		func (e *Endpoints) Count(ctx context.Context, text string, symbol string) (count int) {
//			req := CountRequest{
//				Symbol: symbol,
//				Text:   text,
//			}
//			resp, err := e.CountEndpoint(ctx, &req)
//			if err != nil {
//				return
//			}
//			return resp.(*CountResponse).Count
//		}
//
//		func CountEndpoint(svc StringService) endpoint.Endpoint {
//			return func(ctx context.Context, request interface{}) (interface{}, error) {
//				req := request.(*CountRequest)
//				count := svc.Count(ctx, req.Text, req.Symbol)
//				return &CountResponse{Count: count}, nil
//			}
//		}
//
func (EndpointsTemplate) Render(i *parser.Interface) *File {
	f := NewFile(i.PackageName)

	f.Type().Id("Endpoints").StructFunc(func(g *Group) {
		for _, signature := range i.FuncSignatures {
			g.Id(signature.Name+"Endpoint").Qual(PackageAliasGoKit, "Endpoint")
		}
	})

	for _, signature := range i.FuncSignatures {
		f.Add(endpointFunc(signature))
		f.Line()
	}
	f.Line()
	for _, signature := range i.FuncSignatures {
		f.Add(newEndpointFunc(signature, i))
		f.Line()
	}

	return f
}

func (EndpointsTemplate) Path() string {
	return "./endpoints.go"
}

// Render full endpoints method
//
//		func (e *Endpoints) Count(ctx context.Context, text string, symbol string) (count int) {
//			req := CountRequest{
//				Symbol: symbol,
//				Text:   text,
//			}
//			resp, err := e.CountEndpoint(ctx, &req)
//			if err != nil {
//				return
//			}
//			return resp.(*CountResponse).Count
//		}
//
func endpointFunc(signature *parser.FuncSignature) *Statement {
	return methodDefinition("Endpoints", signature).
		BlockFunc(endpointBody(signature))
}

// Render interface method body
//
//		req := CountRequest{
//			Symbol: symbol,
//			Text:   text,
//		}
//		resp, err := e.CountEndpoint(ctx, &req)
//		if err != nil {
//			return
//		}
//		return resp.(*CountResponse).Count
//
func endpointBody(signature *parser.FuncSignature) func(g *Group) {
	req := "req"
	resp := "resp"
	return func(g *Group) {
		g.Id(req).Op(":=").Id(signature.Name + "Request").Values(mapInitByFuncFields(signature.Params))
		g.List(Id(resp), Err()).Op(":=").Id(util.FirstLowerChar("Endpoint")).Dot(signature.Name+"Endpoint").Call(Id(Context), Op("&").Id(req))
		g.If(Err().Op("!=").Nil()).Block(
			Return(),
		)
		g.ReturnFunc(func(group *Group) {
			for _, field := range signature.Results {
				group.Add(typeCasting(resp, signature.Name+"Response")).Op(".").Add(structFieldName(field))
			}
		})
	}
}

// Render new Endpoint body
//
//		return func(ctx context.Context, request interface{}) (interface{}, error) {
//			req := request.(*CountRequest)
//			count := svc.Count(ctx, req.Text, req.Symbol)
//			return &CountResponse{Count: count}, nil
//		}
//
func newEndpointFuncBody(signature *parser.FuncSignature) *Statement {
	return Return(Func().Params(
		Id("ctx").Qual("context", "Context"),
		Id("request").Interface(),
	).Params(
		Interface(),
		Error(),
	).BlockFunc(func(g *Group) {
		g.Id("req").Op(":=").Add(typeCasting("request", signature.Name+"Request"))
		g.Add(fullServiceMethodCall("svc", "req", signature))
		g.Return(
			Op("&").Id(signature.Name+"Response").Values(mapInitByFuncFields(signature.Results)),
			Nil(),
		)
	}))
}

// Render full new Endpoint function
//
//		func CountEndpoint(svc StringService) endpoint.Endpoint {
//			return func(ctx context.Context, request interface{}) (interface{}, error) {
//				req := request.(*CountRequest)
//				count := svc.Count(ctx, req.Text, req.Symbol)
//				return &CountResponse{Count: count}, nil
//			}
//		}
//
func newEndpointFunc(signature *parser.FuncSignature, svcInterface *parser.Interface) *Statement {
	return Func().
		Id(signature.Name + "Endpoint").Params(Id("svc").Id(svcInterface.Name)).Params(Qual(PackageAliasGoKit, "Endpoint")).
		Block(newEndpointFuncBody(signature))
}
