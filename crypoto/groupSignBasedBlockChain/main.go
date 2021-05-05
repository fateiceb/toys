package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"log"
)

//Register Center
type RC struct {
}

func NewRc() *RC {
	return &RC{}
}

//Register Center根据mobile user的id生成公私钥并返回给MU
func (rc *RC) Register(mu *MU) bool {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		log.Fatal("RC---Register---", err)
	}
	mu.PK = &priv.PublicKey
	mu.Sk = priv
	functionLog("RC", "register", mu.PK)
	return true
}

//mobile user
type MU struct {
	ID    []byte
	IDNum string
	PK    *ecdsa.PublicKey
	Sk    *ecdsa.PrivateKey
}

func NewMu(id []byte, IDNum string) *MU {
	functionLog("NewMu", "新建用户", "id:", id)
	return &MU{
		ID:    id,
		IDNum: IDNum,
	}
}

//Authentication Server
type AS struct {
}

func NewAs() *AS {
	return &AS{}
}
func (as *AS) Authentication(mu *MU, MSID []byte) []byte {
	return []byte{'1', '2', '3'}
}

//MEC Server
type MS struct {
}

func NewMS() *MS {
	return &MS{}
}

func main() {
	phaseLog("初始化MS、AS、RC")
	rc := NewRc()
	ms := NewMS()
	as := NewAs()
	phaseLog("新建用户")
	alice := NewMu(generateId("1"), "1")
	phaseLog("用户经由RC生成公私钥")
	rc.Register(alice)
	phaseLog("认证过程")
	MSID := []byte("1")
	ticket := as.Authentication(alice, MSID)

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
	log.Println("-----" + phase + "-----")
}

//格式化打印各对象及调用的函数结果
func functionLog(entity, function string, content ...interface{}) {
	log.Println(entity, ":", "***"+function+"***", content)
}
