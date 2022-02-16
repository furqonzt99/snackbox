package partner

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/furqonzt99/snackbox/constants"
	"github.com/furqonzt99/snackbox/delivery/common"
	"github.com/furqonzt99/snackbox/delivery/controllers/product"
	"github.com/furqonzt99/snackbox/delivery/controllers/rating"
	"github.com/furqonzt99/snackbox/delivery/middlewares"
	"github.com/furqonzt99/snackbox/helper"
	"github.com/furqonzt99/snackbox/models"
	"github.com/furqonzt99/snackbox/repositories/partner"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/labstack/echo/v4"
	"github.com/leekchan/accounting"
)

type PartnerController struct {
	Repo partner.PartnerInterface
}

func NewPartnerController(partner partner.PartnerInterface) *PartnerController {
	return &PartnerController{partner}
}

func (p PartnerController) ApplyPartner() echo.HandlerFunc {
	return func(c echo.Context) error {
		userJwt, _ := middlewares.ExtractTokenUser(c)

		partnerReq := PartnerUserRequestFormat{}
		c.Bind(&partnerReq)

		if err := c.Validate(partnerReq); err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		var res models.Partner

		user, err := p.Repo.FindUserId(userJwt.UserID)

		var partnerDocument string
		if user.LegalDocument != "" {
			partnerDocument = fmt.Sprintf(constants.LINK_TEMPLATE, constants.S3_BUCKET, constants.S3_REGION, user.LegalDocument)
		}

		if err != nil {
			partnerData := models.Partner{
				UserID:        uint(userJwt.UserID),
				BussinessName: partnerReq.BussinessName,
				Description:   partnerReq.Description,
				Latitude:      partnerReq.Latitude,
				Longtitude:    partnerReq.Longtitude,
				Address:       partnerReq.Address,
				City:          partnerReq.City,
			}

			res, _ = p.Repo.ApplyPartner(partnerData)

			responseFormat := PartnerResponse{
				BussinessName: res.BussinessName,
				Description:   res.Description,
				Latitude:      res.Latitude,
				Longtitude:    res.Longtitude,
				Address:       res.Address,
				City:          res.City,
				LegalDocument: partnerDocument,
				Status:        res.Status,
			}

			return c.JSON(http.StatusOK, common.SuccessResponse(responseFormat))

		}

		if user.Status == "reject" {
			user.BussinessName = partnerReq.BussinessName
			user.Description = partnerReq.Description
			user.Latitude = partnerReq.Latitude
			user.Longtitude = partnerReq.Longtitude
			user.Address = partnerReq.Address
			user.City = partnerReq.City
			user.Status = "pending"

			res, _ = p.Repo.ApplyPartner(user)

			responseFormat := PartnerResponse{
				BussinessName: res.BussinessName,
				Description:   res.Description,
				Latitude:      res.Latitude,
				Longtitude:    res.Longtitude,
				Address:       res.Address,
				City:          res.City,
				LegalDocument: partnerDocument,
				Status:        res.Status,
			}
			return c.JSON(http.StatusOK, common.SuccessResponse(responseFormat))
		}

		responseFormat := PartnerResponse{
			BussinessName: user.BussinessName,
			Description:   user.Description,
			Latitude:      user.Latitude,
			Longtitude:    user.Longtitude,
			Address:       user.Address,
			City:          user.City,
			LegalDocument: partnerDocument,
			Status:        user.Status,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responseFormat))
	}
}

func (p PartnerController) GetAllPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		res, err := p.Repo.GetAllPartner()
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		responseFormat := []GetPartnerResponse{}

		for _, data := range res {
			var partnerDocument string
			if data.LegalDocument != "" {
				partnerDocument = fmt.Sprintf(constants.LINK_TEMPLATE, constants.S3_BUCKET, constants.S3_REGION, data.LegalDocument)
			}
			responseFormat = append(responseFormat, GetPartnerResponse{
				ID:            int(data.ID),
				BussinessName: data.BussinessName,
				Description:   data.Description,
				Latitude:      data.Latitude,
				Longtitude:    data.Longtitude,
				Address:       data.Address,
				City:          data.City,
				LegalDocument: partnerDocument,
				Status:        data.Status,
			})
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responseFormat))
	}
}

