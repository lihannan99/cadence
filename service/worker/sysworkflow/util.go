// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sysworkflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgryski/go-farm"
	"github.com/uber/cadence/.gen/go/shared"
	"github.com/uber/cadence/common/blobstore/blob"
	"strings"
)

const (
	// HistoryBlobKeyExt is the blob key extension on all history blobs
	HistoryBlobKeyExt = "history"
)

type (
	// HistoryBlobHeader is the header attached to all history blobs
	HistoryBlobHeader struct {
		DomainName           *string `json:"domain_name,omitempty"`
		DomainID             *string `json:"domain_id,omitempty"`
		WorkflowID           *string `json:"workflow_id,omitempty"`
		RunID                *string `json:"run_id,omitempty"`
		CurrentPageToken     *string `json:"current_page_token,omitempty"`
		NextPageToken        *string `json:"next_page_token,omitempty"`
		FirstFailoverVersion *string `json:"first_failover_version,omitempty"`
		LastFailoverVersion  *string `json:"last_failover_version,omitempty"`
		FirstEventID         *int64  `json:"first_event_id,omitempty"`
		LastEventID          *int64  `json:"last_event_id,omitempty"`
		UploadDateTime       *string `json:"upload_date_time,omitempty"`
		UploadCluster        *string `json:"upload_cluster,omitempty"`
		EventCount           *int64  `json:"event_count,omitempty"`
	}

	// HistoryBlob is the serializable data that forms the body of a blob
	HistoryBlob struct {
		Header *HistoryBlobHeader `json:"header"`
		Body   *shared.History    `json:"body"`
	}
)

// NewHistoryBlobKey returns a key for history blob
func NewHistoryBlobKey(domainID, workflowID, runID, pageToken, failoverVersion string) (blob.Key, error) {
	if len(domainID) == 0 || len(workflowID) == 0 || len(runID) == 0 || len(pageToken) == 0 || len(failoverVersion) == 0 {
		return nil, errors.New("all inputs required to be non-empty")
	}
	domainIDHash := fmt.Sprintf("%v", farm.Fingerprint64([]byte(domainID)))
	workflowIDHash := fmt.Sprintf("%v", farm.Fingerprint64([]byte(workflowID)))
	runIDHash := fmt.Sprintf("%v", farm.Fingerprint64([]byte(runID)))
	combinedHash := strings.Join([]string{domainIDHash, workflowIDHash, runIDHash}, "")
	return blob.NewKey(HistoryBlobKeyExt, combinedHash, pageToken, failoverVersion)
}

// ConvertHeaderToTags converts header into metadata tags for blob
func ConvertHeaderToTags(header *HistoryBlobHeader) (map[string]string, error) {
	var tempMap map[string]interface{}
	bytes, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &tempMap); err != nil {
		return nil, err
	}
	result := make(map[string]string, len(tempMap))
	for k, v := range tempMap {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
