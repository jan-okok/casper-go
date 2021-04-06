package deploy

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github/casper-go/client"
	"github/casper-go/keys"
	"github/casper-go/keys/blake2b"
	"math/big"
	"testing"
)

func TestDeployPut(t *testing.T) {
	const (
		eventStoreApi = "https://event-store-api-clarity-delta.make.services"
		RpcUrl        = "https://node-clarity-delta.make.services/rpc"
	)
	casper := client.New(RpcUrl, eventStoreApi)
	deploy, err := mockMakeDeploy()
	if err != nil {
		t.Fatal(err)
	}
	sender, err := mockSender()
	if err != nil {
		t.Fatal(err)
	}

	err = deploy.Sign(sender)
	if err != nil {
		t.Fatal(err)
	}

	marshal, err := json.Marshal(deploy)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(marshal))

	result, err := casper.PutDeploy(deploy)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

}

func TestStandardPayment(t *testing.T) {
	//03000000 020004 08
	//长度 / 数值 / tag值
	payment, err := NewStandardPayment(big.NewInt(1024))
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Println(hex.EncodeToString((*payment.ItemModuleBytes.Args.Args)[0].Value.ToBytes()))
	bytess := payment.ItemModuleBytes.ToBytes()
	fmt.Println(len(bytess))
	fmt.Println(bytess)

}

func TestNewTransfer(t *testing.T) {

	recipient, _ := mockRecipient()
	tx, _ := mockTransferSession(recipient)

	ms, _ := json.Marshal(tx)
	fmt.Println(string(ms))
}

func TestMainTrx(t *testing.T) {
	recipient, err := mockRecipient()
	if err != nil {
		t.Fatal(err)
	}
	session, err := mockTransferSession(recipient)
	if err != nil {
		t.Fatal(err)
	}
	payment, err := mockPayment()
	if err != nil {
		t.Fatal(err)
	}
	body := serializeBody(payment, session)
	fmt.Println(hex.EncodeToString(body))
	fmt.Println(blake2b.Hash(body))
	fmt.Println(len(body))

}

func TestDeploy_Sign(t *testing.T) {
	deploy, err := mockMakeDeploy()
	if err != nil {
		t.Fatal(err)
	}

	sender, err := mockSender()
	if err != nil {
		t.Fatal(err)
	}

	err = deploy.Sign(sender)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("hash:%s\n", hex.EncodeToString(deploy.Hash))
	fmt.Printf("body hash:%s\n", hex.EncodeToString(deploy.Header.BodyHash))
	jsonData, _ := json.Marshal(deploy)
	fmt.Println(string(jsonData))
}

func mockMakeDeploy() (*Deploy, error) {
	sender, err := mockSender()
	if err != nil {
		return nil, err
	}
	recipient, err := mockRecipient()
	if err != nil {
		return nil, err
	}
	session, err := mockTransferSession(recipient)
	if err != nil {
		return nil, err
	}
	payment, err := mockPayment()
	if err != nil {
		return nil, err
	}
	bb := payment.ItemModuleBytes.ToBytes()

	fmt.Println(len(bb))
	fmt.Println(bb)

	deploy, err := MakeDeploy(NewParams(sender.RawPublicKey()), session, payment)
	if err != nil {
		return nil, err
	}
	return deploy, nil
}

func mockPayment() (*ExecDeployItem, error) {
	return NewStandardPayment(big.NewInt(1024))
}

func mockRecipient() (keys.KeyHolder, error) {
	//accountHash 23400bdd68d63ffbd3446c4563bf1dd3c7648282ec19b12f0504c6d905bc816d
	//accountHex 01a027ac95925adf648e1a8902dab39e7899f919644c625f21cf4eec9d1b2f158f
	ds, err := hex.DecodeString("a027ac95925adf648e1a8902dab39e7899f919644c625f21cf4eec9d1b2f158f") //secret_key.pem
	if err != nil {
		return nil, err
	}
	holder, err := keys.NewKeyHolder(nil, ds, "ed25519")
	if err != nil {
		return nil, err
	}
	return holder, nil
}

func mockSender() (keys.KeyHolder, error) {
	p, err := hex.DecodeString("a7883a8bf29480a7448a45fec442830200e3135a0fd5bd1e4ff3424de772383ed74e5088891f2c938a38e4dbd37d18157bb65ef97a5cdef1aea44a2293d8d2b2")
	if err != nil {
		return nil, err
	}
	//accountHash c9c6301513d4cb3e71fade128734d484a849c902941be02bb5601de5bdd17310
	//accountHex  01d74e5088891f2c938a38e4dbd37d18157bb65ef97a5cdef1aea44a2293d8d2b2
	ds, err := hex.DecodeString("d74e5088891f2c938a38e4dbd37d18157bb65ef97a5cdef1aea44a2293d8d2b2") //jan_secret_key.pem
	if err != nil {
		return nil, err
	}
	holder, err := keys.NewKeyHolder(p, ds, "ed25519")
	if err != nil {
		return nil, err
	}
	return holder, nil
}

func mockTransferSession(recipient keys.KeyHolder) (*ExecDeployItem, error) {
	session, err := NewTransfer(big.NewInt(2500000000), recipient.AccountHash(), nil)
	if err != nil {
		return nil, err
	}
	return session, err

}
