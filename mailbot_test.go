package mailbot

import (
	"testing"
)

func TestUpgradeTLSSend_Local(t *testing.T) {

	var err error

	// Assuming that exim is running
	t.Log("ASSUMPTION: exim is running")

	sc := ServerConfig{Host: "test.mailbot.net",
		Port:      26,
		User:      "fromster@mailbot.net",
		Password:  "from me from me",
		EmailFrom: "fromster@mailbot.net"}

	t.Logf("send to %s:%d with no From: header\n", sc.Host, sc.Port)
	err = UpgradeTLSSend(sc, "", []string{"toster@mailbot.net"}, []byte("Subject: subby\n\nI am your body."))
	if err != nil {
		if err.Error() != "550 Missing a sender (From, Reply-To or Sender) header" {
			t.Fatalf("Expected to get a missing-header error, but instead got:%s", err)
		}
	} else {
		t.Fatal("UpgradeTLSSend should have failed but did not")
	}

	t.Logf("send to %s:%d with a From: header\n", sc.Host, sc.Port)
	err = UpgradeTLSSend(sc, "", []string{"toster@mailbot.net"}, []byte("From: fromster@mailbot.net\nSubject: subby\n\nI am your body."))
	if err != nil {
		t.Fatalf("UpgradeTLSSend failed with err:%s", err)
	}

	// TODO read the spool dir: /tmp/mailbot.boxes/toster

}

func TestUpgradeTLSSend_Remote(t *testing.T) {

	var err error
	sc := ServerConfig{Host: "mail.grantmurray.com",
		Port:      26,
		User:      "no-reply@grantmurray.com",
		Password:  "!!ProperPasswordNeeded!!<---------------------------",
		EmailFrom: "no-reply@grantmurray.com"}

	// bluehost insists on a From: header, and gmail wants a To: header otherwise it is treated as spam
	t.Logf("send to %s:%d with From: header\n", sc.Host, sc.Port)
	err = UpgradeTLSSend(sc, "", []string{"gmurray1966@gmail.com"}, []byte("To: gmurray1966@gmail.com\nFrom: no-reply@grantmurray.com\nSubject: mailbot test\n\nMailbot remote test - This is the email body."))
	if err != nil {
		t.Fatalf("UpgradeTLSSend failed with err:%s", err)
	}
}
