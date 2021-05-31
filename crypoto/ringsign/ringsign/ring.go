package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

type Ring []*ecdsa.PublicKey //公钥数组声明

type RingSign struct {
	Size  int              // 环的大小
	M     [32]byte         // 要签名的消息
	C     *big.Int         // 环签名值
	S     []*big.Int       // 环签名值
	Ring  Ring             // 公钥数组
	I     *ecdsa.PublicKey // 选取的公钥镜像
	Curve elliptic.Curve
}

//ecdsa.Publickey
/*
type PublicKey struct {
	elliptic.Curve
	X, Y *big.Int
}
*/
/*
type PrivateKey struct {
	PublicKey
	D *big.Int
}
*/
// 辅助函数，返回v的类型
func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// Bytes以字节片的形式返回公钥环。
func (r Ring) Bytes() (b []byte) {
	for _, pub := range r {
		b = append(b, pub.X.Bytes()...)
		b = append(b, pub.Y.Bytes()...)
	}
	return
}

func PadTo32Bytes(in []byte) (out []byte) {
	out = append(out, in...)
	for {
		if len(out) == 32 {
			return
		}
		out = append([]byte{0}, out...)
	}
}

//将签名转换为字节数组

func (r *RingSign) Serialize() ([]byte, error) {
	sig := []byte{}
	//添加大小和消息
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(r.Size))
	sig = append(sig, b[:]...)                      // 8 bytes
	sig = append(sig, PadTo32Bytes(r.M[:])...)      // 32 bytes
	sig = append(sig, PadTo32Bytes(r.C.Bytes())...) // 32 bytes

	// 每次迭代96bytes
	for i := 0; i < r.Size; i++ {
		sig = append(sig, PadTo32Bytes(r.S[i].Bytes())...)
		sig = append(sig, PadTo32Bytes(r.Ring[i].X.Bytes())...)
		sig = append(sig, PadTo32Bytes(r.Ring[i].Y.Bytes())...)
	}

	// 64 bytes
	sig = append(sig, PadTo32Bytes(r.I.X.Bytes())...)
	sig = append(sig, PadTo32Bytes(r.I.Y.Bytes())...)

	if len(sig) != 32*(3*r.Size+4)+8 {
		return []byte{}, errors.New("Could not serialize ring signature")
	}

	return sig, nil
}

//将字节化签名反序列化为RingSign结构体
func Deserialize(r []byte) (*RingSign, error) {
	sig := new(RingSign)
	size := r[0:8]

	if len(r) < 72 {
		return nil, errors.New("incorrect ring size")
	}

	m := r[8:40]

	var m_byte [32]byte
	copy(m_byte[:], m)

	size_uint := binary.BigEndian.Uint64(size)
	size_int := int(size_uint)

	sig.Size = size_int
	sig.M = m_byte
	sig.C = new(big.Int).SetBytes(r[40:72])

	bytelen := size_int * 96

	if len(r) < bytelen+136 {
		return nil, errors.New("incorrect ring size")
	}

	j := 0
	sig.S = make([]*big.Int, size_int)
	sig.Ring = make([]*ecdsa.PublicKey, size_int)

	for i := 72; i < bytelen; i += 96 {
		s_i := r[i : i+32]
		x_i := r[i+32 : i+64]
		y_i := r[i+64 : i+96]

		sig.S[j] = new(big.Int).SetBytes(s_i)
		sig.Ring[j] = new(ecdsa.PublicKey)
		sig.Ring[j].X = new(big.Int).SetBytes(x_i)
		sig.Ring[j].Y = new(big.Int).SetBytes(y_i)
		sig.Ring[j].Curve = crypto.S256()
		j++
	}

	sig.I = new(ecdsa.PublicKey)
	sig.I.X = new(big.Int).SetBytes(r[bytelen+72 : bytelen+104])
	sig.I.Y = new(big.Int).SetBytes(r[bytelen+104 : bytelen+136])
	sig.Curve = crypto.S256()

	return sig, nil
}

//获取公钥环并将与' privkey '相对应的公钥放在环的索引s中
//返回一个类型为[]*ecdsa.PublicKey环
func GenKeyRing(ring []*ecdsa.PublicKey, privkey *ecdsa.PrivateKey, s int) ([]*ecdsa.PublicKey, error) {
	size := len(ring) + 1
	new_ring := make([]*ecdsa.PublicKey, size)
	pubkey := privkey.Public().(*ecdsa.PublicKey)

	if s > len(ring) {
		return nil, errors.New("index s out of bounds")
	}

	new_ring[s] = pubkey
	for i := 1; i < size; i++ {
		idx := (i + s) % size
		new_ring[idx] = ring[i-1]
	}

	return new_ring, nil
}

