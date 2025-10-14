package services

import (
	"context"

	"github.com/lucky-lbc/commons/ctxs"
	"github.com/lucky-lbc/commons/errs"
	"github.com/lucky-lbc/commons/tools"
	apimodels "github.com/lucky-lbc/jugglechat-server/apis/models"
	"github.com/lucky-lbc/jugglechat-server/services/pbobjs"
	"github.com/lucky-lbc/jugglechat-server/storages"
	"github.com/lucky-lbc/jugglechat-server/storages/models"
)

func AddFeedback(ctx context.Context, feedback *apimodels.Feedback) errs.IMErrorCode {
	appkey := ctxs.GetAppKeyFromCtx(ctx)
	userId := ctxs.GetRequesterIdFromCtx(ctx)
	fbContent := &pbobjs.FeedbackContent{
		Text:   feedback.Text,
		Images: feedback.Images,
		Videos: feedback.Videos,
	}
	contentBs, _ := tools.PbMarshal(fbContent)
	storage := storages.NewFeedbackStorage()
	err := storage.Create(models.Feedback{
		AppKey:   appkey,
		UserId:   userId,
		Category: feedback.Category,
		Content:  contentBs,
	})
	if err != nil {
		return errs.IMErrorCode_APP_DEFAULT
	}
	return errs.IMErrorCode_SUCCESS
}
