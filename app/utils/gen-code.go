package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/uptrace/bun"
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
	MemberNoPrefix  = "PKK-"
	MemberNoDigits  = 6
	ProductNoPrefix = "PD-"
	ProductNoDigits = 6
	OrderNoPrefix   = "ORD-"
	OrderNoDigits   = 6
)

func GenerateMemberNo(ctx context.Context, db bun.IDB) (string, error) {
	var last struct {
		MemberNo string `bun:"member_no"`
	}

	// err := db.NewSelect().
	// 	Model((*ent.MemberEntity)(nil)).
	// 	Column("member_no").
	// 	Where("member_no LIKE ?", MemberNoPrefix+"%").
	// 	OrderExpr("member_no DESC").
	// 	Limit(1).
	// 	Scan(ctx, &last)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return fmt.Sprintf("%s%0*d", MemberNoPrefix, MemberNoDigits, 1), nil
	// 	}
	// 	return "", err
	// }

	seqStr := strings.TrimPrefix(last.MemberNo, MemberNoPrefix)
	seq, err := strconv.Atoi(seqStr)
	if err != nil || seq < 0 {
		return fmt.Sprintf("%s%0*d", MemberNoPrefix, MemberNoDigits, 1), nil
	}

	return fmt.Sprintf("%s%0*d", MemberNoPrefix, MemberNoDigits, seq+1), nil
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
