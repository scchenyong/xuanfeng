ETL工具集
====

## 工具说明

ETL工具集是包含了文件处理、数据抽取、数据导出、数据转换、数据过滤等一系列数据处理工具。

工具集统一实现了接受参数：
    -t 执行时间，支持：YYYYMMDD和
    -f
    -n
    -c
    -o


## 工具清单
### 数据抽取
#### oracle2file

#### ftp2file

### 文件处理
#### zip2file

#### file2zip
=======
## 说明

ETL工具集是包含了文件处理、数据抽取、数据导出、数据转换、数据过滤等一系列数据处理工具。

工具集统一实现的执行参数：
- `-t` 执行时间，支持：年月日（YYYYMMDD）和 年月日时分秒（YYYYMMDDHH24MISS）两种格式
- `-f` 执行流程序号，调度时自动带入
- `-n` 执行任务节点编号，调度时自动带入
- `-c` 执行程序配置项ID
- `-o` 是否开启命令行显示日志


## 清单
### 数据抽取
#### oracle2file
从ORACLE数据库中导出数据到文件

##### 配置示例
```json
[
	{
		"id": "123",
		"datatime": {
			"unit": "day",
			"offset": -1
		},
		"source": {
			"host": "127.0.0.1",
			"port": 1521,
			"username": "test",
			"password": "123456",
			"dbname": "ORCL",
			"script": "test.sql"
		},
		"target": {
			"path": "/data/oracle",
			"file": "test.txt" 
		},
		"content": {
			"prefix": "[",
			"suffix": "]",
			"empty": "-",
			"gap": "\\005"
		}
	}
]
```

##### 配置项说明
- id: 配置项ID，用作于执行参数中的 `-c` 值
- datatime: 数据对处理时间配置
    - unit: 时间单位，支持：`year` `month` `day` `hour` `minute` `second`
    - offset: 偏移量，无：0，前一天：-1
- source: 数据源（ORACLE）配置
    - host: 服务器地址
    - port: 数据库端口
    - username: 用户名
    - password: 密码
    - dbname: 数据库名
    - script: 执行脚本路径（非绝对对径）
              对应在`$ETL_HOME/scripts`下的子路径
              内容支持日期替换
- target: 数据文件存放目标配置
    - path: 存放文件路径
    - file: 存放文件名，支持日期替换
- content: 内容加工处理配置
    - prefix: 字段前缀，支持ASCII字符
    - suffix: 字段后缀，支持ASCII字符
    - empty: 空代替，支持ASCII字符
    - gap: 分列字符，支持ASCII字符

#### ftp2file
从FTP服务器中下载文件

##### 配置示例
```json
[
	{	
		"id": "123",
		"datatime": {
			"unit": "day",
			"offset": -1
		},
		"source": {
			"host": "127.0.0.1",
			"port": 21,
			"username": "root",
			"password": "123456",
			"path": "/ftproot/data",
			"file": "test.tar.gz" 
		},
		"target": {
			"path": "/data",
			"file": "test.tar.gz" 
		}
	}
]
```

##### 配置项说明
- id: 配置项ID，用作于执行参数中的 `-c` 值
- datatime: 数据对处理时间配置
    - unit: 时间单位，支持：`year` `month` `day` `hour` `minute` `second`
    - offset: 偏移量，无：0，前一天：-1
- source: 数据文件源（FTP）配置
    - host: 服务器地址
    - port: 服务端口
    - username: 用户名
    - password: 密码
    - path: FTP服务器存放文件路径
    - file: FTP服务器存放文件名，支持日期替换
- target: 下载后文件存放目标配置
    - path: 存放文件路径
    - file: 存放文件名，支持日期替换

### 文件处理
#### zip2file
解压文件

```json
[
	{	
		"id": "123",
		"datatime": {
			"unit": "day",
			"offset": -1
		},
		"datazip": {
		    "type": ".tar.gz",
		    "sourcePath": "/data/test.tar.gz",
		    "targetPath": "/data/test-YYYYMMDD",
		    "clean": false
		}
	}
]
```

##### 配置项说明
- id: 配置项ID，用作于执行参数中的 `-c` 值
- datatime: 数据对处理时间配置
    - unit: 时间单位，支持：`year` `month` `day` `hour` `minute` `second`
    - offset: 偏移量，无：0，前一天：-1
- datazip: 数据解压配置
    - type: 压缩文件类型，支持：`.zip` `.tar` `.gz` `.tar.gz`
    - sourcePath: 压缩文件存放路径，支持日期替换
    - targetPath: 解压文件存放路径，支持日期替换
    - clean: 是否清理源文件（压缩文件）

#### file2zip
压缩文件

```json
[
	{	
		"id": "123",
		"datatime": {
			"unit": "day",
			"offset": -1
		},
		"datazip": {
		    "type": ".tar.gz",
		    "sourcePath": "/data/test.tar.gz",
		    "targetPath": "/data/test-YYYYMMDD",
		    "clean": false
		}
	}
]
```

##### 配置项说明
- id: 配置项ID，用作于执行参数中的 `-c` 值
- datatime: 数据对处理时间配置
    - unit: 时间单位，支持：`year` `month` `day` `hour` `minute` `second`
    - offset: 偏移量，无：0，前一天：-1
- datazip: 数据压缩配置
    - type: 压缩文件类型，支持：`.zip` `.tar` `.gz` `.tar.gz`
    - hasSelf: 是否包含文件目录自己，只对目录有效
    - sourcePath: 压缩文件存放路径，支持日期替换
    - targetPath: 压缩后文件存放路径，支持日期替换
    - clean: 是否清理源文件（压缩前的文件）
