package main

import "testing"

const (
	AddrTrdbPB   = "localhost:8087"     // riak PB port
	AddrDummy    = "12.345.67.890:8087" // dummy IP
	AddrTrdbHTTP = "localhost:8098"     // http port
)

func TestSetLogger(t *testing.T) {
	var err error
	Logger, err = SetLogger()
	if err != nil {
		t.Errorf("SetLogger is failure status.")
	}
}

func TestCheckConnect_TCP_Success(t *testing.T) {
	caseAddr := AddrTrdbPB
	caseProtocol := "TCP"
	myRiak := SetRiakClient(caseAddr)
	err := CheckConnect(myRiak, caseProtocol)
	if err != nil {
		t.Errorf("CheckConnect(%s) is failure status, want success.(caseAddr=%s)", caseProtocol, caseAddr)
	}
}

func TestCheckConnect_PB_Success(t *testing.T) {
	caseAddr := AddrTrdbPB
	caseProtocol := "PB"
	myRiak := SetRiakClient(caseAddr)
	err := CheckConnect(myRiak, caseProtocol)
	if err != nil {
		t.Errorf("CheckConnect(%s) is failure status, want success.(caseAddr=%s)", caseProtocol, caseAddr)
	}
}

func TestCheckConnect_TCP_UnknownIP(t *testing.T) {
	caseAddr := AddrDummy
	caseProtocol := "TCP"
	myRiak := SetRiakClient(caseAddr)
	err := CheckConnect(myRiak, caseProtocol)
	if err == nil {
		t.Errorf("CheckConnect(%s) is success status, want failure.(caseAddr=%s)", caseProtocol, caseAddr)
	}
}

func TestCheckConnect_PB_UnknownIP(t *testing.T) {
	caseAddr := AddrDummy
	caseProtocol := "PB"
	myRiak := SetRiakClient(caseAddr)
	err := CheckConnect(myRiak, caseProtocol)
	if err == nil {
		t.Errorf("CheckConnect(%s) is success status, want failure.(caseAddr=%s)", caseProtocol, caseAddr)
	}
}

func TestCheckConnect_TCP_PORT8098(t *testing.T) {
	caseAddr := AddrTrdbHTTP
	caseProtocol := "TCP"
	myRiak := SetRiakClient(caseAddr)
	err := CheckConnect(myRiak, caseProtocol)
	if err != nil {
		t.Errorf("CheckConnect(%s) is failure status, want success.(caseAddr=%s)", caseProtocol, caseAddr)
	}
}

func TestCheckConnect_PB_PORT8098(t *testing.T) {
	caseAddr := AddrTrdbHTTP
	caseProtocol := "PB"
	myRiak := SetRiakClient(caseAddr)
	err := CheckConnect(myRiak, caseProtocol)
	if err == nil {
		t.Errorf("CheckConnect(%s) is success status, want failure.(caseAddr=%s)", caseProtocol, caseAddr)
	}
}

func TestCheckConnect_UnknownProtocol(t *testing.T) {
	caseAddr := AddrTrdbPB
	caseProtocol := "unknown"
	myRiak := SetRiakClient(caseAddr)
	err := CheckConnect(myRiak, caseProtocol)
	if err == nil {
		t.Errorf("CheckConnect(%s) is success status, want failure.(caseAddr=%s)", caseProtocol, caseAddr)
	}
}
