package lists

import (
	"context"
	"fmt"
	"strconv"

	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListEmbedCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListEmbedCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListEmbedCodeLogic {
	return &GetListEmbedCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListEmbedCodeLogic) GetListEmbedCode(req *types.GetEmbedCodeRequest) (resp *types.EmbedCodeResponse, err error) {
	listID, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid list ID: %w", err)
	}

	list, err := l.svcCtx.DB.GetEmailList(l.ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("list not found: %w", err)
	}

	baseURL := l.svcCtx.Config.App.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:9888"
	}

	formURL := fmt.Sprintf("%s/s/%s", baseURL, list.Slug)

	html := fmt.Sprintf(`<!-- Outlet.sh Subscribe Form for "%s" -->
<form action="%s" method="POST" style="max-width: 400px; font-family: system-ui, sans-serif;">
  <div style="margin-bottom: 12px;">
    <input type="email" name="email" placeholder="Email address" required
      style="width: 100%%; padding: 10px; border: 1px solid #ccc; border-radius: 4px; font-size: 14px;">
  </div>
  <div style="margin-bottom: 12px;">
    <input type="text" name="name" placeholder="Your name (optional)"
      style="width: 100%%; padding: 10px; border: 1px solid #ccc; border-radius: 4px; font-size: 14px;">
  </div>
  <button type="submit"
    style="width: 100%%; padding: 12px; background: #0066ff; color: white; border: none; border-radius: 4px; font-size: 14px; cursor: pointer;">
    Subscribe
  </button>
</form>`, list.Name, formURL)

	return &types.EmbedCodeResponse{
		Html:    html,
		ListId:  strconv.FormatInt(list.ID, 10),
		Slug:    list.Slug,
		BaseUrl: baseURL,
	}, nil
}
