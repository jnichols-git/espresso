package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jakenichols2719/simpleblog/pkg/listings"
	"github.com/jakenichols2719/simpleblog/pkg/posts"
	"golang.org/x/exp/slices"
)

func HandleGetListing(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	pID := params.Get("post_id")
	if pID == "" {
		HandleError(w, fmt.Errorf("get request must include post_id"))
		return
	}
	// Connect to AWS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		HandleError(w, err)
		return
	}
	// Read listings
	l, err := listings.ReadOne(pID, listings.Connect(cfg))
	if err != nil {
		HandleError(w, err)
		return
	}
	// Write out
	raw, err := json.Marshal(l)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Write(raw)
}

func HandleGetListings(w http.ResponseWriter, req *http.Request) {
	// Get request params
	params := req.URL.Query()
	pc, ok := strconv.Atoi(params.Get("page_count"))
	if ok != nil {
		HandleError(w, ok)
		return
	}
	pn, ok := strconv.Atoi(params.Get("page_number"))
	if ok != nil {
		pn = 0
	}
	// compile filter
	// tag match (optional)
	tags := params["tags"]

	filter := func(listing *listings.Listing) bool {
		match := true
		if len(tags) != 0 {
			for _, tag := range tags {
				if slices.Contains(listing.Tags, tag) {
					goto tagFound
				}
			}
			match = false
		}
	tagFound:
		return match
	}
	// Connect to AWS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		HandleError(w, err)
		return
	}
	// Read listings
	ls, err := listings.ReadMany(pc, pn, filter, listings.Connect(cfg))
	if err != nil {
		HandleError(w, err)
		return
	}
	// Write out
	raw, err := json.Marshal(ls)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Write(raw)
}

func HandleGetPost(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	pID := params.Get("post_id")
	vID := params.Get("version_id")
	if pID == "" {
		HandleError(w, fmt.Errorf("get request must include post_id"))
		return
	}
	if vID == "" {
		HandleError(w, fmt.Errorf("get request must include version_id"))
		return
	}
	// Connect to AWS
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		HandleError(w, err)
		return
	}

	// if vID == "latest", get the listing for pID and swap out vID for the live version id
	if vID == "latest" {
		l, err := listings.ReadOne(pID, listings.Connect(cfg))
		if err != nil {
			HandleError(w, err)
			return
		}
		vID = l.LiveVersionID
	}
	// Read post
	p, err := posts.Read(pID, vID, posts.Connect(cfg))
	if err != nil {
		HandleError(w, err)
		return
	}
	// Write out
	raw, err := json.Marshal(p.Content)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Write(raw)
}
