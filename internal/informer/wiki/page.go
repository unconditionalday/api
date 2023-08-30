package wiki

import (
	"errors"
	"maps"
	"slices"

	mapsx "github.com/unconditionalday/server/internal/x/maps"
	netx "github.com/unconditionalday/server/internal/x/net"

	"reflect"
	"strconv"
	"strings"
)

// Result after we parse the response of Wikipedia API.
// Some attributes must be get manually using the WikipediaPage methods
type WikipediaPage struct {
	PageID         int              `json:"pageid"`
	Title          string           `json:"title"`
	OriginalTitle  string           `json:"originaltitle"`
	Content        string           `json:"content"`
	HTML           string           `json:"html"`
	URL            string           `json:"fullurl"`
	RevisionID     float64          `json:"revid"`
	ParentID       float64          `json:"parentid"`
	Summary        string           `json:"summary"`
	CheckedImage   bool             `json:"checkedimage"`
	CheckedThumb   bool             `json:"checkedthumb"`
	Thumbnail      string           `json:"thumbnail"`
	Images         []string         `json:"images"`
	Coordinate     []float64        `json:"coordinates"`
	Language       string           `json:"lang"`
	Reference      []string         `json:"references"`
	Link           []string         `json:"links"`
	Category       []string         `json:"categories"`
	Section        []string         `json:"sections"`
	SectionOffset  map[string][]int `json:"sectionoffset"`
	Disambiguation []string         `json:"disambiguation"`
}

/*
Return true if the 2 pages are the same
*/
func (page WikipediaPage) Equal(other WikipediaPage) bool {
	return page.PageID == other.PageID
}

/*
Get the string content of the page. Save it into the page.Content for later use
*/
func (page *WikipediaPage) GetContent(client Client, lang string) (string, error) {
	if page.Content != "" {
		return page.Content, nil
	}
	pageid := strconv.Itoa(page.PageID)
	args := map[string]string{
		"action":      "query",
		"prop":        "extracts|revisions",
		"explaintext": "",
		"rvprop":      "ids",
		"titles":      page.Title,
	}
	res, err := client.RequestWikiApi(args, lang)
	if err != nil {
		return "", err
	}
	if res.Error.Code != "" {
		return "", errors.New(res.Error.Info)
	}
	page.Content = res.Query.Page[pageid].Extract
	page.RevisionID = res.Query.Page[pageid].Revision[0]["revid"].(float64)
	page.ParentID = res.Query.Page[pageid].Revision[0]["parentid"].(float64)

	return page.Content, nil
}

/*
Get the html of the page. Save it into the page.HTML for later use\

**Warning:: This can get pretty slow on long pages.
*/
func (page *WikipediaPage) GetHTML(client Client, lang string) (string, error) {
	if page.HTML != "" {
		return page.HTML, nil
	}
	args := map[string]string{
		"action":  "query",
		"prop":    "revisions",
		"rvprop":  "content",
		"rvlimit": strconv.Itoa(1),
		"rvparse": "",
		"titles":  page.Title,
	}
	res, err := client.RequestWikiApi(args, lang)
	if err != nil {
		return "", err
	}
	if res.Error.Code != "" {
		return "", errors.New(res.Error.Info)
	}
	page.HTML = res.Query.Page[strconv.Itoa(page.PageID)].Revision[0]["*"].(string)
	return page.HTML, nil
}

/*
Get the revid of the page. Save it into the page.HTML for later use

The revision ID is a number that uniquely identifies the current version of the page.
It can be used to create the permalink or for other direct API calls. See Help:Page history <http://en.wikipedia.org/wiki/Wikipedia:Revision>
for more information.
*/
func (page *WikipediaPage) GetRevisionID(client Client, lang string) (float64, error) {
	if page.RevisionID != 0 {
		return page.RevisionID, nil
	}
	_, err := page.GetContent(client, lang)
	if err != nil {
		return -1, err
	}
	return page.RevisionID, nil
}

