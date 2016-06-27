package main

import (
	util "github.com/woanware/goutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"html/template"
    "path"
    "strings"
)

// ##### Methods #############################################################

//
func routeIndex (c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
	})
}

func loadData (dataType int, c *gin.Context) {
	currentPageNumber := 0
	numRecsPerPage := 10

	temp :=  c.PostForm("num_recs_per_page")
	if len(temp) > 0 {
		if util.IsNumber(temp) == true {
			numRecsPerPage = util.ConvertStringToInt(temp)
		}
	}

	mode, hasMode := c.GetPostForm("mode")

	// Appears to be the first request to send the initial set of data
	if (mode != "first" &&
	    mode != "next" &&
	    mode != "previous") || hasMode == false {

        loadEventData(c, currentPageNumber, numRecsPerPage)
		return
	}

	temp =  c.PostForm("current_page_num")
	if len(temp) == 0 {
		currentPageNumber = 0
	}

	if util.IsNumber(temp) == false {
		currentPageNumber = 0
	}

	currentPageNumber = util.ConvertStringToInt(temp)

	if mode == "first"{
		currentPageNumber = 0
	}

	if mode == "next" {
		currentPageNumber += 1
	}

	if mode == "previous" {
		currentPageNumber -= 1
	}

	if currentPageNumber < 0 {
		currentPageNumber = 0
	}

    loadEventData(c, currentPageNumber, numRecsPerPage)
}

//
func routeEvents (c *gin.Context) {
	loadData(TYPE_EVENT, c)
}

//
func loadEventData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getEvents(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "events", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getEvents(numRecsPerPage int, currentPageNumber int) (bool, bool, []*Event) {
	var data []*Event

	err := db.
		Select("id, domain, host, utc_time, type, message_html").
		From("event").
		OrderBy("utc_time DESC").
		Offset(uint64(numRecsPerPage * currentPageNumber)).
		Limit(uint64(numRecsPerPage + 1)).
		QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
		v.Data = template.HTML(v.MessageHtml)
		v.UtcTimeStr = v.UtcTime.Format("15:04:05 02/01/2006")
	}

	noMoreRecords := false
	if len(data) < numRecsPerPage + 1 {
		noMoreRecords = true
	} else {
		// Remove the last item in the slice/array
		data = data[:len(data) - 1]
	}

	return false, noMoreRecords, data
}

//
func routeExport (c *gin.Context) {

    exportType := 0

    temp :=  c.PostForm("export_type")
    if len(temp) > 0 {
        if util.IsNumber(temp) == true {
            exportType = util.ConvertStringToInt(temp)
        }
    }

    if exportType == 0 {
        c.HTML(http.StatusOK, "export", gin.H{
            "has_data": false,
            "export_type": 0,
            "data": nil,
        })
        return
    }

    errored, data := getExports(exportType)
    if errored == true {
        c.String(http.StatusInternalServerError, "")
        return
    }

    hasData := true
    if len(data) == 0 {
        hasData = false
    }

    c.HTML(http.StatusOK, "export", gin.H{
        "has_data": hasData,
        "export_type": exportType,
        "data": data,
    })
}

//
func getExports(exportType int) (bool, []*Export) {

    var data []*Export

    err := db.
        Select(`*`).
        From("export").
        Where("data_type = $1", exportType).
        Limit(10).
        OrderBy("updated").
        QueryStructs(&data)

    if err != nil {
        logger.Errorf("Error querying for exports: %v (%d)", err, exportType)
        return true, data
    }

    // Perform some cleaning of the data, so that it displays better in the HTML
    for _, v := range data {
        v.OtherData = template.HTML(`<a href="/export/`+ util.ConvertInt64ToString(v.Id) + `">` + v.Updated.Format("15:04:05 02/01/2006") + `</a>`)
    }

    return false, data
}

