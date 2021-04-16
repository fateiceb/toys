package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	easyrand "math/rand"
	"time"
)

type SP struct {
	Record   map[string]Machine
	Privkeys map[string][]byte
}

func NewSP() *SP {
	return &SP{
		Record:   make(map[string]Machine, 1000),
		Privkeys: make(map[string][]byte, 1000),
	}
}

func (sp *SP) addRecord(machine Machine) error {
	sp.Record[machine.ID] = machine
	return nil
}
func (sp *SP) addPrivkey(gid string) {
	v := make([]byte, 32)
	_, err := rand.Read(v)
	if err != nil {
		log.Printf("Register Machine:error generate privkey %v", err)
	}

	sp.Privkeys[gid] = v
	fmt.Println("Sp add new key for gateway ", gid)
}
func (sp *SP) RegisterMachine(id string, gid string) *Machine {
	v, ok := sp.Privkeys[gid]
	if !ok {
		sp.addPrivkey(gid)
	}
	v = sp.Privkeys[gid]

	machine := &Machine{
		ID:      id,
		GID:     gid,
		PrivKey: v,
	}
	sp.addRecord(*machine)
	return machine
}

type Machine struct {
	ID      string
	GID     string
	PrivKey []byte
	buffer  *bytes.Buffer //信息存放
}

func (ma *Machine) ParseMessageTime() int64 {
	content := ma.buffer.Bytes()
	if len(content) < 42 {
		log.Printf("ParseMessage error MessageTime")
	}
	t1 := byteArrToInt64(content[33:41])
	return t1
}
func (ma *Machine) ParseMessageRandomA() int64 {
	content := ma.buffer.Bytes()
	if len(content) < 43 {
		log.Printf("ParseMessage error randomA")
	}
	randoma := content[41]
	return int64(randoma)
}
func (ma *Machine) ParseMessageC() []byte {
	content := ma.buffer.Bytes()
	if len(content) < 43 {
		log.Println("Parsemessage error C")
	}

	return content[1:33]
}

// Commucate 发送消息
/*
	A生成一个一次性随机数a，并计算Q_AB=h(K_HG∥id_A∥id_B )
	和C=MAC[Q_AB,Gid∥id_A∥T_1∥a]，
	然后发送请求消息{id_A,C,T_1,a,id_B }到B（这里的T1表示当前设备A的时间戳）。
*/
func (ma *Machine) Commucate(target *Machine) {
	randnuma, err := rand.Int(rand.Reader, big.NewInt(100))
	// fmt.Println("func() Commucate:randumaa", randnuma)
	t := time.Now().Unix()
	if err != nil {
		log.Fatal(err)
	}
	hash := sha256.New()
	hash.Write(ma.PrivKey)
	hash.Write([]byte(ma.ID))
	hash.Write([]byte(target.ID))
	Qab := hash.Sum(nil)
	// log.Println(Qab)
	hash.Write(Qab)
	hash.Write([]byte(ma.GID))
	hash.Write([]byte(ma.ID))
	hash.Write(int64ToByteArr(t))
	hash.Write(randnuma.Bytes())
	// log.Println("commucate:t:", int64ToByteArr(t))
	// log.Println("comucate:randnuma:", randnuma.Bytes())
	C := hash.Sum(nil)
	//然后发送请求消息{id_A,C,T_1,a,id_B }
	// length := len(ma.ID) + len(C) + len(int64ToByteArr(t)) + len(randnuma.Bytes()) + len(target.ID)
	// fmt.Println(length)
	buf := make([]byte, 0)
	// fmt.Println([]byte(ma.ID)[:])
	buf = append(buf, []byte(ma.ID)...)
	buf = append(buf, C...)
	buf = append(buf, int64ToByteArr(t)...)
	buf = append(buf, randnuma.Bytes()...)
	// fmt.Println("functionCommucate:lenrand", len(randnuma.Bytes()))
	buf = append(buf, []byte(target.ID)...)
	// fmt.Println("functionCommucate:targetid", len([]byte(target.ID)))
	// fmt.Println(buf, "-------", len(buf))
	target.buffer = bytes.NewBuffer(buf)
	log.Println("Commucate time1stamp:", t)

}

