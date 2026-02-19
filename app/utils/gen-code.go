package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

func NextAlphaCode(last string) string {
	if last == "" {
		return "A"
	}
	// ถ้า last เป็นชุดเดียวกัน เช่น "A", "B", ..., "Z", "AA", "BB", ..., "ZZ", "AAA", ...
	// ให้เพิ่มความยาวอีก 1 ถ้า last เป็น Z, ZZ, ZZZ, ...
	upper := strings.ToUpper(last)
	n := len(upper)
	for i := 0; i < n; i++ {
		if upper[i] != 'Z' {
			break
		}
		if i == n-1 {
			// ถ้าเป็น Z, ZZ, ZZZ, ... ให้เพิ่มความยาวอีก 1
			return strings.Repeat("A", n+1)
		}
	}
	// ถ้าไม่ใช่ Z, ZZ, ... ให้เพิ่มตัวอักษรตัวแรก เช่น A->B, B->C, ..., Y->Z, AA->BB, BB->CC
	first := upper[0]
	next := first + 1
	return strings.Repeat(string(next), n)
}

const (
	MemberNoPrefixMember = "MEM-"
	MemberNoPrefixAdmin  = "ADM-"
	ProductNoPrefix      = "PD-"
	ProductNoDigits      = 6
	OrderNoPrefix        = "ORD-"
	OrderNoDigits        = 6
)

func GenerateMemberNo(memberID uuid.UUID, role string) string {
	normalized := strings.ToUpper(strings.ReplaceAll(memberID.String(), "-", ""))
	prefix := MemberNoPrefixMember
	if strings.EqualFold(strings.TrimSpace(role), "admin") {
		prefix = MemberNoPrefixAdmin
	}
	if len(normalized) <= 8 {
		return fmt.Sprintf("%s%s", prefix, normalized)
	}
	return fmt.Sprintf("%s%s", prefix, normalized[len(normalized)-8:])
}

func GenerateProductNo() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%0*d", ProductNoPrefix, ProductNoDigits, n.Int64()), nil
}

func GenerateOrderNo() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%0*d", OrderNoPrefix, OrderNoDigits, n.Int64()), nil
}
