// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mock

import (
	"context"
	"sync"

	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/bee/pkg/topology"
)

type mock struct {
	peers          []swarm.Address
	closestPeer    swarm.Address
	closestPeerErr error
	addPeerErr     error
	mtx            sync.Mutex
}

func WithAddPeerErr(err error) Option {
	return optionFunc(func(d *mock) {
		d.addPeerErr = err
	})
}

func WithClosestPeer(addr swarm.Address) Option {
	return optionFunc(func(d *mock) {
		d.closestPeer = addr
	})
}

func WithClosestPeerErr(err error) Option {
	return optionFunc(func(d *mock) {
		d.closestPeerErr = err
	})
}

func NewTopologyDriver(opts ...Option) topology.Driver {
	d := new(mock)
	for _, o := range opts {
		o.apply(d)
	}
	return d
}

func (d *mock) AddPeer(_ context.Context, addr swarm.Address) error {
	if d.addPeerErr != nil {
		return d.addPeerErr
	}

	d.mtx.Lock()
	d.peers = append(d.peers, addr)
	d.mtx.Unlock()
	return nil
}

func (d *mock) Peers() []swarm.Address {
	return d.peers
}

func (d *mock) ClosestPeer(addr swarm.Address) (peerAddr swarm.Address, err error) {
	return d.closestPeer, d.closestPeerErr
}

type Option interface {
	apply(*mock)
}

type optionFunc func(*mock)

func (f optionFunc) apply(r *mock) { f(r) }