// CheckMessage 接收到信息后检查步骤
/*
	一旦B接收到请求消息，首先检验(T_2-T_1 )≤T是否成立，如果成立则继续下一步
	：计算Q_AB=h(K_HG∥id_A∥id_B )和C=MAC[Q_AB,Gid∥id_A∥T_1∥a]，
	然后验证是否C=C^*，如果相等那么设备B就生成一个随机的秘密b并计算N_B=E_(K_HG )
	 [id_B,b,a,T_2 ]和tag=HMAC[Q_AB,id_A∥id_B∥b∥a∥T_2 ]，
	最后发送回应消息{N_B,tag,T_2 }给A（这里的T2表示当前设备B的时间戳）。
	若在此期间出现验证失败则设备B生成一个错误消息并终止通信。
*/
func (ma *Machine) CheckMessage(checkma *Machine) bool {
	t2 := time.Now().Unix()
	if ma.buffer.Len() < 0 {
		log.Println(ma.ID, "：未接收到消息缓冲区为空")
		return false
	}
	if ma.ParseMessageTime() >= t2 {
		log.Println(ma.ID, ":时间戳验证错误")
		log.Println("matime:", ma.ParseMessageTime())
		log.Println("t2", t2)
		return false
	}
	hash := sha256.New()
	hash.Write(checkma.PrivKey)
	hash.Write([]byte(checkma.ID))
	hash.Write([]byte(ma.ID))
	Qab := hash.Sum(nil)
	hash.Write(Qab)
	// log.Println(Qab)
	hash.Write([]byte(checkma.GID))
	hash.Write([]byte(checkma.ID))
	hash.Write(int64ToByteArr(ma.ParseMessageTime()))
	// log.Println("checkMessage:t", int64ToByteArr(ma.ParseMessageTime()))
	hash.Write(int64ToByteArr(ma.ParseMessageRandomA())[0:1])
	// log.Println("checkMessage:randomA", int64ToByteArr(ma.ParseMessageRandomA())[0:1])
	C := hash.Sum(nil)
	//DEBUG
	// log.Println("ccccc", C)
	// log.Println("c parse", ma.ParseMessageC())
	//从缓存中获取C
	cparse := ma.ParseMessageC()
	//校验C是否与缓存中C相同
	for i := 0; i < len(C); i++ {
		if cparse[i] != C[i] {
			log.Println(ma.ID, ":C校验错误")
			log.Println("c:", C)
			log.Println(ma.ID, "buferr", ma.ParseMessageC())
			return false
		}
	}

	//生成秘密
	secret := []byte{8}
	checkmaBuffer := make([]byte, 0)
	checkmaBuffer = append(checkmaBuffer, int64ToByteArr(t2)...)
	checkmaBuffer = append(checkmaBuffer, int64ToByteArr(ma.ParseMessageRandomA())[0:1]...)
	checkmaBuffer = append(checkmaBuffer, secret...)
	//发送信息
	checkma.buffer = bytes.NewBuffer(checkmaBuffer)
	log.Println("checkmessage time1stamp:", ma.ParseMessageTime())
	log.Println("checkmessage time2stamp:", t2)
	return true
}

// ReplyPhaseone A检查收到的消息
/*
	A首先检验(T_3-T_2 )≤T是否成立，
	如果成立则用存在自己存储器中的KHG解密NB来获得id_B^*,b,a^*,T_2^*，
	并验证是否id_B=id_B^*,a=a^*,T_2=T_2^*，
	如果全部通过则计算tag^*=HMAC[Q_AB,id_A∥id_B^*∥b∥a^*∥T_2^* ]
	并继续验证是否有tag=tag^*。
	若成立则可生成会话密钥sk=h(id_A∥id_B∥b∥T_3∥Q_AB )并计算N_AB=E_(K_HG ) [sk,b,T_3 ]，
	最后发送通知消息{N_AB,T_3 }给设备B（这里的T3表示当前设备A的时间戳）
	。其间任何验证不通过A都会发送错误信息并终止通信。*/
