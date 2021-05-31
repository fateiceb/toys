package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

// sha256[ c Condition ]
type ConditionHash []byte

func (h ConditionHash) MarshalJSON() (string, error) {
	return hex.EncodeToString(h), nil
}

func (h ConditionHash) String() string {
	s, _ := h.MarshalJSON()
	return s
}

// hmac
func (c ConditionHash) Approve(k PrivKey) ConditionProof {
	h := sha256.New()
	h.Write(k)
	h.Write(c)
	h.Write(k)
	v := h.Sum(nil)
	return ConditionProof(v[:])
}

var one = []byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 1,
}

func byteNeg(a []byte) []byte {
	// 翻转所有bit
	c := make([]byte, 32)
	for i := 31; i >= 0; i-- {
		c[i] = ^a[i]
	}
	// 添加 one
	return byteAdd(c, one)
}

func byteAdd(a []byte, b []byte) []byte {
	c := make([]byte, 32)
	carry := int(0)
	for i := 31; i >= 0; i-- {
		v := (int(a[i]) + int(b[i]) + carry) % 256
		carry = (int(a[i]) + int(b[i])) / 256
		c[i] = byte(v)
	}
	return c
}

func byteSub(a []byte, b []byte) []byte {
	return byteAdd(a, byteNeg(b))
}

// sha256[ ch ConditionHash, caPrivate ]
type ConditionProof []byte

func (h ConditionProof) MarshalJSON() (string, error) {
	return hex.EncodeToString(h), nil
}

func (h ConditionProof) String() string {
	s, _ := h.MarshalJSON()
	return s
}

type PrivKey []byte

func NewPrivKey() PrivKey {
	v := make([]byte, 32)
	_, err := rand.Read(v)
	if err != nil {
		log.Printf("error reading random: %v", err)
		return nil
	}
	return PrivKey(v)
}

func (p PrivKey) MarshalJSON() (string, error) {
	return hex.EncodeToString(p), nil
}

func (h PrivKey) String() string {
	s, _ := h.MarshalJSON()
	return s
}

func (p PrivKey) PubKey() PubKey {
	v := sha256.Sum256(p)
	return PubKey(v[:])
}

type PubKey []byte

func (p PubKey) MarshalJSON() (string, error) {
	return hex.EncodeToString(p), nil
}

func (h PubKey) String() string {
	s, _ := h.MarshalJSON()
	return s
}

// 获取条件
type Condition struct {
	Name  string
	Value string
}

// 条件hash
func (c *Condition) Hash() ConditionHash {
	v := sha256.Sum256([]byte(fmt.Sprintf("%s: %s", c.Name, c.Value)))
	return ConditionHash(v[:])
}

type KeyPair struct {
	Pub  PubKey
	Priv PrivKey
}

type Certificate map[Condition]ConditionProof

type PolicyCase struct {
	Required []ConditionHash
	// 减去condition proof的key
	Target PrivKey
}

type Policy struct {
	Pub PubKey
	// 匹配的情况
	Cases []PolicyCase
}

// 给ca policprivkey
func NewPolicy(capriv PrivKey, policypriv PrivKey, cases []PolicyCase) Policy {
	p := Policy{
		Pub:   policypriv.PubKey(),
		Cases: cases,
	}
	for i := range cases {
		// 针对case，计算
		// target = policyPriv - H[ policypub + proofa0 + proofa1 + ...]
		total := make([]byte, 32)
		total = byteAdd(total, p.Pub)
		//log.Printf("newpolicy init %s", PrivKey(total))
		c := cases[i]
		for j := range c.Required {
			a := c.Required[j]
			proof := []byte(a.Approve(capriv))
			total = byteAdd(total, proof)
			//log.Printf("newpolicy %d %d: %s",i,j,PrivKey(proof))
		}
		v := sha256.Sum256(total)
		cases[i].Target = byteSub(policypriv, v[:]) //byteSub([]byte(policypriv), v[:])
	}
	return p
}

