package template

import (
	"fmt"

	"github.com/devimteam/microgen/generator/write_strategy"
	"github.com/vetcher/godecl/types"
	. "github.com/vetcher/jennifer/jen"
)

type gRPCClientTemplate struct {
	Info *GenerationInfo
}

func NewGRPCClientTemplate(info *GenerationInfo) Template {
	return &gRPCClientTemplate{
		Info: info.Duplicate(),
	}
}

func (t *gRPCClientTemplate) grpcConverterPackagePath() string {
	return t.Info.ServiceImportPath + "/transport/converter/protobuf"
}

// Render whole grpc client file.
//
//		// This file was automatically generated by "microgen" utility.
//		// Please, do not edit.
//		package transportgrpc
//
//		import (
//			svc "github.com/devimteam/microgen/test/svc"
//			protobuf "github.com/devimteam/microgen/test/svc/transport/converter/protobuf"
//			grpc1 "github.com/go-kit/kit/transport/grpc"
//			stringsvc "gitlab.devim.team/protobuf/stringsvc"
//			grpc "google.golang.org/grpc"
//		)
//
//		func NewGRPCClient(conn *grpc.ClientConn, opts ...grpc1.ClientOption) svc.StringService {
//			return &svc.Endpoints{CountEndpoint: grpc1.NewClient(
//				conn,
//				"devim.string.protobuf.StringService",
//				"Count",
//				protobuf.EncodeCountRequest,
//				protobuf.DecodeCountResponse,
//				stringsvc.CountResponse{},
//				opts...,
//			).Endpoint()}
//		}
//
func (t *gRPCClientTemplate) Render() write_strategy.Renderer {
	f := NewFile(t.Info.ServiceImportPackageName)
	f.PackageComment(FileHeader)
	f.PackageComment(`Please, do not edit.`)

	f.Func().Id("NewGRPCClient").
		Params(
			Id("conn").Op("*").Qual(PackagePathGoogleGRPC, "ClientConn"),
			Id("opts").Op("...").Qual(PackagePathGoKitTransportGRPC, "ClientOption"),
		).Qual(t.Info.ServiceImportPath, t.Info.Iface.Name).
		BlockFunc(func(g *Group) {
			g.Return().Op("&").Qual(t.Info.ServiceImportPath, "Endpoints").Values(DictFunc(func(d Dict) {
				for _, m := range t.Info.Iface.Methods {
					d[Id(endpointStructName(m.Name))] = Qual(PackagePathGoKitTransportGRPC, "NewClient").Call(
						Line().Id("conn"),
						Line().Lit(t.Info.GRPCRegAddr),
						Line().Lit(m.Name),
						Line().Qual(pathToConverter(t.Info.ServiceImportPath), requestEncodeName(m)),
						Line().Qual(pathToConverter(t.Info.ServiceImportPath), responseDecodeName(m)),
						Line().Add(t.replyType(m)),
						Line().Id("opts").Op("...").Line(),
					).Dot("Endpoint").Call()
				}
			}))
		})
	return f
}

// Renders reply type argument
// 		stringsvc.CountResponse{}
func (t *gRPCClientTemplate) replyType(signature *types.Function) *Statement {
	return Qual(t.Info.ProtobufPackage, responseStructName(signature)).Values()
}

func (gRPCClientTemplate) DefaultPath() string {
	return "./transport/grpc/client.go"
}

func (t *gRPCClientTemplate) Prepare() error {
	if t.Info.GRPCRegAddr == "" {
		return fmt.Errorf("grpc server address is empty")
	}
	if t.Info.ProtobufPackage == "" {
		return fmt.Errorf("protobuf package is empty")
	}
	return nil
}

func (t *gRPCClientTemplate) ChooseStrategy() (write_strategy.Strategy, error) {
	return write_strategy.NewFileMethod(t.Info.AbsOutPath, t.DefaultPath()), nil
}
