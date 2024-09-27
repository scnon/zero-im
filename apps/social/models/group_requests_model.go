package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupRequestsModel = (*customGroupRequestsModel)(nil)

type (
	// GroupRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupRequestsModel.
	GroupRequestsModel interface {
		groupRequestsModel
		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *GroupRequests) error
		FindByGroupIdAndReqId(ctx context.Context, groupId, reqId string) (*GroupRequests, error)
		ListNoHandler(ctx context.Context, groupId string) ([]*GroupRequests, error)
	}

	customGroupRequestsModel struct {
		*defaultGroupRequestsModel
	}
)

// NewGroupRequestsModel returns a model for the database table.
func NewGroupRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupRequestsModel {
	return &customGroupRequestsModel{
		defaultGroupRequestsModel: newGroupRequestsModel(conn, c, opts...),
	}
}

func (c *customGroupRequestsModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return c.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (c *customGroupRequestsModel) UpdateWithSession(ctx context.Context, session sqlx.Session, data *GroupRequests) error {
	groupRequestsIdKey := fmt.Sprintf("%s%v", cacheGroupRequestsIdPrefix, data.Id)

	_, err := c.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", c.table, groupRequestsRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.ReqId, data.GroupId, data.ReqMsg, data.ReqTime, data.JoinSource,
			data.InviterUserId, data.HandleUserId, data.HandleTime, data.HandleResult, data.Id)
	}, groupRequestsIdKey)

	return err
}

func (c *customGroupRequestsModel) FindByGroupIdAndReqId(ctx context.Context, groupId, reqId string) (*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `req_id` = ? and `group_id` = ?", groupRequestsRows, c.table)

	var resp GroupRequests
	err := c.QueryRowNoCacheCtx(ctx, &resp, query, reqId, groupId)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (c *customGroupRequestsModel) ListNoHandler(ctx context.Context, groupId string) ([]*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `handle_result` = 1 ", groupRequestsRows, c.table)

	var resp []*GroupRequests
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
