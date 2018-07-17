package checkmail

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

var emailRexp = regexp.MustCompile("^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,64}$")
var noContext = context.Background()
// Add Timeout 5 second
var defaultDialer = net.Dialer{
	Timeout: 5 * time.Second,
}

func Check(checkEmail string) (result Result, err error) {
	if !CheckSyntax(checkEmail) {
		return InvalidSyntax, nil
	}

	if CheckDisposable(checkEmail) {
		return Disposable, nil
	}
	return CheckMailbox("info@e-mas.com", checkEmail)
	//return Valid, nil
}

// CheckSyntax returns true for a valid email, false otherwise
func CheckSyntax(checkEmail string) bool {
	return emailRexp.Match([]byte(checkEmail))
}

// CheckDisposable returns true if the mail is a disposal mail, false otherwise
func CheckDisposable(checkEmail string) bool {
	host := strings.ToLower(hostname(checkEmail))
	return DisposableDomains[host]
}

// CheckMailbox checks the checkEmail by connecting to the target mailbox and returns the result.
func CheckMailbox(fromEmail, checkEmail string) (result Result, err error) {
	mxList, err := net.LookupMX(hostname(checkEmail))

	if err != nil || len(mxList) == 0 {
		return InvalidDomain, nil
	}

	return Valid, nil
	//return checkMailbox(noContext, fromEmail, checkEmail, mxList, 25)
}

type checkRv struct {
	res Result
	err error
}

func checkMailbox(ctx context.Context, fromEmail, checkEmail string, mxList []*net.MX, port int) (result Result, err error) {
	// try to connect to one mx
	var c *smtp.Client
	for _, mx := range mxList {
		var conn net.Conn
		conn, err = defaultDialer.DialContext(ctx, "tcp", fmt.Sprintf("%v:%v", mx.Host, port))
		// Try other port if blocked (port 587)
		if err != nil {
			conn, err = defaultDialer.DialContext(ctx, "tcp", fmt.Sprintf("%v:%v", mx.Host, 587))
			// Try other port if blocked (port 2525)
			if err != nil {
				conn, err = defaultDialer.DialContext(ctx, "tcp", fmt.Sprintf("%v:%v", mx.Host, 2525))
				// Throw error bad email
				if err != nil {
					if err.(*net.OpError).Timeout() == true {
						return TimeoutError, err
					}
					return NetworkError, err
				}
			}
		}

		c, err = smtp.NewClient(conn, mx.Host)

		if err == nil {
			break
		}
	}
	if err != nil {
		return MailserverError, err
	}
	if c == nil {
		// just to get very sure, that we have a connection
		// this code line should never be reached!
		return MailserverError, fmt.Errorf("can't obtain connection for %v", checkEmail)
	}

	resChan := make(chan checkRv, 1)

	go func() {
		defer c.Close()
		defer c.Quit() // defer ist LIFO
		// HELO
		err = c.Hello(hostname(fromEmail))
		if err != nil {
			resChan <- checkRv{MailserverError, err}
			return
		}

		// MAIL FROM
		err = c.Mail(fromEmail)
		if err != nil {
			resChan <- checkRv{MailserverError, err}
			return
		}

		// RCPT TO
		id, err := c.Text.Cmd("RCPT TO:<%s>", checkEmail)
		if err != nil {
			resChan <- checkRv{MailserverError, err}
			return
		}
		c.Text.StartResponse(id)
		code, _, err := c.Text.ReadResponse(25)
		c.Text.EndResponse(id)
		if code == 550 {
			resChan <- checkRv{MailboxUnavailable, nil}
			return
		}

		if err != nil {
			resChan <- checkRv{MailserverError, err}
			return
		}

		resChan <- checkRv{Valid, nil}

	}()
	select {
	case <-ctx.Done():
		return TimeoutError, ctx.Err()
	case q := <-resChan:
		return q.res, q.err
	}
}

func hostname(mail string) string {
	return mail[strings.Index(mail, "@")+1:]
}