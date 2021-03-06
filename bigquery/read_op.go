// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigquery

import "golang.org/x/net/context"

// RecordsPerRequest returns a ReadOption that sets the number of records to fetch per request when streaming data from BigQuery.
func RecordsPerRequest(n int64) ReadOption { return recordsPerRequest(n) }

type recordsPerRequest int64

func (opt recordsPerRequest) customizeRead(conf *pagingConf) {
	conf.recordsPerRequest = int64(opt)
	conf.setRecordsPerRequest = true
}

// StartIndex returns a ReadOption that sets the zero-based index of the row to start reading from.
func StartIndex(i uint64) ReadOption { return startIndex(i) }

type startIndex uint64

func (opt startIndex) customizeRead(conf *pagingConf) {
	conf.startIndex = uint64(opt)
}

func (c *Client) readTable(src *Table, options []ReadOption) (*Iterator, error) {
	conf := &readTabledataConf{}
	src.customizeReadSrc(conf)

	for _, o := range options {
		o.customizeRead(&conf.paging)
	}

	pageFetcher := func(ctx context.Context, token string) (*readDataResult, error) {
		conf.paging.pageToken = token
		return c.service.readTabledata(ctx, conf)
	}

	it := &Iterator{
		pf: pageFetcher,
	}
	return it, nil
}
