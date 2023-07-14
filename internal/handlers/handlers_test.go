package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/AdiAkhileshSingh15/bookmyroom/internal/driver"
	"github.com/AdiAkhileshSingh15/bookmyroom/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{
		"home", "/", "GET", http.StatusOK,
	},
	{
		"about", "/about", "GET", http.StatusOK,
	},
	{
		"generals-quarters", "/generals-quarter", "GET", http.StatusOK,
	},
	{
		"majors-suite", "/majors-suite", "GET", http.StatusOK,
	},
	{
		"search-availability", "/search-availability", "GET", http.StatusOK,
	},
	{
		"contact", "/contact", "GET", http.StatusOK,
	},

	// {
	// 	"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK,
	// },
	// {
	// 	"post-reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK,
	// },
	// {
	// 	"post-search-availability", "/search-availability", "POST", []postData{
	// 		{key: "start_date", value: "2020-01-01"},
	// 		{key: "end_date", value: "2020-01-02"},
	// 	}, http.StatusOK,
	// },
	// {
	// 	"post-search-availability-json", "/search-availability-json", "POST", []postData{
	// 		{key: "start_date", value: "2020-01-01"},
	// 		{key: "end_date", value: "2020-01-02"},
	// 	}, http.StatusOK,
	// },
	// {
	// 	"make-reservation-post", "/make-reservation", "POST", []postData{
	// 		{key: "first_name", value: "Adi"},
	// 		{key: "last_name", value: "Singh"},
	// 		{key: "email", value: "adads@ad.com"},
	// 		{key: "phone", value: "1234567890"},
	// 	}, http.StatusOK,
	// },
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

var reservationTests = []struct {
	name               string
	reservation        models.Reservation
	expectedStatusCode int
	expectedLocation   string
	expectedHTML       string
}{
	{
		name: "reservation-in-session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "General's Quarters",
			},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "reservation-not-in-session",
		reservation:        models.Reservation{},
		expectedStatusCode: http.StatusTemporaryRedirect,
		expectedLocation:   "/",
	},
	{
		name: "non-existent-room",
		reservation: models.Reservation{
			RoomID: 100,
			Room: models.Room{
				ID:       100,
				RoomName: "General's Quarters",
			},
		},
		expectedStatusCode: http.StatusTemporaryRedirect,
		expectedLocation:   "/",
	},
}

