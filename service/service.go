package service

import (
	"lchat/service/conf"
	"lchat/service/router"
	"lchat/service/logger"
	"lchat/service/entity"
)

func Run() {
	// 初始化日志
	logger.New()

	// 加载配置信息
	err := conf.Load()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// 初始化数据和仓库
	err = entity.NewStore()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer entity.Close()

	// 初始化管理员信息
	err = initAdmin()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	//开启路由
	err = router.Run()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

// 初始化管理员信息
func initAdmin() error {
	if (!entity.HasAdmin()) {
		admin := &entity.User{
			Email: conf.Get().Admin.Email,
			Phone: conf.Get().Admin.Phone,
			Password: conf.Get().Admin.Password,
			NickName: conf.Get().Admin.NickName,
			IsAdmin: true,
		}
		err := admin.Save()
		if err != nil {
			return err
		}
	}
	return nil
}
