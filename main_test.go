package main

import (
	"testing"
	"time"
)

func TestConvertToBackupPath1(t *testing.T) {
	actual := convertToBackupPath("07270128.nma")
	expected := "07270128.bak"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestConvertToBackupPath2(t *testing.T) {
	actual := convertToBackupPath(".\\07270128.nma")
	expected := ".\\07270128.bak"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestConvertToBackupPath3(t *testing.T) {
	actual := convertToBackupPath("C:\\foo\\bar\\07270128.nma")
	expected := "C:\\foo\\bar\\07270128.bak"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestConvertToBackupPath4(t *testing.T) {
	actual := convertToBackupPath("07270128")
	expected := "07270128.bak"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestCorrectGprmc1(t *testing.T) {
	actual1, actual2 := correctGprmc("$GPRMC,012806.000,A,3627.2071,N,13922.7285,E,0.01,0.00,270700,,,A*67")
	expected1 := "$GPRMC,012806.000,A,3627.2071,N,13922.7285,E,0.01,0.00,120320,,,A*67"
	expected2 := time.Date(2020, 3, 12, 1, 28, 6, 0, time.UTC)
	if actual1 != expected1 {
		t.Errorf("actual %v / expected: %v", actual1, expected1)
	}
	if actual2 != expected2 {
		t.Errorf("actual %v / expected: %v", actual2, expected2)
	}
}

func TestCorrectGprmc2(t *testing.T) {
	actual1, actual2 := correctGprmc("$GPRMC,001237.000,A,3537.4722,N,13900.5570,E,0.04,0.00,190999,,,A*61")
	expected1 := "$GPRMC,001237.000,A,3537.4722,N,13900.5570,E,0.04,0.00,050519,,,A*68"
	expected2 := time.Date(2019, 5, 5, 0, 12, 37, 0, time.UTC)
	if actual1 != expected1 {
		t.Errorf("actual %v / expected: %v", actual1, expected1)
	}
	if actual2 != expected2 {
		t.Errorf("actual %v / expected: %v", actual2, expected2)
	}
}

func TestCorrectDate1(t *testing.T) {
	actual, _ := correctDate("270700")
	expected := "120320"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestCorrectDate2(t *testing.T) {
	actual, _ := correctDate("190999")
	expected := "050519"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestCalcCheckSum1(t *testing.T) {
	actual := calcCheckSum("GPRMC,012806.000,A,3627.2071,N,13922.7285,E,0.01,0.00,270700,,,A")
	expected := 0x67
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestCalcCheckSum2(t *testing.T) {
	actual := calcCheckSum("GPRMC,001237.000,A,3537.4722,N,13900.5570,E,0.04,0.00,190999,,,A")
	expected := 0x61
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestCorrectCheckSum1(t *testing.T) {
	actual := correctCheckSum("$GPRMC,012806.000,A,3627.2071,N,13922.7285,E,0.01,0.00,120320,,,A*67")
	expected := "$GPRMC,012806.000,A,3627.2071,N,13922.7285,E,0.01,0.00,120320,,,A*67"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestCorrectCheckSum2(t *testing.T) {
	actual := correctCheckSum("$GPRMC,001237.000,A,3537.4722,N,13900.5570,E,0.04,0.00,050519,,,A*61")
	expected := "$GPRMC,001237.000,A,3537.4722,N,13900.5570,E,0.04,0.00,050519,,,A*68"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestCreateFileName1(t *testing.T) {
	actual := createFileName(time.Date(2000, 7, 27, 1, 28, 6, 0, time.UTC))
	expected := "07270128"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}

func TestCreateFileName2(t *testing.T) {
	actual := createFileName(time.Date(1999, 9, 19, 0, 12, 37, 0, time.UTC))
	expected := "09190012"
	if actual != expected {
		t.Errorf("actual %v / expected: %v", actual, expected)
	}
}
