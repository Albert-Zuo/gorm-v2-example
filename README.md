## Use Gorm V2 Encrypt/Decrypt MySQL Data
基于 Golang 的 Gorm V2 库加解密 MySQL 数据库数据字段案例。

语言与 Gorm 框架都并不是关键，只需要能支持运行原生 SQL 语句即可。

思路：将AES加密后的二进制数据转为16进制存储，解密时将数据转回去后AES解密。

- 加解密代码实现路径：[gorm-v2-example/internal/repo/user.go](https://github.com/Albert-Zuo/gorm-v2-example/blob/164f3a40a60e95f97344512025d82e395c33c5e6/internal/repo/user.go#L1)
- 测试代码路径：[gorm-v2-example/test/user_test.go](https://github.com/Albert-Zuo/gorm-v2-example/blob/164f3a40a60e95f97344512025d82e395c33c5e6/test/user_test.go#L20)

**首先MySQL的数据库的编码格式不能是utf8，需要是utf8mb4**



#### 加密 Decrypt

代码实现：

```go
// CreateUser CreateUser
func (ur userRepo) CreateUser(user *domain.User) (bool, error) {
	ur.db.Model(domain.User{}).Create(map[string]interface{}{
		"UserName": clause.Expr{SQL: "HEX(AES_ENCRYPT(?, ?))", Vars: []interface{}{user.UserName, key}},
		"Password": clause.Expr{SQL: "HEX(AES_ENCRYPT(?, ?))", Vars: []interface{}{user.Password, key}},
		"Email":    clause.Expr{SQL: "HEX(AES_ENCRYPT(?, ?))", Vars: []interface{}{user.Email, key}},
		"Mobile":   clause.Expr{SQL: "HEX(AES_ENCRYPT(?, ?))", Vars: []interface{}{user.Mobile, key}},
	})

	return true, nil
}
```


等价于原生Sql语句：

```sql
INSERT INTO
    `table_name` (`email`,`mobile`,`password`,`user_name`)
VALUES(
    HEX(AES_ENCRYPT('xiaozuo1221@gmail.com', 'xiaozuo1221@gmail.com')),
    HEX(AES_ENCRYPT('12345678900', 'xiaozuo1221@gmail.com')),
    HEX(AES_ENCRYPT('00000', 'xiaozuo1221@gmail.com')),
    HEX(AES_ENCRYPT('test4', 'xiaozuo1221@gmail.com')))
```

测试运行后的控制台打印：
![测试用例控制台打印](http://images.gxuwzapp.top/2021/2/test.png)

数据库对应数据：
![数据库数据](http://images.gxuwzapp.top/2021/2/data.png)

#### 解密 Encrypt

代码实现：
```go
// ListUser ListUser
func (ur userRepo) ListUser() ([]*domain.User, error) {
	users := make([]*domain.User, 0) 
    // 使用CAST( )函数将查询字段作为一个整体查询，避免乱码
	rows, err := ur.db.Raw(`SELECT 
			user_id,
			CAST(AES_DECRYPT(UNHEX(user_name), ?) AS CHAR) AS user_name,
			CAST(AES_DECRYPT(UNHEX(password), ?) AS CHAR) AS password,
			CAST(AES_DECRYPT(UNHEX(email), ?) AS CHAR) AS email,
			CAST(AES_DECRYPT(UNHEX(mobile), ?) AS CHAR) AS mobile
		FROM blog_user`, key, key, key, key).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user domain.User
		ur.db.ScanRows(rows, &user)
		users = append(users, &user)
	}
	return users, nil
}
```


等价于原生Sql语句
```sql
SELECT 
	user_id,
	CAST(AES_DECRYPT(UNHEX(user_name), 'key') AS CHAR) AS user_name,
	CAST(AES_DECRYPT(UNHEX(password), 'key') AS CHAR) AS password,
	CAST(AES_DECRYPT(UNHEX(email), 'key') AS CHAR) AS email,
	CAST(AES_DECRYPT(UNHEX(mobile), 'key') AS CHAR) AS mobile
FROM table_name 
```



#### IO 执行效率
当前机器 i5 4核 16GB, MySQL 8.0

- 添加100个User对象耗时7s
- 查询100个User对象耗时0.01s
