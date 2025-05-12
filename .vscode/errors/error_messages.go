package errors

import "errors"

var (
	ErrInvalidJSONInput    = errors.New("ข้อมูล JSON ไม่ถูกต้อง")
	ErrBookingNotFound     = errors.New("ไม่พบข้อมูลการจอง")
	ErrInternalServer      = errors.New("เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์")
	ErrNotfound            = errors.New("ไม่พบข้อมูล")
	ErrInvalidUID          = errors.New("รหัสคำขอธุรกรรมไม่ถูกต้อง")
	ErrCreateRequest       = errors.New("ไม่สามารถสร้างคำขอได้")
	ErrBookingCannotUpdate = errors.New("ไม่สามารถอัปเดตการจองได้")

	ErrInvalidCostCenter          = errors.New("รหัสศูนย์ต้นทุนไม่ถูกต้อง")
	ErrInvalidWbsNo               = errors.New("รหัส WBS ไม่ถูกต้อง")
	ErrInvalidNetworkNo           = errors.New("รหัสเครือข่ายไม่ถูกต้อง")
	ErrInvalidProjectNo           = errors.New("รหัสโครงการไม่ถูกต้อง")
	ErrInvalidVehicleOwnerDeptSAP = errors.New("รหัสแผนกเจ้าของยานพาหนะใน SAP ไม่ถูกต้อง")
	ErrInvalidRequestNo           = errors.New("รหัสคำขอไม่ถูกต้อง")
)
