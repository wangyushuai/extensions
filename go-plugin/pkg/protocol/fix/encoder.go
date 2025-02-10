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

	"mosn.io/api"
	"mosn.io/pkg/log"
)

func encodeRequest(ctx context.Context, request *Request) (api.IoBuffer, error) {
	// This message can be parsed on demand, if there is an observability or governance requirement.
	log.DefaultLogger.Infof("[fix] send req:%s", request.Data.String())
	return request.Data.Clone(), nil
}

func encodeResponse(ctx context.Context, response *Response) (api.IoBuffer, error) {
	//This message can be parsed on demand, if there is an observability or governance requirement.
	if _, ok := response.Get(ExtChunkedFinish); ok {
		setChunkFinished(ctx, true)
	}
	log.DefaultLogger.Infof("[fix] send resp:%s", response.Data.String())
	return response.Data.Clone(), nil
}
