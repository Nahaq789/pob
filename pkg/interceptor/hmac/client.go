package hmac

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func HmacClientInterceptor(secret string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		message, err := proto.Marshal(req.(proto.Message))
		if err != nil {
			return err
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(message))
		signature := mac.Sum(nil)
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", base64.StdEncoding.EncodeToString(signature))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
