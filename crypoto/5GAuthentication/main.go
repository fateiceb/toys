package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"log"
	"math/big"
	"math/rand"
	"time"
)

//Register Center
type RC struct {
	Mus map[string][]byte
}

func NewRc() *RC {
	return &RC{
		Mus: make(map[string][]byte),
	}
}
//公开的公共参数变量
var params *Params
//公共参数
type Params struct {
	G elliptic.Curve
	Gt elliptic.Curve
	E elliptic.Curve
	Ppub ecdsa.PublicKey
	Priv *ecdsa.PrivateKey
	H0   hash.Hash
	H1   hash.Hash
	H2   hash.Hash
	H3   hash.Hash
	H4   hash.Hash
	H5   hash.Hash
}
//打印params公共参数
func (pa *Params) String () string{
	return fmt.Sprintf("Curve:%+v\n Pub:%+v\n h0:%+v\n h2:%+vv\n h5%+v\n",
		pa.G.Params(),pa.Ppub,pa.H0,pa.H2,pa.H5)
}


func (rc *RC)InitParams() *Params{
	h0 := sha256.New()
	h1 := sha256.New()
	h2 := sha256.New()
	h3 := sha256.New()
	h4 := sha256.New()
	h5 := sha256.New()
	priv,err := ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		log.Fatal("生成私钥失败")
	}
	pub := priv.PublicKey;
	return &Params{
		G:    elliptic.P256(),
		Gt:   elliptic.P256(),
		E:    elliptic.P256(),
		Ppub: pub,
		Priv: priv,
		H0:   h0,
		H1:   h1,
		H2:   h2,
		H3:   h3,
		H4:   h4,
		H5:   h5,
	}
}

//Register Center根据mobile user的id生成公私钥并返回给MU
func (rc *RC) Register(mu *MU) bool {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		log.Fatal("RC---Register---", err)
	}
	mu.PK = &priv.PublicKey
	mu.Sk = priv
	privRsa,err := rsa.GenerateKey(crand.Reader,256)
	if err != nil {
		log.Fatal("RC---Register---",err)
	}
	mu.PkRsa = &privRsa.PublicKey
	mu.SKRsa = privRsa

	//select randomnum
	rand.Seed(time.Now().Unix())
	rn := rand.Int()
	params.H0.Write(mu.ID)
	params.H0.Write(int64ToByteArr(int64(rn)))
	log.Println("用户",string(mu.ID),"注册")
	functionLog("rc","registerMu 选择随机数","random number: ",rn)
	functionLog("rc","registerMu 生成Hu","Hu",params.H0.Sum(nil))
	functionLog("rc", "registerMu 生成私钥", fmt.Sprintf("%+v",mu.Sk))
	//重置h0
	params.H0.Reset()
	params.H0.Write(mu.ID)
	//添加可信用户到rc可新用户列表
	rc.Mus[string(mu.ID)] = params.H0.Sum(nil)
	functionLog("rc","rclist",fmt.Sprintf("%v",rc.Mus))
	//重置hash偏移位置
	params.H0.Reset()
	return true
}
func (rc *RC)RegisterMs(ms *MS) bool{
	priv,err:= rsa.GenerateKey(crand.Reader,256)
	if err != nil {
		log.Fatal("MS generate key priv  failed",err)
	}
	priv2,err := ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		log.Fatal("MS generate key priv  failed",err)

	}
	ms.Sk = priv2
	ms.PK = priv2.PublicKey
	ms.PKGRsa = &priv.PublicKey
	ms.SkGRsa = priv
	//h1 hash Id
	params.H1.Write(ms.ID)
	log.Println("MEC",string(ms.ID),"注册")
	functionLog("rc","registerMs Id的Hash值","Result of Ms's Id Hash value:",params.H1.Sum(nil))
	functionLog("rc","registerMs 生成私钥",fmt.Sprintf("%+v",priv2))
	//重置Hash偏移位置
	params.H0.Reset()
	return true
}

func (rc *RC) Compare(hashu []byte, mu *MU, ms *MS) []byte {
	if v,ok := rc.Mus[string(mu.ID)] ;ok {
		functionLog("rc","step3","存在用户",v)
		c,err :=rsa.EncryptPKCS1v15(crand.Reader,ms.PKGRsa,mu.ID)
		if err != nil {
			log.Fatal("rc compare,encrypt error")
		}
		functionLog("rc","step3","ciphertext",c)
		return c
	}
	functionLog("rc","step3","不存在用户")
	return nil
}

//mobile user
type MU struct {
	ID    []byte
	IDNum string
	PK    *ecdsa.PublicKey
	Sk    *ecdsa.PrivateKey
	PkRsa *rsa.PublicKey
	SKRsa *rsa.PrivateKey
}

