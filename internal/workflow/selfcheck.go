package workflow

import "github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"

func (s *Service) SelfCheck() error {
	b, e := s.CreateBatch(CreateBatchRequest{"演示数据集", "管理员", "伦理-1", []string{"机构内部"}}, "selfcheck-create")
	if e != nil {
		return e
	}
	src, e := s.AddSource(b.BatchID, AddSourceRequest{"实验室", "匿名样本", "consent-demo", "restricted", "checksum-demo"}, b.Version, "selfcheck-source")
	if e != nil {
		return e
	}
	rv, e := s.Review(b.BatchID, src.SourceID, ReviewRequest{"审查员", model.DecisionApprove, nil, []string{"evidence-demo"}, b.Version + 1}, "selfcheck-review")
	if e != nil {
		return e
	}
	_ = rv
	snap, e := s.Freeze(b.BatchID, ActorRequest{"发布负责人", b.Version + 2}, "selfcheck-freeze")
	if e != nil {
		return e
	}
	_, e = s.Close(b.BatchID, ActorRequest{"发布负责人", b.Version + 3}, "selfcheck-close")
	if e != nil {
		return e
	}
	if snap.SourceCount != 1 {
		return model.Invalid("自检快照来源数错误")
	}
	if _, e = s.Export(b.BatchID); e != nil {
		return e
	}
	if e = s.Verify(b.BatchID); e != nil {
		return e
	}
	return nil
}
