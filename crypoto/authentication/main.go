package main

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"log"
	"math/rand"
	"time"
)

// Params 公共参数结构体
type Params struct {
	Curve elliptic.Curve
	H hash.Hash
	H1 hash.Hash
	H2 hash.Hash
	MAC hash.Hash
	PublicKey crypto.PublicKey
}
// params 公共参数结
var params *Params
type TTP struct {
	//G1,G2,P.q
	Curve elliptic.Curve
	H hash.Hash
	H1 hash.Hash
	H2 hash.Hash
	MAC hash.Hash
	//密钥对
	Priv *ecdsa.PrivateKey
}
func NewTTP() *TTP {
	curve := elliptic.P256()
	priv,err := ecdsa.GenerateKey(curve,crand.Reader)
	if err != nil {
		err.Error()
	}
	return &TTP{Curve: curve,H: sha256.New(),H1:sha256.New(),H2: sha256.New(),MAC: sha256.New(),Priv: priv}
}

func (ttp *TTP) Init() *Params  {
	return  &Params{
		Curve: ttp.Curve,
		H: ttp.H,
		H1:ttp.H1,
		H2: ttp.H2,
		MAC: ttp.MAC,
		PublicKey: ttp.Priv.PublicKey,
	}
}
func (ttp *TTP)RegisterUser(user *User, hashVal []byte) bool{
	user.IsRegister = true
	priv,err := ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		fmt.Println("RegisterUser:",err)
		return false
	}
	user.Priv = priv
	return true
}
//格式化输出TTP包含信息
func (ttp *TTP) String() string  {
	return fmt.Sprintf("******生成TTP******\nCurve:%+v\nPriv:%+v\nHash0:%+v",ttp.Curve.Params(),ttp.Priv,ttp.H)

}


type Sensor struct {
	Id []byte
	SensorPirv  *ecdsa.PrivateKey
	UserPublicKey ecdsa.PublicKey
	UserRequestRandom int
	UserReplyTime int64
	UserReplyH1val []byte
}
func NewSensor(Id []byte) *Sensor {
	priv,err:= ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		 log.Fatal("new sensor fail:",err)
	}
	return &Sensor{
		Id: Id,
		SensorPirv: priv,
		UserPublicKey: ecdsa.PublicKey{},
	}
}
func (s Sensor) receiveCheckAndReply(u *User)  {
	if s.UserRequestRandom != 0{
		log.Println("sensor:receive request from ",string(u.UserId))
	}
	mac := params.MAC
	rand.Seed(80)
	r := rand.Int()
	t := time.Now().Unix()
	log.Println("sensor:时间戳 ",t)
	if mac == nil {
		println(111)
	}
	mac.Write(int64ToByteArr(t))
	mact := mac.Sum(nil)
	if u != nil {
		u.Mact = mact
		u.RandomNumFromSensor = r
		u.Time = t
		log.Println("sensor:reply",string(u.UserId),"message:",mact)
	}else {
		log.Fatal("user is not alive")
	}
}

func (s Sensor) CheckAndReplyPhaseTwo(user *User) {
	h1 := params.H1
	h1.Reset()
	//checkmessage is receive
	if s.UserReplyTime == 0 {
		log.Println("not receive message from",string(user.UserId))
		return
	}
	//checktime
	if t3 := time.Now().Unix();s.UserReplyTime > t3 {
		log.Println("time stamp is error:user timestamp:",s.UserReplyTime,"now:",t3)
		return
	}
	log.Println("sensor:user timestamp:",s.UserReplyTime,"now:",time.Now().Unix())
	//checkmessage content
	h1.Write(int64ToByteArr(s.UserReplyTime))
	if message := h1.Sum(nil);bytes.Compare(s.UserReplyH1val,message) != 0 {
		log.Println("message check error")
		return
	}
	log.Println("timestamp:ok",s.UserReplyTime,"message:ok",h1.Sum(nil))
	h1.Reset()
	h1.Write(s.Id)
	h1.Write(int64ToByteArr(s.UserReplyTime))
	sk := h1.Sum(nil)
	log.Println("sensor sk",sk)
}
type User struct {
	UserId []byte
	UserPassword []byte
	Priv *ecdsa.PrivateKey
	IsRegister bool
	Time int64
	Mact []byte
	RandomNumFromSensor int
}