func NewMu(id []byte, IDNum string) *MU {
	functionLog("NewMu", "新建用户", "id:", id)
	return &MU{
		ID:    id,
		IDNum: IDNum,
	}
}
//Step1
func (mu *MU)CommuciateWithMs(ms *MS) ([]byte, []byte, []byte, int64) {
	rand.Seed(time.Now().Unix())
	x,err:= crand.Int(crand.Reader,big.NewInt(int64(big.MaxPrec)))
	if err != nil {
		log.Fatal("随机数生成错误")
	}
	functionLog(string(mu.ID),"step1","random number x ",x)
	P := params.Ppub.Params().P
	X := P.Mul(P,x)
	gx := params.Ppub.Params().Gx
	functionLog(string(mu.ID),"step1","X",X)
	functionLog(string(mu.ID),"step1","Gx",gx)
	params.H1.Write(ms.ID)
	M := params.H1.Sum(nil)
	defer params.H1.Reset()
	functionLog(string(mu.ID),"step1","M",M)
	params.H2.Write(gx.Bytes())
	defer params.H2.Reset()
	hgx := params.H2.Sum(nil)
	params.H0.Write(mu.ID)
	hmuid := params.H0.Sum(nil)
	params.H0.Reset()
	N := []byte{}
	//xor
	for i,v := range hgx{
		N = append(N, v ^ hmuid[i])
	}
	functionLog(string(mu.ID),"step1","N",N)
	params.H0.Write(mu.ID)
	W := params.H0.Sum(nil)
	params.H0.Reset()
	functionLog(string(mu.ID),"step1","W",W)
	Tu := time.Now().Unix();
	//Unix时间戳
	functionLog(string(mu.ID),"step1","Tu",Tu)
	return M,N,W,Tu
}

func (mu *MU) CheckAndSetSessionKey(a []byte, t []byte, y *big.Int, tms int64) {
	params.H0.Write(mu.ID)
	if bytes.Equal(a,params.H0.Sum(nil)){
		functionLog(string(mu.ID),"step5","CheckAndSetSessionKey","step5","a is Ok")
	}else {
		log.Fatal("a error")
	}
	params.H0.Reset()
	params.H5.Write(a)
	params.H5.Write(y.Bytes())
	params.H5.Write(t)
	key := params.H5.Sum(nil)
	functionLog("ms","step6","SetsessionKey",key)

}



//MEC Server
type MS struct {
	ID   []byte
	PKGRsa *rsa.PublicKey
	SkGRsa *rsa.PrivateKey
	PK ecdsa.PublicKey
	Sk *ecdsa.PrivateKey
}

func (ms MS) Receive(m []byte, n []byte, w []byte, t int64, mu *MU) []byte {
	//check time
	timenow := time.Now().Unix()
	if timenow < t {
		log.Fatal("tu is not freshness")
	}
	functionLog("ms","step2","当前时间戳",timenow,)
	params.H0.Write(mu.ID)
	otherW := params.H0.Sum(nil)
	params.H0.Reset()
	if !bytes.Equal(w,otherW) {
		log.Fatal("verify w fail,terminate session",w,otherW)
	}
	functionLog("ms","step2","w",otherW,"验证相等")
	return otherW

}

func (ms MS) CommunicateWithmu(mu *MU) ([]byte, []byte, *big.Int, int64) {
	y, err := crand.Int(crand.Reader,big.NewInt(big.MaxExp))
	if err != nil {
		log.Fatal(err)
	}
	Tms := time.Now().Unix()
	Y := y.Mul(y,ms.PK.Params().P)

	params.H3.Write(Y.Bytes())
	t := params.H3.Sum(nil)
	params.H3.Reset()
	params.H1.Reset()
	params.H1.Write(mu.ID)
	a := params.H1.Sum(nil)
	params.H1.Reset()
	functionLog("ms","step4","time",Tms,"Y",Y)
	return a,t,Y,Tms
}

func NewMS(id []byte) *MS {
	priv,err:= rsa.GenerateKey(crand.Reader,256)
	if err != nil {
		log.Fatal("MS generate key priv  failed",err)
	}
	priv2,err := ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		log.Fatal("MS generate key priv  failed",err)

	}
	return &MS{
		id,
		&priv.PublicKey,
		priv,
		priv2.PublicKey,
		priv2,
	}
}





func main() {
	phaseLog("Initialization phase")
	rc := NewRc();
	params = rc.InitParams()
	log.Println(params)
	phaseLog("Registration phase")
	ms := NewMS([]byte("ms1"))
	alice :=NewMu([]byte("Alice"),"1")
	bob := NewMu([]byte("bob"),"2");
	//rc注册MEC server
	rc.RegisterMs(ms)
	rc.Register(alice)
	rc.Register(bob)
	//认证
	phaseLog("Authentication phase")
	M,N,W,T :=alice.CommuciateWithMs(ms)
	//time.Sleep(1 * time.Second)
	result := ms.Receive(M,N,W,T,alice)
	rc.Compare(result,alice,ms)
	a,t,y,Tms := ms.CommunicateWithmu(alice)
	alice.CheckAndSetSessionKey(a,t,y,Tms)
}

//Help Functions
//generateId 输入num，返回sha256生成的Id
func generateId(num string) []byte {
	hash := sha256.New()
	hash.Write([]byte(num))
	return hash.Sum(nil)
}

//格式化打印每阶段标题
func phaseLog(phase string) {
	log.Println("")
	log.Println("")
	log.Println("-----" + phase + "-----")
}

//格式化打印各对象及调用的函数结果
func functionLog(entity, function string, content ...interface{}) {
	log.Println(entity, ":", "***"+function+"***", content)
}
func int64ToByteArr(num int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(num))
	return buf
}
func byteArrToInt64(bytes []byte) int64 {
	num := binary.LittleEndian.Uint64(bytes)
	return int64(num)
}