func (p PartnerController) AcceptPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		partnerId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		res, _ := p.Repo.FindPartnerId(partnerId)

		if res.Status == "reject" {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		err = p.Repo.AcceptPartner(res)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (p PartnerController) RejectPartner() echo.HandlerFunc {
	return func(c echo.Context) error {

		partnerId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		res, _ := p.Repo.FindPartnerId(partnerId)

		if res.Status == "active" {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		err = p.Repo.RejectPartner(res)
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}

func (p PartnerController) GetPartnerProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		partnerId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		partner, _ := p.Repo.GetPartner(partnerId)

		productItems := []product.ProductResponse{}
		for _, item := range partner.Products {
			productItems = append(productItems, product.ProductResponse{
				Title:       item.Title,
				Type:        item.Type,
				Description: item.Description,
				Price:       item.Price,
			})
		}

		response := GetPartnerProductResponse{
			ID:            int(partner.ID),
			BussinessName: partner.BussinessName,
			Description:   partner.Description,
			Latitude:      partner.Latitude,
			Longtitude:    partner.Longtitude,
			Address:       partner.Address,
			City:          partner.City,
			Rating:        helper.CalculateRating(partner.Ratings),
			Products:      productItems,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}

}

func (p PartnerController) GetPartnerRating() echo.HandlerFunc {
	return func(c echo.Context) error {

		partnerId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		partner, _ := p.Repo.GetPartner(partnerId)

		ratingItems := []rating.RatingResponse{}
		for _, item := range partner.Ratings {
			ratingItems = append(ratingItems, rating.RatingResponse{
				PartnerID: partnerId,
				UserID:    item.Rating,
				Username:  item.User.Name,
				Rating:    item.Rating,
				Comment:   item.Comment,
			})
		}

		response := GetPartnerRatingResponse{
			ID:            int(partner.ID),
			BussinessName: partner.BussinessName,
			Description:   partner.Description,
			Latitude:      partner.Latitude,
			Longtitude:    partner.Longtitude,
			Address:       partner.Address,
			City:          partner.City,
			Rating:        helper.CalculateRating(partner.Ratings),
			Ratings:       ratingItems,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}

}

func (pc PartnerController) Upload(c echo.Context) error {
	var requestUpload UploadDocumentRequest

	c.Bind(&requestUpload)

	user, _ := middlewares.ExtractTokenUser(c)

	partner, err := pc.Repo.FindUserId(user.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	if partner.Status == "active" || partner.Status == "pending" {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	file, err := c.FormFile("legal_document")
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	defer src.Close()

	head := make([]byte, 261)
	src.Read(head)

	kind, _ := filetype.Match(head)

	const ACCEPTED_FILE_TYPE = "pdf"

	if kind.Extension != ACCEPTED_FILE_TYPE {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "extension must .pdf"))
	}

	prefix := "legal-documents/"

	fileID := strings.ReplaceAll(uuid.New().String(), "-", "")
	file.Filename = fmt.Sprint(prefix, fileID, ".", kind.Extension)

	if partner.LegalDocument != "" {
		if err := helper.GetObjectS3(partner.LegalDocument); err == nil {
			_ = helper.DeleteObjectS3(partner.LegalDocument)
		}
	}

	if err := helper.UploadObjectS3(file.Filename, src); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	const PENDING_STATUS = "pending"
	partnerData := models.Partner{
		LegalDocument: file.Filename,
		Status:        PENDING_STATUS,
	}

	_, err = pc.Repo.UploadDocument(int(partner.ID), partnerData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (p PartnerController) Report() echo.HandlerFunc {
	return func(c echo.Context) error {

		userJwt, _ := middlewares.ExtractTokenUser(c)
		transactions, err := p.Repo.Report(userJwt.PartnerID)
		if err != nil {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}
		////////////////////////////////////////

		contents := [][]string{}
		ac := accounting.Accounting{Symbol: "Rp", Precision: 0}
		for i := 0; i < len(transactions); i++ {
			temp := []string{}

			date := fmt.Sprint(transactions[i].CreatedAt)
			invoice := transactions[i].InvoiceID
			totalPrice := ac.FormatMoney(transactions[i].TotalPrice)
			quantity := strconv.Itoa(transactions[i].Quantity)
			paymentChannel := transactions[i].PaymentChannel
			status := transactions[i].Status
			var items string
			for _, item := range transactions[i].Products {
				items += item.Title + ", "
			}
			product := fmt.Sprint(items)[:len(items)-2]

			temp = append(temp, date[:16])
			temp = append(temp, invoice)
			temp = append(temp, totalPrice)
			temp = append(temp, product)
			temp = append(temp, quantity)
			temp = append(temp, paymentChannel)
			temp = append(temp, status)

			contents = append(contents, temp)
		}
		m := pdf.NewMaroto(consts.Landscape, consts.A4)
		m.SetPageMargins(10, 10, 10)

		m.RegisterHeader(func() {
			m.Row(20, func() {
				m.Col(12, func() {
					m.Text("Tabel List Transaction", props.Text{
						Top:    2,
						Size:   14,
						Align:  consts.Center,
						Family: consts.Arial,
					})
				})
			})
		})

		m.SetBackgroundColor(color.NewWhite())

		tableHeadings := []string{"Transaction Date", "Invoice ID", "Total Transaction", "Product", "Quantity", "Payment", "Status"}

		m.TableList(tableHeadings, contents, props.TableList{
			HeaderProp: props.TableListContent{
				Size:      12,
				Style:     consts.Bold,
				GridSizes: []uint{3, 3, 2, 1, 1, 1, 1},
			},

			ContentProp: props.TableListContent{
				Size:      10,
				GridSizes: []uint{3, 3, 2, 1, 1, 1, 1},
			},
			Align:                consts.Center,
			AlternatedBackground: &color.Color{Red: 230, Blue: 230, Green: 230},
			HeaderContentSpace:   2,
			Line:                 true,
		})

		prefix := "reports/"

		fileID := strings.ReplaceAll(uuid.New().String(), "-", "")
		filename := fmt.Sprint(prefix, fileID, ".", "pdf")

		path := fmt.Sprint("./", filename)

		m.OutputFileAndClose(path)

		file, err := os.OpenFile(path, os.O_RDWR, 0755)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, common.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
		defer os.Remove(path)

		if err := helper.UploadObjectS3(filename, file); err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
		}

		reportLink := fmt.Sprintf(constants.LINK_TEMPLATE, constants.S3_BUCKET, constants.S3_REGION, filename)

		responses := ReportResponse{
			ReportLink: reportLink,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responses))
	}

}
