package services

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services/fileengine"
	"github.com/juggleim/jugglechat-server/storages/dbs"
	"github.com/juggleim/jugglechat-server/utils"
	"github.com/juggleim/jugglechat-server/utils/caches"
)

func GetFileCred(ctx context.Context, req *apimodels.QryFileCredReq) (errs.IMErrorCode, *apimodels.QryFileCredResp) {
	appkey := GetAppKeyFromCtx(ctx)
	fileConf := GetFileConf(ctx, appkey)
	if fileConf == nil || fileConf == notExistFileConf {
		return errs.IMErrorCode_APP_FILE_NOOSS, nil
	}
	dir := fileTypeToDir(req.FileType)
	switch fileConf.FileEngine {
	case fileengine.ChannelQiNiu:
		if fileConf.QiNiu == nil {
			return errs.IMErrorCode_APP_FILE_NOOSS, nil
		}
		uploadToken, domain := fileConf.QiNiu.UploadToken(req.Ext)
		return errs.IMErrorCode_SUCCESS, &apimodels.QryFileCredResp{
			OssType: apimodels.OssType_QiNiu,
			QiNiuCredResp: &apimodels.QiNiuCredResp{
				Domain: domain,
				Token:  uploadToken,
			},
		}
	case fileengine.ChannelMinio:
		if fileConf.Minio == nil {
			return errs.IMErrorCode_APP_FILE_NOOSS, nil
		}
		signedURL, err := fileConf.Minio.PreSignedURL(req.Ext, dir)
		if err != nil {
			return errs.IMErrorCode_APP_FILE_SIGNERR, nil
		}
		return errs.IMErrorCode_SUCCESS, &apimodels.QryFileCredResp{
			OssType: apimodels.OssType_Minio,
			PreSignResp: &apimodels.PreSignResp{
				Url: signedURL,
			},
		}
	case fileengine.ChannelAws:
		if fileConf.S3 == nil {
			return errs.IMErrorCode_APP_FILE_NOOSS, nil
		}

		signedURL, err := fileConf.S3.PreSignedURL(req.Ext, dir)
		if err != nil {
			return errs.IMErrorCode_APP_FILE_SIGNERR, nil
		}
		return errs.IMErrorCode_SUCCESS, &apimodels.QryFileCredResp{
			OssType: apimodels.OssType_S3,
			PreSignResp: &apimodels.PreSignResp{
				Url: signedURL,
			},
		}
	case fileengine.ChannelOss:
		if fileConf.Oss == nil {
			return errs.IMErrorCode_APP_FILE_NOOSS, nil
		}
		signedURL, err := fileConf.Oss.PreSignedURL(req.Ext, dir)
		if err != nil {
			return errs.IMErrorCode_APP_FILE_SIGNERR, nil
		}
		resp := fileConf.Oss.PostSign(req.Ext, dir)
		return errs.IMErrorCode_SUCCESS, &apimodels.QryFileCredResp{
			OssType: apimodels.OssType_Oss,
			PreSignResp: &apimodels.PreSignResp{
				Url:         signedURL,
				ObjKey:      resp.ObjKey,
				Policy:      resp.Policy,
				SignVersion: resp.SignVersion,
				Credential:  resp.Credential,
				Date:        resp.Date,
				Signature:   resp.Signature,
			},
		}
	default:
		return errs.IMErrorCode_APP_FILE_NOOSS, nil
	}
}

type FileConfItem struct {
	AppKey     string `json:"app_key,omitempty"`
	FileEngine string `json:"file_engine"`
	//QiniuConfig *QiniuFileConfig `json:"qiniu,omitempty"`

	QiNiu *fileengine.QiNiuStorage
	Oss   *fileengine.OssStorage
	Minio *fileengine.MinioStorage
	S3    *fileengine.S3Storage
}

var fileConfCache *caches.LruCache
var fileLock *sync.RWMutex
var notExistFileConf *FileConfItem

func init() {
	fileConfCache = caches.NewLruCacheWithAddReadTimeout(1000, nil, 5*time.Minute, 10*time.Minute)
	fileLock = &sync.RWMutex{}
	notExistFileConf = &FileConfItem{}
}

func GetFileConf(ctx context.Context, appkey string) *FileConfItem {
	if obj, exist := fileConfCache.Get(appkey); exist {
		return obj.(*FileConfItem)
	} else {
		fileLock.Lock()
		defer fileLock.Unlock()

		if obj, exist := fileConfCache.Get(appkey); exist {
			return obj.(*FileConfItem)
		} else { //load from db
			fileConf, err := loadFileConfFromDb(appkey)
			if err != nil {
				fileConf = notExistFileConf
			}
			fileConfCache.Add(appkey, fileConf)
			return fileConf
		}
	}
}

func loadFileConfFromDb(appkey string) (*FileConfItem, error) {
	dao := dbs.FileConfDao{}
	conf, err := dao.FindEnableFileConf(appkey)
	if err != nil {
		return nil, err
	}
	fileConf := &FileConfItem{
		AppKey:     appkey,
		FileEngine: conf.Channel,
	}
	var confData = make(map[string]interface{})
	_ = json.Unmarshal([]byte(conf.Conf), &confData)

	switch conf.Channel {
	case fileengine.ChannelQiNiu:
		c := utils.MapToStruct[fileengine.QiNiuConfig](confData)
		fileConf.QiNiu = fileengine.NewQiNiu(c)
	case fileengine.ChannelMinio:
		c := utils.MapToStruct[fileengine.MinioConfig](confData)
		fileConf.Minio = fileengine.NewMinio(c)
	case fileengine.ChannelOss:
		c := utils.MapToStruct[fileengine.OssConfig](confData)
		fileConf.Oss = fileengine.NewOss(c)
	case fileengine.ChannelAws:
		c := utils.MapToStruct[fileengine.S3Config](confData)
		fileConf.S3 = fileengine.NewS3Storage(fileengine.WithConf(c))
	}
	return fileConf, nil
}

func fileTypeToDir(fileType apimodels.FileType) string {
	switch fileType {
	case apimodels.FileType_Image:
		return "images"
	case apimodels.FileType_Video:
		return "videos"
	case apimodels.FileType_Audio:
		return "audios"
	case apimodels.FileType_File:
		return "files"
	case apimodels.FileType_Log:
		return "logs"
	default:
		return "files"
	}
}
