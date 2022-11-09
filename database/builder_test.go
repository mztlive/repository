package database

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

type testStruct struct {
	Identity  string     `db:"identity"`
	Name      string     `db:"name"`
	Phone     string     `db:"phone"`
	CreatedAt *time.Time `db:"created_at"`
}

type allTypeStruct struct {
	Identity  string     `db:"identity"`
	Age       int        `db:"age"`
	isDeleted bool       `db:"is_deleted"`
	Amount    float64    `db:"amount"`
	CreatedAt *time.Time `db:"created_at"`
}

type noFieldStruct struct {
	Name  string
	Phone string
}

func TestBuilder(t *testing.T) {
	convey.Convey("测试所有字段都有值", t, func() {
		nowTime := time.Now()
		obj := testStruct{
			Identity:  "test",
			Name:      "test",
			Phone:     "test",
			CreatedAt: &nowTime,
		}

		sql, args, err := BuilderInsertSQL(&obj)
		convey.So(err, convey.ShouldBeNil)
		convey.So(sql, convey.ShouldEqual, "insert into test_structs(identity,name,phone,created_at) values(?,?,?,?)")
		convey.So(args, convey.ShouldResemble, []interface{}{"test", "test", "test", &nowTime})
	})

	convey.Convey("测试部分字段都有值", t, func() {
		obj := testStruct{
			Identity: "test",
			Name:     "test",
			Phone:    "test",
		}

		sql, args, err := BuilderInsertSQL(&obj)
		convey.So(err, convey.ShouldBeNil)
		convey.So(sql, convey.ShouldEqual, "insert into test_structs(identity,name,phone) values(?,?,?)")
		convey.So(args, convey.ShouldResemble, []interface{}{"test", "test", "test"})
	})

	convey.Convey("测试所有类型的字段", t, func() {
		nowTime := time.Now()
		obj := allTypeStruct{
			Identity:  "test",
			Age:       1,
			isDeleted: true,
			Amount:    1.1,
			CreatedAt: &nowTime,
		}

		sql, _, err := BuilderInsertSQL(&obj)
		convey.So(err, convey.ShouldBeNil)
		convey.So(sql, convey.ShouldEqual, "insert into all_type_structs(identity,age,is_deleted,amount,created_at) values(?,?,?,?,?)")
	})

	convey.Convey("测试没有tag的结构体", t, func() {
		obj := noFieldStruct{
			Name:  "123123",
			Phone: "123123",
		}

		_, _, err := BuilderInsertSQL(&obj)
		convey.So(err, convey.ShouldBeError)
		convey.So(err.Error(), convey.ShouldEqual, "No Builder. Check your struct")
	})
}
