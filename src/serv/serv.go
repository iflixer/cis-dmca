package serv

import (
	"dmca/telegram"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/cloudmailin/cloudmailin-go"
)

type Service struct {
	// mux                   *http.ServeMux
	port                  string
	dleApiURL             string
	dleBasicLogin         string
	dleBasicPassword      string
	telegramReportGroupID int64
	telegramService       *telegram.Service
}

func (s *Service) Run() {
	log.Println("Start http server on :" + s.port)
	http.HandleFunc("/incoming", func(w http.ResponseWriter, req *http.Request) {
		log.Println("new email!")
		message, err := cloudmailin.ParseIncoming(req.Body)
		if err != nil {
			http.Error(w, "Error parsing message: "+err.Error(), http.StatusUnprocessableEntity)
			return
		}
		body := message.Plain

		body = strings.ReplaceAll(body, "[.]", ".")
		body = strings.ReplaceAll(body, "hxxps", "https")

		links := s.parseLinks(body)

		if os.Getenv("DEBUG_EMAIL_BODY") == "1" {
			log.Printf("BODY:%s", body)
		}

		log.Println("links:", links)

		answer := ""
		dleError := error(nil)
		if len(links) > 0 {
			msg := ""
			for _, link := range links {
				msg = msg + link + "\n"
			}
			answer, dleError = s.sendLinks(links)
		} else {
			log.Println("no links body:", body)
			answer = "no links in email"
		}

		if dleError == nil {
			s.telegramService.Send(s.telegramReportGroupID, fmt.Sprintf("DMCA complain!\nDLE answer:%s", answer))
		} else {
			s.telegramService.Send(s.telegramReportGroupID, fmt.Sprintf("DMCA complain!\nError sending to DLE:%s", dleError))
		}

		// send to dle to change links

	})

	http.HandleFunc("/alive", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintf(w, "OK")
	})

	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		panic(err)
	}
}

func (s *Service) sendLinks(links []string) (answer string, err error) {

	client := &http.Client{}
	// v := url.Values{}
	// for _, link := range links {
	// 	v.Add("links[]", link)
	// }

	linksStr := strings.Join(links, "\n")

	req, _ := http.NewRequest("POST", s.dleApiURL, strings.NewReader(linksStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.dleBasicLogin, s.dleBasicPassword)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Printf("RESPONCE:%s\n", body)
	// s.telegramService.Send(s.telegramReportGroupID, fmt.Sprintf("DMCA DLE RESPONCE:\n%s", body))
	return string(body), nil
}

func (s *Service) parseLinks(text string) (res []string) {
	r := regexp.MustCompile(`(?m:^https.*\.html)`)
	res = r.FindAllString(text, -1)
	return
}

func NewService(port string, telegramService *telegram.Service, dleApiURL string, telegramReportGroupID int64, dleBasicLogin, dleBasicPassword string) (*Service, error) {
	s := &Service{
		port:                  port,
		dleApiURL:             dleApiURL,
		dleBasicLogin:         dleBasicLogin,
		dleBasicPassword:      dleBasicPassword,
		telegramReportGroupID: telegramReportGroupID,
		telegramService:       telegramService,
	}
	return s, nil
}
