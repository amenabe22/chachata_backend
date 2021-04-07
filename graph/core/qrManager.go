package core

import (
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

// type Export struct {
// }

func GenerateQrOnSignup(uid string) (string, error) {
	path := "/home/anonny/projects/fun/chachata/chachata_backend/static/"
	storagePath := fmt.Sprintf("%s%s.png", path, uid)
	err := qrcode.WriteFile(uid, qrcode.Medium, 256, storagePath)
	if err != nil {
		return "error", err
	}
	return fmt.Sprintf("%s.png", uid), nil
}
