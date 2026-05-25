package interceptor

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func ServerInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// リクエスト内容からMAC計算
		message, err := proto.Marshal(req.(proto.Message))
		if err != nil {
			return nil, err
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(message))
		inputSignature := mac.Sum(nil)

		// メタデータからMAC取り出し
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization")
		}
		validateSignature, err := base64.StdEncoding.DecodeString(values[0])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization")
		}

		// 検証
		if !hmac.Equal(inputSignature, validateSignature) {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		res, err := handler(ctx, req)
		return res, err
	}
}
