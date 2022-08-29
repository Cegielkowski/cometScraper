package entity

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0

var RequestIDHeader = "X-Request-Id"

const (
	Start             string = "PROCESS STARTED, WILL LOGIN"
	FailedCredentials        = "WRONG CREDENTIALS"
	Logged                   = "LOGGED SUCCESSFULLY, GOING TO CRAWL THE BASIC DATA"
	Fail                     = "INTERNAL ERROR, CONTACT ADMIN"
	Basic                    = "CRAWLED BASIC DATA, STARTING TO CRAWL EXPERIENCES AND SKILLS"
	Success                  = "SUCCESS"
	TimeOut                  = "THE OPERATION TOOK LONGER THAN EXPECTED, PLEASE TRY AGAIN"
)