/*
Revision ID of the parent version of the current revision of this page.

See “revision_id“ for more information.
*/
func (page *WikipediaPage) GetParentID(client Client, lang string) (float64, error) {
	if page.RevisionID != 0 {
		return page.ParentID, nil
	}
	_, err := page.GetContent(client, lang)
	if err != nil {
		return -1, err
	}
	return page.ParentID, nil
}

/*
String summary of a page
*/
func (page *WikipediaPage) GetSummary(client Client, lang string) (string, error) {
	if page.Summary != "" {
		return page.Summary, nil
	}

	pageid := strconv.Itoa(page.PageID)
	args := map[string]string{
		"action":      "query",
		"prop":        "extracts",
		"explaintext": "",
		"exintro":     "",
		"exsentences": "3",
		"exlimit":     "1",
		"titles":      page.Title,
	}
	res, err := client.RequestWikiApi(args, lang)
	if err != nil {
		return "", err
	}

	page.Summary = res.Query.Page[pageid].Extract
	return page.Summary, nil
}

/*
Based on <https://www.mediawiki.org/wiki/API:Query#Continuing_queries>
*/
func (page *WikipediaPage) ContinuedQuery(args map[string]string, client Client, lang string) ([]interface{}, error) {
	// args["pageids"] = strconv.Itoa(page.PageID)
	args["titles"] = page.Title
	last := map[string]interface{}{}
	prop := args["prop"]
	result := make([]interface{}, 0, 7)
	for {
		new_args := maps.Clone(args)
		mapsx.Update(new_args, last)

		res, err := client.RequestWikiApi(args, lang)
		if err != nil {
			return result, err
		}
		if res.Error.Code != "" {
			return result, errors.New(res.Error.Info)
		}

		if reflect.DeepEqual(RequestQuery{}, res.Query) {
			break
		}

		if _, ok := args["generator"]; ok {
			for _, v := range res.Query.Page {
				result = append(result, v)
			}
		} else {
			if prop == "extlinks" {
				temp := res.Query.Page[strconv.Itoa(page.PageID)].Extlink
				for _, v := range temp {
					result = append(result, v["*"])
				}
			} else {
				temp := []map[string]interface{}{}
				switch prop {
				case "links":
					temp = res.Query.Page[strconv.Itoa(page.PageID)].Link
				case "categories":
					temp = res.Query.Page[strconv.Itoa(page.PageID)].Category
				}
				for _, v := range temp {
					result = append(result, v["title"].(string))
				}
			}

		}

		if len(res.Continue) == 0 {
			break
		}

		last = res.Continue
	}
	return result, nil
}

/*
List of URLs of images on the page.
*/
func (page *WikipediaPage) GetImagesURL(client Client, lang string) ([]string, error) {
	if page.CheckedImage {
		return page.Images, nil
	}
	args := map[string]string{
		"action":    "query",
		"generator": "images",
		"gimlimit":  "max",
		"prop":      "imageinfo",
		"iiprop":    "url",
	}

	res, err := page.ContinuedQuery(args, client, lang)
	if err != nil && len(res) == 0 {
		return []string{}, err
	}
	result := make([]string, 0, 7)
	for _, v := range res {
		temp := v.(InnerPage).ImageInfo
		if len(temp) > 0 {
			result = append(result, temp[0]["url"])
		}
	}
	page.CheckedImage = true
	page.Images = result
	return page.Images, nil
}

/*
Get Thumbnail URL of the page
*/
func (page *WikipediaPage) GetThumbURL(client Client, lang string) (string, error) {
	if page.CheckedThumb {
		return page.Thumbnail, nil
	}

	args := map[string]string{
		"action":      "query",
		"prop":        "pageimages",
		"titles":      page.Title,
		"piprop":      "thumbnail",
		"pithumbsize": "500",
	}

	res, err := client.RequestWikiApi(args, lang)
	if err != nil {
		return "", err
	}

	page.CheckedThumb = true
	page.Thumbnail = res.Query.Page[strconv.Itoa(page.PageID)].Thumbnail.Source

	return page.Thumbnail, nil
}

