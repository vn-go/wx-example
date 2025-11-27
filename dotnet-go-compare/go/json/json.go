package main

import (
	"fmt"
	"time"

	json "github.com/json-iterator/go"
)

type Address struct {
	City     string `json:"city"`
	District string `json:"district"`
}

type Person struct {
	Name        string   `json:"name"`
	Age         int      `json:"age"`
	IsDeveloper bool     `json:"isDeveloper"`
	Skills      []string `json:"skills"`
	Address     Address  `json:"address"`
	BigArray    []int    `json:"bigArray"`
}

func main() {
	fmt.Println("Đang tạo JSON lớn ~1.9MB...")

	// Tạo mảng 100_000 phần tử dạng [0,1,2,...,99999]
	bigArray := make([]int, 100_000)
	for i := range bigArray {
		bigArray[i] = i
	}
	bigArrayBytes, _ := json.Marshal(bigArray)

	// Template JSON
	template := `{
  "name": "Nguyễn Văn A",
  "age": 28,
  "isDeveloper": true,
  "skills": ["C#", "Go", "Kubernetes", ".NET 10"],
  "address": {
    "city": "Hà Nội",
    "district": "Cầu Giấy"
  },
  "bigArray": %s
}`

	jsonStr := fmt.Sprintf(template, string(bigArrayBytes))

	fmt.Printf("JSON size: %.1f KB\n", float64(len(jsonStr))/1024)

	// Warmup
	var p Person
	json.Unmarshal([]byte(jsonStr), &p)

	// Đo tốc độ thật sự: parse 100 lần
	start := time.Now()

	const loops = 100
	for i := 0; i < loops; i++ {
		var person Person
		err := json.Unmarshal([]byte(jsonStr), &person)
		if err != nil {
			panic(err)
		}
		_ = person // tránh bị optimize bỏ
	}

	elapsed := time.Since(start)
	avg := float64(elapsed.Milliseconds()) / loops

	fmt.Printf("Go parse %d lần: %d ms\n", loops, elapsed.Milliseconds())
	fmt.Printf("Trung bình mỗi lần: %.3f ms\n", avg)
	fmt.Printf("Name: %s, BigArray count: %d\n", p.Name, len(p.BigArray))
}