func (ma *Machine) ReplyPhaseone(target *Machine) {
	content := ma.buffer.Bytes()
	if len(content) < 9 {
		log.Fatal("message receive wrong")
	}
	t2 := content[0:8]
	randnum := content[8]
	se := content[9]
	time2 := byteArrToInt64(t2)
	t3 := time.Now().Unix()
	if time2 > t3 {
		log.Println("t2 is wrong")
	}
	if !bytes.Equal(ma.PrivKey, target.PrivKey) {
		log.Println(ma.ID, ":", target.ID, "的信息错误")
	}
	discussionkey := make([]byte, 32)
	easyrand.Seed(23)
	rand := easyrand.Int()
	// log.Println("ReplyPhaseone:rand:", rand)
	k := int64ToByteArr(int64(rand))
	copy(discussionkey, k)
	// log.Println("discussionkey:key", discussionkey)
	target.buffer.Reset()
	// log.Println("discussionkey:targetbuffer", target.buffer.Bytes())
	//添加随机数消息，sercert，time3等信息
	discussionkey = append(discussionkey, se, randnum)
	discussionkey = append(discussionkey, int64ToByteArr(t3)...)
	_, err := target.buffer.Write(discussionkey)
	if err != nil {
		log.Println("ReplyPhaseone", "sendto", target.ID, "wrong")
	}
	log.Println("phaseone time2stamp:", byteArrToInt64(t2))
	log.Println("phaseone time3stamp:", t3)
}

/*
	一旦设备B接收到来自设备A的通知消息，设备B先检验(T_4-T_3 )≤T是否成立，
	如果成立则用KHG解密NAB来获得sk^*,b^*,T_3^*。验证是否T_3=T_3^*,b=b^*，
	若全部通过则承认会话密钥sk=h(id_A∥id_B∥b∥T_3∥T_2∥Q_AB )，
	至此设备B和设备A成功建立安全的会话密钥（这里的T4表示当前设备B的时间戳）。
*/
func (ma *Machine) ReplyPhasetwo(target *Machine) {
	message := ma.buffer.Bytes()
	if len(message) == 0 {
		log.Println(ma.ID, "ReplyPhasetwo:message is wrong")
	}
	// log.Println("ReplyPhasetwo:lenmessage", len(message))
	t4 := time.Now().Unix()
	t3 := byteArrToInt64(message[34:])
	log.Println("phasetwo time3stamp:", t3)
	log.Println("phasetwo time4stamp:", t4)
	if t3 > t4 {
		log.Println("reply phasetwo:", "t3 is wrong")
	}
	log.Println("Machine", ma.ID, "get the discussionkey:", message[0:32])

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
	//生成SP
	sp := NewSP()
	//向sp注册a,b两个设备，第一个参数为设备名必须为单字节，第二个参数为它们所在的网关名
	a := sp.RegisterMachine("1", "A")
	fmt.Println(a)
	b := sp.RegisterMachine("2", "A")
	fmt.Println(b)
	// c := sp.RegisterMachine("2", "B")
	// fmt.Println(c)
	//设备a向b发起通信
	a.Commucate(b)
	time.Sleep(1 * time.Second)
	//b检查a发送到它缓冲区的消息，并回复消息
	b.CheckMessage(a)
	//time.Sleep是为了等待b向a的缓存写入结束
	time.Sleep(1 * time.Second)
	//a回复b，检验安全参数
	a.ReplyPhaseone(b)
	//b检验a发回的消息，得到会话密钥
	time.Sleep(1 * time.Second)
	b.ReplyPhasetwo(a)
	//随机数和时间检验DEBUG
	// fmt.Println("parserandoma", b.ParseMessageRandomA())
	// fmt.Println("parsetime", b.ParseMessageTime())

}