func TestRepository_Reservation(t *testing.T) {
	for _, e := range reservationTests {
		req, _ := http.NewRequest("GET", "/make-reservation", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.Reservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var postReservationTests = []struct {
	name                 string
	reservationSession   bool
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name:               "valid-data",
		reservationSession: true,
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"Eddy"},
			"last_name":  {"Singh"},
			"email":      {"eddy@singh.com"},
			"phone":      {"1234567890"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/reservation-summary",
	},
	{
		name:               "valid-data but reservation not in session",
		reservationSession: false,
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"Eddy"},
			"last_name":  {"Singh"},
			"email":      {"eddy@singh.com"},
			"phone":      {"1234567890"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
	{
		name:                 "missing-post-body",
		reservationSession:   true,
		postedData:           nil,
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
	{
		name:               "invalid-data",
		reservationSession: true,
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"Ed"},
			"last_name":  {"Singh"},
			"email":      {"eddy@singh.com"},
			"phone":      {"1234567890"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:               "database-insert-fails-reservation",
		reservationSession: true,
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"Eddy"},
			"last_name":  {"Singh"},
			"email":      {"eddy@singh.com"},
			"phone":      {"123456789"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
	{
		name:               "database-insert-fails-restriction",
		reservationSession: true,
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"first_name": {"RoomRestriction"},
			"last_name":  {"Singh"},
			"email":      {"eddy@singh.com"},
			"phone":      {"1234567890"},
			"room_id":    {"1"},
		},
		expectedResponseCode: http.StatusTemporaryRedirect,
		expectedLocation:     "/",
	},
}

func TestRepository_PostReservation(t *testing.T) {
	for _, e := range postReservationTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/make-reservation", nil)

		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		if e.reservationSession {
			session.Put(ctx, "reservation", models.Reservation{})
		}
		handler := http.HandlerFunc(Repo.PostReservation)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedResponseCode)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

var testPostAvailabilityData = []struct {
	name               string
	postedData         url.Values
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "rooms not available",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
		},
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name: "rooms are available",
		postedData: url.Values{
			"start_date": {"2040-01-01"},
			"end_date":   {"2040-01-02"},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "empty post body",
		postedData:         nil,
		expectedStatusCode: http.StatusTemporaryRedirect,
	},
	{
		name: "start date wrong format",
		postedData: url.Values{
			"start_date": {"invalid"},
			"end_date":   {"2040-01-02"},
		},
		expectedStatusCode: http.StatusTemporaryRedirect,
	},
	{
		name: "end date wrong format",
		postedData: url.Values{
			"start_date": {"2040-01-01"},
			"end_date":   {"invalid"},
		},
		expectedStatusCode: http.StatusTemporaryRedirect,
	},
	{
		name: "database query fails",
		postedData: url.Values{
			"start_date": {"2060-01-01"},
			"end_date":   {"2060-01-02"},
		},
		expectedStatusCode: http.StatusTemporaryRedirect,
	},
}

func TestRepository_PostAvailability(t *testing.T) {
	for _, e := range testPostAvailabilityData {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability", nil)

		}

		// get the context with session
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the request header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// make our handler a http.HandlerFunc and call
		handler := http.HandlerFunc(Repo.PostAvailability)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s gave wrong status code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

var testAvailabilityJSONData = []struct {
	name            string
	postedData      url.Values
	expectedOK      bool
	expectedMessage string
}{
	{
		name: "rooms not available",
		postedData: url.Values{
			"start_date": {"2050-01-01"},
			"end_date":   {"2050-01-02"},
			"room_id":    {"1"},
		},
		expectedOK: false,
	}, {
		name: "rooms are available",
		postedData: url.Values{
			"start_date": {"2040-01-01"},
			"end_date":   {"2040-01-02"},
			"room_id":    {"1"},
		},
		expectedOK: true,
	},
	{
		name:            "empty post body",
		postedData:      nil,
		expectedOK:      false,
		expectedMessage: "Internal Server Error",
	},
	{
		name: "database query fails",
		postedData: url.Values{
			"start_date": {"2060-01-01"},
			"end_date":   {"2060-01-02"},
			"room_id":    {"1"},
		},
		expectedOK:      false,
		expectedMessage: "Error querying database",
	},
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	for _, e := range testAvailabilityJSONData {
		// create request, get the context with session, set header, create recorder
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/search-availability-json", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// make our handler a http.HandlerFunc and call
		handler := http.HandlerFunc(Repo.AvailabilityJSON)
		handler.ServeHTTP(rr, req)

		var j jsonResponse
		err := json.Unmarshal([]byte(rr.Body.String()), &j)
		if err != nil {
			t.Error("failed to parse json!")
		}

		if j.OK != e.expectedOK {
			t.Errorf("%s: expected %v but got %v", e.name, e.expectedOK, j.OK)
		}
	}
}

var reservationSummaryTests = []struct {
	name               string
	reservation        models.Reservation
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "res-in-session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "General's Quarters",
			},
		},
		url:                "/reservation-summary",
		expectedStatusCode: http.StatusOK,
		expectedLocation:   "",
	},
	{
		name:               "res-not-in-session",
		reservation:        models.Reservation{},
		url:                "/reservation-summary",
		expectedStatusCode: http.StatusTemporaryRedirect,
		expectedLocation:   "/",
	},
}

func TestRepository_ReservationSummary(t *testing.T) {
	for _, e := range reservationSummaryTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ReservationSummary)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var chooseRoomTests = []struct {
	name               string
	reservation        models.Reservation
	url                string
	expectedStatusCode int
	expectedLocation   string
}{
	{
		name: "reservation-in-session",
		reservation: models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "General's Quarters",
			},
		},
		url:                "/choose-room/1",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/make-reservation",
	},
	{
		name:               "reservation-not-in-session",
		reservation:        models.Reservation{},
		url:                "/choose-room/1",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
	{
		name:               "malformed-url",
		reservation:        models.Reservation{},
		url:                "/choose-room/fish",
		expectedStatusCode: http.StatusSeeOther,
		expectedLocation:   "/",
	},
}

func TestRepository_ChooseRoom(t *testing.T) {
	for _, e := range chooseRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		// set the RequestURI on the request so that we can grab the ID from the URL
		req.RequestURI = e.url

		rr := httptest.NewRecorder()
		if e.reservation.RoomID > 0 {
			session.Put(ctx, "reservation", e.reservation)
		}

		handler := http.HandlerFunc(Repo.ChooseRoom)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}

		if e.expectedLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}
	}
}

var bookRoomTests = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{
		name:               "database-works",
		url:                "/book-room?s=2050-01-01&e=2050-01-02&id=1",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "database-fails",
		url:                "/book-room?s=2040-01-01&e=2040-01-02&id=4",
		expectedStatusCode: http.StatusTemporaryRedirect,
	},
}

func TestRepository_BookRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	for _, e := range bookRoomTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		session.Put(ctx, "reservation", reservation)

		handler := http.HandlerFunc(Repo.BookRoom)

		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s failed: returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, _ := session.Load(req.Context(), req.Header.Get("X-Session"))
	return ctx
}
