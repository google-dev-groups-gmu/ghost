package main

// scraper for GMU courses
// performs guest handshake to get session cookie + synchronizer token
// then fetches course data for a given subject and term
// outputs the raw JSON response to stdout
//

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"cmd/internal/firestore"
	"cmd/internal/types"

	"github.com/joho/godotenv"
)

const (
	BaseURL = "https://ssbstureg.gmu.edu/StudentRegistrationSsb"
	Term    = "202610" // spring 2026 term
)

func main() {
	// initialize HTTP client with cookie jar
	// NOTE: handles JSESSIONID automatically
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	// initialize firestore
	// NOTE: this is not an API endpoint, so we init firebase here
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load")
	}

	if err := firestore.Init(); err != nil {
		log.Printf("Failed to initialize Firestore: %v", err)
	}
	defer firestore.Close()

	// guest handshake to get X-Synchronizer-Token
	fmt.Println("== 1 == visiting Search Page to get Token...")

	// visit the main page just to parse the token from the HTML
	targetURL := BaseURL + "/ssb/classSearch/classSearch"

	req, _ := http.NewRequest("GET", targetURL, nil)
	setHeaders(req, "")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		log.Fatalf("server returned %d", resp.StatusCode)
	}

	// regex to find window.synchronizerToken
	re := regexp.MustCompile(`name="synchronizerToken"\s+content="([^"]+)"`)
	matches := re.FindStringSubmatch(bodyString)

	if len(matches) < 2 {
		// debugging html dump
		os.WriteFile("debug_page.html", bodyBytes, 0644)
		log.Fatal("could not find X-Synchronizer-Token in HTML.")
	}
	token := matches[1]
	fmt.Printf("	> Token Found: %s\n", token)

	// setting the term
	// all requests will happen after setting the term
	fmt.Println("== 2 == setting term to", Term, "...")

	formData := url.Values{}
	formData.Set("term", Term)
	formData.Set("studyPath", "")
	formData.Set("studyPathText", "")
	formData.Set("startDatepicker", "")
	formData.Set("endDatepicker", "")

	// NOTE: use the "uniqueSessionId" param from the cookies
	// to mimic a real user session
	uniqueID := fmt.Sprintf("guest%d", time.Now().Unix())
	termUrl := fmt.Sprintf("%s/ssb/term/search?mode=search&uniqueSessionId=%s", BaseURL, uniqueID)

	req, _ = http.NewRequest("POST", termUrl, strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// attach the token header
	setHeaders(req, token)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// search for classes
	// use this to get all subjects in prod
	// subjects, err := GetSubjects(client, token)
	// in local we are only testing CS and MATH
	var subjects = []string{"CS", "MATH"}

	fmt.Printf("== 3 == fetching classes for %s...\n", subjects)

	// for concurrency
	// using worker pool pattern
	var wg sync.WaitGroup

	// room aggregation
	// sync.Mutex is used to protect concurrent map writes
	// what sync.Mutex does is it allows only one goroutine
	// to access the critical section of code at a time
	rooms := make(map[string]*types.Room)
	var roomMu sync.Mutex

	for _, subj := range subjects {
		// NOTE: banner api is stateful. which means we need to reset the search
		// every time we change the subject
		// otherwise it will keep appending to the previous search
		resetSearch(client, token)

		offset := 0
		maxSize := 50

		for {
			fmt.Printf("fetching %s (offset %d)...\n", subj, offset)

			// build URL with dynamic offset
			apiURL := fmt.Sprintf(
				"%s/ssb/searchResults/searchResults?txt_subject=%s&txt_term=%s&pageOffset=%d&pageMaxSize=%d",
				BaseURL, subj, Term, offset, maxSize,
			)

			// fetch
			req, _ = http.NewRequest("GET", apiURL, nil)

			// attach the token header
			setHeaders(req, token)

			resp, err = client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			bodyBytes, _ = io.ReadAll(resp.Body)

			if resp.StatusCode != 200 {
				log.Fatalf("error %d: %s", resp.StatusCode, string(bodyBytes))
			}

			fmt.Printf("received %d bytes of JSON.\n", len(bodyBytes))
			fmt.Println(string(bodyBytes[:1000]))

			// unmarshal into struct
			var response types.BannerResponse
			json.Unmarshal(bodyBytes, &response)

			if len(response.Data) == 0 {
				break
			}

			// process and save
			for _, rawSec := range response.Data {
				// get all meetings for this section
				meetings := parseBannerMeetings(rawSec)

				roomMu.Lock()
				for _, meeting := range meetings {

					// filter unknown locations
					if strings.Contains(meeting.Location, "Online") || strings.Contains(meeting.Location, "TBA") {
						continue
					}

					roomID := strings.ReplaceAll(meeting.Location, " ", "_")

					if _, exists := rooms[roomID]; !exists {
						parts := strings.Split(meeting.Location, " ")
						number := ""
						building := meeting.Location
						if len(parts) > 1 {
							number = parts[len(parts)-1]
							building = strings.Join(parts[:len(parts)-1], " ")
						}

						rooms[roomID] = &types.Room{
							ID:       roomID,
							Building: building,
							Number:   number,
							Schedule: []types.Meeting{},
						}
					}

					// add to schedule
					rooms[roomID].Schedule = append(rooms[roomID].Schedule, meeting)
				}
				roomMu.Unlock()
			}

			// pagination: increment offset
			offset += maxSize
			if offset >= response.TotalCount {
				break
			}

			// no classes found, skip to next subject
			if response.TotalCount == 0 {
				fmt.Printf("   > No classes found for %s. Skipping.\n", subj)
				break
			}

			// this is just to be polite
			// and not get rate limited by the server lol
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Printf("sleeping 2s before next subject\n")
		time.Sleep(2 * time.Second)
	}
	// wait for all goroutines to finish
	wg.Wait()

	// save rooms
	fmt.Println("== 4 == Saving Room Schedules...")
	var roomWg sync.WaitGroup
	// semaphore to limit concurrency
	sem := make(chan struct{}, 20)

	for _, r := range rooms {
		roomWg.Add(1)
		sem <- struct{}{}
		go func(room *types.Room) {
			defer roomWg.Done()
			defer func() { <-sem }()

			if err := firestore.SaveRoom(context.Background(), *room); err != nil {
				log.Printf("Error saving room %s: %v", room.ID, err)
			} else {
				fmt.Printf("   > Saved Room: %s\n", room.ID)
			}
		}(r)
	}
	roomWg.Wait()

	fmt.Println("== DONE == all subjects processed.")
}

// clears the search criteria in the session
func resetSearch(client *http.Client, token string) {
	resetURL := fmt.Sprintf("%s/ssb/classSearch/resetDataForm", BaseURL)
	req, _ := http.NewRequest("POST", resetURL, nil)
	setHeaders(req, token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Warning: Failed to reset search: %v", err)
	} else {
		resp.Body.Close()
	}
}

// set the headers
func setHeaders(req *http.Request, token string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	if token != "" {
		req.Header.Set("X-Synchronizer-Token", token)
	}
}

// parsing time: 1330 -> 13 + 30
// return 0 if nil or invalid
func parseTimeStr(t *string) int {
	if t == nil {
		return 0
	}
	val := *t
	if len(val) < 4 {
		return 0
	}
	// "1330" -> 13, 30
	hh, _ := strconv.Atoi(val[:2])
	mm, _ := strconv.Atoi(val[2:])
	return (hh * 60) + mm
}

// safely dereference string pointer
func getStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// returns a list of meetings with the course info
func parseBannerMeetings(raw types.BannerSection) []types.Meeting {
	var meetings []types.Meeting

	// extract info
	profName := "Unknown"
	if len(raw.Faculty) > 0 {
		profName = raw.Faculty[0].DisplayName
	}

	info := types.MeetingInfo{
		ID:        raw.CRN,
		CourseID:  raw.Subject + raw.CourseNumber,
		Section:   raw.SequenceNumber,
		Professor: profName,
	}

	for _, mf := range raw.MeetingsFaculty {
		mt := mf.MeetingTime

		startMin := parseTimeStr(mt.BeginTime)
		endMin := parseTimeStr(mt.EndTime)

		if startMin == 0 || endMin == 0 {
			continue
		}

		bldg := getStr(mt.Building)
		room := getStr(mt.Room)
		location := fmt.Sprintf("%s %s", bldg, room)
		if bldg == "" || room == "" {
			location = "TBA"
		}

		daysMap := map[int]bool{
			0: mt.Sunday, 1: mt.Monday, 2: mt.Tuesday,
			3: mt.Wednesday, 4: mt.Thursday, 5: mt.Friday, 6: mt.Saturday,
		}

		for dayCode, isActive := range daysMap {
			if isActive {
				meetings = append(meetings, types.Meeting{
					Day:       dayCode,
					StartTime: startMin,
					EndTime:   endMin,
					Location:  location,
					Label:     []types.MeetingInfo{info},
				})
			}
		}
	}
	return meetings
}

// fetch all subjects from banner
func GetSubjects(client *http.Client, token string) ([]string, error) {
	fmt.Println("== 0 == Fetching Subject List...")

	// fetch all subjects (max=500 should cover it)
	apiURL := fmt.Sprintf(
		"%s/ssb/classSearch/get_subject?searchTerm=&term=%s&offset=1&max=500",
		BaseURL, Term,
	)

	req, _ := http.NewRequest("GET", apiURL, nil)
	setHeaders(req, token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var subjects []types.BannerSubject
	if err := json.Unmarshal(body, &subjects); err != nil {
		return nil, err
	}

	var codes []string
	for _, s := range subjects {
		codes = append(codes, s.Code)
	}

	fmt.Printf("   > Found %d active subjects.\n", len(codes))
	return codes, nil
}