//创建一个大小为' size '的环，并将与' privkey '相对应的公钥放在环的索引s中
//返回一个ecdsa公钥数组
func GenNewKeyRing(size int, privkey *ecdsa.PrivateKey, s int) ([]*ecdsa.PublicKey, error) {
	ring := make([]*ecdsa.PublicKey, size)
	pubkey := privkey.Public().(*ecdsa.PublicKey)

	if s > len(ring) {
		return nil, errors.New("index s out of bounds")
	}

	ring[s] = pubkey
	//生成公钥填满环
	for i := 1; i < size; i++ {
		idx := (i + s) % size
		priv, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}

		pub := priv.Public()
		ring[idx] = pub.(*ecdsa.PublicKey)
	}

	return ring, nil
}

// 计算key image I = x * H_p(P)，其中H_p是一个返回一个点的哈希函数
// H_p(P) = sha3(P) * G
func GenKeyImage(privkey *ecdsa.PrivateKey) *ecdsa.PublicKey {
	pubkey := privkey.Public().(*ecdsa.PublicKey)
	image := new(ecdsa.PublicKey)

	//计算sha3 (P)
	h_x, h_y := HashPoint(pubkey)

	//计算H_p(p) = x * sha3(p) * G
	i_x, i_y := privkey.Curve.ScalarMult(h_x, h_y, privkey.D.Bytes())

	image.X = i_x
	image.Y = i_y
	return image
}

func HashPoint(p *ecdsa.PublicKey) (*big.Int, *big.Int) {
	hash := sha3.Sum256(append(p.X.Bytes(), p.Y.Bytes()...))
	return p.Curve.ScalarBaseMult(hash[:])
}

