package enums

type AppointmentStatus string

const (
	AppointmentProgrammed AppointmentStatus = "Programada"
	AppointmentConfirmed  AppointmentStatus = "Confirmada"
	AppointmentDone       AppointmentStatus = "Realizada"
	AppointmentCancelled  AppointmentStatus = "Cancelada"
	AppointmentNoShow     AppointmentStatus = "NoAsistio"
)
