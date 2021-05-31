package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"time"
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
	privRsa,err := rsa.GenerateKey(crand.Reader,256)
	if err != nil {
		log.Fatal("RC---Register---",err)
	}
	mu.PkRsa = &privRsa.PublicKey
	mu.SKRsa = privRsa
	functionLog("RC", "register", mu.PK)
	log.Printf("%+v",mu.Sk)
	return true
}

//mobile user
type MU struct {
	ID    []byte
	IDNum string
	PK    *ecdsa.PublicKey
	Sk    *ecdsa.PrivateKey
	PkRsa *rsa.PublicKey
	SKRsa *rsa.PrivateKey
	GPK *ecdsa.PublicKey
	GSK *ecdsa.PrivateKey
	GPkRsa *rsa.PublicKey
	GSKRsa *rsa.PrivateKey
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
	PkRsa *rsa.PublicKey
	SkRsa *rsa.PrivateKey
}

func NewAs() *AS {
	priv,err := rsa.GenerateKey(crand.Reader,256)
	if err != nil {
		log.Fatal("As init err",err)
	}
	functionLog("AS","init")
	log.Printf("priv%+v",priv)
	return &AS{
		&priv.PublicKey,
		priv,
	}
}
func (as *AS) Authentication(mu *MU, ms *MS,group *Group) []byte {
	//设置随机数种子
	rand.Seed(time.Now().Unix())
	//生成随机数
	random := rand.Int()
	functionLog("AS","Authentication","随机数",random)
	//加密随机数
	encodeMsg,err := rsa.EncryptPKCS1v15(crand.Reader,mu.PkRsa,[]byte(strconv.Itoa(random)))
	//debug
	if err != nil {
		log.Fatal("加密失败",encodeMsg)
	}
	functionLog("AS","Authentication","密文",encodeMsg)
	//用户解密随机数
	decodeMsg,err := rsa.DecryptPKCS1v15(crand.Reader,mu.SKRsa,encodeMsg)
	result,err:= strconv.Atoi(string(decodeMsg))
	if err != nil {
		log.Fatal(err)
	}
	functionLog("MU","Authentication","解密",result)
	//用户使用AS公钥加密随机数
	encodeMuMsg,err := rsa.EncryptPKCS1v15(crand.Reader,as.PkRsa,[]byte(strconv.Itoa(random)))
	if err != nil {
		log.Fatal("MU EncryPt failed",err)
	}
	functionLog("MU","使用AS公钥加密随机数R","密文",encodeMuMsg)
	//AS解密并验证随机数数值r == original r
	decodeMuMsg,err := rsa.DecryptPKCS1v15(crand.Reader,as.SkRsa,encodeMuMsg)
	if err != nil {
		log.Fatal("Decrypt MU message failed",err)
	}
	r,err:= strconv.Atoi(string(decodeMuMsg))
	if err != nil {
		log.Fatal(err)
	}
	if r != random {
		log.Fatal("r != original r,认证失败")
	}
	functionLog("AS","解密并验证随机数数值r == original r,认证成功","解密后数值",r)
	//认证成功后为mu生成群组公钥和私钥
	Gpriv,err :=ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		log.Fatal("Authentication generate mu's group priv fail",err)
	}
	mu.GPK = &Gpriv.PublicKey
	mu.GSK = Gpriv
	GprivRsa,err := rsa.GenerateKey(crand.Reader,256)
	if err != nil {
		log.Fatal("Authentication generate mu's group priv fail",err)
	}
	mu.GPkRsa = &GprivRsa.PublicKey
	mu.GSKRsa = GprivRsa
	//生成ticket并加密返回
	functionLog("Mu","IdNum"+mu.IDNum,"GetTicket")
	ticket := map[string]interface{}{
		"Pk of Mobile user in Group":mu.GPK,
		"Sk of Mobile user in Group":mu.GSK,
		"Pk of mobile user":mu.PK,
		"Pk of Group Manager":group.GroupManager.GPK,
		"Pk of Group":group.GroupPk,

	}
	for key,v := range ticket {
		log.Printf("%s:%+v",key,v)
	}
	jsonticket,err:= json.Marshal(ticket)
	if err != nil {
		log.Fatal("encode to json failed")
	}
	//将mu加入群组
	group.addMu(mu)
	//debug
	//log.Println(len(group.MUS))
	return jsonticket
}

//MEC Server
type MS struct {
	ID   []byte
	PKGRsa *rsa.PublicKey
	SkGRsa *rsa.PrivateKey
}

