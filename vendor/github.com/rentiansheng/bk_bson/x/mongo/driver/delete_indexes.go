// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"

	"time"

	"github.com/rentiansheng/bk_bsonbson"
	"github.com/rentiansheng/bk_bsonmongo/options"
	"github.com/rentiansheng/bk_bson/x/bsonx"
	"github.com/rentiansheng/bk_bsonx/mongo/driver/session"
	"github.com/rentiansheng/bk_bsonx/mongo/driver/topology"
	"github.com/rentiansheng/bk_bsonx/mongo/driver/uuid"
	"github.com/rentiansheng/bk_bsonx/network/command"
	"github.com/rentiansheng/bk_bsonx/network/description"
)

// DropIndexes handles the full cycle dispatch and execution of a dropIndexes
// command against the provided topology.
func DropIndexes(
	ctx context.Context,
	cmd command.DropIndexes,
	topo *topology.Topology,
	selector description.ServerSelector,
	clientID uuid.UUID,
	pool *session.Pool,
	opts ...*options.DropIndexesOptions,
) (bson.Raw, error) {

	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return nil, err
	}

	conn, err := ss.Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	dio := options.MergeDropIndexesOptions(opts...)
	if dio.MaxTime != nil {
		cmd.Opts = append(cmd.Opts, bsonx.Elem{"maxTimeMS", bsonx.Int64(int64(*dio.MaxTime / time.Millisecond))})
	}

	// If no explicit session and deployment supports sessions, start implicit session.
	if cmd.Session == nil && topo.SupportsSessions() {
		cmd.Session, err = session.NewClientSession(pool, clientID, session.Implicit)
		if err != nil {
			return nil, err
		}
		defer cmd.Session.EndSession()
	}

	return cmd.RoundTrip(ctx, ss.Description(), conn)
}