// NewUser 新建未注册的User，密钥为空，未注册
func NewUser(id string ,password string) *User {
	return &User{
		UserId: []byte(id),
		UserPassword: []byte(password),
		Priv: nil,
		IsRegister: false,
	}
}
func (u *User) HashIdandPassword()  []byte{
	hash := sha256.New()
	hash.Write(u.UserId)
	hash.Write(u.UserPassword)
	return hash.Sum(nil)
}

// Login 用户U插入智能卡，输入身份ID和密码PW登录系统。智能卡根据用户
func (u *User) Login(s string) {
	hash := sha256.New()
	log.Println(s," Verify:")
	//e
	hval := u.HashIdandPassword()
	//e'
	hash.Write(u.UserId)
	hash.Write(u.UserPassword)
	smartval := hash.Sum(nil)
	if bytes.Compare(hval,smartval) == 0 {
		log.Println(s,"验证成功")
	}else {
		log.Println("验证失败")
	}
}

// Request 将请求消息 通过一个公共信道发送给SN
func (u *User) Request(sensor *Sensor) {
	//产生随机数
	rand.Seed(23)
	r := rand.Int()
	log.Println("request random number",r)
	if sensor != nil {
		sensor.UserRequestRandom = r
		sensor.UserPublicKey = u.Priv.PublicKey
	}else {
		log.Fatal("no sensor")
	}

}

func (u *User) CheckAndReply(sensor *Sensor) {
	mac := params.MAC
	mac.Reset()
	h1 := params.H1
	h1.Reset()
	//check
	if u.RandomNumFromSensor == 0 {
		log.Println("not receive message")
		return
	}
	if t := time.Now().Unix();u.Time > t {
		log.Println("time stamp is wrong\nsensor timestamp:",u.Time,"now:",t)
	}
	mac.Write(int64ToByteArr(u.Time))
	if mact := mac.Sum(nil); bytes.Compare(mact,u.Mact) != 0 {
		log.Println("check mac message wrong")
	}
	t2 := time.Now().Unix()
	log.Println("sensor timestamp:",u.Time,"now:",t2)
	h1.Write(int64ToByteArr(t2))
	//disscussion key
	h1hashval := h1.Sum(nil)
	//Reply
	sensor.UserReplyTime = t2
	sensor.UserReplyH1val = h1hashval
	h1.Reset()
	h1.Write(sensor.Id)
	h1.Write(int64ToByteArr(t2))
	sk := h1.Sum(nil)
	log.Println("user sk",sk)

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

func main()  {
	//生成TTP
	ttp := NewTTP()
	sensor := NewSensor([]byte("sensor1"))
	log.Printf("%+v\n",ttp)
	//ttp初始化公共参数
	params =ttp.Init()
	log.Printf("******初始化公共params******\n%+v\n",params)
	//新建未注册用户
	log.Println("******User Register******")
	alice := NewUser("alice","alice999")
	aliceHashVal := alice.HashIdandPassword()
	//ttp注册用户
	ttp.RegisterUser(alice,aliceHashVal)
	log.Println("alice ID and passwod hashvalue",aliceHashVal)
	//用户登录
	log.Println("******User Login******")
	alice.Login("smartcard")
	//用户发送请求
	log.Println("******Request******")
	alice.Request(sensor)
	//sensor检查请求并返回消息
	log.Println("******sensor******")
	//时间间隔
	time.Sleep(1 * time.Second)
	sensor.receiveCheckAndReply(alice)
	log.Println("******alice*******")
	time.Sleep(1 * time.Second)
	alice.CheckAndReply(sensor)
	time.Sleep(1 * time.Second)
	log.Println("******sensor phase two*******")
	sensor.CheckAndReplyPhaseTwo(alice)

}


