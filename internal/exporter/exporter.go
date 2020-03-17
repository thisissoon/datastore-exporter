package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/oauth2/google"
)

var (
	datastoreApi = "https://datastore.googleapis.com/v1/"
	scopes       = []string{
		"https://www.googleapis.com/auth/datastore",
		"https://www.googleapis.com/auth/cloud-platform",
	}
)

const (
	STATE_UNSPECIFIED = "STATE_UNSPECIFIED"
	INITIALIZING      = "INITIALIZING"
	PROCESSING        = "PROCESSING"
	CANCELLING        = "CANCELLING"
	FINALIZING        = "FINALIZING"
	SUCCESSFUL        = "SUCCESSFUL"
	FAILED            = "FAILED"
	CANCELLED         = "CANCELLED"
)

type entityFilter struct {
	Kinds        []string `json:"kinds"`
	NamespaceIDs []string `json:"namespaceIds"`
}

type datastoreReqBody struct {
	OutputURLPrefix string       `json:"outputUrlPrefix,omitempty"`
	EntityFilter    entityFilter `json:"entityFilter,omitempty"`
}

type opError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

func (e opError) String() string {
	return fmt.Sprintf("code: %d: %s. Reason: %s", e.Code, e.Status, e.Message)
}

type operationResp struct {
	Name     string `json:"name,omitempty"`
	Metadata struct {
		OutputURLPrefix string `json:"outputUrlPrefix,omitempty"`
		Common          struct {
			StartTime time.Time `json:"startTime,omitempty"`
			EndTime   time.Time `json:"endTime,omitempty"`
			State     string    `json:"state,omitempty"`
		} `json:"common,omitempty"`
	} `json:"metadata,omitempty"`
	Done  bool     `json:"done,omitempty"`
	Error *opError `json:"error,omitempty"`
}

type Exporter struct {
	client     *http.Client
	log        zerolog.Logger
	ProjectID  string
	BucketName string
}

func NewExporter(ctx context.Context, log zerolog.Logger, projectID, bucketName string) (*Exporter, error) {
	c, err := google.DefaultClient(ctx, scopes...)
	if err != nil {
		return nil, fmt.Errorf("error getting default client: %v", err)
	}
	return &Exporter{
		client:     c,
		log:        log,
		ProjectID:  projectID,
		BucketName: bucketName,
	}, nil
}

// Export Initiates the datastore export operation exporting all kinds and namespaces
// Checks to see if the export operation was successful
func (e *Exporter) Export(ctx context.Context) error {
	apiUrl := fmt.Sprintf("%sprojects/%s:export", datastoreApi, e.ProjectID)

	b, err := json.Marshal(datastoreReqBody{
		OutputURLPrefix: "gs://" + e.BucketName,
		EntityFilter: entityFilter{
			Kinds:        []string{},
			NamespaceIDs: []string{},
		},
	})
	if err != nil {
		return fmt.Errorf("error marshalling body: %v", err)
	}
	body := bytes.NewReader(b)
	res, err := e.client.Post(apiUrl, "application/json", body)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var r operationResp
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	if res.StatusCode != 200 {
		if r.Error != nil {
			return fmt.Errorf("error making export call: %v", r.Error)
		}
		return fmt.Errorf("error making export call: %d: %s", res.StatusCode, res.Status)
	}
	e.log.Info().Msgf("Datastore export started. Output location will be: %s", r.Metadata.OutputURLPrefix)
	return e.watchOp(ctx, r.Name)
}

func (e *Exporter) watchOp(ctx context.Context, name string) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			resp, err := e.getOp(ctx, name)
			if err != nil {
				return err
			}
			state := resp.Metadata.Common.State
			switch state {
			case FAILED, CANCELLED:
				return fmt.Errorf("Datastore operation failed: %v", state)
			case SUCCESSFUL:
				e.log.Info().Msg("Datastore export success")
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (e *Exporter) getOp(ctx context.Context, name string) (*operationResp, error) {
	url := fmt.Sprintf("%s%s", datastoreApi, name)
	resp, err := e.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error making get operation request: %s", resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r operationResp
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
