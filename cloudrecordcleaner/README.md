# cloudrecordcleaner
云端录像清理服务，负责清理 minio 中的过期的文件
1. 与 cloudrecord 服务使用同一个数据库
1. 定期检查数据库中过期的数据，先打上软删除标记
1. 从 minio 中删除文件，或者检查是否存在，成功后删除数据库的数据