func NewMS(id []byte) *MS {
	priv,err:= rsa.GenerateKey(crand.Reader,256)
	if err != nil {
		log.Fatal("MS generate key priv  failed",err)
	}
	return &MS{
		id,
		&priv.PublicKey,
		priv,
	}
}
//MS生成区块
func (ms *MS)GnerateBlock(tx []byte) *Block{
	return generateBlock(tx)
}
//Group
type Group struct {
	id []byte
	//群管理
	GroupManager MU
	//MS
	GroupMS MS
	//用户组
	MUS []*MU
	GroupPk *ecdsa.PublicKey
	GroupSK *ecdsa.PrivateKey
}
func (g *Group)AggregatePk() *ecdsa.PublicKey{
	if g.GroupPk != nil {
		return g.GroupPk
	}
	log.Fatal("AggregatePk failed")
	return nil
}
func NewGroup(GroupManagerId []byte,GroupManagerNum string) *Group{
	priv,err:=ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		log.Fatal("group generate priv failed",err)
	}
	return &Group{
		GroupManager:*newGroupManager(GroupManagerId,GroupManagerNum),
		MUS: []*MU{
		},
		GroupPk: &priv.PublicKey,
		GroupSK:priv,

	}
}
func (g *Group) addMu(mobileUser *MU)bool{
	g.MUS = append(g.MUS,mobileUser)
	return true
}
//Sign 组签名返回签名结果
func (g *Group)Sign(b *Block) (r,s *big.Int){
	if b == nil {
		log.Fatal("block not generated")
	}
	r,s,err := ecdsa.Sign(crand.Reader,g.GroupSK,b.BlockHash)
	if err != nil {
		log.Fatal("group",g.id,err)
	}
	return r,s
}
func (g *Group)Validation(SignGroup *Group,block *Block,r,s*big.Int) bool{
	return ecdsa.Verify(SignGroup.GroupPk,block.BlockHash,r,s)
}
func newGroupManager(id []byte,idNum string) *MU{
	manager := NewMu(id,idNum)
	priv,err := ecdsa.GenerateKey(elliptic.P256(),crand.Reader)
	if err != nil {
		log.Fatal("init Group manager failed",err)
	}
	manager.GPK = &priv.PublicKey
	manager.GSK = priv
	privRsa, err := rsa.GenerateKey(crand.Reader,256)
	if err != nil {
		log.Fatal("init Group manager failed",err)

	}
	manager.GSKRsa = privRsa
	manager.GPkRsa = &privRsa.PublicKey
	return manager
}
//区块结构
type Block struct {
	Tx []byte
	BlockHash []byte
}
func generateBlock(tx []byte) *Block {
	hash := sha256.New()
	hash.Write(tx)
	return &Block{
		tx,hash.Sum(nil),
	}
}

func main() {
	phaseLog("初始化MS、AS、RC")
	//GroupMnager 也是Mu类型，GroupMnaagerNum是mu类型的IdNum字段,生成的Group需要注册Groupmanager
	group1 := NewGroup([]byte("groupmanager1"),"1")
	group2 := NewGroup([]byte("groupmanager2"),"2")
	group3 := NewGroup([]byte("groupmanager3"),"3")
	rc := NewRc()
	//生成ms
	ms1 := NewMS([]byte("ms1"))
	ms2 := NewMS([]byte("ms2"))
	as := NewAs()
	phaseLog("新建用户")
	//群组一用户
	alice := NewMu(generateId("1"), "1")
	bob := NewMu(generateId("2"),"2")
	//群组二用户
	eve := NewMu(generateId("3"),"3")
	alex := NewMu(generateId("4"),"4")
	phaseLog("用户经由RC生成公私钥")
	//用户在register center注册
	rc.Register(alice)
	rc.Register(bob)
	rc.Register(eve)
	rc.Register(alex)
	//注册group管理员
	phaseLog("group manager register")
	rc.Register(&group1.GroupManager)
	rc.Register(&group2.GroupManager)
	rc.Register(&group3.GroupManager)
	phaseLog("认证过程")
	//认证和交流ticket过程
	as.Authentication(alice, ms1,group1)
	as.Authentication(bob,ms1,group1)
	as.Authentication(eve,ms2,group2)
	as.Authentication(alex,ms2,group2)
	//聚合公钥
	phaseLog("公钥聚合")
	Gpk1 := group1.AggregatePk()
	Gpk2 := group2.AggregatePk()
	log.Println("AggregatePk:","pk1",Gpk1,"\n","pk2",Gpk2)
	phaseLog("区块生成和签名")
	//ms1生成block1
	block1 :=ms1.GnerateBlock([]byte("transaction1"))
	//group1对block1签名
	r,s := group1.Sign(block1)
	log.Println("区块hash",block1.BlockHash)
	log.Println("签名后的区块hash",r)
	phaseLog("广播并进行第二轮有效性验证")
	//group2接受到广播并验证有效性，如果要展示失效效果，可以生成新的未签名block进行展示
	testBlock := ms1.GnerateBlock([]byte("transction3"))
	validation2 :=group2.Validation(group1,block1,r,s)
	validation3 :=group3.Validation(group1,block1,r,s)
	validationfail := group2.Validation(group1,testBlock,r,s);
	log.Println("group2验证结果:",validation2)
	log.Println("group3验证结果：",validation3)
	log.Println("group2验证错误blockhash:",validationfail)
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
