package table

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yaoapp/gou/fs"
	"github.com/yaoapp/gou/model"
	gouProcess "github.com/yaoapp/gou/process"
	"github.com/yaoapp/gou/types"
	"github.com/yaoapp/kun/any"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/log"
	"github.com/yaoapp/kun/maps"
	"github.com/yaoapp/yao/helper"
	"github.com/yaoapp/yao/widgets/app"
)

// Export process

func exportProcess() {
	gouProcess.Register("yao.table.setting", processSetting)
	gouProcess.Register("yao.table.xgen", processXgen)
	gouProcess.Register("yao.table.component", processComponent)
	gouProcess.Register("yao.table.upload", processUpload)
	gouProcess.Register("yao.table.download", processDownload)
	gouProcess.Register("yao.table.search", processSearch)
	gouProcess.Register("yao.table.get", processGet)
	gouProcess.Register("yao.table.find", processFind)
	gouProcess.Register("yao.table.save", processSave)
	gouProcess.Register("yao.table.create", processCreate)
	gouProcess.Register("yao.table.insert", processInsert)
	gouProcess.Register("yao.table.update", processUpdate)
	gouProcess.Register("yao.table.updatewhere", processUpdateWhere)
	gouProcess.Register("yao.table.updatein", processUpdateIn)
	gouProcess.Register("yao.table.delete", processDelete)
	gouProcess.Register("yao.table.deletewhere", processDeleteWhere)
	gouProcess.Register("yao.table.deletein", processDeleteIn)
	gouProcess.Register("yao.table.export", processExport)
	gouProcess.Register("yao.table.load", processLoad)
}

func processXgen(process *gouProcess.Process) interface{} {

	tab := MustGet(process)
	data := process.ArgsMap(1, map[string]interface{}{})
	excludes := app.Permissions(process, "tables", tab.ID)
	setting, err := tab.Xgen(data, excludes)
	if err != nil {
		exception.New(err.Error(), 500).Throw()
	}
	return setting
}

func processDownload(process *gouProcess.Process) interface{} {

	process.ValidateArgNums(4)
	tab := MustGet(process)
	field := process.ArgsString(1)
	file := process.ArgsString(2)
	tokenString := process.ArgsString(3)

	// checking
	ext := fs.ExtName(file)
	if _, has := fs.DownloadWhitelist[ext]; !has {
		exception.New("%s.%s .%s file does not allow", 403, tab.ID, field, ext).Throw()
	}

	// Auth
	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
	if tokenString == "" {
		exception.New("%s.%s No permission", 403, tab.ID, field).Throw()
	}
	claims := helper.JwtValidate(tokenString)

	// Get Process name
	name := "fs.system.Download"
	if tab.Action.Download.Process != "" {
		name = tab.Action.Download.Process
	}

	// Create process
	p, err := gouProcess.Of(name, file)
	if err != nil {
		log.Error("[downalod] %s.%s %s", tab.ID, field, err.Error())
		exception.New("[downalod] %s.%s %s", 400, tab.ID, field, err.Error()).Throw()
	}

	// Excute process
	res, err := p.WithGlobal(process.Global).WithSID(claims.SID).Exec()
	if err != nil {
		log.Error("[downalod] %s.%s %s", tab.ID, field, err.Error())
		exception.New("[downalod] %s.%s %s", 500, tab.ID, field, err.Error()).Throw()
	}

	return res
}

func processUpload(process *gouProcess.Process) interface{} {

	process.ValidateArgNums(4)
	tab := MustGet(process)
	xpath := process.ArgsString(1)
	method := process.ArgsString(2)
	key := fmt.Sprintf("%s.$%s", xpath, method)

	// get cloud props
	cProp, has := tab.CProps[key]
	if !has {
		exception.New("%s does not exist", 400, key).Throw()
	}

	// $file.file
	tmpfile, ok := process.Args[3].(types.UploadFile)
	if !ok {
		exception.New("parameters error: %v", 400, process.Args[3]).Throw()
	}

	// execute upload
	res, err := cProp.ExecUpload(process, tmpfile)
	if err != nil {
		exception.New(err.Error(), 500).Throw()
	}

	if file, ok := res.(string); ok {
		field := strings.TrimSuffix(xpath, ".edit.props")
		file = fmt.Sprintf("/api/__yao/table/%s/download/%s?name=%s", tab.ID, url.QueryEscape(field), file)
		return file
	}

	return res
}

