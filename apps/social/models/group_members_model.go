package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupMembersModel = (*customGroupMembersModel)(nil)

type (
	// GroupMembersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupMembersModel.
	GroupMembersModel interface {
		groupMembersModel
		InsertWithSession(ctx context.Context, session sqlx.Session, data *GroupMembers) (sql.Result, error)
		FindByGroudIdAndUserId(ctx context.Context, userId, groupId string) (*GroupMembers, error)
		ListByUserId(ctx context.Context, userId string) ([]*GroupMembers, error)
		ListByGroupId(ctx context.Context, groupId string) ([]*GroupMembers, error)
	}

	customGroupMembersModel struct {
		*defaultGroupMembersModel
	}
)

// NewGroupMembersModel returns a model for the database table.
func NewGroupMembersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupMembersModel {
	return &customGroupMembersModel{
		defaultGroupMembersModel: newGroupMembersModel(conn, c, opts...),
	}
}

func (c *customGroupMembersModel) InsertWithSession(ctx context.Context, session sqlx.Session, data *GroupMembers) (sql.Result, error) {
	groupMembersIdKey := fmt.Sprintf("%s%v", cacheGroupMembersIdPrefix, data.Id)

	ret, err := c.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", c.table, groupMembersRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.GroupId, data.UserId, data.RoleLevel, data.JoinTime, data.JoinSource,
				data.InviterUid, data.OperatorUid)
		}
		return conn.ExecCtx(ctx, query, data.GroupId, data.UserId, data.RoleLevel, data.JoinTime, data.JoinSource,
			data.InviterUid, data.OperatorUid)
	}, groupMembersIdKey)

	return ret, err
}

func (c *customGroupMembersModel) FindByGroudIdAndUserId(ctx context.Context, userId, groupId string) (*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `group_id` = ?", groupMembersRows, c.table)

	var resp GroupMembers
	err := c.QueryRowNoCacheCtx(ctx, &resp, query, userId, groupId)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (c *customGroupMembersModel) ListByUserId(ctx context.Context, userId string) ([]*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", groupMembersRows, c.table)

	var resp []*GroupMembers
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, userId)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (c *customGroupMembersModel) ListByGroupId(ctx context.Context, groupId string) ([]*GroupMembers, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ?", groupMembersRows, c.table)

	var resp []*GroupMembers
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
