package emailconfig

import (
	"context"
	"fmt"

	"outlet/internal/services/email"
	"outlet/internal/svc"
	"outlet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetectSESQuotaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetectSESQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetectSESQuotaLogic {
	return &DetectSESQuotaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetectSESQuotaLogic) DetectSESQuota(req *types.DetectSESQuotaRequest) (resp *types.SESQuotaResponse, err error) {
	accessKey := req.AWSAccessKey
	secretKey := req.AWSSecretKey
	region := req.AWSRegion

	// If credentials not provided, try to use org's stored credentials
	if accessKey == "" || secretKey == "" {
		config, err := email.GetOrgEmailConfig(l.ctx, l.svcCtx.DB, req.OrgId)
		if err != nil {
			return nil, fmt.Errorf("failed to get org config: %w", err)
		}
		if config.HasOwnAWSCredentials() {
			accessKey = config.AWSAccessKey
			secretKey = config.AWSSecretKey
			if region == "" {
				region = config.AWSRegion
			}
		}
	}

	// Query AWS SES for quota
	quota, err := email.GetSESQuota(l.ctx, region, accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get SES quota: %w", err)
	}

	return &types.SESQuotaResponse{
		Max24HourSend:   quota.Max24HourSend,
		MaxSendRate:     quota.MaxSendRate,
		SentLast24Hours: quota.SentLast24Hours,
		RemainingQuota:  quota.RemainingQuota,
	}, nil
}
