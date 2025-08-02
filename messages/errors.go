package messages

import "errors"

var (
	ErrTryAgain            = errors.New("เกิดความผิดพลาดโปรดลองใหม่")
	ErrInvalidJSONInput    = errors.New("ข้อมูล JSON ไม่ถูกต้อง")
	ErrBookingNotFound     = errors.New("ไม่พบข้อมูลคำขอ")
	ErrInternalServer      = errors.New("เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์")
	ErrNotfound            = errors.New("ไม่พบข้อมูล")
	ErrInvalidUID          = errors.New("รหัสคำขอธุรกรรมไม่ถูกต้อง")
	ErrCreateRequest       = errors.New("ไม่สามารถสร้างคำขอได้")
	ErrBookingCannotUpdate = errors.New("ไม่สามารถอัปเดตคำขอได้")
	ErrAnnualCannotUpdate  = errors.New("ไม่สามารถอัปเดตคำขอได้")
	ErrForbidden           = errors.New("ไม่สามารถเข้าถึงข้อมูลนี้ได้")
	ErrBadRequest          = errors.New("คำขอไม่ถูกต้อง")
	ErrInvalidDate         = errors.New("วันที่ไม่ถูกต้อง")
	ErrInvalidFileType     = errors.New("ประเภทไฟล์ไม่ถูกต้อง")
	ErrUnauthorized        = errors.New("ไม่มีสิทธิ์เข้าถึง")
	ErrInvalidRequest      = errors.New("คำขอไม่ถูกต้อง")
	ErrAlreadyExist        = errors.New("ข้อมูลนี้มีอยู่ในระบบแล้ว")
)
