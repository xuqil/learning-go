package session

import (
	"github.com/google/uuid"
	"leanring-go/web"
)

type Manager struct {
	Propagator
	Store
	CtxSessKey string
}

// GetSession 尝试从 ctx 中拿到 Session
// 如果成功了，那么它会将 Session 实例缓存到 ctx 的UserValues 里面
func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	val, ok := ctx.UserValues[m.CtxSessKey]
	//ctx.Req.Context().Value(m.CtxSessKey)
	if ok {
		return val.(Session), nil
	}
	// 尝试缓存 session
	sessID, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	sess, err := m.Get(ctx.Req.Context(), sessID)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.CtxSessKey] = sess
	// 使用 Context 缓存 session 会使用复制，性能会有所损失，同时父 Context 无法访问子 Context 的内容
	//ctx.Req = ctx.Req.WithContext(context.WithValue(ctx.Req.Context(), m.CtxSessKey, sess))
	return sess, err
}

// InitSession 初始化一个 session，并且注入到 http response 里面
func (m *Manager) InitSession(ctx *web.Context) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	// 注入进去 HTTP 响应里面
	m.Inject(id, ctx.Resp)
	return sess, nil
}

// RefreshSession 刷新 Session
func (m *Manager) RefreshSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Refresh(ctx.Req.Context(), sess.ID())
}

// RemoveSession 删除 Session
func (m *Manager) RemoveSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
