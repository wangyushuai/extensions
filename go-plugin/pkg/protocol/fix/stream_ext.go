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
	"os"
	"strconv"
	"sync"

	"mosn.io/extensions/go-plugin/pkg/config"
	"mosn.io/pkg/log"
)

var enableFixOneway bool = true

const enableFixOnewayFlag = "ENABLE_FIX_ONEWAY"

func init() {
	if v := os.Getenv(enableFixOnewayFlag); len(v) > 0 {
		enableFixOneway, _ = strconv.ParseBool(v)
		log.DefaultLogger.Infof("[fix]enableFixOneway=%t", enableFixOneway)
	}
}

const (
	ExtContentLength string = "x-mosn-content-length"
	ExtChunkedLength string = "x-mosn-chunked-length"
	ExtChunkedFinish string = "x-mosn-chunked-finish"
)

func setContentLength(ctx context.Context, length uint32) {
	if connCtx, ok := config.GetSharedConnCtx(ctx); ok {
		connCtx.LoadOrStore(ExtContentLength, length)
		log.DefaultLogger.Infof("[fix] ctx=%p, set %s=%d", &ctx, ExtContentLength, length)
		return
	}
	log.DefaultLogger.Warnf("[fix] failed to set cotent length, shared conn context not found.ctx=%p", &ctx)
}

func setChunkedSize(ctx context.Context, length uint32) {
	var connCtx *sync.Map
	var ok bool
	var currentLength uint32
	if connCtx, ok = config.GetSharedConnCtx(ctx); !ok {
		log.DefaultLogger.Warnf("[fix] failed to set chunk length, shared conn context not found.ctx=%p", &ctx)
		return
	}
	if val, ok := connCtx.Load(ExtChunkedLength); ok {
		currentLength = val.(uint32)
	}
	newLength := currentLength + length
	connCtx.Store(ExtChunkedLength, newLength)
	log.DefaultLogger.Infof("[fix] ctx=%p, set %s=%d", &ctx, ExtChunkedLength, newLength)
}

func setChunkFinished(ctx context.Context, finished bool) {
	if connCtx, ok := config.GetSharedConnCtx(ctx); ok {
		connCtx.Store(ExtChunkedFinish, finished)
		log.DefaultLogger.Infof("[fix] ctx=%p, set %s=%t", &ctx, ExtChunkedFinish, finished)
		return
	}
	log.DefaultLogger.Warnf("[fix] failed to set chunk finished, shared conn context not found.ctx=%p", &ctx)
}

func getContentLength(ctx context.Context) (uint32, bool) {
	if connCtx, ok := config.GetSharedConnCtx(ctx); ok {
		if v, ok := connCtx.Load(ExtContentLength); ok {
			return v.(uint32), true
		}
	}
	return 0, false
}

func getChunkLength(ctx context.Context) (uint32, bool) {
	if connCtx, ok := config.GetSharedConnCtx(ctx); ok {
		if v, ok := connCtx.Load(ExtChunkedLength); ok {
			return v.(uint32), true
		}
	}
	return 0, false
}

func isChunkComplete(ctx context.Context) bool {
	contentLen, ok := getContentLength(ctx)
	if !ok {
		return false
	}
	chuckLen, ok := getChunkLength(ctx)
	if !ok {
		return false
	}
	return contentLen == chuckLen
}

func isLastAck(ctx context.Context) bool {
	if !enableFixOneway {
		return false
	}
	if connCtx, ok := config.GetSharedConnCtx(ctx); ok {
		if _, ok := connCtx.Load(ExtChunkedFinish); ok {
			return true
		}
	}
	return false
}
