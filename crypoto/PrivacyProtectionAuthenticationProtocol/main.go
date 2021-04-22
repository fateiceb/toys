package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"log"
	"math/big"
	"time"
)

type U interface {
	// SendMessage U发送服务请求信息M1给网关节点GWN，请求认证
	SendMessage(message string,gwn GWN) bool
	// Verify 用户U通过验证信息M4的合法性，进而验证S与GWN的合法性
	Verify(t int64,hashval []byte) bool
}
type GWN interface {
	// Init 初始化，产生公共参数
	Init() Params
	// SendMessageToSensor GWN完成U认证之后，GWN发送认证信息M2给传感器节点S
	SendMessageToSensor(message string,sensor *Sensor) (int64,[]byte)
	SendMessageToUser(message string, user U)
	// Verify 此时GWN收到M3后，验证其合法性，并且发送信息M4给用户U
	VerifyUser(message string,t int64,hashval []byte,randnum *big.Int,user User)
	VerifySensor(message string,t int64,hashval []byte,randnum *big.Int,sensor Sensor)

}
type S interface {
	// 当S收到GWN发来的信息M2之后，便发送有关于会话密钥建立的信息M3给GWN
	BuildDiscussionKey(message string, gwn GWN)
	sendMessage(message string, gwn GWN) bool
}
//全局公共参数
var publicParams Params
//定义公共参数需要的属性，curve代表椭圆曲线
type Params struct {
	Curve elliptic.Curve `curve`
	H	hash.Hash
	H1 	hash.Hash
	H2 	hash.Hash
	PublicKey ecdsa.PublicKey
}

func (p Params) String() string {
	return fmt.Sprintf("params\nCurve:%+v\npublicKey%+v\nHash1%+v\n,",p.Curve.Params(),p.PublicKey,p.H1)
}
//Gateway
type Gateway struct {
	Priv ecdsa.PrivateKey
}



func (g *Gateway) SendMessageToSensor(message string, sensor *Sensor) (int64,[]byte){
	//验证t2时间戳
	t2 := time.Now().Unix()
	fmt.Println("sensor:时间戳t2",t2)
	//时间间隔
	time.Sleep(1 * time.Second)
	//网关发送t3时间戳消息
	hash := publicParams.H1
	t3 := time.Now().Unix()
	hash.Reset()
	hash.Write([]byte(sensor.SensorId))
	hashval := hash.Sum(nil)
	random, _ := crand.Int(crand.Reader,big.NewInt(10000000))
	//网关验证传感器消息
	g.VerifySensor("sensor",t3,hashval,random,*sensor)
	fmt.Println("gatway: send to sensor t3时间戳",t3)
	time.Sleep(1 * time.Second)
	t4 := time.Now().Unix()
	hash.Reset()
	hash.Write(int64ToByteArr(t4))
	val := hash.Sum(nil)
	return t4,val
}

func (g *Gateway) SendMessageToUser(message string, user U) {
	hash := publicParams.H
	hash.Write([]byte(message))
	mu := hash.Sum(nil)
	fmt.Println("Mu",mu)
}

func (g *Gateway) VerifyUser(message string,t int64,hashval []byte,randnum *big.Int,user User) {
	hash := publicParams.H
		t11 := time.Now().Unix()
		hash.Reset()
		hash.Write([]byte(user.Username))
		verfiyhash := hash.Sum(nil)
		if t>t11 {
			log.Fatal(t,t11,"gateway：时间错误")
		}
		if bytes.Compare(hashval,verfiyhash) != 0 {
			log.Fatal("gateway：验证错误")
		}
		fmt.Println("gateway:用户验证成功")
}
func (g *Gateway) VerifySensor(message string,t int64,hashval []byte,randnum *big.Int,sensor Sensor) {
	hash := publicParams.H1
	hash.Reset()
	hash.Write([]byte(sensor.SensorId))
	verifyhash := hash.Sum(nil)
	if bytes.Compare(hashval,verifyhash) != 0{
		fmt.Println("gateway:传感器请求错误")
	}
	fmt.Println("gateway:时间戳t3",t)
	fmt.Println("gateway:传感器请求验证成功")
}

