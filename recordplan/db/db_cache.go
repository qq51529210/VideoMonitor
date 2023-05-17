package db

import "sync"

type mapCacheKey[K int64 | string | WeekPlanStreamKey] interface {
	key() K
}

// mapCache 用于缓存数据
type mapCache[K int64 | string | WeekPlanStreamKey, M mapCacheKey[K]] struct {
	sync.Mutex
	// 数据
	d map[K]M
	// gorm.Model 要的模型
	m M
	// 标记数据 d 是否需要重新加载
	ok bool
	// 创建函数
	new func() M
}

func newMapCache[K int64 | string | WeekPlanStreamKey, M mapCacheKey[K]](nf func() M) *mapCache[K, M] {
	c := new(mapCache[K, M])
	c.d = make(map[K]M)
	c.new = nf
	c.m = nf()
	return c
}

// check 检查内存数据是否需要重新加载
func (c *mapCache[K, M]) check() error {
	// 看看标记是否有效
	if !c.ok {
		// 加载
		err := c.loadAll()
		if err != nil {
			// 失败标记
			c.ok = false
			return err
		}
		// 成功标记
		c.ok = true
	}
	return nil
}

// load 尝试加载，添加和修改时候调用
func (c *mapCache[K, M]) load(k K) {
	// 读取
	m := c.new()
	err := _db.First(m).Error
	// 失败
	if err != nil {
		c.ok = false
		return
	}
	// 成功
	c.d[k] = m
	c.ok = true
}

func (c *mapCache[K, M]) loadAll() error {
	c.d = make(map[K]M)
	//
	var models []M
	// 数据库
	err := _db.Find(&models).Error
	if err != nil {
		return err
	}
	// 内存
	for _, model := range models {
		c.d[model.key()] = model
	}
	return nil
}

// Load 尝试加载，添加和修改时候调用
func (c *mapCache[K, M]) Load(k K) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 读取
	c.load(k)
}

// All 返回所有
func (c *mapCache[K, M]) All() ([]M, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.check()
	if err != nil {
		return nil, err
	}
	// 列表
	var models []M
	for _, v := range c.d {
		models = append(models, v)
	}
	// 返回
	return models, nil
}

// Get 返回
func (c *mapCache[K, M]) Get(k K) (m M, err error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保加载
	err = c.check()
	if err != nil {
		return
	}
	// 返回
	m = c.d[k]
	// 返回
	return
}

// Add 添加
func (c *mapCache[K, M]) Add(m M) (int64, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 数据库
	db := _db.Create(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		c.load(m.key())
	}
	// 返回
	return db.RowsAffected, nil
}

// Update 更新
func (c *mapCache[K, M]) Update(m M) (int64, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 数据库
	db := _db.Updates(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		c.load(m.key())
	}
	// 返回
	return db.RowsAffected, nil
}

// Save 保存
func (c *mapCache[K, M]) Save(k K, m *M) (int64, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 数据库
	db := _db.Save(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		c.load(k)
	}
	return db.RowsAffected, nil
}

// Delete 删除
func (c *mapCache[K, M]) Delete(k K) (int64, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 数据库
	db := _db.Delete(c.m, k)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		delete(c.d, k)
	}
	// 返回
	return db.RowsAffected, nil
}

// BatchDelete 批量删除
func (c *mapCache[K, M]) BatchDelete(ks []K) (int64, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 数据库
	db := _db.Delete(c.m, ks)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if db.RowsAffected > 0 {
		for _, k := range ks {
			delete(c.d, k)
		}
	}
	// 返回
	return db.RowsAffected, nil
}