//通过公钥数组创建环签名:
// m:字节数组，要签名的消息
// ring:ecdsa公钥数组，公钥包含在环中
// privkey: 签名人的私钥
// s:签名人在环中的索引
func Sign(m [32]byte, ring []*ecdsa.PublicKey, privkey *ecdsa.PrivateKey, s int) (*RingSign, error) {
	// 检查 ringsize > 1
	ringsize := len(ring)
	if ringsize < 2 {
		return nil, errors.New("环中少于两个成员")
	} else if s >= ringsize || s < 0 {
		return nil, errors.New("签名人密钥索引超出环大小")
	}

	// 开始签名
	//pubkey := privkey.Public().(*ecdsa.PublicKey)
	pubkey := &privkey.PublicKey
	curve := pubkey.Curve
	sig := new(RingSign)
	sig.Size = ringsize
	sig.M = m
	sig.Ring = ring
	sig.Curve = curve

	// 检查ring在s位置的公钥确实是签名人公钥
	if ring[s] != pubkey {
		return nil, errors.New("不是签名人公钥")
	}

	// 复刻privkey
	image := GenKeyImage(privkey)
	sig.I = image

	//从c[1]开始
	//取随机标量u (glue value)，计算c[1] = H(m, u*G)，其中H是哈希函数，G是曲线的基点
	C := make([]*big.Int, ringsize)
	S := make([]*big.Int, ringsize)

	//取随机标量u
	u, err := rand.Int(rand.Reader, curve.Params().P)
	if err != nil {
		return nil, err
	}

	//从秘密索引s开始
	//计算L_s = u*G
	l_x, l_y := curve.ScalarBaseMult(u.Bytes())
	//计算R s = u*H p(p [s])
	h_x, h_y := HashPoint(pubkey)
	r_x, r_y := curve.ScalarMult(h_x, h_y, u.Bytes())

	l := append(l_x.Bytes(), l_y.Bytes()...)
	r := append(r_x.Bytes(), r_y.Bytes()...)

	//连接m和u*G，计算c[s+1] = H(m, L_s, R_s)
	C_i := sha3.Sum256(append(m[:], append(l, r...)...))
	idx := (s + 1) % ringsize
	C[idx] = new(big.Int).SetBytes(C_i[:])

	// start loop at s+1
	for i := 1; i < ringsize; i++ {
		idx := (s + i) % ringsize

		// 选取随机标量 s_i
		s_i, err := rand.Int(rand.Reader, curve.Params().P)
		S[idx] = s_i
		if err != nil {
			return nil, err
		}

		if curve == nil {
			return nil, errors.New(fmt.Sprintf("No curve at index %d", idx))
		}
		if ring[idx] == nil {
			return nil, errors.New(fmt.Sprintf("No public key at index %d", idx))
		}

		// 计算 L_i = s_i*G + c_i*P_i
		px, py := curve.ScalarMult(ring[idx].X, ring[idx].Y, C[idx].Bytes()) // px, py = c_i*P_i
		sx, sy := curve.ScalarBaseMult(s_i.Bytes())                          // sx, sy = s[n-1]*G
		l_x, l_y := curve.Add(sx, sy, px, py)

		// 计算 R_i = s_i*H_p(P_i) + c_i*I
		px, py = curve.ScalarMult(image.X, image.Y, C[idx].Bytes()) // px, py = c_i*I
		hx, hy := HashPoint(ring[idx])
		sx, sy = curve.ScalarMult(hx, hy, s_i.Bytes()) // sx, sy = s[n-1]*H_p(P_i)
		r_x, r_y := curve.Add(sx, sy, px, py)

		// 计算 c[i+1] = H(m, L_i, R_i)
		l := append(l_x.Bytes(), l_y.Bytes()...)
		r := append(r_x.Bytes(), r_y.Bytes()...)
		C_i = sha3.Sum256(append(m[:], append(l, r...)...))

		if i == ringsize-1 {
			C[s] = new(big.Int).SetBytes(C_i[:])
		} else {
			C[(idx+1)%ringsize] = new(big.Int).SetBytes(C_i[:])
		}
	}
	//查找S[S] = (u - c[S]*k[S]) mod P，其中k[S]是私钥，P是曲线的阶数
	S[s] = new(big.Int).Mod(new(big.Int).Sub(u, new(big.Int).Mul(C[s], privkey.D)), curve.Params().N)

	//检查 u*G = S[s]*G + c[s]*P[s]
	ux, uy := curve.ScalarBaseMult(u.Bytes()) // u*G
	px, py := curve.ScalarMult(ring[s].X, ring[s].Y, C[s].Bytes())
	sx, sy := curve.ScalarBaseMult(S[s].Bytes())
	l_x, l_y = curve.Add(sx, sy, px, py)

	// 检查 u*H_p(P[s]) = S[s]*H_p(P[s]) + C[s]*I
	px, py = curve.ScalarMult(image.X, image.Y, C[s].Bytes()) // px, py = C[s]*I
	hx, hy := HashPoint(ring[s])
	tx, ty := curve.ScalarMult(hx, hy, u.Bytes())
	sx, sy = curve.ScalarMult(hx, hy, S[s].Bytes()) // sx, sy = S[s]*H_p(P[s])
	r_x, r_y = curve.Add(sx, sy, px, py)

	l = append(l_x.Bytes(), l_y.Bytes()...)
	r = append(r_x.Bytes(), r_y.Bytes()...)

	// 检查  H(m, L[s], R[s]) == C[s+1]
	C_i = sha3.Sum256(append(m[:], append(l, r...)...))

	if !bytes.Equal(ux.Bytes(), l_x.Bytes()) || !bytes.Equal(uy.Bytes(), l_y.Bytes()) || !bytes.Equal(tx.Bytes(), r_x.Bytes()) || !bytes.Equal(ty.Bytes(), r_y.Bytes()) { //|| !bytes.Equal(C[(s+1)%ringsize].Bytes(), C_i[:]) {
		return nil, errors.New("error closing ring")
	}

	// 添加签名值
	sig.S = S
	sig.C = C[0]

	return sig, nil
}

//验证RingSign结构体中包含的环签名
//如果签名有效则返回true，否则返回false
func Verify(sig *RingSign) bool {
	//开始
	ring := sig.Ring
	ringsize := sig.Size
	S := sig.S
	C := make([]*big.Int, ringsize)
	C[0] = sig.C
	curve := sig.Curve
	image := sig.I

	// 计算 c[i+1] = H(m, s[i]*G + c[i]*P[i])
	// c[0] = H(m, s[n-1]*G + c[n-1]*P[n-1]) n是环的size大小
	for i := 0; i < ringsize; i++ {
		// 计算 L_i = s_i*G + c_i*P_i
		px, py := curve.ScalarMult(ring[i].X, ring[i].Y, C[i].Bytes()) // px, py = c_i*P_i
		sx, sy := curve.ScalarBaseMult(S[i].Bytes())                   // sx, sy = s[i]*G
		l_x, l_y := curve.Add(sx, sy, px, py)

		// 计算 R_i = s_i*H_p(P_i) + c_i*I
		px, py = curve.ScalarMult(image.X, image.Y, C[i].Bytes()) // px, py = c[i]*I
		hx, hy := HashPoint(ring[i])
		sx, sy = curve.ScalarMult(hx, hy, S[i].Bytes()) // sx, sy = s[i]*H_p(P[i])
		r_x, r_y := curve.Add(sx, sy, px, py)

		// 计算 c[i+1] = H(m, L_i, R_i)
		l := append(l_x.Bytes(), l_y.Bytes()...)
		r := append(r_x.Bytes(), r_y.Bytes()...)
		C_i := sha3.Sum256(append(sig.M[:], append(l, r...)...))

		if i == ringsize-1 {
			C[0] = new(big.Int).SetBytes(C_i[:])
		} else {
			C[i+1] = new(big.Int).SetBytes(C_i[:])
		}
	}

	return bytes.Equal(sig.C.Bytes(), C[0].Bytes())
}