func NewGateWay() *Gateway {
	return &Gateway{}
}
func (g *Gateway) Init() Params {
	curve :=elliptic.P256()
	h := sha256.New()
	h1 := sha256.New()
	h2 := sha256.New()
	privKey,err := ecdsa.GenerateKey(curve,crand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	return Params{curve,h,h1,h2,privKey.PublicKey}
}


// User
type User struct {
	Username string
	Priv *ecdsa.PrivateKey
}

func (u *User) SendMessage(message string, gwn GWN) bool {
	switch message {
	case "register":
		//使用公共参数中的hash函数
		h := publicParams.H
		h.Write([]byte(u.Username))
		midu := h.Sum(nil)
		fmt.Println("MIDU",midu)
		gwn.SendMessageToUser(string(midu),u)
		if len(midu) == 0 {
			return  false
		}
		return true
	case "gateway":
		h := publicParams.H
		rand,err := crand.Int(crand.Reader,big.NewInt(1000000000000))
		if err != nil {
			log.Fatal(err)
		}
		t1 :=time.Now().Unix()
		fmt.Println("user随机数",rand,"时间戳t1",t1)
		h.Reset()
		h.Write([]byte(u.Username))
		hashval := h.Sum(nil)
		//网关验证
		gwn.VerifyUser("user",t1,hashval,rand,*u)
		return  true
	}
	return false
}

func (u *User) Verify(t int64,hashval []byte) bool {
	h := publicParams.H1
	h.Reset()
	time.Sleep(1 * time.Second)
	t4 := time.Now().Unix()
	fmt.Println("user:t4时间戳",t,"t4'时间戳",t4)
	if t4 < t {
		fmt.Println("user:时间验证错误")
		return false
	}
	h.Write(int64ToByteArr(t))
	verifyval := h.Sum(nil)
	if bytes.Compare(hashval,verifyval) != 0 {
		fmt.Println("user:hash值错误")
		return false
	}
	fmt.Println("user:验证成功计算会话密钥")
	h2 := publicParams.H2
	h2.Reset()
	h2.Write(u.Priv.D.Bytes())
	sk := h2.Sum(nil)
	fmt.Println("会话密钥",sk)
	return true
}


func NewUser(name string) *User  {
	priv,err := ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		 log.Fatal(err)
	}
	return &User{Username: name,Priv: priv}
}

//Sensor
type Sensor struct {
	SensorId string
	priv *ecdsa.PrivateKey
}
func NewSensor(Id string) *Sensor{
	priv,err := ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		log.Fatal(priv)
	}
	return &Sensor{SensorId: Id,priv: priv}
}
func (s *Sensor) BuildDiscussionKey(message string, gwn GWN) {
	panic("implement me")
}

func (s *Sensor) sendMessage(message string, gwn GWN) bool {
	switch message {
	case "rigister":
		hash := publicParams.H
		hash.Write([]byte(s.SensorId))
		me := hash.Sum(nil)
		if len(me) == 0 {
			return false
		}
		fmt.Println("传感器注册信息：",me)
		//返回信息给网关
		gwn.SendMessageToSensor(string(me),s)
		return true
	case "gateway":
		//fmt.Println("sensor to gateway")
	}
	return false
}

//辅助函数
func int64ToByteArr(num int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(num))
	return buf
}
func byteArrToInt64(bytes []byte) int64 {
	num := binary.LittleEndian.Uint64(bytes)
	return int64(num)
}

func main() {
	//占位符
	indent := "-------------"
	//新建gateWay
	gwn := NewGateWay()
	//新建User
	user := NewUser("user1")
	//新建SenSor
	sensor := NewSensor("s1")
	//gatWay公布公共参数
	publicParams = gwn.Init()
	fmt.Println(indent,"公共参数",indent)
	fmt.Printf("%+v",publicParams)
	//发送注册信息
	fmt.Println(indent,"用户注册",indent)
	user.SendMessage("register", gwn)
	fmt.Println(indent,"传感器注册",indent)
	sensor.sendMessage("rigister",gwn)

	//认证与密钥协商
	//用户选择随机数，发送给网关节点
	fmt.Println(indent,"认证过程",indent)
	//时间间隔
	time.Sleep(1*time.Second)
	//用户发送请求，网关验证
	succes := user.SendMessage("gateway",gwn)
	if !succes{
		log.Fatal("用户请求错误")
	}
	//网关发送信息到传感器，并验证传感器回传信息
	t4,val := gwn.SendMessageToSensor("sensor",sensor)
	//用户验证传感器参数
	user.Verify(t4,val)

}
