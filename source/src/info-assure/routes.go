package main

import (
	util "github.com/woanware/goutil"
	"github.com/gin-gonic/gin"
	"net/http"
	//"strings"
	"html/template"
	"fmt"
    "path"
)

// ##### Methods #############################################################

//
func index (c *gin.Context) {
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

		switch dataType {
		case TYPE_EVENT:
			loadEventData(c, currentPageNumber, numRecsPerPage)
		case TYPE_PROCESS_CREATE:
			loadProcessCreateData(c, currentPageNumber, numRecsPerPage)
		case TYPE_PROCESS_TERMINATE:
			loadProcessTerminateData(c, currentPageNumber, numRecsPerPage)
		case TYPE_FILE_CREATION_TIME:
			loadFileCreationTimeData(c, currentPageNumber, numRecsPerPage)
		case TYPE_NETWORK_CONNECTION:
			loadNetworkConnectionData(c, currentPageNumber, numRecsPerPage)
		case TYPE_DRIVER_LOADED:
			loadDriverLoadedData(c, currentPageNumber, numRecsPerPage)
		case TYPE_IMAGE_LOADED:
			loadImageLoadedData(c, currentPageNumber, numRecsPerPage)
		case TYPE_RAW_ACCESS_READ:
			loadRawAccessReadData(c, currentPageNumber, numRecsPerPage)
		case TYPE_CREATE_REMOTE_THREAD:
			loadCreateRemoteThreadData(c, currentPageNumber, numRecsPerPage)
		}

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

	switch dataType {
	case TYPE_EVENT:
		loadEventData(c, currentPageNumber, numRecsPerPage)
	case TYPE_PROCESS_CREATE:
		loadProcessCreateData(c, currentPageNumber, numRecsPerPage)
	case TYPE_PROCESS_TERMINATE:
		loadProcessTerminateData(c, currentPageNumber, numRecsPerPage)
	case TYPE_FILE_CREATION_TIME:
		loadFileCreationTimeData(c, currentPageNumber, numRecsPerPage)
	case TYPE_NETWORK_CONNECTION:
		loadNetworkConnectionData(c, currentPageNumber, numRecsPerPage)
	case TYPE_DRIVER_LOADED:
		loadDriverLoadedData(c, currentPageNumber, numRecsPerPage)
	case TYPE_IMAGE_LOADED:
		loadImageLoadedData(c, currentPageNumber, numRecsPerPage)
	case TYPE_RAW_ACCESS_READ:
		loadRawAccessReadData(c, currentPageNumber, numRecsPerPage)
	case TYPE_CREATE_REMOTE_THREAD:
		loadCreateRemoteThreadData(c, currentPageNumber, numRecsPerPage)
	}
}

