package model

// GetAllModels 维护了系统中所有需要自动建表的实体模型
// 以后每次新增一张表（比如 User, Device），只需要在这里加一行即可
func GetAllModels() []interface{} {
	return []interface{}{
		&Log{},
		// &User{},
		// &Device{},
		// &DeviceLog{},
	}
}
