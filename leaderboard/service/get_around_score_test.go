package service_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/podium/leaderboard/database"
	"github.com/topfreegames/podium/leaderboard/model"
	"github.com/topfreegames/podium/leaderboard/service"
)

var _ = Describe("Service GetAroundScore", func() {
	var ctrl *gomock.Controller
	var mock *database.MockDatabase
	var svc *service.Service

	var leaderboard string = "leaderboardTest"
	var totalMembers int = 10
	var pageSize int = 3
	var score int64 = 1
	var member string = "member"
	var order string = "asc"

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mock = database.NewMockDatabase(ctrl)

		svc = &service.Service{mock}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("Should return members slice if all is OK", func() {
		rank := 6
		start := 5
		stop := 7

		membersDatabaseReturn := []*database.Member{
			&database.Member{
				Member: "member1",
				Score:  float64(1),
				Rank:   5,
			},
			&database.Member{
				Member: "member2",
				Score:  float64(2),
				Rank:   6,
			},
			&database.Member{
				Member: "member3",
				Score:  float64(3),
				Rank:   7,
			},
		}

		membersReturn := []*model.Member{
			&model.Member{
				PublicID: "member1",
				Score:    1,
				Rank:     6,
			},
			&model.Member{
				PublicID: "member2",
				Score:    2,
				Rank:     7,
			},
			&model.Member{
				PublicID: "member3",
				Score:    3,
				Rank:     8,
			},
		}

		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{member}, nil)
		mock.EXPECT().GetRank(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(member), gomock.Eq(order)).Return(rank, nil)
		mock.EXPECT().GetTotalMembers(gomock.Any(), gomock.Eq(leaderboard)).Return(totalMembers, nil)
		mock.EXPECT().GetOrderedMembers(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(start), gomock.Eq(stop), gomock.Eq(order)).Return(membersDatabaseReturn, nil)

		membersFromService, err := svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
		Expect(err).NotTo(HaveOccurred())

		Expect(membersFromService).To(Equal(membersReturn))
	})

	It("Should return member not found if GetMemberIDsWithScoreInsideRange return no member", func() {
		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{}, nil)

		_, err := svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
		Expect(err).To(Equal(service.NewMemberNotFoundError(leaderboard, fmt.Sprint(score))))
	})

	It("Should return error if GetMemberIDsWithScoreInsideRange return in error", func() {
		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return(nil, fmt.Errorf("database error"))

		_, err := svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
		Expect(err).To(Equal(service.NewGeneralError("get around score", "database error")))
	})

	It("Should return error if getRank return member not found", func() {
		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{member}, nil)
		mock.EXPECT().GetRank(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(member), gomock.Eq(order)).Return(-1, database.NewMemberNotFoundError(leaderboard, member))

		_, err := svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
		Expect(err).To(Equal(service.NewMemberNotFoundError(leaderboard, member)))
	})

	It("Should return error if getRank return in error", func() {
		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{member}, nil)
		mock.EXPECT().GetRank(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(member), gomock.Eq(order)).Return(-1, fmt.Errorf("database error"))

		_, err := svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
		Expect(err).To(Equal(service.NewGeneralError("get around score", "database error")))
	})

	It("Should return error if TotalMembers return in error", func() {
		rank := 6

		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{member}, nil)
		mock.EXPECT().GetRank(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(member), gomock.Eq(order)).Return(rank, nil)
		mock.EXPECT().GetTotalMembers(gomock.Any(), gomock.Eq(leaderboard)).Return(-1, fmt.Errorf("database error"))

		_, err := svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
		Expect(err).To(Equal(service.NewGeneralError("get around score", "database error")))
	})

	It("Should return error if GetOrderedMembers return in error", func() {
		rank := 6
		start := 5
		stop := 7

		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{member}, nil)
		mock.EXPECT().GetRank(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(member), gomock.Eq(order)).Return(rank, nil)
		mock.EXPECT().GetTotalMembers(gomock.Any(), gomock.Eq(leaderboard)).Return(totalMembers, nil)
		mock.EXPECT().GetOrderedMembers(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(start), gomock.Eq(stop), gomock.Eq(order)).Return(nil, fmt.Errorf("database error"))

		_, err := svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
		Expect(err).To(Equal(service.NewGeneralError("get around score", "database error")))
	})

	It("Should ask for first 3 members if user is the first one", func() {
		rank := 0
		start := 0
		stop := 2

		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{member}, nil)

		mock.EXPECT().GetRank(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(member), gomock.Eq(order)).Return(rank, nil)
		mock.EXPECT().GetTotalMembers(gomock.Any(), gomock.Eq(leaderboard)).Return(totalMembers, nil)

		//this is the assertation relevant to this test
		mock.EXPECT().GetOrderedMembers(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(start), gomock.Eq(stop), gomock.Eq(order)).Times(1).Return(nil, fmt.Errorf("database error"))

		svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
	})

	It("Should ask for last members if user is the last one", func() {
		rank := 10
		start := 8
		stop := 10

		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{member}, nil)

		mock.EXPECT().GetRank(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(member), gomock.Eq(order)).Return(rank, nil)
		mock.EXPECT().GetTotalMembers(gomock.Any(), gomock.Eq(leaderboard)).Return(totalMembers, nil)

		//this is the assertation relevant to this test
		mock.EXPECT().GetOrderedMembers(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(start), gomock.Eq(stop), gomock.Eq(order)).Times(1).Return(nil, fmt.Errorf("database error"))

		svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
	})

	It("Should ask for all members if totalMembers is less than pageSize", func() {
		var rank int = 0
		var totalMembers int = 2
		start := 0
		stop := 2

		mock.EXPECT().GetMemberIDsWithScoreInsideRange(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq("-inf"), gomock.Eq(fmt.Sprint(score)), gomock.Eq(0), gomock.Eq(1)).Return([]string{member}, nil)
		mock.EXPECT().GetRank(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(member), gomock.Eq(order)).Return(rank, nil)
		mock.EXPECT().GetTotalMembers(gomock.Any(), gomock.Eq(leaderboard)).Return(totalMembers, nil)

		//this is the assertation relevant to this test
		mock.EXPECT().GetOrderedMembers(gomock.Any(), gomock.Eq(leaderboard), gomock.Eq(start), gomock.Eq(stop), gomock.Eq(order)).Times(1).Return(nil, fmt.Errorf("database error"))

		svc.GetAroundScore(context.Background(), leaderboard, pageSize, score, order)
	})
})
