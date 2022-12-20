package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm: 只支持指向结构体的一级指针")

	ErrNoRows = errors.New("orm: 没有数据")
	// ErrInsertZeroRow 代表插入 0 行
	ErrInsertZeroRow = errors.New("orm: 插入 0 行")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm: 不支持的表达式类型 %v", expr)
}
func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm: 未知字段 %s", name)
}

func NewErrUnknownColumn(name string) error {
	return fmt.Errorf("orm: 未知列 %s", name)
}

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm: 非法标签 %s", pair)
}

func NewErrUnsupportedAssignable(expr any) error {
	return fmt.Errorf("orm: 不支持的赋值表达式类型 %v", expr)
}
