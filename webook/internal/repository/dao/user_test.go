package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string
		// 为什么不用 ctrl ?
		// 因为这里是sqlmock 不是gomock
		mock    func(t *testing.T) *sql.DB
		ctx     context.Context
		user    User
		wantErr error
		//wantId  int64
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				// 这边预期的是正则表达式
				// 这个写法的意思就是 只要是 INSERT 到 users的语句
				// 这个写法的意思就是 只要是 INSERT 到users的语句
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
			//wantId: 3,
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 这边预期的是正则表达式
				// 这个写法的意思就是 只要是 INSERT 到users的语句
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				require.NoError(t, err)
				return mockDB
			},
			user:    User{},
			wantErr: ErrUserDuplicate,
			//wantId: 3,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				// 这边预期的是正则表达式
				// 这个写法的意思就是 只要是 INSERT 到users的语句
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDB
			},
			user:    User{},
			wantErr: errors.New("数据库错误"),
			//wantId: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//mockDB, mock, err := sqlmock.New()
			//require.NoError(t, err)
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// mock DB 不需要ping
				DisableAutomaticPing: true,
				// 这个是：
				SkipDefaultTransaction: true,
			})
			d := NewUserDAO(db)
			u := tc.user
			err = d.Insert(tc.ctx, u)
			assert.Equal(t, tc.wantErr, err)
			// 可以比较一下
			//assert.Equal(t, tc.wantId, u.Id)
		})

		// TODO: 理论上让GORM执行
		// Insert xxx
		// 实际上 GORM
		// BEGIN;
		// INSERT;
		// COMMIT;
	}
}
