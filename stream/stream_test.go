package stream

import (
	"fmt"
	"testing"
)

type User struct {
	Name   string
	Age    int
	Money  float64
	Gender string
}

var users = []User{
	{
		Name:   "zhangsan",
		Age:    18,
		Money:  1.1,
		Gender: "man",
	},
	{
		Name:   "lisi",
		Age:    19,
		Money:  1.2,
		Gender: "woman",
	},
	{
		Name:   "wangwu",
		Age:    20,
		Money:  1.3,
		Gender: "man",
	},
	{
		Name:   "zhaoliu",
		Age:    21,
		Money:  1.4,
		Gender: "man",
	},
}

func TestStream(t *testing.T) {
	list := []int64{1, 2, 3, 4, 5}
	fmt.Printf("list: %v\n", list)
	fmt.Printf("Filter: %v\n", Filter(list, func(v int64) bool {
		return v > 3
	}))

	fmt.Printf("Map: %v\n", Map(list, func(v int64) string {
		return fmt.Sprintf("str_%d", v)
	}))

	fmt.Printf("FlatMap: %v\n", FlatMap(list, func(v int64) []string {
		return Map([]string{"a", "b", "c"}, func(str string) string {
			return fmt.Sprintf("%s#%d", str, v)
		})
	}))

	fmt.Printf("Summing: %v\n", Summing(Map(list, func(v int64) int64 {
		return v * 10
	})))

	fmt.Printf("AnyMatch: %v\n", AnyMatch(list, func(v int64) bool {
		return v == 5
	}))

	fmt.Printf("AllMatch: %v\n", AllMatch(list, func(v int64) bool {
		return v > 0
	}))

}

func TestGroupBy(t *testing.T) {
	list := users

	fmt.Printf("GroupBy: %v\n", GroupBy(list, func(u User) string {
		return u.Gender
	}))

	fmt.Printf("GroupByMapping: %v\n", GroupByMapping(list, func(u User) string {
		return u.Gender
	}, func(u User) float64 {
		return u.Money
	}))

	fmt.Printf("GroupBySumming: %v\n", GroupBySumming(list, func(u User) string {
		return u.Gender
	}, func(u User) float64 {
		return u.Money
	}))

	fmt.Printf("GroupByReduce: %v\n", GroupByReduce(list, func(u User) string {
		return u.Gender
	}, func(u User) string {
		return u.Name
	}, func(u1, u2 string) string {
		if u1 == "" || u2 == "" {
			return u1 + u2
		}
		return u1 + " and " + u2
	}))

}
