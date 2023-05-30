package util

import (
	"sync"

	"gorm.io/gorm"
)

// GORMCache 用于缓存数据
type GORMCache[K comparable, M any] struct {
	sync.Mutex
	// 数据库
	db *gorm.DB
	// 是否开启缓存
	cache bool
	// 数据
	D map[K]M
	// gorm.Model 要的模型
	m M
	// 标记数据 d 是否需要重新加载
	ok bool
	// 创建函数
	new func() M
	// 返回 M 的主键
	key func(M) K
	// 返回 M 的主键，用于单个，查询，修改，删除
	whereKey func(*gorm.DB, K) *gorm.DB
	// 返回 M 的主键列表，用于批量删除
	whereKeys func(*gorm.DB, []K) *gorm.DB
}

// NewGORMCache 返回新的缓存，enable 为 false 则不开启缓存
func NewGORMCache[K comparable, M any](
	db *gorm.DB,
	cache bool,
	newFunc func() M,
	keyFunc func(M) K,
	whereKeyFunc func(*gorm.DB, K) *gorm.DB,
	whereKeysFunc func(*gorm.DB, []K) *gorm.DB,
) *GORMCache[K, M] {
	c := new(GORMCache[K, M])
	c.Init(db, cache, newFunc, keyFunc, whereKeyFunc, whereKeysFunc)
	return c
}

func (c *GORMCache[K, M]) Init(
	db *gorm.DB,
	cache bool,
	newFunc func() M,
	keyFunc func(M) K,
	whereKeyFunc func(*gorm.DB, K) *gorm.DB,
	whereKeysFunc func(*gorm.DB, []K) *gorm.DB,
) {
	c.db = db
	c.cache = cache
	c.D = make(map[K]M)
	c.new = newFunc
	c.key = keyFunc
	c.whereKey = whereKeyFunc
	c.whereKeys = whereKeysFunc
	c.m = newFunc()
}

func (c *GORMCache[K, M]) IsCache() bool {
	return c.cache
}

func (c *GORMCache[K, M]) DB() *gorm.DB {
	return c.db
}

func (c *GORMCache[K, M]) Model() *gorm.DB {
	return c.db.Model(c.m)
}

// Check 检查内存数据是否需要重新加载，同步
func (c *GORMCache[K, M]) LoadAll() error {
	// 不启用
	if !c.cache {
		return nil
	}
	//
	c.Lock()
	err := c.loadAll()
	c.Unlock()
	return err
}

// loadAll 检查内存数据是否需要重新加载
func (c *GORMCache[K, M]) loadAll() error {
	// 数据是否有效
	if !c.ok {
		var models []M
		// 数据库
		err := c.db.Find(&models).Error
		if err != nil {
			c.ok = false
			return err
		}
		// 重置
		c.D = make(map[K]M)
		// 添加
		for _, model := range models {
			c.D[c.key(model)] = model
		}
		c.ok = true
	}
	return nil
}

// Load 尝试加载，添加和修改时候调用，同步
func (c *GORMCache[K, M]) Load(k K) {
	// 不启用
	if !c.cache {
		return
	}
	// 启用
	c.Lock()
	c.load(k)
	c.Unlock()
}

// load 尝试加载，添加和修改时候调用
func (c *GORMCache[K, M]) load(k K) error {
	// 读取
	m := c.new()
	err := c.whereKey(c.db, k).First(m).Error
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
	c.D[k] = m
	c.ok = true
	//
	return nil
}

// All 返回所有内存，不要修改返回的指针，同步
func (c *GORMCache[K, M]) All() ([]M, error) {
	// 不开启
	if !c.cache {
		var models []M
		err := c.db.Find(&models).Error
		if err != nil {
			return nil, err
		}
		return models, nil
	}
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll()
	if err != nil {
		return nil, err
	}
	// 列表
	var models []M
	for _, v := range c.D {
		models = append(models, v)
	}
	// 返回
	return models, nil
}

// Get 返回指定，不要修改返回的指针，同步
func (c *GORMCache[K, M]) Get(k K) (m M, err error) {
	// 不开启
	if !c.cache {
		mm := c.new()
		err = c.whereKey(c.db, k).First(mm).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = nil
			}
			return
		}
		m = mm
		return
	}
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保加载
	err = c.loadAll()
	if err != nil {
		return
	}
	// 返回
	m = c.D[k]
	// 返回
	return
}

