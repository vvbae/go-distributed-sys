// DEPRECATED: use DistributedLog instead
package log

// import (
// 	"context"
// 	"sync"

// 	api "github.com/vvbae/proglog/api/v1"
// 	"go.uber.org/zap"
// 	"google.golang.org/grpc"
// )

// type Replicator struct {
// 	DialOptions []grpc.DialOption
// 	LocalServer api.LogClient

// 	logger *zap.Logger

// 	mu      sync.Mutex
// 	servers map[string]chan struct{} // map of server addr, use to stop replicate when server leave/fail
// 	closed  bool
// 	close   chan struct{}
// }

// func (r *Replicator) Join(name, addr string) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	r.init()
// 	if r.closed {
// 		return nil
// 	}

// 	if _, ok := r.servers[name]; ok {
// 		// already replicating so skip
// 		return nil
// 	}

// 	r.servers[name] = make(chan struct{})
// 	go r.replicate(addr, r.servers[name])
// 	return nil
// }

// func (r *Replicator) replicate(addr string, leave chan struct{}) {
// 	cc, err := grpc.Dial(addr, r.DialOptions...)
// 	if err != nil {
// 		r.logError(err, "failed to dial", addr)
// 		return
// 	}
// 	defer cc.Close()

// 	// client and open up a stream to consume all logs on the server
// 	client := api.NewLogClient(cc)
// 	ctx := context.Background()
// 	stream, err := client.ConsumeStream(ctx,
// 		&api.ConsumeRequest{
// 			Offset: 0,
// 		})

// 	if err != nil {
// 		r.logError(err, "failed to consume", addr)
// 		return
// 	}

// 	// create records channe;=l
// 	records := make(chan *api.Record)
// 	go func() {
// 		for {
// 			recv, err := stream.Recv()
// 			if err != nil {
// 				r.logError(err, "failed to receive", addr)
// 				return
// 			}
// 			records <- recv.Record
// 		}
// 	}()

// 	// consumes the logs from the discovered server in a stream and then
// 	// produces to the local server to save a copy
// 	for {
// 		select {
// 		// it that server fails and replicator closes the channel for that server
// 		case <-r.close:
// 			return
// 			// if that server left
// 		case <-leave:
// 			return
// 			// if get record on records channel
// 		case record := <-records:
// 			_, err = r.LocalServer.Produce(ctx,
// 				&api.ProduceRequest{
// 					Record: record,
// 				})
// 			if err != nil {
// 				r.logError(err, "failed to produce", addr)
// 				return
// 			}
// 		}
// 	}
// }

// // close channel and remove contacts
// func (r *Replicator) Leave(name string) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	r.init()
// 	if _, ok := r.servers[name]; !ok {
// 		return nil
// 	}
// 	close(r.servers[name])
// 	delete(r.servers, name)
// 	return nil
// }

// func (r *Replicator) init() {
// 	if r.logger == nil {
// 		r.logger = zap.L().Named("replicator")
// 	}
// 	if r.servers == nil {
// 		r.servers = make(map[string]chan struct{})
// 	}
// 	if r.close == nil {
// 		r.close = make(chan struct{})
// 	}
// }

// // Close() closes the replicator so it doesn’t replicate new servers that join the cluster
// // and it stops replicating existing servers
// func (r *Replicator) Close() error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	r.init()
// 	if r.closed {
// 		return nil
// 	}
// 	r.closed = true
// 	close(r.close)
// 	return nil
// }

// func (r *Replicator) logError(err error, msg, addr string) {
// 	r.logger.Error(
// 		msg,
// 		zap.String("addr", addr), zap.Error(err),
// 	)
// }
