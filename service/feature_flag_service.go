package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/status"

	"github.com/rafaelmaier/featureflags/domain"
	pb "github.com/rafaelmaier/featureflags/proto/v1"
	"github.com/rafaelmaier/featureflags/repository"
	"google.golang.org/grpc/codes"
)

type FeatureFlagService struct {
	pb.UnimplementedFeatureAdminServiceServer

	repo repository.FeatureFlagRepository
}

func NewFeatureFlagService(repo repository.FeatureFlagRepository) *FeatureFlagService{
	return &FeatureFlagService{
		repo: repo,
	}
}

func (s *FeatureFlagService) CreateFeature(ctx context.Context, req *pb.CreateFeatureRequest) (*pb.FeatureResponse, error){
	if req.Feature == nil {
		return nil, status.Error(codes.InvalidArgument, "feature cannot be nil")
	}
	
	domainFlag, err := protoToDomain(req.Feature)

	if err != nil{
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.repo.Create(ctx, domainFlag); err != nil {
		return nil, status.Error(codes.Internal, "failed to create feature: "+err.Error())
	}

	return &pb.FeatureResponse{
		Feature: domainToProto(domainFlag),
	}, nil
}

func (s *FeatureFlagService) GetFeature(ctx context.Context, req *pb.GetFeatureRequest) (*pb.FeatureResponse, error){
	if req.Key == ""{
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}
	domainFlag, err := s.repo.Get(ctx, req.Key)
	if err != nil {
		return nil, status.Error(codes.NotFound, "feature not found: " +err.Error())
	}

	return &pb.FeatureResponse{
		Feature: domainToProto(domainFlag),
	}, nil
}

func (s *FeatureFlagService) UpdateFeature(ctx context.Context, req *pb.UpdateFeatureRequest) (*pb.FeatureResponse, error){
	if req.Feature == nil{
		return nil, status.Error(codes.InvalidArgument, "Feature cannot be nil")
	}
	domainFlag, err := protoToDomain(req.Feature)
	if err != nil{
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := domainFlag.Validate(); err != nil{
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.repo.Update(ctx, domainFlag); err != nil{
		return nil, status.Error(codes.Internal, "failed to update feature: "+err.Error())
	}
	return &pb.FeatureResponse{
		Feature: domainToProto(domainFlag),
	}, nil
}

func (s *FeatureFlagService) DeleteFeature(ctx context.Context, req *pb.DeleteFeatureRequest) (*pb.DeleteFeatureResponse, error){
	if req.Key == ""{
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}
	if err := s.repo.Delete(ctx, req.Key); err != nil{
		return nil, status.Error(codes.Internal, "failed to delete features: "+err.Error())
	}
	return &pb.DeleteFeatureResponse{
		Success: true,
	}, nil
}

func (s *FeatureFlagService) ListFeatures(ctx context.Context, req *pb.ListFeaturesRequest) (*pb.ListFeaturesResponse, error){
	domainFlags, err := s.repo.List(ctx)
	if err != nil{
		return nil, status.Error(codes.Internal, "error listing feature")
	}
	protoFlags := make([]*pb.FeatureFlag, len(domainFlags))
	for i, domainFlag := range domainFlags{
		protoFlags[i] = domainToProto(domainFlag)
	}
	return &pb.ListFeaturesResponse{
		Features: protoFlags,
	}, nil
}

func protoToDomain(proto *pb.FeatureFlag) (*domain.FeatureFlag, error){
	if proto == nil{
		return nil, errors.New("proto feature flag cannot be nil")
	}

	domainRules := make([]domain.Rule, len(proto.Rules))
	for i, protoRule := range proto.Rules {
		domainRule, err := protoRuleToDomain(protoRule)
		if err != nil {
			return nil, err
		}
		domainRules[i] = domainRule
	}
	return &domain.FeatureFlag{
		Key: proto.Key,
		Enabled: proto.Enabled,
		Rules: domainRules,
		Version: proto.Version,
	}, nil
}

func protoRuleToDomain(proto *pb.Rule) (domain.Rule, error){
	if proto == nil{
		return domain.Rule{}, errors.New("proto rule cannot be nil")
	}

	switch rule := proto.Rule.(type){
	case *pb.Rule_Percentage:
		return domain.Rule{
			Evaluator: &domain.PercentageRule{Percentage: rule.Percentage.Percentage},
		}, nil
	case *pb.Rule_UserId:
		return domain.Rule{
			Evaluator: &domain.UserIDRule{ UserIDs: rule.UserId.UserIds},
		}, nil
	default:
			return domain.Rule{}, errors.New("unknown rule type")
	}
}

func domainToProto(d *domain.FeatureFlag) *pb.FeatureFlag{
	if d == nil {
		return nil
	}

	protoRules := make([]*pb.Rule, len(d.Rules))
	for i, domainRule := range d.Rules{
		protoRules[i] = domainRuleToProto(domainRule)
	}
	return &pb.FeatureFlag{
		Key: d.Key,
		Enabled: d.Enabled,
		Rules: protoRules,
		Version: d.Version,
	}
}

func domainRuleToProto(d domain.Rule) *pb.Rule{
	if d.Evaluator == nil{
		return nil
	}
	switch eval := d.Evaluator.(type){
	case *domain.PercentageRule:
		return &pb.Rule{
			Rule: &pb.Rule_Percentage{
				Percentage: &pb.PercentageRule{
					Percentage: eval.Percentage,
				},
			},
		}
	case *domain.UserIDRule:
		return &pb.Rule{
			Rule: &pb.Rule_UserId{
				UserId: &pb.UserIDRule{
					UserIds: eval.UserIDs,
				},
			},
		}
	default:
		return nil
	}
}