//
func routeExportData(c *gin.Context) {

    id, successful := processInt64Parameter(c.Param("id"))
    if successful == false {
        c.String(http.StatusInternalServerError, "")
        return
    }

    if id < 1 {
        c.String(http.StatusInternalServerError, "")
        return
    }

    errored, export := getExport(id)

    if errored == true {
        c.String(http.StatusInternalServerError, "")
        return
    }

    // Load file contents
    if util.DoesFileExist(path.Join(config.ExportDir, export.FileName)) == false {
        logger.Errorf("Export file does not exist: %s", export.FileName)
        c.String(http.StatusInternalServerError, "")
        return
    }

    // Return file contents as download
    data, err := util.ReadTextFromFile(path.Join(config.ExportDir, export.FileName))
    if err != nil {
        logger.Errorf("Error reading export file: %v (%s)", err, export.FileName)
        c.String(http.StatusInternalServerError, "")
        return
    }

    c.Header("Content-Disposition", "attachment; filename=\"" + export.FileName)
    c.Data(http.StatusOK, "text/csv", []byte(data))
}

//
func getExport(id int64) (bool, Export) {

    var e Export

    err := db.
        Select(`id, data_type, file_name, updated`).
        From("export").
        Where("id = $1", id).
        QueryStruct(&e)

    if err != nil {
        logger.Errorf("Error querying for export: %v", err)
        return true, e
    }

    return false, e
}

//
func routeSearch (c *gin.Context) {

    currentPageNumber := 0

    numRecsPerPage, successful := processIntParameter(c.PostForm("num_recs_per_page"))
    if successful == false {
        numRecsPerPage = 10
    }

    mode, hasMode := c.GetPostForm("mode")

    // Appears to be the first request to send the initial set of data
    if (mode != "first" &&
        mode != "next" &&
        mode != "previous") || hasMode == false {

        loadSearchData(c, "", currentPageNumber, numRecsPerPage)
        return
    }

    searchValue :=  c.PostForm("search_value")
    if len(searchValue) == 0 {
        c.String(http.StatusInternalServerError, "")
        return
    }

    currentPageNumber = processCurrentPageNumber(c.PostForm("current_page_num"), mode)

    loadSearchData(c,  searchValue, currentPageNumber, numRecsPerPage)
}

//
func loadSearchData(
    c *gin.Context,
    searchValue string,
    currentPageNumber int,
    numRecsPerPage int) {

    if len(searchValue) == 0 {
        c.HTML(http.StatusOK, "search", gin.H{
            "current_page_num": currentPageNumber,
            "num_recs_per_page": numRecsPerPage,
            "no_more_records": true,
            "events": nil,
            "has_data": false,
            "search_value": searchValue,
        })
        return
    }

    errored, noMoreRecords, events := getSearch(searchValue, numRecsPerPage, currentPageNumber)
    if errored == true {
        c.String(http.StatusInternalServerError, "")
        return
    }

    hasData := true
    if len(events) == 0 {
        hasData = false
    }

    c.HTML(http.StatusOK, "search", gin.H{
        "current_page_num": currentPageNumber,
        "num_recs_per_page": numRecsPerPage,
        "no_more_records": noMoreRecords,
        "events": events,
        "has_data": hasData,
        "search_value": searchValue,
    })
}

//
func getSearch(searchValue string, numRecsPerPage int, currentPageNumber int) (bool, bool, []*Event) {

    var data []*Event

    err := db.
        Select(`id, domain, host, utc_time, type, message_html`).
        From("event").
        OrderBy("utc_time DESC").
        Offset(uint64(numRecsPerPage * currentPageNumber)).
        Limit(uint64(numRecsPerPage + 1)).
        Where("LOWER(message) LIKE $1", "%" + strings.ToLower(searchValue) + "%").
        QueryStructs(&data)

    if err != nil {
        logger.Errorf("Error querying for search: %v", err)
        return true, false, data
    }

    // Perform some cleaning of the data, so that it displays better in the HTML
    for _, v := range data {
        v.Data = template.HTML(v.MessageHtml)
        v.UtcTimeStr = v.UtcTime.Format("15:04:05 02/01/2006")
    }

    noMoreRecords := false
    if len(data) < numRecsPerPage + 1 {
        noMoreRecords = true
    } else {
        // Remove the last item in the slice/array
        data = data[:len(data) - 1]
    }

    return false, noMoreRecords, data
}