func UnlockPolicy(cert Certificate, policy Policy) PrivKey {
	for c := range policy.Cases {
		// 针对每一个case计算
		// policyPriv = target + H[ policypub + proofa0 + proofa1 + ...]
		policyCase := policy.Cases[c]
		hasAll := true
		total := make([]byte, 32)
		total = byteAdd(total, policy.Pub)
		//log.Printf("unlock init %s", PrivKey(total))
		for r := range policyCase.Required {
			required := policyCase.Required[r]
			foundIt := false
			for a := range cert {
				if bytes.Compare(a.Hash(), required) == 0 {
					foundIt = true

					total = byteAdd(total, cert[a])
					//log.Printf("unlock %d %d %s", c, r, PrivKey(cert[a]))
				}
			}
			if foundIt == false {
				hasAll = false
			}
		}
		if hasAll {
			v := sha256.Sum256(total)
			return PrivKey(byteAdd(policyCase.Target, v[:]))
		}
	}
	return nil
}

//机构集合
type Organzations struct {
	name   []string
	police []Policy
}

func NewOrganzations(name []string, police []Policy) Organzations {

	return Organzations{
		name:   name,
		police: police,
	}
}
func (org *Organzations) Unlock(crt Certificate) []PrivKey {
	key := make([]PrivKey, len(org.name))
	for i := 0; i < len(org.name); i++ {
		if key[i] = UnlockPolicy(crt, org.police[i]); key[i] == nil {
			return nil
		}
	}
	return key
}
func main() {
	//设置属性
	isEmailAdmin := Condition{Name: "email", Value: "admin@foo.com"}
	isEmailAdminHash := isEmailAdmin.Hash()

	isAgeAdult := Condition{Name: "age", Value: "adult"}
	isAgeAdultHash := isAgeAdult.Hash()

	isCitizenUS := Condition{Name: "citizen", Value: "US"}
	isCitizenUSHash := isCitizenUS.Hash()

	isCitizenUK := Condition{Name: "citizen", Value: "UK"}
	isCitizenUKHash := isCitizenUK.Hash()

	isCitizenNL := Condition{Name: "citizen", Value: "NL"}
	//isCitizenNLHash := isCitizenNL.Hash()

	capriv := NewPrivKey()
	//capub := capriv.PubKey()
	//个人证书生成
	alice := Certificate{
		isAgeAdult:  isAgeAdult.Hash().Approve(capriv),
		isCitizenUK: isCitizenUK.Hash().Approve(capriv),
	}

	bob := Certificate{
		isAgeAdult:  isAgeAdult.Hash().Approve(capriv),
		isCitizenNL: isCitizenNL.Hash().Approve(capriv),
	}

	charles := Certificate{
		isAgeAdult:  isAgeAdult.Hash().Approve(capriv),
		isCitizenUS: isCitizenUS.Hash().Approve(capriv),
	}

	dave := Certificate{
		isEmailAdmin: isEmailAdmin.Hash().Approve(capriv),
	}

	policyPriv := NewPrivKey()

	//获取机构私钥的策略
	policy := NewPolicy(
		capriv,
		policyPriv,
		[]PolicyCase{
			PolicyCase{
				Required: []ConditionHash{isAgeAdultHash, isCitizenUKHash},
			},
			PolicyCase{
				Required: []ConditionHash{isAgeAdultHash, isCitizenUSHash},
			},
			PolicyCase{
				Required: []ConditionHash{isEmailAdminHash},
			},
		},
	)

	policyPrivb := NewPrivKey()
	policyb := NewPolicy(
		capriv,
		policyPrivb,
		[]PolicyCase{
			PolicyCase{
				Required: []ConditionHash{isAgeAdultHash},
			},
		},
	)
	//机构名称
	names := []string{"orga", "orgb"}
	//使用策略
	policys := []Policy{policy, policyb}
	//新建双机构
	orgs := NewOrganzations(names, policys)
	_ = policyPriv

	log.Printf(
		"expect: %s", policyPriv)
	log.Printf("alice unlock: %s", UnlockPolicy(alice, policy))
	log.Printf("bob unlock: %s", UnlockPolicy(bob, policy))
	log.Printf("charles unlock: %s", UnlockPolicy(charles, policy))
	log.Printf("dave unlock: %s", UnlockPolicy(dave, policy))
	log.Printf("多机构:")
	log.Println("alice unlock", orgs.name, orgs.Unlock(alice))
	log.Println("bob unlock", orgs.name, orgs.Unlock(bob))
}