//比较两次环签名的公钥中的x，y，相等返回true
func Link(sig_a *RingSign, sig_b *RingSign) bool {
	return sig_a.I.X.Cmp(sig_b.I.X) == 0 && sig_a.I.Y.Cmp(sig_b.I.Y) == 0
}

//测试
func main() {
	TestLinkabilityTrue()
	TestLinkabilityFalse()
}

func TestLinkabilityTrue() {
	/* 从16进制字符串生成公私钥对 */
	privkey, _ := crypto.HexToECDSA("358be44145ad16a1add8622786bef07e0b00391e072855a5667eb3c78b9d3803")

	/* 生成消息摘要 */
	msg1 := "helloworld"
	msgHash1 := sha3.Sum256([]byte(msg1))

	/* 生成环 */
	keyring1, err := GenNewKeyRing(2, privkey, 0)
	log.Printf("环中公钥")
	for k,v := range keyring1{
		log.Printf("第%d个公钥",k+1)
		log.Printf("%+v",v)
	}
	if err != nil {
		log.Fatal(err)
	}
	/*签名*/
	sig1, err := Sign(msgHash1, keyring1, privkey, 0)
	if err != nil {
		log.Println("error sign")
	} else {
		log.Println("签名成功")
		log.Printf("%v\n", sig1)
		log.Println("-----------------------------")
		spew.Dump(sig1.I)
	}
	/* 第二次签名过程 */
	/* 生成摘要 */
	msg2 := "hello world"
	msgHash2 := sha3.Sum256([]byte(msg2))

	/* 生成环 */
	keyring2, err := GenNewKeyRing(2, privkey, 0)
	if err != nil {
		log.Println(err)
	}

	sig2, err := Sign(msgHash2, keyring2, privkey, 0)
	if err != nil {
		log.Println("error sign")
	} else {
		log.Println(" 第二次签名成功")
		log.Println("-----------------------------")
		log.Printf("%v\n", sig2)
	}
	/*验证是否可连接*/
	link := Link(sig1, sig2)
	if link {
		log.Println("签名是可连接")
	} else {
		log.Println("不可链接")
	}
}

func TestLinkabilityFalse() {
	privkey1, _ := crypto.HexToECDSA("358be44145ad16a1add8622786bef07e0b00391e072855a5667eb3c78b9d3803")
	/* 生成消息摘要 */
	msg1 := "helloworld"
	msgHash1 := sha3.Sum256([]byte(msg1))

	/* 生成环 */
	keyring1, err := GenNewKeyRing(2, privkey1, 0)
	if err != nil {
		log.Fatal(err)
	}
	/*签名*/
	sig1, err := Sign(msgHash1, keyring1, privkey1, 0)
	if err != nil {
		log.Println("error sign")
	} else {
		log.Println("签名成功")
		log.Printf("%v", sig1)
		log.Println("-----------------------------")
		//打印公钥
		spew.Dump(sig1.I)
	}

	privkey2, _ := crypto.HexToECDSA("01ad23ee4fbabbcf31dda1270154a623f5f7c07433193ff07395b33ac5bf2bea")
	/* 第二次签名过程 */
	/* 生成摘要 */
	msg2 := "hello world"
	msgHash2 := sha3.Sum256([]byte(msg2))

	/* 生成环 */
	keyring2, err := GenNewKeyRing(2, privkey2, 0)
	if err != nil {
		log.Println(err)
	}

	sig2, err := Sign(msgHash2, keyring2, privkey2, 0)
	if err != nil {
		log.Println("error sign")
	} else {
		log.Println(" 第二次签名成功")
		log.Println("-----------------------------")
		log.Printf("%v\n", sig2)

	}
	/*验证是否可连接*/
	link := Link(sig1, sig2)
	if link {
		log.Println("签名是可连接")
	} else {
		log.Println("签名使用了不同的私钥，不可连接")
	}
}