/*
Slice of float64 in the form of (lat, lon)
*/
func (page *WikipediaPage) GetCoordinate(client Client, lang string) ([]float64, error) {
	if len(page.Coordinate) == 2 {
		return page.Coordinate, nil
	}
	args := map[string]string{
		"action":  "query",
		"prop":    "coordinates",
		"colimit": "max",
		"titles":  page.Title,
	}

	res, err := client.RequestWikiApi(args, lang)
	if err != nil {
		return []float64{}, err
	}
	if res.Error.Code != "" {
		return []float64{}, errors.New(res.Error.Info)
	}

	if reflect.DeepEqual(RequestQuery{}, res.Query) {
		page.Coordinate = []float64{-1, -1}
		return page.Coordinate, nil
	} else {
		temp := res.Query.Page[strconv.Itoa(page.PageID)].Coordinate[0]
		page.Coordinate = []float64{temp["lat"].(float64), temp["lon"].(float64)}
	}
	return page.Coordinate, nil
}

/*
		List of URLs of external links on a page.
	    May include external links within page that aren't technically cited anywhere.
*/
func (page *WikipediaPage) GetReference(client Client, lang string) ([]string, error) {
	if len(page.Reference) > 0 {
		return page.Reference, nil
	}
	args := map[string]string{
		"action":  "query",
		"prop":    "extlinks",
		"ellimit": "max",
	}
	res, err := page.ContinuedQuery(args, client, lang)
	if err != nil && len(res) == 0 {
		return []string{}, err
	}
	for _, v := range res {
		page.Reference = append(page.Reference, netx.HelpAddURL(v.(string)))
	}

	return page.Reference, nil
}

/*
		List of titles of Wikipedia page links on a page.
	    **Note:: Only includes articles from namespace 0, meaning no Category, User talk, or other meta-Wikipedia pages.
*/
func (page *WikipediaPage) GetLink(client Client, lang string) ([]string, error) {
	if len(page.Link) > 0 {
		return page.Link, nil
	}
	args := map[string]string{
		"action":      "query",
		"prop":        "links",
		"plnamespace": "0",
		"pllimit":     "max",
	}
	res, err := page.ContinuedQuery(args, client, lang)
	if err != nil && len(res) == 0 {
		return []string{}, err
	}
	for _, v := range res {
		page.Link = append(page.Link, v.(string))
	}
	return page.Link, nil
}

/*
List of categories of a page.
*/
func (page *WikipediaPage) GetCategory(client Client, lang string) ([]string, error) {
	if len(page.Category) > 0 {
		return page.Category, nil
	}
	args := map[string]string{
		"action":  "query",
		"prop":    "categories",
		"cllimit": "max",
	}
	res, err := page.ContinuedQuery(args, client, lang)
	if err != nil && len(res) == 0 {
		return []string{}, err
	}
	for _, v := range res {
		page.Category = append(page.Category, strings.Replace(v.(string), "Category:", "", 1))
	}
	return page.Category, nil
}

/*
List of section titles from the table of contents on the page.
*/
func (page *WikipediaPage) GetSectionList(client Client, lang string) ([]string, error) {
	if len(page.Section) > 0 {
		return page.Section, nil
	}
	args := map[string]string{
		"action": "parse",
		"prop":   "sections",
	}
	if page.Title != "" {
		args["page"] = page.Title
	}
	res, err := client.RequestWikiApi(args, lang)
	if err != nil {
		return []string{}, err
	}
	if res.Error.Code != "" {
		return []string{}, errors.New(res.Error.Info)
	}
	for _, v := range res.Parse["sections"].([]interface{}) {
		page.Section = append(page.Section, v.(map[string]interface{})["line"].(string))
	}
	return page.Section, nil
}

