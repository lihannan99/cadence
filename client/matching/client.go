package matching

import (
	"sync"
	"time"

	"golang.org/x/net/context"

	m "code.uber.internal/devexp/minions/.gen/go/matching"
	workflow "code.uber.internal/devexp/minions/.gen/go/shared"
	"code.uber.internal/devexp/minions/common"
	"code.uber.internal/devexp/minions/common/membership"
	tchannel "github.com/uber/tchannel-go"
	"github.com/uber/tchannel-go/thrift"
)

var _ Client = (*clientImpl)(nil)

type clientImpl struct {
	connection      *tchannel.Channel
	resolver        membership.ServiceResolver
	thriftCacheLock sync.RWMutex
	thriftCache     map[string]m.TChanMatchingService
}

// NewClient creates a new history service TChannel client
func NewClient(ch *tchannel.Channel, monitor membership.Monitor) (Client, error) {
	sResolver, err := monitor.GetResolver(common.MatchingServiceName)
	if err != nil {
		return nil, err
	}

	client := &clientImpl{
		connection:  ch,
		resolver:    sResolver,
		thriftCache: make(map[string]m.TChanMatchingService),
	}
	return client, nil
}

func (c *clientImpl) AddActivityTask(addRequest *m.AddActivityTaskRequest) error {
	client, err := c.getHostForRequest(addRequest.GetTaskList().GetName())
	if err != nil {
		return err
	}
	ctx, cancel := c.createContext()
	defer cancel()
	return client.AddActivityTask(ctx, addRequest)
}

func (c *clientImpl) AddDecisionTask(addRequest *m.AddDecisionTaskRequest) error {
	client, err := c.getHostForRequest(addRequest.GetTaskList().GetName())
	if err != nil {
		return err
	}
	ctx, cancel := c.createContext()
	defer cancel()
	return client.AddDecisionTask(ctx, addRequest)
}

func (c *clientImpl) PollForActivityTask(pollRequest *workflow.PollForActivityTaskRequest) (*workflow.PollForActivityTaskResponse, error) {
	client, err := c.getHostForRequest(pollRequest.GetTaskList().GetName())
	if err != nil {
		return nil, err
	}
	ctx, cancel := c.createLongPollContext()
	defer cancel()
	return client.PollForActivityTask(ctx, pollRequest)
}

func (c *clientImpl) PollForDecisionTask(pollRequest *workflow.PollForDecisionTaskRequest) (*workflow.PollForDecisionTaskResponse, error) {
	client, err := c.getHostForRequest(pollRequest.GetTaskList().GetName())
	if err != nil {
		return nil, err
	}
	ctx, cancel := c.createLongPollContext()
	defer cancel()
	return client.PollForDecisionTask(ctx, pollRequest)
}

func (c *clientImpl) getHostForRequest(key string) (m.TChanMatchingService, error) {
	host, err := c.resolver.Lookup(key)
	if err != nil {
		return nil, err
	}
	return c.getThriftClient(host.GetAddress()), nil
}

func (c *clientImpl) createContext() (thrift.Context, context.CancelFunc) {
	// TODO: make timeout configurable
	return thrift.NewContext(time.Second * 10)
}

func (c *clientImpl) createLongPollContext() (thrift.Context, context.CancelFunc) {
	// TODO: make timeout configurable
	return thrift.NewContext(time.Minute * 3)
}

func (c *clientImpl) getThriftClient(hostPort string) m.TChanMatchingService {
	c.thriftCacheLock.RLock()
	client, ok := c.thriftCache[hostPort]
	c.thriftCacheLock.RUnlock()
	if ok {
		return client
	}

	c.thriftCacheLock.Lock()
	defer c.thriftCacheLock.Unlock()

	// check again if in the cache cause it might have been added
	// before we acquired the lock
	client, ok = c.thriftCache[hostPort]
	if !ok {
		tClient := thrift.NewClient(c.connection, common.MatchingServiceName, &thrift.ClientOptions{
			HostPort: hostPort,
		})

		client = m.NewTChanMatchingServiceClient(tClient)
		c.thriftCache[hostPort] = client
	}
	return client
}