//
func events (c *gin.Context) {
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
		Select("id, domain, host, utc_time, type, message").
		From("unified").
		OrderBy("utc_time DESC").
		Offset(uint64(numRecsPerPage * currentPageNumber)).
		Limit(uint64(numRecsPerPage + 1)).
		QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
		v.MessageHtml = template.HTML(v.Message)
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
func processCreate (c *gin.Context) {
	loadData(TYPE_PROCESS_CREATE, c)
}

//
func loadProcessCreateData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getProcessCreate(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "process_create", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getProcessCreate(numRecsPerPage int, currentPageNumber int) (bool, bool, []*ProcessCreate) {
	var data []*ProcessCreate

	err := db.
	Select("id, domain, host, utc_time, process_id, image, command_line, sha256, md5, parent_process_id, parent_image, parent_command_line, process_user").
	From("process_create").
	OrderBy("utc_time DESC").
	Offset(uint64(numRecsPerPage * currentPageNumber)).
	Limit(uint64(numRecsPerPage + 1)).
	QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
		v.OtherData = template.HTML(fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Process User:</strong> %s<br><strong>Image:</strong> %s<br><strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s<br><strong>Parent Process ID:</strong> %d<br><strong>Parent Image:</strong> %s<br><strong>Parent Command Line:</strong> %s<br>`,
			v.ProcessId, v.ProcessUser, v.Image, v.Md5, v.Sha256, v.ParentProcessId, v.ParentImage, v.ParentCommandLine))
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
func processTerminate (c *gin.Context) {
	loadData(TYPE_PROCESS_TERMINATE, c)
}

//
func loadProcessTerminateData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getProcessTerminate(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "process_terminate", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getProcessTerminate(numRecsPerPage int, currentPageNumber int) (bool, bool, []*ProcessTerminate) {
	var data []*ProcessTerminate

	err := db.
	Select("id, domain, host, utc_time, process_id, image").
	From("process_terminate").
	OrderBy("utc_time DESC").
	Offset(uint64(numRecsPerPage * currentPageNumber)).
	Limit(uint64(numRecsPerPage + 1)).
	QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
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
func fileCreationTime (c *gin.Context) {
	loadData(TYPE_FILE_CREATION_TIME, c)
}

//
func loadFileCreationTimeData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getFileCreationTimeData(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "file_creation_time", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getFileCreationTimeData(numRecsPerPage int, currentPageNumber int) (bool, bool, []*FileCreationTime) {
	var data []*FileCreationTime

	err := db.
	Select("id, domain, host, utc_time, process_id, image, target_file_name, creation_utc_time, previous_creation_utc_time").
	From("file_creation_time").
	OrderBy("utc_time DESC").
	Offset(uint64(numRecsPerPage * currentPageNumber)).
	Limit(uint64(numRecsPerPage + 1)).
	QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
		v.UtcTimeStr = v.UtcTime.Format("15:04:05 02/01/2006")
		v.OtherData = template.HTML(fmt.Sprintf(`<strong>Target File Name:</strong> %s<br><strong>Creation UTC Time:</strong> %s<br><strong>Previous Creation UTC Time:</strong> %s<br>`,
									v.TargetFileName, v.CreationUtcTime.Format("15:04:05 02/01/2006"), v.PreviousCreationUtcTime.Format("15:04:05 02/01/2006")))
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
func networkConnection (c *gin.Context) {
	loadData(TYPE_NETWORK_CONNECTION, c)
}

//
func loadNetworkConnectionData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getNetworkConnectionData(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "network_connection", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getNetworkConnectionData(numRecsPerPage int, currentPageNumber int) (bool, bool, []*NetworkConnection) {
	var data []*NetworkConnection

	err := db.
        Select("id, domain, host, utc_time, process_id, image, process_user, protocol, initiated, source_ip, source_port, destination_ip, destination_port").
        From("network_connection").
        OrderBy("utc_time DESC").
        Offset(uint64(numRecsPerPage * currentPageNumber)).
        Limit(uint64(numRecsPerPage + 1)).
        QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
		v.UtcTimeStr = v.UtcTime.Format("15:04:05 02/01/2006")
		v.OtherData = template.HTML(fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Initiated:</strong> %v<br>`,
			v.ProcessId, v.Image, v.Initiated))
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
func loadHttpRouteDriverLoaded (c *gin.Context) {
	loadData(TYPE_DRIVER_LOADED, c)
}

//
func loadDriverLoadedData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getDriverLoadedData(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "driver_loaded", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getDriverLoadedData(numRecsPerPage int, currentPageNumber int) (bool, bool, []*DriverLoaded) {
	var data []*DriverLoaded

	err := db.
        Select("id, domain, host, utc_time, image_loaded, md5, sha256, signed, signature").
        From("driver_loaded").
        OrderBy("utc_time DESC").
        Offset(uint64(numRecsPerPage * currentPageNumber)).
        Limit(uint64(numRecsPerPage + 1)).
        QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
		v.UtcTimeStr = v.UtcTime.Format("15:04:05 02/01/2006")
		v.OtherData = template.HTML(fmt.Sprintf(`<strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s`,
			v.Md5, v.Sha256))
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
func loadHttpRouteImageLoaded (c *gin.Context) {
	loadData(TYPE_IMAGE_LOADED, c)
}

//
func loadImageLoadedData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getImageLoadedData(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "image_loaded", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getImageLoadedData(numRecsPerPage int, currentPageNumber int) (bool, bool, []*ImageLoaded) {
	var data []*ImageLoaded

	err := db.
        Select("id, domain, host, utc_time, image_loaded, md5, sha256, signed, signature").
        From("image_loaded").
        OrderBy("utc_time DESC").
        Offset(uint64(numRecsPerPage * currentPageNumber)).
        Limit(uint64(numRecsPerPage + 1)).
        QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
		v.UtcTimeStr = v.UtcTime.Format("15:04:05 02/01/2006")
		v.OtherData = template.HTML(fmt.Sprintf(`<strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s`,
			v.Md5, v.Sha256))
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
func loadHttpRouteRawAccessRead (c *gin.Context) {
	loadData(TYPE_RAW_ACCESS_READ, c)
}

//
func loadRawAccessReadData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getRawAccessReadData(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "raw_access_read", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getRawAccessReadData(numRecsPerPage int, currentPageNumber int) (bool, bool, []*RawAccessRead) {
	var data []*RawAccessRead

	err := db.
        Select("id, domain, host, utc_time, process_id, image, device").
        From("raw_access_read").
        OrderBy("utc_time DESC").
        Offset(uint64(numRecsPerPage * currentPageNumber)).
        Limit(uint64(numRecsPerPage + 1)).
        QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
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
func loadHttpRouteCreateRemoteThread (c *gin.Context) {
	loadData(TYPE_CREATE_REMOTE_THREAD, c)
}

//
func loadCreateRemoteThreadData(
	c *gin.Context,
	currentPageNumber int,
	numRecsPerPage int) {

	errored, noMoreRecords, events := getCreateRemoteThreadData(numRecsPerPage, currentPageNumber)
	if errored == true {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.HTML(http.StatusOK, "create_remote_thread", gin.H{
		"current_page_num": currentPageNumber,
		"num_recs_per_page": numRecsPerPage,
		"no_more_records": noMoreRecords,
		"events": events,
	})
}

func getCreateRemoteThreadData(numRecsPerPage int, currentPageNumber int) (bool, bool, []*CreateRemoteThread) {
	var data []*CreateRemoteThread

	err := db.
	Select("id, domain, host, utc_time, source_process_id, source_image, target_process_id, target_image, new_thread_id, start_address, start_module, start_function").
	From("create_remote_thread").
	OrderBy("utc_time DESC").
	Offset(uint64(numRecsPerPage * currentPageNumber)).
	Limit(uint64(numRecsPerPage + 1)).
	QueryStructs(&data)

	if err != nil {
		logger.Error(err)
	}

	// Perform some cleaning of the data, so that it displays better in the HTML
	for _, v := range data {
		v.UtcTimeStr = v.UtcTime.Format("15:04:05 02/01/2006")
		v.OtherData = template.HTML(fmt.Sprintf(`<strong>Target Process ID:</strong> %d<br><strong>Target Image:</strong> %s<br><strong>New Thread ID:</strong> %d<br><strong>Start Address:</strong> %s<br><strong>Start Module:</strong> %s<br><strong>Start Function:</strong> %s`,
			v.TargetProcessId, v.TargetImage, v.NewThreadId, v.StartAddress, v.StartModule, v.StartFunction))
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
func export (c *gin.Context) {

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
func exportData(c *gin.Context) {

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


