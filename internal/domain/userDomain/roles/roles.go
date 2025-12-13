package roles

type Role string

const (
	RoleDoctor  Role = "DOCTOR"
	RolePatient Role = "PATIENT"
	RoleAdmin   Role = "ADMIN"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleDoctor, RolePatient, RoleAdmin:
		return true
	default:
		return false
	}
}
