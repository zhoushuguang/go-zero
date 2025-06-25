package main

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 使用Mysql语句一张user表，有id、name、age三个字段，id为主键，name唯一
//CREATE TABLE user (
//	id INT AUTO_INCREMENT PRIMARY KEY,
//	name VARCHAR(100) NOT NULL UNIQUE,
//	age INT NOT NULL
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

//CREATE TABLE classroom (
//	id INT AUTO_INCREMENT PRIMARY KEY,
//	name VARCHAR(100) NOT NULL,
//	student_id INT NOT NULL
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

// 往user表插入3条数据
//INSERT INTO user (name, age) VALUES ('Alice_replica_three', 30), ('Bob_replica_three', 25), ('Charlie_replica_three', 35);

type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"` // 用户名
	Age  int64  `db:"age"`  // 年龄
}

func main() {
	mainDSN := "root:123456@tcp(127.0.0.1:3306)/main?parseTime=true&loc=Local"
	replicaOne := "root:123456@tcp(127.0.0.1:3306)/replica_one?parseTime=true&loc=Local"
	replicaTwo := "root:123456@tcp(127.0.0.1:3306)/replica_two?parseTime=true&loc=Local"
	replicaThree := "root:123456@tcp(127.0.0.1:3306)/replica_three?parseTime=true&loc=Local"

	conn := sqlx.NewMysql(mainDSN, sqlx.WithReplicas([]string{replicaOne, replicaTwo, replicaThree}...), sqlx.WithPolicy(sqlx.PolicyRandom))

	ctx := context.Background()
	//ctx = sqlx.WithReadMode(ctx)

	//for i := 0; i < 10; i++ {
	//	var user User
	//	err := conn.QueryRowCtx(ctx, &user, "SELECT * FROM user WHERE age = ?", 30)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("User found: ID=%d, Name=%s, Age=%d\n", user.ID, user.Name, user.Age)
	//	time.Sleep(time.Second)
	//}

	//stmt, err := conn.PrepareCtx(ctx, "SELECT * FROM user WHERE age = ?")
	//if err != nil {
	//	panic(err)
	//}

	//for i := 0; i < 10; i++ {
	//	stmt, err := conn.PrepareCtx(ctx, "SELECT * FROM user WHERE age = ?")
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	var user User
	//	err = stmt.QueryRowCtx(ctx, &user, 30)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("User found: ID=%d, Name=%s, Age=%d\n", user.ID, user.Name, user.Age)
	//	time.Sleep(time.Second)
	//}

	err := conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		r, err := session.ExecCtx(ctx, "insert into user (name, age) values (?, ?)", "light_2", 30)
		if err != nil {
			return err
		}
		id, err := r.LastInsertId()
		if err != nil {
			return err
		}
		r, err = session.ExecCtx(ctx, "insert into classroom (name, student_id) values (?, ?)", "class_one", id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
