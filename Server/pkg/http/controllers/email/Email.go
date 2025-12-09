package email

type EmailController struct {
	services EmailServices
}

func NewEmailController(services EmailServices) *EmailController {
	return &EmailController{services: services}
}
