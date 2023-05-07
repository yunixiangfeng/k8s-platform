package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"k8s-platform/config"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

var Terminal terminal

type terminal struct{}

// wshanlder
func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	//加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		logger.Error("加载k8s配置失败, " + err.Error())
		return
	}
	//解析form入参，获取namespace，pod，container参数
	if err := r.ParseForm(); err != nil {
		logger.Error("解析参数失败, " + err.Error())
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	logger.Info("exec pod: %s, container: %s, namespace: %s\n", podName, containerName, namespace)

	//new一个terminalsession
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		logger.Error("实例化TerminalSession失败, " + err.Error())
		return
	}
	//处理关闭
	defer func() {
		logger.Info("关闭TerminalSession")
		pty.Close()
	}()
	//组装post请求
	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: containerName,
			Command:   []string{"/bin/bash"},
		}, scheme.ParameterCodec)
	logger.Info("exec post request url: ", req)

	//升级SPDY协议
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		logger.Error("建立SPDY连接失败, " + err.Error())
		return
	}
	//与kubelet建立stream连接
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})

	if err != nil {
		logger.Error("执行 pod 命令失败, " + err.Error())
		//将报错返回给web端
		pty.Write([]byte("执行 pod 命令失败, " + err.Error()))
		//标记关闭
		pty.Done()
	}
}

// 消息内容
type terminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

// 交互的结构体，接管输入和输出
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

// 初始化一个websocket.Upgrader类型的对象，用于http协议升级为ws协议
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

// 创建TerminalSession类型的对象并返回
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	//升级ws协议
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, errors.New("升级websocket失败," + err.Error())
	}
	//new
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}

	return session, nil
}

// 读数据的方法
// 返回值int是读成功了多少数据
func (t *TerminalSession) Read(p []byte) (int, error) {
	//从ws中读取消息
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("读取消息错误: %v", err)
		return 0, err
	}
	//反序列化
	var msg terminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("读取消息语法错误: %v", err)
		return 0, err
	}
	//逻辑判断
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("消息类型错误'%s'", msg.Operation)
		return 0, fmt.Errorf("消息类型错误'%s'", msg.Operation)
	}
}

// 写数据的方法,拿到apiserver的返回内容，向web端输出
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("写消息语法错误: %v", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("写消息错误: %v", err)
		return 0, err
	}

	return len(p), nil
}

// 标记关闭的方法
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// 关闭的方法
func (t *TerminalSession) Close() {
	t.wsConn.Close()
}

// resize方法，以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}
