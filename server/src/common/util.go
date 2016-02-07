package common

import (
	"errors"
	"log"
	"regexp"
	"net/smtp"
	"strings"
	"crypto/rand"
	"fmt"
	"html/template"
	"concept-build/server/src/config"
)
func GenSecret() (string, error) {
	c := 24
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
var signup_verify_tmpl *template.Template
var signup_welcome_tmpl *template.Template

var smtp_dsn string

const (
	TMPL_TYPE_SIGNUP_VERIFICATION 	int =	1+iota
	TMPL_TYPE_SIGNUP_WELCOME
)
type EmailData struct{
	Name string
	Email string
	Secret string
}
var template_map map[int]*template.Template

func UtilInit(cfg *config.Config) {
	
	template_map = make(map[int]*template.Template, 0) 
	// sign up verificaton template intialize
	var err error
	tpath := cfg.Base.Path + "views/email_verify.tmpl"
	template_map[TMPL_TYPE_SIGNUP_VERIFICATION], err = template.ParseFiles(tpath)
	if err != nil {
		log.Println("Verification Email template init Failed ", err)
		template_map[TMPL_TYPE_SIGNUP_VERIFICATION] = nil
	}else{
		log.Println("Verification Email template init :SUCCESS")
	}
	
	//welcome mail template initialization

	tpath = cfg.Base.Path + "views/welcome.tmpl"
	template_map[TMPL_TYPE_SIGNUP_WELCOME], err = template.ParseFiles(tpath)
	if err != nil {
		log.Println("Welcome Email template init Failed", err)
		template_map[TMPL_TYPE_SIGNUP_WELCOME] = nil
	}else{
		log.Println("Welcome Email template init SUCCESS")
	}
	if cfg.Smtp.Address == "" {
		smtp_dsn = "127.0.0.1:25"
	} else {
		smtp_dsn = cfg.Smtp.Address
	}

}

func GetTemplateData(email_type int) (map[string]string, error){
	if template_map[email_type] == nil{
		log.Println("template is not initialized for type: ", email_type)
		return nil, errors.New("template is not initialized")
	}
	data := make(map[string]string, 0)
	data["From"]="merchant.helpdesk@paytm.com"

	switch(email_type){
	case TMPL_TYPE_SIGNUP_VERIFICATION:
		data["Subject"]="Verification"
		data["Content-Type"]="text/HTML"
	case TMPL_TYPE_SIGNUP_WELCOME:
		data["Subject"]="User Signup"
	}
	return data, nil
}
func SendEmail(email_data EmailData, email_type int) error {
	var err error
	data := map[string]string{
		"Email":email_data.Email,
		"Name":email_data.Name,
		"Secret":email_data.Secret,
	}
	mdata, err := GetTemplateData(email_type)
	if err != nil{
		log.Println("Template is not initialized")
		return err
	}

	log.Println("sending email email_type", email_type)
	
	for k, v := range data{
		mdata[k] = v
	}
	
	c, err := smtp.Dial(smtp_dsn)
	if err != nil {
		log.Println(err)
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail("care@paytm.com"); err != nil {
		log.Println(err)
		return err
	}
	if err := c.Rcpt(data["To"]); err != nil {
		log.Println(err)
		return err

	}

	wc, err := c.Data()
	if err != nil {
		log.Println("failed to get channel to write data ", err)
		return err
	}
	defer wc.Close()
	err = template_map[email_type].Execute(wc, mdata)
	if err != nil {
		log.Println("error in sending the mail.", err)
		return err
	}
	return nil
}
func EmailRegexValidation(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]+$`)
	if !re.MatchString(email){
	      return false
	}
	if strings.Contains(email, "noreply") || strings.Contains(email, "no-reply"){
	      return false
	}
	return true
}
func PhoneRegexValidation(phno string) bool {
	re := regexp.MustCompile(`[1-9]\d{9}`)
	return re.MatchString(phno)
}
func NameRegexValidation(name  string) bool{
	re := regexp.MustCompile(`^[a-zA-Z0-9 ._]{2,50}$`)
	return re.MatchString(name)
}