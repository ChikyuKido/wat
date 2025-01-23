package wat

var Config struct {
	WatBaseURL           string
	SmtpServer           string
	SmtpHost             string
	SmtpPassword         string
	SmtpEmail            string
	EmailVerification    bool
	EmailVerificationUrl string
	JwtSecret            string
	Debug                bool
	AllowedEmails        []string
	FirstUser            bool
	ResourceVersion      string
}
