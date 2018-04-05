package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

const (
	leetCodeJSON = "leetcode.json"
)

type leetcode struct {
	Username string

	Record   record   // 已解答题目与全部题目的数量，按照难度统计
	Problems problems // 所有问题的集合

	Ranking int

	Updated time.Time
}

func newLeetCode() *leetcode {
	log.Println("开始，获取 LeetCode 数据")

	lc, err := readLeetCode()
	if err != nil {
		log.Println("读取 LeetCode 的记录失败，正在重新生成 LeetCode 记录。失败原因：", err.Error())
		return getLeetCode()
	}

	lc.refresh()

	log.Println("完成，获取 LeetCode 数据")
	return lc
}

func readLeetCode() (*leetcode, error) {
	data, err := ioutil.ReadFile(leetCodeJSON)
	if err != nil {
		return nil, errors.New("读取文件失败：" + err.Error())
	}

	lc := new(leetcode)
	if err := json.Unmarshal(data, lc); err != nil {
		return nil, errors.New("转换成 leetcode 时，失败：" + err.Error())
	}

	return lc, nil
}

func (lc *leetcode) save() {
	raw, err := json.MarshalIndent(lc, "", "\t")
	if err != nil {
		log.Fatal("无法把Leetcode数据转换成[]bytes: ", err)
	}
	if err = ioutil.WriteFile(leetCodeJSON, raw, 0666); err != nil {
		log.Fatal("无法把Marshal后的lc保存到文件: ", err)
	}
	log.Println("最新的 LeetCode 记录已经保存。")
	return
}

func (lc *leetcode) refresh() {
	log.Println("开始，刷新 LeetCode 数据")

	if time.Since(lc.Updated) < 7*time.Minute {
		log.Printf("LeetCode 数据在 %s 前刚刚更新过，跳过此次刷新\n", time.Since(lc.Updated))
		return
	}

	newLC := getLeetCode()
	logDiff(lc, newLC)
	lc = newLC

	lc.save()
}

func logDiff(old, new *leetcode) {
	// 对比 ranking
	str := fmt.Sprintf("当前排名 %d", new.Ranking)
	verb, delta := "进步", old.Ranking-new.Ranking
	if new.Ranking > old.Ranking {
		verb, delta = "后退", new.Ranking-old.Ranking
	}
	str += fmt.Sprintf("，%s了 %d 名", verb, delta)
	log.Println(str)

	// 对比 已完成的问题
	lenOld, lenNew := len(old.Problems), len(new.Problems)
	hasNewFinished := false

	i := 0
	for i < lenOld {
		o, n := old.Problems[i], new.Problems[i]

		if o.IsAccepted == false && n.IsAccepted == true {
			log.Printf("～新完成～ %d.%s", n.ID, n.Title)
			dida("re", n)
			hasNewFinished = true
		}

		i++
	}

	if !hasNewFinished {
		log.Println("～ 没有新完成习题 ～")
	}

	// 检查新添加的习题
	for i < lenNew {
		if new.Problems[i].isAvailble() {
			log.Printf("新题: %d - %s", new.Problems[i].ID, new.Problems[i].Title)
			dida("do", new.Problems[i])
			i++
		}
	}
}
