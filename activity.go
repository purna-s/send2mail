package send2mail

import (

	"fmt"
	"log"
	"strings"
	"net/smtp"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-flogo-send2mail")

// MyActivity is a stub for your Activity implementation
type send2mail struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &send2mail{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *send2mail) Metadata() *activity.Metadata {
	return a.metadata
}


// Eval implements activity.Activity.Eval
func (a *send2mail) Eval(ctx activity.Context) (done bool, err error) {
	
	
	serverport := ctx.GetInput("A_server:port").(string)
	sender := ctx.GetInput("B_sender").(string)
	rcpnt := ctx.GetInput("C_rcpnt").(string)
	msub := ctx.GetInput("E_sub").(string)
	mbody := ctx.GetInput("F_body").(string)
	
	mrcpnt := []string{rcpnt}
	
	mail := SendMail(serverport, sender, msub, mbody, mrcpnt)
	//mail := SendMail("server:port", "CEP_System_Alerts@alert.lta.gov.sg", "TestMailSubject", "Test-SendMailBody", []string{"ikiran@ncs.com.sg"})
	
	fmt.Println(mail)
	
	log.Println("Mail Sent")
	// Set the output as part of the context
	activityLog.Debugf("Activity has sent the mail Successfully")
	fmt.Println("Activity has sent the mail Successfully")

	ctx.SetOutput("output", "Mail_Sent_Successfully")

	return true, nil
}

func SendMail(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	
	t := []string{"To:", rcpnt}
	strings.Join(t, " ")
	
	s := []string{"Subject:", sub}
	strings.Join(s, " ")
	
	msg := []byte(strings.Join(t, " ") + "\r\n" + strings.Join(s, " ") + "\r\n" + mbody + "\r\n")
	//msg := []byte("To: recipient@example.net\r\n" + "Subject: discount Gophers!\r\n" + "\r\n" + "This is the email body.\r\n")
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}