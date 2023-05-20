package db

import (
	"sync"

	"gorm.io/gorm"
)

var (
	enableCache bool
)

// mapCache 用于缓存数据
type mapCache[K comparable, M any] struct {
	sync.Mutex
	// 数据
	d map[K]M
	// gorm.Model 要的模型
	m M
	// 标记数据 d 是否需要重新加载
	ok bool
	// 创建函数
	new func() M
	//
	key       func(M) K
	whereKey  func(*gorm.DB, K) *gorm.DB
	whereKeys func(*gorm.DB, []K) *gorm.DB
}

// newMapCache 返回新的缓存
func newMapCache[K comparable, M any](
	newFunc func() M,
	keyFunc func(M) K,
	whereKeyFunc func(*gorm.DB, K) *gorm.DB,
	whereKeysFunc func(*gorm.DB, []K) *gorm.DB) *mapCache[K, M] {
	c := new(mapCache[K, M])
	c.d = make(map[K]M)
	c.new = newFunc
	c.key = keyFunc
	c.whereKey = whereKeyFunc
	c.whereKeys = whereKeysFunc
	c.m = newFunc()
	return c
}

// check 检查内存数据是否需要重新加载
func (c *mapCache[K, M]) check() error {
	// 没有缓存
	if !enableCache {
		return nil
	}
	// 看看标记是否有效
	if !c.ok {
		// 加载
		return c.loadAll()
	}
	return nil
}

// load 尝试加载，添加和修改时候调用
func (c *mapCache[K, M]) load(k K) error {
	// 没有缓存
	if !enableCache {
		return nil
	}
	// 读取
	m := c.new()
	err := c.whereKey(_db, k).First(m).Error
	// 失败
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.ok = true
			return nil
		}
		c.ok = false
		return err
	}
	// 成功
	c.d[k] = m
	c.ok = true
	return nil
}

func (c *mapCache[K, M]) loadAll() error {
	// 没有缓存
	if !enableCache {
		return nil
	}
	var models []M
	// 数据库
	err := _db.Find(&models).Error
	if err != nil {
		c.ok = false
		return err
	}
	// 重置
	c.d = make(map[K]M)
	// 添加
	for _, model := range models {
		k := c.key(model)
		c.d[k] = model
	}
	c.ok = true
	return nil
}

// Load 尝试加载，添加和修改时候调用
func (c *mapCache[K, M]) Load(k K) {
	// 没有缓存
	if !enableCache {
		return
	}
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 读取
	c.load(k)
}

// All 返回所有
func (c *mapCache[K, M]) All() ([]M, error) {
	// 没有缓存
	if !enableCache {
		var models []M
		err := _db.Find(&models).Error
		if err != nil {
			return nil, err
		}
		return models, nil
	}
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
	// 没有缓存
	if !enableCache {
		mm := c.new()
		if err = c.whereKey(_db, k).First(mm).Error; err != nil {
			return
		}
		m = mm
		return
	}
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
	// 数据库
	db := _db.Create(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if enableCache && db.RowsAffected > 0 {
		// 上锁
		c.Lock()
		defer c.Unlock()
		// 加载
		c.load(c.key(m))
	}
	// 返回
	return db.RowsAffected, nil
}

// Update 更新
func (c *mapCache[K, M]) Update(m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(_db, k).Updates(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if enableCache && db.RowsAffected > 0 {
		// 上锁
		c.Lock()
		defer c.Unlock()
		// 加载
		c.load(k)
	}
	// 返回
	return db.RowsAffected, nil
}

// Save 保存
func (c *mapCache[K, M]) Save(m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(_db, k).Save(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if enableCache && db.RowsAffected > 0 {
		// 上锁
		c.Lock()
		defer c.Unlock()
		// 加载
		c.load(k)
	}
	return db.RowsAffected, nil
}

// Delete 删除
func (c *mapCache[K, M]) Delete(k K) (int64, error) {
	// 数据库
	db := c.whereKey(_db, k).Delete(c.m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if enableCache && db.RowsAffected > 0 {
		// 上锁
		c.Lock()
		defer c.Unlock()
		// 删除
		delete(c.d, k)
	}
	// 返回
	return db.RowsAffected, nil
}

// BatchDelete 批量删除
func (c *mapCache[K, M]) BatchDelete(ks []K) (int64, error) {
	// 数据库
	db := c.whereKeys(_db, ks).Delete(c.m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if enableCache && db.RowsAffected > 0 {
		// 上锁
		c.Lock()
		defer c.Unlock()
		// 删除
		for _, k := range ks {
			delete(c.d, k)
		}
	}
	// 返回
	return db.RowsAffected, nil
}