func (page *WikipediaPage) GetSection(section string, client Client, lang string) (string, error) {
	sections, err := page.GetSectionList(client, lang)
	if err != nil {
		return "", err
	}
	if !slices.Contains(sections, section) {
		return "", errors.New("section not exist")
	}
	content, err := page.GetContent(client, lang)
	if err != nil {
		return "", err
	}
	if page.SectionOffset == nil {
		page.SectionOffset = map[string][]int{}
	}
	if value, ok := page.SectionOffset[section]; ok {
		return content[value[0]:value[1]], nil
	}
	sectiontitle := "== " + section + " =="
	start := strings.Index(content, sectiontitle) + len(sectiontitle)
	// If you cannot find the section in the content (but it's there in the API for some reason)
	if start < len(sectiontitle) {
		page.SectionOffset[section] = []int{0, 0}
		return "", nil
	}
	end := start + strings.Index(content[start:], "==")
	if end == -1 {
		page.SectionOffset[section] = []int{start, len(content)}
		return content[start:], nil
	}
	page.SectionOffset[section] = []int{start, end}
	return strings.TrimSpace(strings.TrimLeft(content[start:end], "=")), nil
}

/*
		Load basic information from Wikipedia.

	    Confirm that page exists. If it's a disambiguation page, get a list of suggesting
*/
func MakeWikipediaPage(pageid int, title string, originaltitle string, redirect bool, client Client, lang string) (WikipediaPage, error) {
	page := WikipediaPage{}
	args := map[string]string{
		"action":    "query",
		"prop":      "info|pageprops",
		"inprop":    "url",
		"ppprop":    "disambiguation",
		"redirects": "",
	}
	page.Title = title
	page.OriginalTitle = title
	if pageid != -1 {
		args["pageids"] = strconv.Itoa(pageid)
		page.PageID = pageid
	} else {
		args["titles"] = title
	}
	if originaltitle != "" {
		page.OriginalTitle = originaltitle
	}
	res, err := client.RequestWikiApi(args, lang)
	if err != nil {
		return page, err
	}

	target := InnerPage{}
	target.Missing = "false"
	var index string
	for i, v := range res.Query.Page {
		index = i
		target = v
		break
	}

	if pageid == -1 {
		page.PageID = target.PageID
	}
	if title == "" {
		page.Title = target.Title
		page.OriginalTitle = target.Title
	}

	page.Language = target.PageLanguage

	if target.Missing == "" && index == "-1" {
		return page, errors.New("missing")
	}
	// if field redirects exist
	if len(res.Query.Redirect) > 0 {
		if !redirect {
			return page, errors.New("set the redirect argument to true to allow automatic redirects")
		}
		tempstr := page.Title
		if len(res.Query.Normalize) > 0 {
			if res.Query.Normalize[0].From != page.Title {
				return page, errors.New("an unexpected weird error, report me if it happened")
			}
			tempstr = res.Query.Normalize[0].To
		}
		if tempstr != res.Query.Redirect[0].From {
			return page, errors.New("an unexpected weird error, report me if it happened")
		}
		return MakeWikipediaPage(-1, res.Query.Redirect[0].To, "", redirect, client, lang)
	}

	// If the page is a disambiguation page
	if _, ok := target.PageProps["disambiguation"]; ok {
		return WikipediaPage{}, errors.New("disambiguation")
		// args = map[string]string{
		// 	"action":  "query",
		// 	"prop":    "revisions",
		// 	"rvprop":  "content",
		// 	"rvparse": "",
		// 	"rvlimit": strconv.Itoa(1),
		// 	"titles":  page.Title,
		// }
		// res, err := client.RequestWikiApi(args, lang)
		// if err != nil {
		// 	return page, err
		// }

		// html := res.Query.Page[strconv.Itoa(page.PageID)].Revision[0]["*"].(string)
		// doc := soup.HTMLParse(html)
		// links := doc.FindAll("li")
		// disa := make([]string, 0, 10)
		// for _, link := range links {
		// 	li := link.FindAll("a")
		// 	for _, l := range li {
		// 		if ref, ok := l.Attrs()["title"]; ok {
		// 			if len(ref) >= 1 && !slices.Contains(disa, ref) {
		// 				disa = append(disa, ref)
		// 			}
		// 		}
		// 	}
		// }
		// page.Disambiguation = disa
		// return page, nil
	}

	page.URL = target.FullURL

	return page, nil
}
