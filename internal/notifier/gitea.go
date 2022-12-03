package notifier

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	eventv1 "github.com/fluxcd/pkg/apis/event/v1beta1"
	"github.com/fluxcd/pkg/apis/meta"
	"net/http"
	"net/url"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

type Gitea struct {
	BaseURL string
	Token   string
	Owner   string
	Repo    string
	client  *gitea.Client
	debug   bool
}

var _ Interface = &Gitea{}

func NewGitea(addr string, token string, certPool *x509.CertPool) (*Gitea, error) {
	if len(token) == 0 {
		return nil, errors.New("github token cannot be empty")
	}

	host, id, err := parseGitAddress(addr)
	if err != nil {
		return nil, err
	}

	if _, err := url.Parse(host); err != nil {
		return nil, err
	}

	comp := strings.Split(id, "/")
	if len(comp) != 2 {
		return nil, fmt.Errorf("invalid repository id %q", id)
	}

	client, err := gitea.NewClient(host, gitea.SetToken(token))
	if err != nil {
		return nil, err
	}

	if certPool != nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		}
		client.SetHTTPClient(&http.Client{Transport: tr})
	}

	debug := false
	if os.Getenv("NOTIFIER_GITEA_DEBUG") == "true" {
		debug = true
	}

	return &Gitea{
		BaseURL: host,
		Token:   token,
		Owner:   comp[0],
		Repo:    comp[1],
		client:  client,
		debug:   debug,
	}, err
}

func (g *Gitea) Post(ctx context.Context, event eventv1.Event) error {
	revString, ok := event.Metadata[eventv1.MetaRevisionKey]
	if !ok {
		return errors.New("missing revision metadata")
	}
	rev, err := parseRevision(revString)
	if err != nil {
		return err
	}
	state, err := toGiteaState(event)
	if err != nil {
		return err
	}
	name, desc := formatNameAndDescription(event)

	status := gitea.CreateStatusOption{
		State:       state,
		TargetURL:   "",
		Description: desc,
		Context:     name,
	}

	listStatusesOpts := gitea.ListStatusesOption{
		ListOptions: gitea.ListOptions{
			Page:     0,
			PageSize: 50,
		},
	}
	statuses, _, err := g.client.ListStatuses(g.Owner, g.Repo, rev, listStatusesOpts)
	if err != nil {
		return fmt.Errorf("could not list commit statuses: %v", err)
	}
	if duplicateGiteaStatus(statuses, &status) {
		if g.debug {
			ctrl.Log.Info("gitea skip posting duplicate status",
				"owner", g.Owner, "repo", g.Repo, "commit_hash", rev, "status", status)
		}
		return nil
	}

	if g.debug {
		ctrl.Log.Info("gitea create commit begin",
			"base_url", g.BaseURL, "token", g.Token, "event", event, "status", status)
	}

	st, rsp, err := g.client.CreateStatus(g.Owner, g.Repo, rev, status)
	if err != nil {
		if g.debug {
			ctrl.Log.Error(err, "gitea create commit failed", "status", status)
		}
		return err
	}

	if g.debug {
		ctrl.Log.Info("gitea create commit ok", "response", rsp, "response_status", st)
	}

	return nil
}

func toGiteaState(event eventv1.Event) (gitea.StatusState, error) {
	// progressing events
	if event.HasReason(meta.ProgressingReason) {
		// pending
		return gitea.StatusPending, nil
	}
	switch event.Severity {
	case eventv1.EventSeverityInfo:
		return gitea.StatusSuccess, nil
	case eventv1.EventSeverityError:
		return gitea.StatusFailure, nil
	default:
		return gitea.StatusError, errors.New("can't convert to gitea state")
	}
}

// duplicateStatus return true if the latest status
// with a matching context has the same state and description
func duplicateGiteaStatus(statuses []*gitea.Status, status *gitea.CreateStatusOption) bool {
	if status == nil || statuses == nil {
		return false
	}

	for _, s := range statuses {
		if s.Context == "" || s.State == "" || s.Description == "" {
			continue
		}

		if s.Context == status.Context {
			if s.State == status.State && s.Description == status.Description {
				return true
			}

			return false
		}
	}

	return false
}
