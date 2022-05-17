package container

import (
	"docker_demo/common"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type ContainerInfo struct {
	Pid         string   `json:"pid"`     //容器init进程在宿主机上的pid
	Id          string   `json:"id"`      //容器ID
	Command     string   `json:"command"` //容器内init进程的运行命令
	Name        string   `json:"name"`
	CreateTime  string   `json:"create_time"`
	Status      string   `json:"status"`
	Volume      string   `json:"volume"`       //容器的数据卷
	PortMapping []string `json:"port_mapping"` //端口映射
}

func RecordContainerInfo(containerPID int, cmdArray []string, containerName, containerID string) error {
	// 容器信息封装
	info := &ContainerInfo{
		Pid:        strconv.Itoa(containerPID),
		Id:         containerID,
		Command:    strings.Join(cmdArray, ""),
		Name:       containerName,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     common.Running,
	}
	//创建容器保存文件夹
	dir := path.Join(common.DefaultContainerInfoPath, containerName)
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logrus.Errorf("mkdir container dir: %v, err: %v", dir, err)
			return err
		}
	}
	//创建容器信息文件
	fileName := fmt.Sprintf("%s/%s", dir, common.ContainerInfoFileName)
	logrus.Info("file name=======", fileName)
	file, err := os.Create(fileName)
	if err != nil {
		logrus.Errorf("create config.json, fileName: %s, err: %v", fileName, err)
		return err
	}
	// 保存容器信息到文件
	bs, _ := json.Marshal(info)
	_, err = file.WriteString(string(bs))
	if err != nil {
		logrus.Errorf("write config.json, fileName: %s, err: %v", fileName, err)
		return err
	}
	return nil
}

// ListContainerInfo 容器列表信息
func ListContainerInfo() {
	files, err := ioutil.ReadDir(common.DefaultContainerInfoPath)
	if err != nil {
		logrus.Errorf("read info dir, err: %v", err)
	}

	var infos []*ContainerInfo
	for _, file := range files {
		info, err := getContainerInfo(file.Name())
		if err != nil {
			logrus.Errorf("get container info, name: %s, err: %v", file, err)
			continue
		}
		infos = append(infos, info)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 2, ' ', 0)
	_, _ = fmt.Fprint(w, "ID\tName\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, info := range infos {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", info.Id, info.Name, info.Pid, info.Status, info.Command, info.CreateTime)
	}

	//刷新标准输出流缓存区，打印容器列表
	if err := w.Flush(); err != nil {
		logrus.Errorf("flush info, err: %v", err)
	}
}

// 通过容器名称读取容器信息
func getContainerInfo(containerName string) (*ContainerInfo, error) {
	filePath := path.Join(common.DefaultContainerInfoPath, containerName, common.ContainerInfoFileName)
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("read file, path: %s, err: %v", filePath, err)
		return nil, err
	}

	info := &ContainerInfo{}
	err = json.Unmarshal(bs, info)
	return info, err
}

func GenContainerId(n int) string {
	letterBytes := "0123456789qwertyuiopasdfghjklzxcvbnm"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func DeleteContainerInfo(containerName string) {
	dir := path.Join(common.DefaultContainerInfoPath, containerName)
	err := os.RemoveAll(dir)
	if err != nil {
		logrus.Errorf("remove container info, err: %v", err)
	}
}