func processComponent(process *gouProcess.Process) interface{} {

	process.ValidateArgNums(3)
	tab := MustGet(process)
	xpath := process.ArgsString(1)
	method := process.ArgsString(2)
	key := fmt.Sprintf("%s.$%s", xpath, method)

	// get cloud props
	cProp, has := tab.CProps[key]
	if !has {
		exception.New("%s does not exist", 400, key).Throw()
	}

	// :query
	query := map[string]interface{}{}
	if process.NumOfArgsIs(4) {
		query = process.ArgsMap(3)
	}

	// execute query
	res, err := cProp.ExecQuery(process, query)
	if err != nil {
		exception.New(err.Error(), 500).Throw()
	}

	return res
}

func processSetting(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	process.Args = append(process.Args, process.Args[0]) // table name
	return tab.Action.Setting.MustExec(process)
}

func processSearch(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.Search.MustExec(process)
}

func processGet(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.Get.MustExec(process)
}

func processSave(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.Save.MustExec(process)
}

func processCreate(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.Create.MustExec(process)
}

func processFind(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.Find.MustExec(process)
}

func processInsert(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.Insert.MustExec(process)
}

func processUpdate(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.Update.MustExec(process)
}

func processUpdateWhere(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.UpdateWhere.MustExec(process)
}

func processUpdateIn(process *gouProcess.Process) interface{} {
	process.ValidateArgNums(3)
	tab := MustGet(process)
	ids := strings.Split(process.ArgsString(1), ",")
	process.Args[1] = model.QueryParam{
		Wheres: []model.QueryWhere{
			{Column: tab.Layout.Primary, OP: "in", Value: ids},
		},
	}
	return tab.Action.UpdateIn.MustExec(process)
}

func processDelete(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.Delete.MustExec(process)
}

func processDeleteWhere(process *gouProcess.Process) interface{} {
	tab := MustGet(process)
	return tab.Action.DeleteWhere.MustExec(process)
}

func processDeleteIn(process *gouProcess.Process) interface{} {
	process.ValidateArgNums(2)
	tab := MustGet(process)
	ids := strings.Split(process.ArgsString(1), ",")
	process.Args[1] = model.QueryParam{
		Wheres: []model.QueryWhere{
			{Column: tab.Layout.Primary, OP: "in", Value: ids},
		},
	}
	return tab.Action.DeleteIn.MustExec(process)
}

// processExport yao.table.Export (:table, :queryParam, :chunkSize)
func processExport(process *gouProcess.Process) interface{} {
	process.ValidateArgNums(1)
	tab := MustGet(process) // 0
	params := process.ArgsQueryParams(1, types.QueryParam{})
	pagesize := process.ArgsInt(2, 50)

	// Filename
	hash := md5.Sum([]byte(time.Now().Format("20060102-15:04:05")))
	fingerprint := string(hex.EncodeToString(hash[:]))
	fingerprint = strings.ToUpper(fingerprint)
	dir := time.Now().Format("20060102")
	filename := filepath.Join(string(os.PathSeparator), dir, fmt.Sprintf("%s.xlsx", fingerprint))

	// Create Data Path
	fs := fs.MustGet("system")
	if has, _ := fs.Exists(dir); !has {
		fs.MkdirAll(dir, uint32(os.ModePerm))
	}

	// Query
	page := 1
	for page > 0 {
		process.Args = []interface{}{tab.ID, params, page, pagesize}
		data, err := tab.Action.Search.Exec(process)
		if err != nil {
			log.Error("[table] export error %s", err.Error())
			page = -1
			continue
		}

		res, ok := data.(map[string]interface{})
		if !ok {
			res, ok = data.(maps.MapStrAny)
			if !ok {
				log.Error("[table] export Search Action response data error %#v", data)
				page = -1
				continue
			}
		}

		if _, ok := res["next"]; !ok {
			page = -1
			continue
		}

		size := pagesize
		if _, ok := res["pagesize"]; ok {
			size = any.Of(res["pagesize"]).CInt()
		}

		// Export
		err = tab.Export(filename, res["data"], page, size)
		if err != nil {
			log.Error("Export %s %s", tab.ID, err.Error())
		}

		page = any.Of(res["next"]).CInt()
	}

	return filename
}

// processLoad yao.table.Load (:file)
func processLoad(process *gouProcess.Process) interface{} {
	process.ValidateArgNums(1)
	file := process.ArgsString(0)
	if file == "" {
		exception.New("file is required", 400).Throw()
	}

	file = strings.TrimPrefix(file, string(os.PathSeparator))
	return LoadFile("tables", file)
}
