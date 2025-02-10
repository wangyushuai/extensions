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
	"strconv"
	"strings"

	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
)

func decodeRequest(ctx context.Context, data api.IoBuffer) (cmd interface{}, err error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. 判断最小协议头部长度，是否可以计算出完整报文长度
	if bytesLen < RequestHeaderLen {
		return
	}

	var frameLen int
	var headerLen uint16
	var contentLen uint16
	var isRequestAck bool

	if bytesLen == 4 && strings.EqualFold(string(bytes[0:4]), RequestACKFlag) {
		frameLen = 4
		isRequestAck = true
	} else {
		//TODO:  判断报文是否完整, 流式首次请求可能会超过 128位，易导致请求没有请求完整就被发出
		//frameLen = RequestHeaderLen + int(headerLen) + int(contentLen)
		frameLen = bytesLen
		if bytesLen < frameLen {
			log.DefaultLogger.Infof("[fix] continue read req:%d", bytesLen)
			return
		}
	}
	log.DefaultLogger.Infof("[fix] recevie req:%s", string(bytes))
	data.Drain(frameLen)

	request := &Request{}
	request.ProtocolHeader = ProtocolHeader{
		Flag:       getStreamType(bytes),
		HeaderLen:  headerLen,
		ContentLen: contentLen,
	}
	request.Data = buffer.GetIoBuffer(frameLen)
	//3. 完整报文复制到Data字段中
	request.Data.Write(bytes[:frameLen])
	request.rawData = request.Data.Bytes()

	//oriRemoteAddr := mosnctx.Get(ctx, types.ContextOriRemoteAddr)
	// Mock ServiceName
	request.Set(ServiceNameKey, "cmbc_stream_server")
	//ctx = context.WithValue(ctx, ServiceNameKey, "cmbc_stream_server")
	//variable.SetString(ctx, types.VarDirection, "100.88.140.206:9090")

	// Oneway逻辑处理
	if isRequestAck {
		if isLastAck(ctx) {
			request.Flag = CmdRequestOneway
			log.DefaultLogger.Infof("[fix] lask ack, request type:%d", request.Flag)
		}
		return request, err
	}
	//TODO: 如果有需要，这里可以进一步解析 Header和Body
	return request, err
}

func decodeResponse(ctx context.Context, data api.IoBuffer) (cmd interface{}, err error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. 判断最小协议头部长度，是否可以计算出完整报文长度
	if bytesLen < ResponseHeaderLen {
		return
	}

	var (
		headerLen  uint16
		contentLen int // body长度
		frameLen   int
	)

	if bytesLen == 8 {
		// 首次报文
		frameLen = 8
		if _, err = strconv.Atoi(strings.TrimSpace(string(bytes[0:8]))); err != nil {
			err = fmt.Errorf("pares total length failed,err:%v", err)
			return
		}
	} else {
		// Chunk报文
		headerLen = 0
		contentLen, err = strconv.Atoi(string(bytes[0:8]))
		frameLen = ResponseHeaderLen + int(headerLen) + contentLen
		if err != nil {
			err = fmt.Errorf("pares total length failed,err:%v", err)
			return
		}
	}

	if bytesLen < frameLen {
		log.DefaultLogger.Infof("[fix] continue read resp: %d", bytesLen)
		return
	}
	log.DefaultLogger.Infof("[fix] receive resp:%s", string(bytes))
	// 非常重要: 丢弃tcp连接解码后的数据，防止内核重复推送重复数据
	data.Drain(frameLen)

	response := &Response{}
	response.ProtocolHeader = ProtocolHeader{
		Flag:       getStreamType(bytes),
		HeaderLen:  headerLen,
		ContentLen: uint16(contentLen),
	}

	response.Data = buffer.GetIoBuffer(frameLen)
	response.Data.Write(bytes[:frameLen])
	response.rawData = response.Data.Bytes()

	// Oneway逻辑处理
	if frameLen == 8 {
		//首次报文
		var totalLen int
		if totalLen, err = strconv.Atoi(strings.TrimSpace(string(bytes[0:8]))); err == nil {
			setContentLength(ctx, uint32(totalLen))
		}
	} else {
		//Chunk报文
		setChunkedSize(ctx, uint32(contentLen))
	}
	if isChunkComplete(ctx) {
		response.Set(ExtChunkedFinish, "true")
	}

	return response, err
}
