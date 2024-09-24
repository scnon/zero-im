package models

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel

		FindOneByPhone(ctx context.Context, phone *string) (*Users, error)
		ListByName(ctx context.Context, name *string) ([]*Users, error)
		ListByIds(ctx context.Context, ids []string) ([]*Users, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c, opts...),
	}
}

func (m *customUsersModel) FindOneByPhone(ctx context.Context, phone *string) (*Users, error) {
	zeroImUsersIdKey := fmt.Sprintf("%s%v", cacheZeroImUsersIdPrefix, phone)
	var resp Users
	err := m.QueryRowCtx(ctx, &resp, zeroImUsersIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", usersRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, phone)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) ListByName(ctx context.Context, name *string) ([]*Users, error) {
	query := fmt.Sprintf("select %s from %s where `nickname` like ?", usersRows, m.table)
	var resp []*Users
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, fmt.Sprint("%", name, "%"))
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customUsersModel) ListByIds(ctx context.Context, ids []string) ([]*Users, error) {
	query := fmt.Sprintf("select %s from %s where `id` in ('%s')", usersRows, m.table, strings.Join(ids, "','"))
	var resp []*Users
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, ids)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
