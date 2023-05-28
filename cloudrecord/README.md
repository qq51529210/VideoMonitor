# cloudrecord
云端录像管理服务，负责管理录像文件的存储信息。
1. 接受 recordassist 提交的信息，保存到数据库
1. 进行数据库的查询和删除

文档生成：swag init -o api/internal/docs --parseDependency