/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fix

import (
	"context"
	"fmt"
	"sync/atomic"

	"mosn.io/api"

	"mosn.io/pkg/log"
)

type FixProtocol struct{}

func (proto FixProtocol) Name() api.ProtocolName {
	return ProtocolName
}

func (proto FixProtocol) Encode(ctx context.Context, model interface{}) (api.IoBuffer, error) {
	switch frame := model.(type) {
	case *Request:
		return encodeRequest(ctx, frame)
	case *Response:
		return encodeResponse(ctx, frame)
	default:
		log.DefaultLogger.Errorf("[protocol][fix] encode with unknown command : %+v", model)
		return nil, api.ErrUnknownType
	}
}

func (proto FixProtocol) Decode(ctx context.Context, data api.IoBuffer) (interface{}, error) {
	if data.Len() >= LessLen {
		flag := getStreamType(data.Bytes())

		switch flag {
		case CmdRequest, CmdRequestHeartbeat:
			return decodeRequest(ctx, data)
		case CmdResponse, CmdResponseHeartbeat:
			return decodeResponse(ctx, data)
		default:
			// unknown cmd type
			return nil, fmt.Errorf("decode Error, type = %s, value = %d", UnKnownCmdType, flag)
		}
	}

	return nil, nil
}

// Heartbeater
func (proto FixProtocol) Trigger(ctx context.Context, requestId uint64) api.XFrame {
	// TODO: 构造心跳请求，如果不支持心跳能力，删除TODO和以下注释
	//return &Request{
	//	ProtocolHeader: ProtocolHeader{
	//		Magic:     0xbcbc,
	//		Flag:      CmdRequestHeartbeat,
	//		RequestId: uint32(requestId),
	//		Codec:     0, // not used now.
	//	},
	//	Timeout: 3000,
	//}

	return nil
}

func (proto FixProtocol) Reply(ctx context.Context, request api.XFrame) api.XRespFrame {
	// TODO: 构造心跳响应，如果不支持心跳能力，删除TODO和以下注释
	//return &Response{
	//	ProtocolHeader: ProtocolHeader{
	//		Flag:      CmdResponseHeartbeat,
	//		RequestId: uint32(request.GetRequestId()),
	//		Codec:     0, // not used now.
	//	},
	//	Status: ResponseStatusSuccess,
	//}
	return nil
}

// Hijacker
func (proto FixProtocol) Hijack(ctx context.Context, request api.XFrame, statusCode uint32) api.XRespFrame {
	panic("实现: 根据原始请求和私有协议异常码，构造请求被mosn劫持的响应报文")
	// TODO: example 构造请求被mosn劫持的响应报文
	//return &Response{
	//	ProtocolHeader: ProtocolHeader{
	//		Flag:      CmdResponse,
	//		RequestId: uint32(request.GetRequestId()),
	//		Codec:     0, // not used now.
	//	},
	//	Status: statusCode,
	//}
}

func (proto FixProtocol) Mapping(httpStatusCode uint32) uint32 {
	panic("实现: 根据http状态码转换成私有协议状态码")
	// TODO: example http状态码转换成私有协议状态码
	// switch httpStatusCode {
	// case http.StatusOK:
	// 	return ResponseStatusSuccess
	// case api.RouterUnavailableCode:
	// 	return ResponseStatusNoProcessor
	// case api.NoHealthUpstreamCode:
	// 	return ResponseStatusConnectionClosed
	// case api.UpstreamOverFlowCode:
	// 	return ResponseStatusServerThreadPoolBusy
	// case api.CodecExceptionCode:
	// 	//Decode or Encode Error
	// 	return ResponseStatusCodecException
	// case api.DeserialExceptionCode:
	// 	//Hessian Exception
	// 	return ResponseStatusServerDeserializeException
	// case api.TimeoutExceptionCode:
	// 	//Response Timeout
	// 	return ResponseStatusTimeout
	// default:
	// 	return ResponseStatusUnknown
	// }
}

// PoolMode returns whether ping-pong or multiplex
func (proto FixProtocol) PoolMode() api.PoolMode {
	return api.PingPong
}

func (proto FixProtocol) EnableWorkerPool() bool {
	return false
}

func (proto FixProtocol) GenerateRequestID(streamID *uint64) uint64 {
	return 0
}
