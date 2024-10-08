package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendRequestsModel = (*customFriendRequestsModel)(nil)

type (
	// FriendRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendRequestsModel.
	FriendRequestsModel interface {
		friendRequestsModel
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		UpdateWithReq(ctx context.Context, session sqlx.Session, data *FriendRequests) error
		ListNoHandler(ctx context.Context, userId string) ([]*FriendRequests, error)
		FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error)
	}

	customFriendRequestsModel struct {
		*defaultFriendRequestsModel
	}
)

// NewFriendRequestsModel returns a model for the database table.
func NewFriendRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendRequestsModel {
	return &customFriendRequestsModel{
		defaultFriendRequestsModel: newFriendRequestsModel(conn, c, opts...),
	}
}

func (c *customFriendRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return c.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (c *customFriendRequestsModel) FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `user_id` = ?", friendRequestsRows, c.table)

	var resp FriendRequests
	err := c.QueryRowNoCacheCtx(ctx, &resp, query, rid, uid)

	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (c *customFriendRequestsModel) UpdateWithReq(ctx context.Context, session sqlx.Session, data *FriendRequests) error {
	friendRequestsIdKey := fmt.Sprintf("%s%v", cacheFriendRequestsIdPrefix, data.Id)

	_, err := c.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", c.table, friendRequestsRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime, data.HandleResult, data.HandleMsg, data.HandledAt, data.Id)
	}, friendRequestsIdKey)

	return err
}

func (c *customFriendRequestsModel) ListNoHandler(ctx context.Context, userId string) ([]*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `handle_result` = 1 and `user_id` = ?", friendRequestsRows, c.table)

	var resp []*FriendRequests
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, userId)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