// Add 添加，同步
func (c *GORMCache[K, M]) Add(m M) (int64, error) {
	// 数据库
	db := c.db.Create(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		c.load(c.key(m))
		c.Unlock()
	}
	// 返回
	return db.RowsAffected, nil
}

// Update 更新，同步
func (c *GORMCache[K, M]) Update(m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(c.db, k).Updates(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		c.load(c.key(m))
		c.Unlock()
	}
	// 返回
	return db.RowsAffected, nil
}

// BatchUpdate 更新内存
func (c *GORMCache[K, M]) BatchUpdateCache(fn func(M)) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保加载
	err := c.loadAll()
	if err != nil {
		return
	}
	// 更新
	for _, m := range c.D {
		fn(m)
	}
}

// Save 保存，同步
func (c *GORMCache[K, M]) Save(m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(c.db, k).Save(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		c.load(c.key(m))
		c.Unlock()
	}
	return db.RowsAffected, nil
}

// BatchSave 事务保存，同步
func (c *GORMCache[K, M]) BatchSave(ms []M) (int64, error) {
	// 数据库
	err := c.db.Transaction(func(tx *gorm.DB) error {
		for _, m := range ms {
			db := tx.Save(m)
			if db.Error != nil {
				return db.Error
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	// 内存
	c.LoadAll()
	//
	return int64(len(ms)), nil
}

// Delete 删除
func (c *GORMCache[K, M]) Delete(k K) (int64, error) {
	// 数据库
	db := c.whereKey(c.db, k).Delete(c.m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		delete(c.D, k)
		c.Unlock()
	}
	// 返回
	return db.RowsAffected, nil
}

// BatchDelete 批量删除
func (c *GORMCache[K, M]) BatchDelete(ks []K) (int64, error) {
	// 数据库
	db := c.whereKeys(c.db, ks).Delete(c.m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		for _, k := range ks {
			delete(c.D, k)
		}
		c.Unlock()
	}
	// 返回
	return db.RowsAffected, nil
}

// BatchDeleteCache 删除内存
func (c *GORMCache[K, M]) BatchDeleteCache(fn func(M) bool) {
	// 缓存
	if !c.cache {
		return
	}
	// 上锁
	c.Lock()
	defer c.Unlock()
	for k, m := range c.D {
		if fn(m) {
			delete(c.D, k)
		}
	}
}

// Search 在内存中查找
func (c *GORMCache[K, M]) Search(match func(M) bool) ([]M, error) {
	var mm []M
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll()
	if err != nil {
		return nil, err
	}
	// 查找
	for _, m := range c.D {
		if match(m) {
			mm = append(mm, m)
		}
	}
	//
	return mm, nil
}

// Search 在内存中查找
func (c *GORMCache[K, M]) SearchOne(match func(M) bool) (m M, err error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err = c.loadAll()
	if err != nil {
		return
	}
	// 查找
	for _, v := range c.D {
		if match(v) {
			m = v
			break
		}
	}
	//
	return
}

// Count 返回内存匹配数量
func (c *GORMCache[K, M]) Count(match func(M) bool) (int64, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll()
	if err != nil {
		return 0, err
	}
	// 查找
	var n int64
	for _, v := range c.D {
		if match(v) {
			n++
		}
	}
	//
	return n, nil
}

// Total 返回内存总量
func (c *GORMCache[K, M]) Total() int64 {
	c.Lock()
	n := int64(len(c.D))
	c.Unlock()
	return n
}

// List 查询数据库
func (c *GORMCache[K, M]) List(page *GORMPage, query GORMQuery, res *GORMList[M]) error {
	return gormList(c.db, page, query, res)
}

// SearchGORMCache 模板化的 Search
func SearchGORMCache[T any, K comparable, M any](c *GORMCache[K, M], match func(M) (bool, T)) ([]T, error) {
	var vv []T
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll()
	if err != nil {
		return nil, err
	}
	// 查找
	for _, m := range c.D {
		o, v := match(m)
		if o {
			vv = append(vv, v)
		}
	}
	//
	return vv, nil
}
