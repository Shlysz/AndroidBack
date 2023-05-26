package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"src/androidBackground/entityClass"
	"src/androidBackground/respo"
	"strings"
	"sync"
)

func getRetMessage(question string, key string) string {

	client := &http.Client{}
	allquestions := makeMessages(key, question)
	var data = strings.NewReader(`{
    "model": "gpt-3.5-turbo",
    "messages": ` + allquestions + `
  }`)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer ${your key}")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)
	//调用chatgpt 3.5的api
	//allquestion := makeMessages(key, question)
	////协程插入新的问题
	//go respo.Message{Role: "user", Content: question}.AddMessage()
	//client := &http.Client{}
	//var data = strings.NewReader(`{
	// "model": "gpt-3.5-turbo",
	// "messages": ` + allquestion + `,
	//}`)
	//req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", data)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer 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")
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//bodyText, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	var retMessage entityClass.Response
	err = json.Unmarshal([]byte(bodyText), &retMessage)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	//提取出message
	ret := retMessage.Choices[0].Message.Content
	//插入数据库
	respo.Message{Role: "assistant", Content: ret}.AddMessage()
	return ret
}

// 构造messages
func makeMessages(key string, question string) string {
	var message respo.Message
	messages := message.GetMessage()
	//将messages转换为json “messages”:[{"role":"user","content":"hello"},{"role":"admin","content":"hi"}]
	message.Role = "user"
	message.Content = question
	messages = append(messages, message)
	messageJson, err := json.Marshal(messages)
	if err != nil {
		log.Fatal(err)
	}
	return string(messageJson)
}

func HandleChat(c *gin.Context) {
	//设置字符集为utf-8
	c.Header("Content-Type", "application/json; charset=utf-8")
	key := c.PostForm("key")
	question := c.PostForm("question")

	//根据key获取username
	username := respo.Account{Key: key}.GetUsernameByKey()
	if username == "" {
		c.JSON(200, gin.H{
			"message": "you have no right to do this",
			"status":  "false",
		})
		return
	}
	var wg sync.WaitGroup
	retMessage := ""
	wg.Add(1)
	go func() {
		retMessage = getRetMessage(question, key)
		wg.Done()
	}()
	wg.Wait()

	//返回结果
	c.JSON(200, gin.H{
		"message": retMessage,
		"status":  "true",
	})
}
