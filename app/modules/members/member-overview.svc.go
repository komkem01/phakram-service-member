package members

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

type MemberOverviewServiceRequest struct {
	base.RequestPaginate
	MemberID uuid.UUID
}

type MemberOverviewPaginate struct {
	Addresses *base.ResponsePaginate `json:"addresses"`
	Banks     *base.ResponsePaginate `json:"banks"`
	Files     *base.ResponsePaginate `json:"files"`
	Payments  *base.ResponsePaginate `json:"payments"`
	Wishlist  *base.ResponsePaginate `json:"wishlist"`
	Orders    *base.ResponsePaginate `json:"orders"`
	Carts     *base.ResponsePaginate `json:"carts"`
}

type MemberOverviewServiceResponse struct {
	Member    *InfoServiceResponse      `json:"member"`
	Addresses []*ent.MemberAddressEntity `json:"addresses"`
	Banks     []*ent.MemberBankEntity    `json:"banks"`
	Files     []*ent.MemberFileEntity    `json:"files"`
	Payments  []*ent.MemberPaymentEntity `json:"payments"`
	Wishlist  []*ent.MemberWishlistEntity `json:"wishlist"`
	Orders    []*ent.OrderEntity         `json:"orders"`
	Carts     []*ent.CartEntity          `json:"carts"`
	Paginate  *MemberOverviewPaginate    `json:"paginate"`
}

func (s *Service) MemberOverviewService(ctx context.Context, req *MemberOverviewServiceRequest) (*MemberOverviewServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.overview.start`)

	member, err := s.InfoService(ctx, req.MemberID)
	if err != nil {
		return nil, err
	}

	addresses, addressesPage, err := s.address.ListMemberAddresses(ctx, &entitiesdto.ListMemberAddressesRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, err
	}

	banks, banksPage, err := s.bank.ListMemberBanks(ctx, &entitiesdto.ListMemberBanksRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, err
	}

	files, filesPage, err := s.file.ListMemberFiles(ctx, &entitiesdto.ListMemberFilesRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, err
	}

	payments, paymentsPage, err := s.payment.ListMemberPayments(ctx, &entitiesdto.ListMemberPaymentsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, err
	}

	wishlist, wishlistPage, err := s.wishlist.ListMemberWishlist(ctx, &entitiesdto.ListMemberWishlistRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, err
	}

	orders, ordersPage, err := s.order.ListOrders(ctx, &entitiesdto.ListOrdersRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, err
	}

	carts, cartsPage, err := s.cart.ListCarts(ctx, &entitiesdto.ListCartsRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        req.MemberID,
	})
	if err != nil {
		return nil, err
	}

	span.AddEvent(`members.svc.overview.success`)
	return &MemberOverviewServiceResponse{
		Member:    member,
		Addresses: addresses,
		Banks:     banks,
		Files:     files,
		Payments:  payments,
		Wishlist:  wishlist,
		Orders:    orders,
		Carts:     carts,
		Paginate: &MemberOverviewPaginate{
			Addresses: addressesPage,
			Banks:     banksPage,
			Files:     filesPage,
			Payments:  paymentsPage,
			Wishlist:  wishlistPage,
			Orders:    ordersPage,
			Carts:     cartsPage,
		},
	}, nil
}
